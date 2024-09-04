// Copyright 2024 Raonsecure

package main

import (
	"did-fabric-contract/chaincode"
	"fmt"

	"github.com/hyperledger/fabric-chaincode-go/shim"
)

func main() {
	if err := shim.Start(chaincode.NewOpenDIDCC()); err != nil {
		fmt.Printf("error starting cc : %s", err)
	}
}
