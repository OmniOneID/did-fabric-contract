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

package service

import (
	"did-fabric-contract/chaincode/data"
	. "did-fabric-contract/chaincode/error"
	"did-fabric-contract/chaincode/repository"
	"did-fabric-contract/chaincode/utility"
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/hyperledger-labs/cckit/state"

	"github.com/hyperledger-labs/cckit/router"
	"gopkg.in/go-playground/validator.v9"
)

// RemoveIndex
/*
   The function removes the state associated with the given index from the ledger.

   * @param ctx The router context used for state management.
   * @param index The index of the state to be removed.

   * @return An error if the removal fails, otherwise nil.
*/
func RemoveIndex(ctx router.Context, index string) error {
	return ctx.Stub().DelState(index)
}

// RemoveAll
/*
   The function removes all states from the ledger that match predefined prefixes.

   * @param ctx The router context used for state management.

   * @return An error if the removal fails, otherwise nil. Returns an error if no data is found.
*/
func RemoveAll(ctx router.Context) error {
	jsonQuery := make(map[string]interface{})

	selector := make(map[string]interface{})
	jsonQuery["selector"] = selector

	id := make(map[string]interface{})
	selector["_id"] = id

	orConditions := []interface{}{
		map[string]interface{}{"$regex": data.DIDDOC_PREFIX},
		map[string]interface{}{"$regex": data.DIDDOC_STATUS_PREFIX},
		map[string]interface{}{"$regex": data.VCMETA_PREFIX},
	}
	id["$or"] = orConditions

	query, _ := json.Marshal(jsonQuery)
	iterator, err := ctx.Stub().GetQueryResult(string(query))
	if err != nil {
		return err
	}
	defer iterator.Close()

	if !iterator.HasNext() {
		return fmt.Errorf("cannot found data")
	}

	for iterator.HasNext() {
		queryResponse, _ := iterator.Next()
		ctx.Stub().DelState(queryResponse.Key)
	}

	return nil
}

// RegisterDidDoc
/*
   The function registers a DID document into the ledger after validating and converting it.

   * @param ctx The router context used for state management.
   * @param invokedDidDoc The invoked DID document to be registered.
   * @param roleType The role type of the DID document.

   * @return An error if the registration fails, otherwise nil.
*/
func RegisterDidDoc(ctx router.Context, invokedDidDoc *data.InvokedDidDoc, roleType data.ROLE_TYPE) error {

	if err := validateInvokedDidDoc(ctx, invokedDidDoc, roleType); err != nil {
		return err
	}
	didDoc, err := convertFromMultiBaseToDocument(invokedDidDoc.DidDoc)
	if err != nil {
		return GetContractError(DIDDOC_CONVERT_ERROR, err)
	}
	documentStatus := data.MakeDocumentStatus(didDoc, roleType)
	return saveOrUpdateDocument(ctx, didDoc, documentStatus)
}

// validateInvokedDidDoc
/*
   The function validates the invoked DID document based on the role type and signature.

   * @param ctx The router context used for state management.
   * @param invokedDidDoc The invoked DID document to be validated.
   * @param roleType The role type of the DID document.

   * @return An error if the validation fails, otherwise nil.
*/
func validateInvokedDidDoc(ctx router.Context, invokedDidDoc *data.InvokedDidDoc, roleType data.ROLE_TYPE) error {
	if roleType == data.TAS {
		return nil
	}
	if isExist, _ := repository.IsExistDidDocLatest(ctx, invokedDidDoc.Controller.Did); !isExist {
		return GetContractError(DIDDOC_PROVIDER_INVALID, fmt.Errorf("cannot found provider"))
	}
	if err := verifyDidDocument(ctx, invokedDidDoc); err != nil {
		return GetContractError(DIDDOC_SIGNATURE_VERIFICATION_ERROR, err)
	}
	return nil
}

// verifyDidDocument
/*
   The function verifies the signature of the DID document.

   * @param ctx The router context used for state management.
   * @param invokedDidDoc The invoked DID document to be verified.

   * @return An error if the verification fails, otherwise nil.
*/
func verifyDidDocument(ctx router.Context, invokedDidDoc *data.InvokedDidDoc) error {

	versionId, keyId, err := parsingKeyUrl(invokedDidDoc.Proof.VerificationMethod)
	if err != nil {
		return err
	}

	provider := invokedDidDoc.Controller.Did
	didDoc, _, err := getDidDocAndDocumentStatus(ctx, provider, versionId)
	if err != nil {
		return err
	}

	verificationMethod, err := didDoc.GetVerificationMethod(keyId)
	if err != nil {
		return err
	}

	compressedKey, err := utility.DecodeMultibase(verificationMethod.PublicKeyMultibase)
	if err != nil {
		return err
	}

	publicKey, err := utility.DecompressPublicKey(compressedKey)
	if err != nil {
		return err
	}
	signature, err := utility.DecodeMultibase(invokedDidDoc.Proof.ProofValue)
	if err != nil {
		return err
	}

	invokedDidDoc.Proof.ProofValue = ""
	documentDtoJson, _ := json.Marshal(invokedDidDoc)
	plainText := utility.SortJson(documentDtoJson)

	log.Printf("plain text : %s\n", plainText)
	if !utility.Verify(plainText, signature[1:], publicKey) {
		return fmt.Errorf("cannot verify signature")
	}

	return nil
}

