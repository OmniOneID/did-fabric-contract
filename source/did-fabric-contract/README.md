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
You can easily set up a Fabric network using the test network provided by Hyperledger Fabric. <br>
Refer to the [Hyperledger Fabric official documentation - Using the Fabric test network](https://hyperledger-fabric.readthedocs.io/en/latest/test_network.html) for more details on the network setup and chaincode deployment process outlined below.
1. **Start the Test Network**<br>
   Use the following command to set up the `test-network` composed of two peer organizations and one orderer organization, and create a channel.
   ```bash
   $ cd fabric-sample/test-network
   $ ./network.sh up createChannel -c [channel name] -ca -s couchdb
   ```
2. **Deploy Chaincode**<br>
   Clone the `did-fabric-contract` project under the `fabric-sample` directory.
   ```bash
   $ cd fabric-sample
   $ git clone http://gitlab.raondevops.com/opensourcernd/source/server/did-fabric-contract.git
   ```
   Go back to the `test-network` directory and execute the following command to deploy the chaincode.
   ```bash
   $ cd ./test-network
   $ ./network.sh deployCC -c [channel name] -ccn [chaincode name] -ccp ../did-fabric-contract/source/did-fabric-contract -ccl go -ccs 1
   ```

<br>

## Example Execution
The OpenDID Chaincode improves issues such as low-level chaincode state operations and complex data processing tasks by utilizing the `CCKit Framework`.  
In the `NewOpenDIDCC` function inside the `did-fabric-contract/source/did-fabric-contract/chaincode/opendid.go` file, you can see the various routing paths and methods defined using the `CCKit router` functionality.  
(For more details on CCKit, refer to the [CCKit Github Repository](https://github.com/hyperledger-labs/cckit).)

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

   r.Init(Init)

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
```
<br>

Below is an example of calling the `document_getDidDoc` API to retrieve a DID document from the blockchain using the `peer` CLI.
1. **Set Environment Variables**  
   Set the environment variables in the `test-network` directory to use the `peer` CLI commands.
   ```bash
      $ export PATH=${PWD}/../bin:$PATH
      $ export FABRIC_CFG_PATH=$PWD/../config/
      
      $ export CORE_PEER_TLS_ENABLED=true
      $ export CORE_PEER_LOCALMSPID=Org1MSP
      $ export CORE_PEER_TLS_ROOTCERT_FILE=${PWD}/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt
      $ export CORE_PEER_MSPCONFIGPATH=${PWD}/organizations/peerOrganizations/org1.example.com/users/Admin@org1.example.com/msp
      $ export CORE_PEER_ADDRESS=localhost:7051
   ```
2. **Invoke Chaincode**  
   The `document_getDidDoc` function is called to retrieve a specific DID document where the DID is `did:open:user` and the versionId is `1`.
   ```bash
   $ peer chaincode query -C [channel name] -n [chaincode name] -c '{"Args":["document_getDidDoc","did:open:user","1"]}'
   ```
3. **Return Result**  
   If the command is successful, the result is returned as a `payload` encoded in Base64. The example output is as follows:
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
   The above example output assumes that there is a previously registered DID document. If no DID document is registered, a `null` value is encoded in Base64 and returned as the `payload`.
   ```bash
   {"status":200,"payload":"bnVsbA=="}
   ```