// Copyright 2024 Raonsecure

package chaincode_test

import (
	"crypto/ecdsa"
	"crypto/x509"
	"did-fabric-contract/chaincode"
	"did-fabric-contract/chaincode/data"
	. "did-fabric-contract/chaincode/utility"
	b64 "encoding/base64"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"log"
	"os"
	"testing"
	"time"

	"github.com/btcsuite/btcutil/base58"
	testcc "github.com/hyperledger-labs/cckit/testing"
	expectcc "github.com/hyperledger-labs/cckit/testing/expect"
	"github.com/hyperledger/fabric-protos-go/peer"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestDocument(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, `OPEN DID TEST CODE IS START`)
}

func EncodingStrByMultiBase(data string) string {
	sEnc := b64.RawStdEncoding.EncodeToString([]byte(data))
	res := "m" + sEnc
	return res
}

var _ = Describe("document function test is start", func() {

	cc := testcc.NewMockStub(`OPEN DID CC`, chaincode.NewOpenDIDCC())

	tasDocument := getTasDocument()
	userDocument := getUserdocument()

	Describe("while document registration", func() {
		tasDocumentDto := makeInvokedDocument(tasDocument, getPrivateKeyInPemFile("key/tas/private-key.pem"))
		userDocumentDto := makeInvokedDocument(userDocument, getPrivateKeyInPemFile("key/tas/private-key.pem"))

		It(`should return success`, func() {
			log.Printf("request document : %v", tasDocumentDto)
			queryResponse := expectcc.ResponseOk(cc.Invoke("document_registDidDoc", tasDocumentDto, data.TAS))
			response := expectcc.PayloadIs(queryResponse, &peer.Response{}).(peer.Response)
			expectcc.ResponseOk(response)
			fmt.Printf("response payload : %s\n", response.Payload)
		})

		It(`should return success while user document registration`, func() {
			userDocumentDtoJson, _ := json.Marshal(userDocumentDto)
			log.Printf("request document : %s", userDocumentDtoJson)
			queryResponse := expectcc.ResponseOk(cc.Invoke("document_registDidDoc", userDocumentDto, data.AppProvider))
			response := expectcc.PayloadIs(queryResponse, &peer.Response{}).(peer.Response)
			expectcc.ResponseOk(response)
			fmt.Printf("response payload : %s\n", response.Payload)
		})

		It(`should return success while user document update - increase versionId`, func() {
			userDocument.VersionId = "2"
			userDocumentDto := makeInvokedDocument(userDocument, getPrivateKeyInPemFile("key/tas/private-key.pem"))

			userDocumentDtoJson, _ := json.Marshal(userDocumentDto)
			log.Printf("request document : %s", userDocumentDtoJson)
			queryResponse := expectcc.ResponseOk(cc.Invoke("document_registDidDoc", userDocumentDto, data.AppProvider))
			response := expectcc.PayloadIs(queryResponse, &peer.Response{}).(peer.Response)
			expectcc.ResponseOk(response)
			fmt.Printf("response payload : %s\n", response.Payload)
		})
	})

	Describe("DidDocAndStatus Get", func() {

		It("Should return tas DidDocAndStatus", func() {
			response := expectcc.PayloadIs(cc.Invoke("document_getDidDoc", "did:opendid:user", "1"), &peer.Response{}).(peer.Response)
			var didDocAndStatus data.DidDocAndStatus
			fmt.Printf("response payload : %s\n", string(response.Payload))

			userDocument.VersionId = "1"

			Expect(json.Unmarshal(response.Payload, &didDocAndStatus)).To(Succeed())
			Expect(didDocAndStatus.DidDoc).To(Equal(*userDocument))
		})

		It("Should return empty", func() {
			response := expectcc.PayloadIs(cc.Invoke("document_getDidDoc", "did:opendid:tas", "2"), &peer.Response{}).(peer.Response)
			fmt.Printf("response payload : %s\n", string(response.Payload))
			Expect(response.Payload).To(Equal([]uint8("null")))
		})
	})

	Describe("Update DidDoc status between In-Service ", func() {

		It("Should return deactivated field true", func() {
			response := expectcc.PayloadIs(cc.Invoke("document_updateDidDocStatusInService", "did:opendid:user", "DEACTIVATED", "1"), &peer.Response{}).(peer.Response)
			fmt.Printf("response payload : %s\n", string(response.Payload))

			var didDoc data.DidDoc
			Expect(json.Unmarshal(response.Payload, &didDoc)).To(Succeed())
			Expect(didDoc.Deactivated).To(Equal(true))
		})

		It("Should return unsupported status error", func() {
			queryResponse := expectcc.ResponseError(cc.Invoke("document_updateDidDocStatusInService", "did:opendid:user", "REVOKED", "1"))
			fmt.Printf(queryResponse.Message)
		})
	})

	Describe("Update DidDoc status to Revocation", func() {

		It("Should return error - invalid status", func() {
			response := expectcc.ResponseError(cc.Invoke("document_updateDidDocStatusRevocation", "did:opendid:user", "TERMINATED", "123"))
			fmt.Println(response.Message)
		})

		It("Should return REVOKED status", func() {
			response := expectcc.PayloadIs(cc.Invoke("document_updateDidDocStatusRevocation", "did:opendid:user", "REVOKED", "123"), &peer.Response{}).(peer.Response)
			fmt.Printf("response payload : %s\n", string(response.Payload))
			response = expectcc.PayloadIs(cc.Invoke("document_getDidDoc", "did:opendid:user", ""), &peer.Response{}).(peer.Response)
			var didDocAndStatus data.DidDocAndStatus
			Expect(json.Unmarshal(response.Payload, &didDocAndStatus)).To(Succeed())
			Expect(didDocAndStatus.Status).To(Equal(data.DIDDOC_STATUS("REVOKED")))
		})
	})
})

