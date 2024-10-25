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
	"did-fabric-contract/chaincode/validate"
	"fmt"
	"strings"

	"github.com/hyperledger-labs/cckit/state"

	"github.com/hyperledger-labs/cckit/router"
)

// RegisterVcMetadata
/*
   The function registers the provided VC metadata into the ledger after validating it.

   * @param ctx The router context used for state management.
   * @param vcMeta The VC metadata to be registered in the ledger.

   * @return An error if any issue occurred during validation or insertion, otherwise nil.
*/
func RegisterVcMetadata(ctx router.Context, vcMeta data.VcMeta) error {
	v := validate.RegisterVcMetaValidator()
	if err := v.Struct(vcMeta); err != nil {
		return GetContractError(VCMETA_CONVERT_ERROR, err)
	}

	if err := repository.InsertVcMeta(ctx, &vcMeta); err != nil {
		return err
	}
	return nil
}

// GetVcMetadata
/*
   The function retrieves the VC metadata for the specified VC ID from the ledger.

   * @param ctx The router context used for state management.
   * @param vcId The VC ID of the VC metadata to retrieve.

   * @return The VC metadata associated with the specified VC ID if found.
   * @return An error if any issue occurred during retrieval or if the key is not found.
*/
func GetVcMetadata(ctx router.Context, vcId string) (*data.VcMeta, error) {
	vcMeta, err := repository.GetVcMeta(ctx, vcId)
	if err != nil && !strings.Contains(err.Error(), state.ErrKeyNotFound.Error()) {
		return nil, err
	}
	return vcMeta, nil
}

// UpdateVcStatus
/*
   The function updates the status of the VC metadata associated with the specified VC ID.

   * @param ctx The router context used for state management.
   * @param vcId The VC ID of the VC metadata to update.
   * @param vcStatus The new status to set for the VC metadata.

   * @return An error if any issue occurred during retrieval or update, or if the status update is invalid.
*/
func UpdateVcStatus(ctx router.Context, vcId string, vcStatus data.VC_STATUS) error {
	vcMeta, err := repository.GetVcMeta(ctx, vcId)
	if err != nil {
		return err
	}
	if vcMeta.Status == vcStatus || vcMeta.Status == data.VC_REVOKED {
		return GetContractError(VCMETA_STATUS_INVALID, fmt.Errorf("cannot update status from %s to %s", vcMeta.Status, vcStatus))
	}
	vcMeta.Status = vcStatus
	if err := repository.PutVcMeta(ctx, vcMeta); err != nil {
		return err
	}
	return nil
}
