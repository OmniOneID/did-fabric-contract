// Copyright 2024 OmniOne.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package utility

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"did-fabric-contract/chaincode/data"
	"encoding/base32"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"math/big"

	"github.com/btcsuite/btcutil/base58"
)

func CreateNewEcdsaPrivateKey() (*ecdsa.PrivateKey, error) {
	return ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
}

func Sign(data []byte, privateKey *ecdsa.PrivateKey) ([]byte, error) {
	digest := sha256.Sum256(data)

	r, s, err := ecdsa.Sign(rand.Reader, privateKey, digest[:])
	if err != nil {
		return nil, err
	}

	params := privateKey.Curve.Params()
	curveOrderByteSize := params.P.BitLen() / 8

	rBytes := r.Bytes()
	sBytes := s.Bytes()

	signature := make([]byte, curveOrderByteSize*2)
	copy(signature[curveOrderByteSize-len(rBytes):], rBytes)
	copy(signature[curveOrderByteSize*2-len(sBytes):], sBytes)

	return signature, nil
}

func Verify(data, signature []byte, publicKey *ecdsa.PublicKey) bool {
	digest := sha256.Sum256(data)

	curveOrderByteSize := publicKey.Curve.Params().P.BitLen() / 8

	r, s := new(big.Int), new(big.Int)
	r.SetBytes(signature[:curveOrderByteSize])
	s.SetBytes(signature[curveOrderByteSize:])

	return ecdsa.Verify(publicKey, digest[:], r, s)
}

func DecompressPublicKeyFromString(publicKeyHexString string) (*ecdsa.PublicKey, error) {
	publicKeyBytes, err := hex.DecodeString(publicKeyHexString)
	if err != nil {
		return nil, err
	}

	return DecompressPublicKey(publicKeyBytes)
}

func DecompressPublicKey(publicKey []byte) (*ecdsa.PublicKey, error) {
	if len(publicKey) != 33 {
		return nil, fmt.Errorf("invalid public key size: %d, expected 33", len(publicKey))
	}

	if publicKey[0] != 0x02 && publicKey[0] != 0x03 {
		return nil, fmt.Errorf("invalid public key format: first byte must be 0x02 or 0x03")
	}

	x := new(big.Int).SetBytes(publicKey[1:])
	xx := new(big.Int).Mul(x, x)
	xxx := new(big.Int).Mul(xx, x)

	ax := new(big.Int).Mul(big.NewInt(3), x)

	yy := new(big.Int).Sub(xxx, ax)
	yy.Add(yy, elliptic.P256().Params().B)

	y1 := new(big.Int).ModSqrt(yy, elliptic.P256().Params().P)
	if y1 == nil {
		return nil, fmt.Errorf("can not recovery public key")
	}

	y2 := new(big.Int).Neg(y1)
	y2.Mod(y2, elliptic.P256().Params().P)

	y := new(big.Int)
	if publicKey[0] == 0x02 {
		if y1.Bit(0) == 0 {
			y = y1
		} else {
			y = y2
		}
	} else {
		if y1.Bit(0) == 1 {
			y = y1
		} else {
			y = y2
		}
	}

	return &ecdsa.PublicKey{X: x, Y: y, Curve: elliptic.P256()}, nil
}

func CompressPublicKey(publicKey *ecdsa.PublicKey) []byte {
	params := publicKey.Curve.Params()
	curveOrderByteSize := params.P.BitLen() / 8

	xBytes := publicKey.X.Bytes()
	signature := make([]byte, curveOrderByteSize+1)

	if publicKey.Y.Bit(0) == 1 {
		signature[0] = 0x03
	} else {
		signature[0] = 0x02
	}

	copy(signature[1+curveOrderByteSize-len(xBytes):], xBytes)
	return signature
}

func RecoveryEcdsa(data []byte, rawSign []byte) (*ecdsa.PublicKey, *ecdsa.PublicKey, error) {
	signLength := len(rawSign)
	halfSignLength := signLength / 2

	r := new(big.Int).SetBytes(rawSign[:halfSignLength])
	s := new(big.Int).SetBytes(rawSign[halfSignLength:])
	curve := elliptic.P256().Params()

	expY := new(big.Int).Sub(curve.N, big.NewInt(2))
	rInv := new(big.Int).Exp(r, expY, curve.N)
	z := new(big.Int).SetBytes(data)

	xx := new(big.Int).Mul(r, r)
	xxx := xx.Mul(xx, r)

	ax := new(big.Int).Mul(big.NewInt(3), r)

	yy := new(big.Int).Sub(xxx, ax)
	yy.Add(yy, elliptic.P256().Params().B)

	y1 := new(big.Int).ModSqrt(yy, curve.P)
	if y1 == nil {
		return nil, nil, fmt.Errorf("cannot recover public key")
	}

	y2 := new(big.Int).Neg(y1)
	y2.Mod(y2, curve.P)

	rInvBytes := rInv.Bytes()

	p1, p2 := elliptic.P256().ScalarMult(r, y1, s.Bytes())
	p3, p4 := elliptic.P256().ScalarBaseMult(z.Bytes())

	p5 := new(big.Int).Neg(p4)
	p5.Mod(p5, curve.P)

	q1, q2 := elliptic.P256().Add(p1, p2, p3, p5)
	q3, q4 := elliptic.P256().ScalarMult(q1, q2, rInvBytes)

	n1, n2 := elliptic.P256().ScalarMult(r, y2, s.Bytes())
	n3, n4 := elliptic.P256().ScalarBaseMult(z.Bytes())

	n5 := new(big.Int).Neg(n4)
	n5.Mod(n5, curve.P)

	q5, q6 := elliptic.P256().Add(n1, n2, n3, n5)
	q7, q8 := elliptic.P256().ScalarMult(q5, q6, rInvBytes)

	key1 := ecdsa.PublicKey{Curve: elliptic.P256(), X: q3, Y: q4}
	key2 := ecdsa.PublicKey{Curve: elliptic.P256(), X: q7, Y: q8}

	return &key1, &key2, nil
}

func ComparePublicKey(key1, key2 *ecdsa.PublicKey) bool {
	x := key1.X.Cmp(key2.X)
	y := key1.Y.Cmp(key2.Y)

	if x == 0 && y == 0 {
		return true
	}

	return false
}

func DecodeMultibase(input data.Multibase) ([]byte, error) {
	multibaseData := string(input)
	if len(multibaseData) < 1 {
		return nil, fmt.Errorf("input data is too short")
	}

	encodingType := multibaseData[0]
	encodingData := multibaseData[1:]

	switch encodingType {
	case 'f', 'F':
		return hex.DecodeString(encodingData)
	case 'z':
		return base58.Decode(encodingData), nil
	case 'b':
		return base32.StdEncoding.DecodeString(encodingData)
	case 'm':
		return base64.RawStdEncoding.DecodeString(encodingData)
	case 'u':
		return base64.URLEncoding.DecodeString(encodingData)
	default:
		return nil, fmt.Errorf("unsupported encoding : %c", encodingType)
	}
}
