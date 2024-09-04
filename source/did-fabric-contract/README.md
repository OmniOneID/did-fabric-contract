# Fabric Contract Guide
This document is a guide for using the OpenDID Chaincode. The Chaincode provides functionality to store and manage the necessary DID Document and Verifiable Credential Metadata(hereafter VC Meta) information for OpenDID on the blockchain. The Chaincode processes requests received from the Blockchain SDK to generate and record transactions.

## S/W Specifications
| Category | Details                |
|----------|------------------------|
| Language | Golang 1.22             |

<br>

## Chaincode Features
The OpenDID Chaincode provides transaction processing functions related to DID documents and VC metadata. The functions called via the routing function NewOpenDIDCC are as follows:

* <b>document_registerDidDoc</b>: Registers/updates a new DID document and saves its state.
* <b>document_getDidDoc</b>: Retrieves a specific DID document and its state.
* <b>document_updateDidDocStatusInService</b>: Changes the status of a DID document to "In-Service".
* <b>document_updateDidDocStatusRevocation</b>: Changes the status of a DID document to "Revocation".
* <b>vcMeta_registerVcMetadata</b>: Registers VC metadata.
* <b>vcMeta_getVcMetadata</b>: Retrieves specific VC metadata.
* <b>vcMeta_updateVcStatus</b>: Changes the status of VC metadata.

<br>

## Installation and Deployment
The following steps guide you on how to easily deploy and execute chaincode using the `network.sh` script in the Hyperledger Fabric test network.
For detailed instructions on setting up the network, installation, endorsement, and more, please refer to the [Hyperledger Fabric official documentation](https://hyperledger-fabric.readthedocs.io/). <br>
1. **Start the Test Network**<br>
   First, use the following command to start the test network and create a channel.
   ```bash
   ./network.sh up createChannel -c [channel name] -ca -s couchdb
   ```
2. **Dependency Management**<br>
   From the root directory of the project, run the `go mod tidy` command to fetch dependencies and keep the `go.mod` and `go.sum` files up to date.
   ```go
   go mod tidy
   ```
   Then, to include external dependencies, run the `go mod vendor` command to create the `vendor` directory.
   ```go
   go mod vendor
   ```
3. **Deploying the Chaincode**<br>
   Finally, run the following command to deploy the chaincode:
   ```bash
   ./network.sh deployCC -ccn [chaincode name] -ccp [project path] -ccl go -ccs 1
   ```

<br>

## Example Execution
The OpenDID Chaincode leverages the CCKit framework to address issues related to low-level chaincode state operations and complex data processing tasks. 
Below is an example of calling a specific method in the CCKit-based chaincode to retrieve a DID document from the blockchain.<br>

For more information on CCKit, refer to the [CCKit GitHub Repository](https://github.com/hyperledger-labs/cckit).

In the `NewOpenDIDCC` function within the `chaincode/opendid.go` file, various routing paths and methods are defined using the `CCKit router` feature. The DID document retrieval function `getDidDoc` can be called externally using the name `document_getDidDoc`.
```go
package chaincode

import (
	"fmt"

	"opendid/chaincode/data"
	"opendid/chaincode/service"

	"github.com/hyperledger-labs/cckit/router"
	"github.com/hyperledger-labs/cckit/router/param"
)


func NewOpenDIDCC() *router.Chaincode {

	r := router.New(`OpenDID`)

	r.Group(`document_`).
		Query(`getDidDoc`, getDidDoc,
			param.String("da"),
			param.String("versionId"))
   ...

	return router.NewChaincode(r)
}

func getDidDoc(ctx router.Context) (interface{}, error) {
	da := ctx.ParamString("da")
	versionId := ctx.ParamString("versionId")

	result, err := service.GetDidDocAndStatus(ctx, da, versionId)
	if err != nil {
		return ctx.Response().Error(err), err
	}
	return ctx.Response().Success(result), nil
}
```

Internally, state processing tasks are performed in `repository/document_repository.go`, and with CCKit's state management features, data stored on the blockchain can be easily retrieved.
```go
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
```

You can invoke the DID document retrieval function using CLI commands.<br> The following command calls the `document_getDidDoc` function to retrieve a specific DID document with `DID` as `did:open:user` and `versionId` as `1`.
```bash
peer chaincode query -C [channel name] -n [chaincode name] -c '{"Args":["document_getDidDoc","did:open:user","1"]"}'
```

If the command is successful, the result is returned as a Base64-encoded `payload`. An example of the output is as follows:
```bash
{"status":200,"payload":"eyJkb2N1...iXX1dfSwic3RhdHVzIjoiQUNUSVZBVEVEIn0"}
```
Decoding the `payload` value will reveal the retrieved specific DID document as shown below:
```json
{
  	"document": {
    	"@context": [
      		"https://www.w3.org/ns/did/v1"
    	],
    	"id": "did:opendid:user",
		...
    	"versionId": "1",
    	"deactivated": false,
		...
  	},
  	"status": "ACTIVATED"
}
```