# Fabric Contract Guide
본 문서는 OpenDID Chaincode 사용을 위한 가이드로, 
Open DID에 필요한 DID Document(DID 문서), Verifiable Credential Metadata(이하 VC Meta) 정보를 블록체인에 저장하고 관리하는 기능을 제공합니다. Chaincode는 Blockchain SDK로부터 받은 요청을 처리하여 트랜잭션을 생성 및 기록합니다.

## S/W 사양
| 구분 | 내용                       |
|------|----------------------------|
| Language  | Golang 1.22           |

<br>

## 체인코드 기능
OpenDID Chaincode는 DID 문서 및 VC Meta 데이터와 관련된 트랜잭션 처리 기능을 제공합니다. 라우팅 함수 `NewOpenDIDCC`를 통해 호출하는 함수명과 기능은 다음과 같습니다:

* <b>document_registerDidDoc</b>: 새로운 DID 문서를 등록/변경하고 상태를 저장합니다.
* <b>document_getDidDoc</b>: 특정 DID 문서와 해당 DID 문서 상태를 조회합니다.
* <b>document_updateDidDocStatusInService</b>: In-Service 간 DID 문서 상태를 변경합니다.
* <b>document_updateDidDocStatusRevocation</b>: DID 문서를 Revocation 상태로 변경합니다.
* <b>vcMeta_registerVcMetadata</b>: VC 메타데이터를 등록합니다.
* <b>vcMeta_getVcMetadata</b>: 특정 VC 메타데이터를 조회합니다.
* <b>vcMeta_updateVcStatus</b>: VC 메타데이터의 상태를 변경합니다.

<br>

## 설치 및 배포
Hyperledger Fabric에서 제공하는 test-network를 사용하여 손쉽게 Fabric 네트워크를 구축할 수 있습니다. <br>
아래 단계에서 안내하고 있는 네트워크 실행 및 체인코드 배포 과정은 [Hyperledger Fabric 공식 문서 - Using the Fabric test network](https://hyperledger-fabric.readthedocs.io/en/latest/test_network.html)를 참조하십시오. <br>

1. **테스트 네트워크 실행**<br>
   다음 명령어를 사용하여 두 개의 peer organization과 하나의 order organization으로 구성된 `test-network`를 구축하고 channel을 생성할 수 있습니다.
   ```bash
   $ cd fabric-sample/test-network
   $ ./network.sh up createChannel -c [channel name] -ca -s couchdb
   ```
2. **체인코드 배포**<br>
   `fabric-sample` 디렉터리 하위에 `did-fabric-contract` 프로젝트를 복제합니다.
   ```bash
   $ cd fabric-sample
   $ git clone http://gitlab.raondevops.com/opensourcernd/source/server/did-fabric-contract.git
   ```
   `test-network` 디렉터리로 돌아와서 체인코드 배포를 위해 다음 명령어를 실행합니다.
   ```bash
   $ cd ./test-network
   $ ./network.sh deployCC -c [channel name] -ccn [chaincode name] -ccp ../did-fabric-contract/source/did-fabric-contract -ccl go -ccs 1
   ```

<br>

## 실행 예시
OpenDID Chaincode는 `CCKit Framework`를 활용하여 기존 낮은 수준의 체인코드 상태 작업, 복잡한 데이터 처리 작업 등의 문제점을 개선하고 있습니다. <br>
`did-fabric-contract/source/did-fabric-contract/chaincode/opendid.go` 파일의 `NewOpenDIDCC` 함수 내부에서 `CCKit router` 기능을 사용하여 정의되어 있는 다양한 라우팅 경로와 메서드를 확인할 수 있습니다. <br>
(CCKit에 대한 자세한 내용은 [CCKit Github Repository](https://github.com/hyperledger-labs/cckit)를 참조하십시오.)

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

다음은 `peer` CLI를 사용하여 DID 문서를 블록체인에서 조회하기 위한  `document_getDidDoc` API를 호출하는 예시입니다.<br>
1. **환경변수 설정**<br>
   `peer` CLI 명령어를 사용하기 위해 `test-network` 디렉터리에서 환경변수를 설정합니다.
   ```bash
   $ export PATH=${PWD}/../bin:$PATH
   $ export FABRIC_CFG_PATH=$PWD/../config/
   
   $ export CORE_PEER_TLS_ENABLED=true
   $ export CORE_PEER_LOCALMSPID=Org1MSP
   $ export CORE_PEER_TLS_ROOTCERT_FILE=${PWD}/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt
   $ export CORE_PEER_MSPCONFIGPATH=${PWD}/organizations/peerOrganizations/org1.example.com/users/Admin@org1.example.com/msp
   $ export CORE_PEER_ADDRESS=localhost:7051
   ```
2. **체인코드 호출**<br>
   `document_getDidDoc` 함수를 호출하여 DID와 versionId가 각각 `did:open:user` ,`1`에 해당하는 특정 DID 문서를 조회하고 있습니다.
   ```bash
   $ peer chaincode query -C [channel name] -n [chaincode name] -c '{"Args":["document_getDidDoc","did:open:user","1"]}'
   ```
3. **결괏값 반환**<br>
   명령이 성공하면 Base64로 인코딩된 `payload`로 결과값이 반환됩니다. 예시 출력은 다음과 같습니다: <br>
   ```bash
   {"status":200,"payload":"eyJkb2N1...iXX1dfSwic3RhdHVzIjoiQUNUSVZBVEVEIn0"}
   ```
   `payload` 값을 디코딩하면, 아래와 같이 조회된 특정 DID 문서를 확인할 수 있습니다.
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
   위의 예시 출력은 이전에 등록된 DID 문서가 있다고 가정합니다. 등록된 DID 문서가 없는 경우, `null`값이 Base64 인코딩 되어 `payload` 값으로 반환됩니다.
   ```bash
   {"status":200,"payload":"bnVsbA=="}
   ```
   