func getTasDocument() *data.DidDoc {
	privateKey := getPrivateKeyInPemFile("key/tas/private-key.pem")
	return makeDocument("did:opendid:tas", *privateKey)
}

func getUserdocument() *data.DidDoc {
	privateKey := getPrivateKeyInPemFile("key/user/private-key.pem")
	return makeDocument("did:opendid:user", *privateKey)
}

func getPrivateKeyInPemFile(path string) *ecdsa.PrivateKey {
	log.Printf("PATH : %s", path)
	privateKeyPem, err := os.ReadFile(path)
	if err != nil {
		log.Fatalf("cannot read file. error is : %v", err)
	}

	privateKeyByte, _ := pem.Decode(privateKeyPem)

	privateKey, err := x509.ParseECPrivateKey(privateKeyByte.Bytes)
	if err != nil {
		log.Fatalf("cannot parse EC PrivateKey. error is : %v", err)
		log.Fatal(err)
	}
	return privateKey
}

func makeDocument(id string, privateKey ecdsa.PrivateKey) *data.DidDoc {
	return &data.DidDoc{
		Context: []data.URL{
			"https://www.w3.org/ns/did/v1",
		},
		Id:          id,
		Controller:  id,
		Created:     "2024-05-19T04:32:03",
		Updated:     "2024-05-19T04:32:03",
		VersionId:   "1",
		Deactivated: false,
		VerificationMethod: []data.VerificationMethod{
			{
				Id:                 "pin",
				Type:               data.R1,
				Controller:         id,
				PublicKeyMultibase: data.Multibase("z" + base58.Encode(CompressPublicKey(&privateKey.PublicKey))),
				AuthType:           data.PIN,
			},
		},
		CapabilityInvocation: []data.DID_KEY_ID{
			"pin",
		},
		Service: []data.Service{
			{
				Id:   "homepage",
				Type: "LinkedDomains",
				ServiceEndpoint: []data.URL{
					"http://www.example.com",
				},
			},
		},
	}
}

func makeInvokedDocument(document *data.DidDoc, privateKey *ecdsa.PrivateKey) *data.InvokedDidDoc {
	documentJson, _ := json.Marshal(document)

	doc := EncodingStrByMultiBase(string(documentJson))

	invokedDocument := &data.InvokedDidDoc{
		Proof: data.InvokeProof{
			Type:               data.R1,
			Created:            time.DateTime,
			VerificationMethod: "did:open:tas?versionId=1#pin",
			ProofPurpose:       data.CapabilityInvocation,
		},
		Controller: data.Provider{
			Did:       "did:opendid:tas",
			CertVcRef: "https://test.das.com",
		},
		Nonce:  "test:nonce:1234567890123452988",
		DidDoc: data.Multibase(doc),
	}

	invokedDocumentJson, _ := json.Marshal(invokedDocument)

	plainText := SortJson(invokedDocumentJson)
	log.Printf("signature plain text : %s", plainText)

	signature, err := Sign(plainText, privateKey)
	if err != nil {
		log.Fatal(err)
	}
	invokedDocument.Proof.ProofValue = data.Multibase("z" + base58.Encode(signature))

	return invokedDocument
}
