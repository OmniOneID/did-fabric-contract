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

// InsertVcMeta
/*
   The function inserts the provided VC metadata into the ledger.

   * @param ctx The router context used for state management.
   * @param vcMeta The VC metadata to be inserted into the ledger.

   * @return An error if the insertion fails, otherwise nil.
*/
func InsertVcMeta(ctx router.Context, vcMeta *data.VcMeta) error {
	if err := ctx.State().Insert(vcMeta); err != nil {
		return GetContractError(VCMETA_INSERT_ERROR, err)
	}
	return nil
}

// PutVcMeta
/*
   The function updates or inserts the VC metadata in the ledger. If the VC metadata already exists, it is updated; otherwise, it is inserted.

   * @param ctx The router context used for state management.
   * @param vcMeta The VC metadata to be updated or inserted into the ledger.

   * @return An error if the update or insertion fails, otherwise nil.
*/
func PutVcMeta(ctx router.Context, vcMeta *data.VcMeta) error {
	if err := ctx.State().Put(vcMeta); err != nil {
		return GetContractError(VCMETA_PUT_ERROR, err)
	}
	return nil
}

// GetVcMeta
/*
   The function retrieves the VC metadata for the specified VC ID from the ledger.

   * @param ctx The router context used for state management.
   * @param vcId The VC ID of the VC metadata to retrieve.

   * @return The VC metadata associated with the specified VC ID if found.
   * @return An error if the retrieval fails or if there is an issue converting the data.
*/
func GetVcMeta(ctx router.Context, vcId string) (*data.VcMeta, error) {
	result, err := ctx.State().Get(&data.VcMeta{Id: vcId})
	var vcMeta *data.VcMeta
	if err != nil {
		return nil, GetContractError(VCMETA_GET_ERROR, err)
	}
	if err = json.Unmarshal(result.([]uint8), &vcMeta); err != nil {
		return nil, GetContractError(VCMETA_CONVERT_ERROR, err)
	}
	return vcMeta, nil
}
