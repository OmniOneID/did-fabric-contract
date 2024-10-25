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

package chaincode

import (
	"fmt"

	"did-fabric-contract/chaincode/data"
	"did-fabric-contract/chaincode/service"

	"github.com/hyperledger-labs/cckit/router"
	"github.com/hyperledger-labs/cckit/router/param"
)

const (
	CC_VERSION = "did-fabric-contract-1.0.0"
)

// NewOpenDIDCC init chaincode function.
func NewOpenDIDCC() *router.Chaincode {

	r := router.New(`OpenDID`)

	r.Init(Init)

	// temp
	r.Group(`remove`).
		Invoke(``, func(ctx router.Context) (interface{}, error) {
			index := string(ctx.GetArgs()[1])
			if err := service.RemoveIndex(ctx, index); err != nil {
				return ctx.Response().Error(err), err
			}
			return ctx.Response().Success(fmt.Sprintf("remove '%s' success", index)), nil
		}).
		Invoke(`All`, func(ctx router.Context) (interface{}, error) {
			if err := service.RemoveAll(ctx); err != nil {
				return ctx.Response().Error(err), err
			}
			return ctx.Response().Success("remove all data success"), nil
		})

	r.Group(`document_`).
		Invoke(`registDidDoc`, registerDidDoc,
			param.Struct("InvokedDidDoc", &data.InvokedDidDoc{}),
			param.String("roleType")).
		Query(`getDidDoc`, getDidDoc,
			param.String("da"),
			param.String("versionId")).
		Invoke(`updateDidDocStatusInService`, updateDidDocStatusInService,
			param.String("da"),
			param.String("status"),
			param.String("versionId")).
		Invoke(`updateDidDocStatusRevocation`, updateDidDocStatusRevocation,
			param.String("da"),
			param.String("status"),
			param.String("terminatedTime"))

	r.Group("vcMeta_").
		Invoke("registVcMetadata", registerVcMetadata,
			param.Struct("vcMeta", &data.VcMeta{})).
		Query("getVcMetadata", getVcMetadata,
			param.String("vcId")).
		Invoke("updateVcStatus", updateVcStatus,
			param.String("vcId"),
			param.String("vcStatus"))

	return router.NewChaincode(r)
}

// register a new document router function definition.
func registerDidDoc(ctx router.Context) (interface{}, error) {
	invokedDidDoc := ctx.Param("InvokedDidDoc").(data.InvokedDidDoc)
	roleType := ctx.ParamString("roleType")

	err := service.RegisterDidDoc(ctx, &invokedDidDoc, data.ROLE_TYPE(roleType))
	if err != nil {
		return ctx.Response().Error(err), err
	}
	return ctx.Response().Success(fmt.Sprintf("document registration successful\n")), nil
}

// get document in ledger router function definition.
func getDidDoc(ctx router.Context) (interface{}, error) {
	da := ctx.ParamString("da")
	versionId := ctx.ParamString("versionId")

	result, err := service.GetDidDocAndStatus(ctx, da, versionId)
	if err != nil {
		return ctx.Response().Error(err), err
	}
	return ctx.Response().Success(result), nil
}

// update document in-service status router function definition.
func updateDidDocStatusInService(ctx router.Context) (interface{}, error) {
	da := ctx.ParamString("da")
	status := ctx.ParamString("status")
	versionId := ctx.ParamString("versionId")

	result, err := service.UpdateDidDocStatusInService(ctx, da, versionId, data.DIDDOC_STATUS(status))
	if err != nil {
		return ctx.Response().Error(err), err
	}
	return ctx.Response().Success(result), nil
}

// update document revocation status router function definition.
func updateDidDocStatusRevocation(ctx router.Context) (interface{}, error) {
	da := ctx.ParamString("da")
	status := ctx.ParamString("status")
	terminatedTime := ctx.ParamString("terminatedTime")

	result, err := service.UpdateDidDocStatusRevocation(ctx, da, data.DIDDOC_STATUS(status), data.UTCDateTime(terminatedTime))
	if err != nil {
		return ctx.Response().Error(err), err
	}
	return ctx.Response().Success(result), nil
}

// register vcMeta router function definition.
func registerVcMetadata(ctx router.Context) (interface{}, error) {
	vcMeta := ctx.Param("vcMeta").(data.VcMeta)
	if err := service.RegisterVcMetadata(ctx, vcMeta); err != nil {
		return ctx.Response().Error(err), err
	}
	return ctx.Response().Success(fmt.Sprintf("vc meta registration successful: vcId = %s\n", vcMeta.Id)), nil
}

// get vcMeta router function definition.
func getVcMetadata(ctx router.Context) (interface{}, error) {
	result, err := service.GetVcMetadata(ctx, ctx.ParamString("vcId"))
	if err != nil {
		return ctx.Response().Error(err), err
	}
	return ctx.Response().Success(result), nil
}

// update vcStatus router function definition.
func updateVcStatus(ctx router.Context) (interface{}, error) {
	vcId := ctx.ParamString("vcId")
	if err := service.UpdateVcStatus(ctx, vcId, data.VC_STATUS(ctx.ParamString("vcStatus"))); err != nil {
		return ctx.Response().Error(err), err
	}
	return ctx.Response().Success(fmt.Sprintf("vc meta update successful: vcId = %s\n", vcId)), nil
}

func Init(ctx router.Context) (interface{}, error) {
	ctx.Logger().Info(fmt.Sprintf("Init open DID CC => version %s", CC_VERSION))
	return fmt.Sprintf("Init open DID CC => version %s", CC_VERSION), nil
}
