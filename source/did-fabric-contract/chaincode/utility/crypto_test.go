// Copyright 2024 Raonsecure

package utility_test

import (
	"did-fabric-contract/chaincode/utility"
	"encoding/hex"
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestVerify(t *testing.T) {
	plainText := "Test"

	assertVerify(t,
		"023888e60a7bf9672b09fad449436939a1b9dff53f18bd9a2383a7d18e3e15c861",
		"20c6f54810e3667e349c67b625d50f3344253e38d5e4876cffc2c54e8635e4527d5216507e7b29c3723263d9dfb16ee292ade9908c5519d60756443255d87c07f3",
		plainText)

	assertVerify(t,
		"035406dba5e8a29dc2d05b42c08f925b95d972786d95f91e86e7d2c0f51c6cef9b",
		"1f8578bd7f8535d3a7cccade7670ee0947e7fc778e82690928c56ff4c8ddfbf0ab739563c6f01c6a1e752c1a824e621bcab94e62a708fe0e88dab20c6646919706",
		plainText)

	assertVerify(t,
		"02eb2044610ba2c3960f9d91196bdb1c5498beebcb04983ed6ebe2329c3612907c",
		"1f4c05e607165e08274a690e42e76e6dd602ff0d97ffb5e1d034aab00e1b4662207229205edeba142f74ef9cfb25071d793ebd70294388d66fbc030a1b41216a98",
		plainText)

	assertVerify(t,
		"03380b593f680637362656171bcb755dc51f86e660710829420c4d2e71c44e95f3",
		"20f57f655090c2ec94bc9b12448bcdcb9438006530d1f8bd33aa928ef71b7d4e1463cccc8fe22497e9f869146df813c883e3017a3cdfb711f05c0c209b38b20a47",
		plainText)

	assertVerify(t,
		"039bbce9612e74b5e294edc2eba256c1a97196ebc356ae55b4786b92b96aee87a3",
		"1f94efe30a62347fb0f9521159f4f394b3fcc39f34af74657127c9482c31b7704d44ed5590fb8ed13340d5de9a0a31094396d799348164a2bcc81c032c8a6a68f1",
		plainText)
}

func assertVerify(t *testing.T, publicKey, signature, plainTest string) {

	publicKeyHexString, _ := hex.DecodeString(publicKey)
	signatureHexString, _ := hex.DecodeString(signature)
	data := []byte(plainTest)

	publicKeyBytes, err := utility.DecompressPublicKey(publicKeyHexString)
	if err != nil {
		log.Fatal(err)
	}

	assert.True(t,
		utility.Verify(data, signatureHexString[1:], publicKeyBytes),
		"signature verify success",
	)
}
