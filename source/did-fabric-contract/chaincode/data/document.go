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
	"fmt"
)

type DID_KEY_ID string
type KEY_TYPE string
type VersionID string
type PROOF_PURPOSE string
type DID_KEY_URL string

const (
	DIDDOC_PREFIX = "open:did:doc:"
	LATEST        = ":latest"
)

const (
	RSASignature KEY_TYPE = "RsaSignature2018"
	Secp256k1    KEY_TYPE = "Secp256k1Signature2018"
	Secp256r1    KEY_TYPE = "Secp256r1Signature2018"
)

const (
	AssertionMethod      PROOF_PURPOSE = "assertionMethod"
	Authentication       PROOF_PURPOSE = "authentication"
	KeyAgreement         PROOF_PURPOSE = "keyAgreement"
	CapabilityInvocation PROOF_PURPOSE = "capabilityInvocation"
	CapabilityDelegation PROOF_PURPOSE = "capabilityDelegation"
)

type DidDoc struct {
	Context              []URL                `validate:"required" json:"@context"`
	Id                   string               `validate:"required" json:"id"`
	Controller           string               `validate:"required" json:"controller"`
	Created              UTCDateTime          `validate:"required" json:"created"`
	Updated              UTCDateTime          `validate:"required" json:"updated"`
	VersionId            string               `validate:"required" json:"versionId"`
	Deactivated          bool                 `json:"deactivated"`
	VerificationMethod   []VerificationMethod `validate:"required" json:"verificationMethod"`
	AssertionMethod      []DID_KEY_ID         `json:"assertionMethod,omitempty"`
	Authentication       []DID_KEY_ID         `json:"authentication,omitempty"`
	KeyAgreement         []DID_KEY_ID         `json:"keyAgreement,omitempty"`
	CapabilityInvocation []DID_KEY_ID         `json:"capabilityInvocation,omitempty"`
	CapabilityDelegation []DID_KEY_ID         `json:"capabilityDelegation,omitempty"`
	Service              []Service            `json:"service,omitempty"`
}

type DidDocWithVersionId struct {
	DidDoc
}

type InvokedDidDoc struct {
	DidDoc     Multibase   `json:"didDoc"`
	Proof      InvokeProof `json:"proof"`
	Controller Provider    `json:"controller"`
	Nonce      Multibase   `json:"nonce"`
}

type DidDocAndStatus struct {
	DidDoc DidDoc        `json:"document"`
	Status DIDDOC_STATUS `json:"status"`
}

type InvokeProof struct {
	Type               KEY_TYPE      `json:"type"`
	Created            UTCDateTime   `json:"created"`
	VerificationMethod DID_KEY_URL   `json:"verificationMethod"`
	ProofPurpose       PROOF_PURPOSE `json:"proofPurpose"`
	ProofValue         Multibase     `json:"proofValue,omitempty"`
}

func (d *DidDoc) GetVerificationMethod(id string) (*VerificationMethod, error) {
	if !d.isCapabilityInvocation(DID_KEY_ID(id)) {
		return nil, fmt.Errorf("cannot found key in capabilityInvocation")
	}

	for i := 0; i < len(d.VerificationMethod); i++ {
		if d.VerificationMethod[i].Id == DID_KEY_ID(id) {
			return &d.VerificationMethod[i], nil
		}
	}

	return nil, fmt.Errorf("cannot found key in verificationMethod")
}

func (d *DidDoc) isCapabilityInvocation(id DID_KEY_ID) bool {
	for i := 0; i < len(d.CapabilityInvocation); i++ {
		if d.CapabilityInvocation[i] == DID_KEY_ID(id) {
			return true
		}
	}
	return false
}

func (d *DidDoc) SwitchStatus(status DIDDOC_STATUS) error {
	switch status {
	case DOC_ACTIVATED:
		d.Deactivated = false
		return nil
	case DOC_DEACTIVATED:
		d.Deactivated = true
		return nil
	default:
		return fmt.Errorf("unsupported status: %s", status)
	}
}

func (d DidDoc) Key() ([]string, error) {
	return []string{DIDDOC_PREFIX, d.Id, LATEST}, nil
}

func (d DidDocWithVersionId) Key() ([]string, error) {
	return []string{DIDDOC_PREFIX, d.Id, ":versionId:" + d.VersionId}, nil
}

func MakeDidDocWithVersionId(didDoc *DidDoc) *DidDocWithVersionId {
	return &DidDocWithVersionId{
		DidDoc: *didDoc,
	}
}

func MakeDidDocAndStatus(didDoc *DidDoc, status DIDDOC_STATUS) *DidDocAndStatus {
	return &DidDocAndStatus{
		DidDoc: *didDoc,
		Status: status,
	}
}