// parsingKeyUrl
/*
   The function parses the key URL to extract the version ID and key ID.

   * @param keyUrl The key URL to be parsed.

   * @return The extracted version ID.
   * @return The extracted key ID.
   * @return An error if parsing fails.
*/
func parsingKeyUrl(keyUrl data.DID_KEY_URL) (string, string, error) {
	versionSplit := strings.Split(string(keyUrl), "versionId=")
	if len(versionSplit) != 2 {
		return "", "", fmt.Errorf("invalid verificationMethod. missing versionId")
	}

	versionPart := versionSplit[1]

	keySplit := strings.Split(versionPart, "#")
	if len(keySplit) != 2 {
		return "", "", fmt.Errorf("invalid verificationMethod. missing keyId")
	}

	versionId := keySplit[0]
	keyId := keySplit[1]

	return versionId, keyId, nil
}

// convertFromMultiBaseToDocument
/*
   The function converts a document hash from multibase format to a DID document.

   * @param documentHash The multibase-encoded document hash to be converted.

   * @return The DID document corresponding to the document hash.
   * @return An error if the conversion fails.
*/
func convertFromMultiBaseToDocument(documentHash data.Multibase) (*data.DidDoc, error) {
	documentByteData, err := utility.DecodeMultibase(documentHash)
	if err != nil {
		return nil, err
	}
	didDoc := new(data.DidDoc)
	json.Unmarshal(documentByteData, didDoc)
	if err := validator.New().Struct(didDoc); err != nil {
		return nil, err
	}
	return didDoc, nil

}

// saveOrUpdateDocument
/*
   The function saves or updates the DID document and its status in the ledger.

   * @param ctx The router context used for state management.
   * @param didDoc The DID document to be saved or updated.
   * @param documentStatus The status of the DID document to be saved or updated.

   * @return An error if saving or updating fails, otherwise nil.
*/
func saveOrUpdateDocument(ctx router.Context, didDoc *data.DidDoc, documentStatus *data.DocumentStatus) error {

	savedDidDoc, savedDocumentStatus, err := getDidDocAndDocumentStatus(ctx, didDoc.Id, "")
	if err != nil && !strings.Contains(err.Error(), state.ErrKeyNotFound.Error()) {
		return err
	}

	if savedDidDoc == nil {
		return saveDidDocAndStatus(ctx, didDoc, documentStatus)
	} else {
		return updateDidDocAndStatus(ctx, didDoc, savedDidDoc, savedDocumentStatus)
	}
}

// saveDidDocAndStatus
/*
   The function saves the DID document and its status into the ledger.

   * @param ctx The router context used for state management.
   * @param didDoc The DID document to be saved.
   * @param documentStatus The status of the DID document to be saved.

   * @return An error if saving fails, otherwise nil.
*/
func saveDidDocAndStatus(ctx router.Context, didDoc *data.DidDoc, documentStatus *data.DocumentStatus) error {
	if err := repository.InsertDidDocLatest(ctx, didDoc); err != nil {
		return err
	}
	return repository.InsertDocumentStatus(ctx, documentStatus)
}

// updateDidDocAndStatus
/*
   The function updates the DID document and its status in the ledger.

   * @param ctx The router context used for state management.
   * @param didDoc The DID document to be updated.
   * @param savedDidDoc The previously saved DID document.
   * @param savedDocumentStatus The previously saved document status.

   * @return An error if updating fails, otherwise nil.
*/
func updateDidDocAndStatus(ctx router.Context, didDoc *data.DidDoc, savedDidDoc *data.DidDoc, savedDocumentStatus *data.DocumentStatus) error {
	if savedDidDoc.VersionId >= didDoc.VersionId {
		return GetContractError(DIDDOC_VERSIONID_INVAILD, fmt.Errorf("cannot update didDoc. saved didDoc version id : %s, new didDoc version id %s", savedDidDoc.VersionId, didDoc.VersionId))
	}
	if err := repository.PutDidDocLatest(ctx, didDoc); err != nil {
		return err
	}
	if err := repository.InsertDidDocWithVersionId(ctx, data.MakeDidDocWithVersionId(savedDidDoc)); err != nil {
		return err
	}
	savedDocumentStatus.Version = didDoc.VersionId
	return repository.PutDocumentStatus(ctx, savedDocumentStatus)
}

