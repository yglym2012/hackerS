/*
Copyright IBM Corp. 2016 All Rights Reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

		 http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main


//清算

import (
	"errors"
	"fmt"
	"strings"
	"strconv"
	"encoding/json"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

// SimpleChaincode example simple Chaincode implementation
type SimpleChaincode struct {
}

//部署，没有传入参数
func (t *SimpleChaincode) Init(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	if len(args) != 0 {
		return nil, errors.New("This is the init step. Incorrect number of arguments. Expecting 0")
	}
	return nil, nil
}

//有3个函数：newSell(新建卖单)newPay(新建买单)
func (t *SimpleChaincode) Invoke(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	if function == "newSell" {
		return t.newSell(stub, args)
	}else if function == "newPay" {
		return t.newPay(stub, args)
	}

	return nil, errors.New("This is invoke step. no such a method on this chaincode")
}

//newSell 新建卖单 
//参数有两个 1.卖单ID 2.卖单信息（json格式）
func (t *SimpleChaincode) newSell(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var sellOrderID string	//卖单ID
	var sellOrderInfo string	//卖单信息（json格式）

	var err error

	if len(args) != 2 {
		return nil, errors.New("This is newSell function. Incorrect number of arguments. Expecting 2")
	}

	// Initialize the chaincode
	sellOrderID = args[0]
	sellOrderInfo = args[1]

	// Write the state to the ledger
	err = stub.PutState(sellOrderID, []byte(sellOrderInfo))
	if err != nil {
		return nil, err
	}

	return nil, nil
}

//newPay 新建买单 
//参数有两个 1.买单ID 2.买单信息（json格式）
func (t *SimpleChaincode) newPay(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var payOrderID string	//买单ID
	var payOrderInfo string	//买单信息（json格式）

	var err error

	if len(args) != 2 {
		return nil, errors.New("This is newPay function. Incorrect number of arguments. Expecting 2")
	}

	// Initialize the chaincode
	payOrderID = args[0]
	payOrderInfo = args[1]

	// Write the state to the ledger
	err = stub.PutState(payOrderID, []byte(payOrderInfo))
	if err != nil {
		return nil, err
	}

	return nil, nil
}

// Query callback representing the query of a chaincode
func (t *SimpleChaincode) Query(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	if function != "query" {
		return nil, errors.New("Invalid query function name. Expecting \"query\"")
	}
	var A string // Entities
	var err error

	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting name of the person to query")
	}

	A = args[0]

	// Get the state from the ledger
	Avalbytes, err := stub.GetState(A)
	if err != nil {
		jsonResp := "{\"Error\":\"Failed to get state for " + A + "\"}"
		return nil, errors.New(jsonResp)
	}

	if Avalbytes == nil {
		jsonResp := "{\"Error\":\"Nil amount for " + A + "\"}"
		return nil, errors.New(jsonResp)
	}

	jsonResp := "{\"Name\":\"" + A + "\",\"Amount\":\"" + string(Avalbytes) + "\"}"
	fmt.Printf("Query Response:%s\n", jsonResp)
	return Avalbytes, nil
}

func main() {
	err := shim.Start(new(SimpleChaincode))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}
}
