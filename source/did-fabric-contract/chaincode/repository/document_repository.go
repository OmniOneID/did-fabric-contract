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

package repository

import (
	"did-fabric-contract/chaincode/data"
	. "did-fabric-contract/chaincode/error"
	"encoding/json"

	"github.com/hyperledger-labs/cckit/router"
)

// InsertDidDocLatest
/*
   The function inserts the provided DID document into the ledger.

   * @param ctx The router context used for state management.
   * @param didDoc The DID document to be inserted into the ledger.

   * @return An error if any issue occurred during the insertion, otherwise nil.
*/
func InsertDidDocLatest(ctx router.Context, didDoc *data.DidDoc) error {
	if err := ctx.State().Insert(didDoc); err != nil {
		return GetContractError(DIDDOC_INSERT_ERROR, err)
	}
	return nil
}

// GetDidDocLatest
/*
   The function retrieves the latest DID document from the ledger.

   * @param ctx The router context used for state management.
   * @param did The DID of the DID document to retrieve.

   * @return The DID document containing the latest data if successful.
   * @return An error if any issue occurred during retrieval or conversion.
*/
func GetDidDocLatest(ctx router.Context, did string) (*data.DidDoc, error) {
	result, err := ctx.State().Get(&data.DidDoc{Id: did})
	var didDoc *data.DidDoc
	if err != nil {
		return nil, GetContractError(DIDDOC_GET_ERROR, err)
	}
	if err := json.Unmarshal(result.([]uint8), &didDoc); err != nil {
		return nil, GetContractError(DIDDOC_CONVERT_ERROR, err)
	}
	return didDoc, nil
}

// PutDidDocLatest
/*
   The function updates the provided DID document in the ledger.

   * @param ctx The router context used for state management.
   * @param didDoc The DID document to be updated in the ledger.

   * @return An error if any issue occurred during the update, otherwise nil.
*/
func PutDidDocLatest(ctx router.Context, didDoc *data.DidDoc) error {
	if err := ctx.State().Put(didDoc); err != nil {
		return GetContractError(DIDDOC_PUT_ERROR, err)
	}
	return nil
}

// IsExistDidDocLatest
/*
   The function checks if a DID document with the specified DID exists in the ledger.

   * @param ctx The router context used for state management.
   * @param da The DID of the DID document to check for existence.

   * @return True if the DID document exists, otherwise false.
   * @return An error if any issue occurred during the check.
*/
func IsExistDidDocLatest(ctx router.Context, da string) (bool, error) {
	return ctx.State().Exists(&data.DidDoc{Id: da})
}

// InsertDidDocWithVersionId
/*
   The function inserts the provided DID document with version ID into the ledger.

   * @param ctx The router context used for state management.
   * @param didDoc The DID document with version ID to be inserted into the ledger.

   * @return An error if any issue occurred during the insertion, otherwise nil.
*/
func InsertDidDocWithVersionId(ctx router.Context, didDoc *data.DidDocWithVersionId) error {
	if err := ctx.State().Insert(didDoc); err != nil {
		return GetContractError(DIDDOC_INSERT_ERROR, err)
	}
	return nil
}

// GetDidDocWithVersionID
/*
   The function retrieves the DID document with the specified version ID from the ledger.

   * @param ctx The router context used for state management.
   * @param did The DID of the DID document to retrieve.
   * @param versionId The version ID of the DID document to retrieve.

   * @return The DID document with the specified version ID if found.
   * @return An error if any issue occurred during retrieval or conversion.
*/
func GetDidDocWithVersionID(ctx router.Context, did string, versionId string) (*data.DidDoc, error) {
	result, err := ctx.State().Get(&data.DidDocWithVersionId{DidDoc: data.DidDoc{Id: did, VersionId: versionId}})
	var didDocWithVersionId *data.DidDocWithVersionId
	if err != nil {
		return nil, GetContractError(DIDDOC_GET_ERROR, err)
	}
	if err := json.Unmarshal(result.([]uint8), &didDocWithVersionId); err != nil {
		return nil, GetContractError(DIDDOC_CONVERT_ERROR, err)
	}
	return &didDocWithVersionId.DidDoc, nil
}

// PutDidDocWithVersionId
/*
   The function updates the provided DID document with version ID in the ledger.

   * @param ctx The router context used for state management.
   * @param didDocWithVersionId The DID document with version ID to be updated in the ledger.

   * @return An error if any issue occurred during the update, otherwise nil.
*/
func PutDidDocWithVersionId(ctx router.Context, didDocWithVersionId *data.DidDocWithVersionId) error {
	if err := ctx.State().Put(didDocWithVersionId); err != nil {
		return GetContractError(DIDDOC_PUT_ERROR, err)
	}
	return nil
}

// InsertDocumentStatus
/*
   The function inserts the provided document status into the ledger.

   * @param ctx The router context used for state management.
   * @param documentStatus The document status to be inserted into the ledger.

   * @return An error if any issue occurred during the insertion, otherwise nil.
*/
func InsertDocumentStatus(ctx router.Context, documentStatus *data.DocumentStatus) error {
	if err := ctx.State().Insert(documentStatus); err != nil {
		return GetContractError(DIDDOC_STATUS_INSERT_ERROR, err)
	}
	return nil
}

// PutDocumentStatus
/*
   The function updates the provided document status in the ledger.

   * @param ctx The router context used for state management.
   * @param documentStatus The document status to be updated in the ledger.

   * @return An error if any issue occurred during the update, otherwise nil.
*/
func PutDocumentStatus(ctx router.Context, documentStatus *data.DocumentStatus) error {
	if err := ctx.State().Put(documentStatus); err != nil {
		return GetContractError(DIDDOC_STATUS_PUT_ERROR, err)
	}
	return nil
}

// GetDocumentStatus
/*
   The function retrieves the document status for the specified DID from the ledger.

   * @param ctx The router context used for state management.
   * @param did The DID of the document status to retrieve.

   * @return The document status associated with the specified DID if found.
   * @return An error if any issue occurred during retrieval or conversion.
*/
func GetDocumentStatus(ctx router.Context, did string) (*data.DocumentStatus, error) {
	result, err := ctx.State().Get(&data.DocumentStatus{Id: did})
	var documentStatus *data.DocumentStatus
	if err != nil {
		return nil, GetContractError(DIDDOC_STATUS_GET_ERROR, err)
	}
	if err := json.Unmarshal(result.([]uint8), &documentStatus); err != nil {
		return nil, GetContractError(DIDDOC_STATUS_CONVERT_ERROR, err)
	}
	return documentStatus, json.Unmarshal(result.([]uint8), &documentStatus)
}
