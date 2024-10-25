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

package data

import (
	"bytes"
	"crypto/sha256"
	"encoding/json"
	"fmt"
)

type AUTH_TYPE int

const (
	Free AUTH_TYPE = 1
	PIN  AUTH_TYPE = 2
	BIO  AUTH_TYPE = 4
)

const (
	RSA KEY_TYPE = "RsaVerificationKey2018"
	K1  KEY_TYPE = "Secp256k1VerificationKey2018"
	R1  KEY_TYPE = "Secp256r1VerificationKey2018"
)

type VerificationMethod struct {
	Id                 DID_KEY_ID `validate:"required" json:"id"`
	Type               KEY_TYPE   `validate:"required" json:"type"`
	Controller         string     `validate:"required" json:"controller"`
	PublicKeyMultibase Multibase  `validate:"required" json:"publicKeyMultibase"`
	AuthType           AUTH_TYPE  `validate:"required" json:"authType"`
}

func (a *VerificationMethod) IsEqual(b *VerificationMethod) (bool, error) {
	hashData, err := a.ToHash()
	if err != nil {
		return false, fmt.Errorf("failed to hash the first verificationMethod : %w", err)
	}

	comparedHashData, err := b.ToHash()
	if err != nil {
		return false, fmt.Errorf("failed to hash the second verificationMethod : %w", err)
	}

	return bytes.Equal(hashData, comparedHashData), nil
}

func (v *VerificationMethod) ToJson() ([]byte, error) {
	jsonData, err := json.Marshal(v)
	if err != nil {
		return nil, err
	}

	return jsonData, nil
}

func (v *VerificationMethod) ToObject(data []byte) error {
	return json.Unmarshal(data, v)
}

func (v *VerificationMethod) ToHash() ([]byte, error) {
	hash := sha256.New()

	jsonData, err := v.ToJson()
	if err != nil {
		return nil, err
	}

	hash.Write(jsonData)
	md := hash.Sum(nil)

	return md, nil
}