// GetDidDocAndStatus
/*
   The function retrieves the DID document and its status based on DID and version ID.

   * @param ctx The router context used for state management.
   * @param did The DID of the document.
   * @param versionId The version ID of the document (empty string to get the latest version).

   * @return The DID document and its status if found.
   * @return An error if retrieval fails or if there are issues during processing.
*/
func GetDidDocAndStatus(ctx router.Context, did, versionId string) (*data.DidDocAndStatus, error) {
	didDoc, documentStatus, err := getDidDocAndDocumentStatus(ctx, did, versionId)
	if err != nil {
		if strings.Contains(err.Error(), state.ErrKeyNotFound.Error()) {
			return nil, nil
		}
		return nil, err
	}
	return data.MakeDidDocAndStatus(didDoc, documentStatus.Status), nil
}

// getDidDocAndDocumentStatus
/*
   The function retrieves the DID document and its status based on DID and version ID.

   * @param ctx The router context used for state management.
   * @param did The DID of the document.
   * @param versionId The version ID of the document (empty string to get the latest version).

   * @return The DID document and its status if found.
   * @return An error if retrieval fails.
*/
func getDidDocAndDocumentStatus(ctx router.Context, did, versionId string) (*data.DidDoc, *data.DocumentStatus, error) {
	var didDoc *data.DidDoc
	var err error

	documentStatus, err := repository.GetDocumentStatus(ctx, did)
	if err != nil {
		return nil, nil, err
	}

	if versionId == "" || versionId == documentStatus.Version {
		didDoc, err = repository.GetDidDocLatest(ctx, did)
	} else {
		didDoc, err = repository.GetDidDocWithVersionID(ctx, did, versionId)
	}
	return didDoc, documentStatus, err
}

// UpdateDidDocStatusRevocation
/*
   The function updates the status of the DID document to revoked and sets the termination time.

   * @param ctx The router context used for state management.
   * @param did The DID of the document to be updated.
   * @param status The new status of the document.
   * @param terminatedTime The time when the document was terminated.

   * @return The updated DID document if successful.
   * @return An error if updating fails.
*/
func UpdateDidDocStatusRevocation(ctx router.Context, did string, status data.DIDDOC_STATUS, terminatedTime data.UTCDateTime) (*data.DidDoc, error) {
	didDoc, documentStatus, err := getDidDocAndDocumentStatus(ctx, did, "")
	if err != nil {
		return nil, err
	}
	if err := updateDocumentStatus(documentStatus, status, terminatedTime); err != nil {
		return nil, GetContractError(DIDDOC_STATUS_INVALID, err)
	}

	return didDoc, repository.PutDocumentStatus(ctx, documentStatus)
}

// UpdateDidDocStatusInService
/*
   The function updates the status of the DID document to in-service.

   * @param ctx The router context used for state management.
   * @param did The DID of the document to be updated.
   * @param versionId The version ID of the document to be updated.
   * @param status The new status of the document.

   * @return The updated DID document if successful.
   * @return An error if updating fails.
*/
func UpdateDidDocStatusInService(ctx router.Context, did string, versionId string, status data.DIDDOC_STATUS) (*data.DidDoc, error) {
	didDoc, documentStatus, err := getDidDocAndDocumentStatus(ctx, did, versionId)
	if err != nil {
		return nil, err
	}
	if err := didDoc.SwitchStatus(status); err != nil {
		return nil, GetContractError(DIDDOC_STATUS_INVALID, err)
	}

	if documentStatus.Version == didDoc.VersionId {
		err = repository.PutDidDocLatest(ctx, didDoc)
	} else {
		err = repository.PutDidDocWithVersionId(ctx, data.MakeDidDocWithVersionId(didDoc))
	}
	return didDoc, err
}

// updateDocumentStatus
/*
   The function updates the document status based on the new status.

   * @param documentStatus The document status to be updated.
   * @param status The new status to set.
   * @param terminatedTime The time when the document was terminated (if applicable).

   * @return An error if updating fails, otherwise nil.
*/
func updateDocumentStatus(documentStatus *data.DocumentStatus, status data.DIDDOC_STATUS, terminatedTime data.UTCDateTime) error {
	switch status {
	case data.DOC_REVOKED:
		return documentStatus.Revoke()
	case data.DOC_TERMINATED:
		return documentStatus.Terminated(terminatedTime)
	default:
		return fmt.Errorf("unsupported status: %s", status)
	}
}
