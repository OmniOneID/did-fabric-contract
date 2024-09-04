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
아래 단계는 Hypereldger Fabric 테스트 네트워크에서 `network.sh` 스크립트를 사용하여 체인코드를 손쉽게 배포 및 실행하는 방법을 안내하고 있습니다.
네트워크 구축, 설치 및 승인 등 자세한 네트워크 구성과 배포 과정은 [Hyperledger Fabric 공식 문서](https://hyperledger-fabric.readthedocs.io/)를 참조하십시오. <br>
1. **테스트 네트워크 실행**<br>
   먼저, 다음 명령어를 사용하여 테스트 네트워크를 실행하고 채널을 생성합니다.
   ```bash
   ./network.sh up createChannel -c [channel name] -ca -s couchdb
   ```
2. **의존성 관리**<br>
   프로젝트 루트 디렉터리에서 `go mod tidy` 명령어를 실행하여 의존성을 가져오고, `go.mod`와 `go.sum` 파일을 최신 상태로 유지할 수 있습니다. 
   ```go
   go mod tidy
   ```
   이후, 외부 의존성을 포함하기 위해 `go mod vendor` 명령어를 실행하여 `vendor` 디렉터리를 생성합니다.
   ```go
   go mod vendor
   ```
3. **체인코드 배포**<br>
   최종적으로 체인코드 배포를 위해 다음 명령어를 실행합니다.
   ```bash
   ./network.sh deployCC -ccn [chaincode name] -ccp [project path] -ccl go -ccs 1
   ```

<br>

## 실행 예시
OpenDID Chaincode는 CCKit 프레임워크를 활용하여 기존 낮은 수준의 체인코드 상태 작업, 복잡한 데이터 처리 작업 등의 문제점을 개선하고 있습니다. 
다음은 DID 문서를 블록체인에서 조회하기 위해 CCKit 프레임워크 기반의 Chaincode의 특정 메서드를 호출하는 예시입니다.<br>

CCKit에 대한 자세한 내용은 [CCKit Github Repository](https://github.com/hyperledger-labs/cckit)를 참조하십시오. 

`chaincode/opendid.go` 파일의 `NewOpenDIDCC` 함수에서 `CCKit router` 기능을 사용하여 다양한 라우팅 경로와 메서드를 정의하고 있습니다. 아래 `document_getDidDoc`이라는 이름을 사용하여 DID 문서 조회 함수 `getdidDoc`을 외부에서 호출할 수 있습니다.
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

내부적으로 `repository/document_repository.go`에서 상태 처리 작업이 수행되며, CCKit의 상태 관리 기능을 통해 블록체인에 저장된 데이터를 손쉽게 조회할 수 있습니다.
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

CLI 명령어을 사용하여 DID 문서를 조회 기능을 호출할 수 있습니다.<br> 다음 명령어는 `document_getDidDoc` 함수를 호출하여 DID와 versionId가 각각 `did:open:user` ,`1`에 해당하는 특정 DID 문서를 조회하고 있습니다.
```bash
peer chaincode query -C [channel name] -n [chaincode name] -c '{"Args":["document_getDidDoc","did:open:user","1"]"}'
```


명령이 성공하면 Base64로 인코딩된 `payload`로 결과값이 반환됩니다. 예시 출력은 다음과 같습니다:
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