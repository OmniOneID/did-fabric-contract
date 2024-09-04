// Copyright 2024 Raonsecure

package chaincode_test

//
//import (
//	"encoding/json"
//	"fmt"
//	"opendid/chaincode"
//	"opendid/chaincode/data"
//	"testing"
//
//	"github.com/hyperledger/fabric-protos-go/peer"
//
//	. "github.com/onsi/ginkgo"
//	. "github.com/onsi/gomega"
//
//	testcc "github.com/hyperledger-labs/cckit/testing"
//	expectcc "github.com/hyperledger-labs/cckit/testing/expect"
//)
//
//var (
//	TEST_VCMETA = data.VcMeta{
//		Id: "test:vc:01",
//		Issuer: data.ProviderDetail{
//			Provider: data.Provider{
//				Did:       "open:did:tas",
//				CertVcRef: "http://w3.org",
//			},
//			Name: "tas",
//		},
//		Subject: "open:did:user",
//		CredentialSchema: data.CredentialSchema{
//			Id:   "http://schema.com",
//			Type: "OsdSchemaCredential",
//		},
//		Status:        data.VC_ACTIVE,
//		IssuanceDate:  "2024-05-19T04:32:03",
//		ValidFrom:     "2024-05-19T04:32:03",
//		ValidUntil:    "2024-05-19T04:32:03",
//		FormatVersion: "1.0.1",
//		Language:      "ko",
//	}
//)
//
//func TestVcMetadata(t *testing.T) {
//	RegisterFailHandler(Fail)
//	RunSpecs(t, `OPEN DID TEST CODE IS START`)
//}
//
//var _ = Describe("VC Metadata function test", func() {
//	cc := testcc.NewMockStub(`open did cc`, chaincode.NewOpenDIDCC())
//	// register
//	Describe("VC Metadata Register", func() {
//
//		It("Should return success", func() {
//			queryResponse := expectcc.ResponseOk(cc.Invoke("vcMeta_registVcMetadata", TEST_VCMETA))
//			response := expectcc.PayloadIs(queryResponse, &peer.Response{}).(peer.Response)
//			fmt.Printf("response payload : %s\n", response.Payload)
//		})
//
//		It("Should return error - struct validate fail", func() {
//			invalidRequest := data.VcMeta{Id: "test:vc:02"}
//			queryResponse := expectcc.ResponseError(cc.Invoke("vcMeta_registVcMetadata", invalidRequest))
//			fmt.Println(queryResponse.Message)
//		})
//
//		duplicateKeyRequest := TEST_VCMETA
//		It("Should return error - state key already exists", func() {
//			queryResponse := expectcc.ResponseError(cc.Invoke("vcMeta_registVcMetadata", duplicateKeyRequest))
//			fmt.Println(queryResponse.Message)
//		})
//	})
//
//	//get
//	Describe("VC Metadata Get", func() {
//
//		It("Should be a vcMetaDao", func() {
//			response := expectcc.PayloadIs(cc.Invoke("vcMeta_getVcMetadata", "test:vc:01"), &peer.Response{}).(peer.Response)
//			var vcMeta data.VcMeta
//			fmt.Printf("response payload : %s\n", string(response.Payload))
//			Expect(json.Unmarshal(response.Payload, &vcMeta)).To(Succeed())
//			Expect(vcMeta).To(Equal(TEST_VCMETA))
//		})
//
//		It("Should be empty", func() {
//			queryResponse := expectcc.ResponseOk(cc.Invoke("vcMeta_getVcMetadata", "test:vc:02"))
//			response := expectcc.PayloadIs(queryResponse, &peer.Response{}).(peer.Response)
//			Expect(response.Payload).To(Equal([]uint8("null")))
//		})
//	})
//
//	// update
//	Describe("VC Metadata Update", func() {
//
//		It("Should return success", func() {
//			queryResponse := expectcc.ResponseOk(cc.Invoke("vcMeta_updateVcStatus", "test:vc:01", data.VC_INACTIVE))
//			response := expectcc.PayloadIs(queryResponse, &peer.Response{}).(peer.Response)
//			fmt.Printf("response payload : %s\n", string(response.Payload))
//		})
//
//		It("Should return error - state entry not found", func() {
//			queryResponse := expectcc.ResponseError(cc.Invoke("vcMeta_updateVcStatus", "test:vc:02", data.VC_INACTIVE))
//			fmt.Println(queryResponse.Message)
//		})
//
//		It("Should return error - cannot update vc status ", func() {
//			cc.Invoke("vcMeta_updateVcStatus", "test:vc:01", data.VC_REVOKED)
//			queryResponse := expectcc.ResponseError(cc.Invoke("vcMeta_updateVcStatus", "test:vc:01", data.VC_ACTIVE))
//			fmt.Println(queryResponse.Message)
//		})
//	})
//})
