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

//医疗救助对象

import (
	"errors"
	"fmt"
	"strings"
	//"strconv"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

// SimpleChaincode example simple Chaincode implementation
type SimpleChaincode struct {
}

//传入一个参数，[0]是操作人ID
func (t *SimpleChaincode) Init(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {

	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting 1")
	}

	return nil, nil
}

func (t *SimpleChaincode) Invoke(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	if function == "add" {
		return t.add(stub, args)
	}

	return nil, errors.New("no such a method on this chaincode")
}

func (t *SimpleChaincode) add(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {

	if len(args) != 3 {
		return nil, errors.New("Incorrect number of arguments. Expecting 3 ")
	}
	NListID := args[0]    //医疗救助人员ID
	NListIDval := args[1] //当次医疗救助金额信息
	var err error

	// 判断更新的ID是否在链上
	ListIDvalTemp, errs := stub.GetState(NListID)
	//第一次则创建切片数组
	if errs != nil {
		var infoSlice []string = make([]string, 0, 999999)
		infoSlice = append(infoSlice, NListIDval)
		byteContent := "\x00" + strings.Join(infoSlice, "\x02\x00")
		err = stub.PutState(NListID, []byte(byteContent))
		if err != nil {
			return nil, err
		}
	} else {
		l := string(ListIDvalTemp)
		l = append(l, NListIDval)
		byteContent := "\x00" + strings.Join(infoSlice, "\x02\x00")
		err = stub.PutState(NListID, []byte(byteContent))
		if err != nil {
			return nil, err
		}
	}
	return nil, nil
}

// Query callback representing the query of a chaincode,1个参数，[0]医疗救助人员ID
func (t *SimpleChaincode) Query(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	if function != "query" {
		return nil, errors.New("Invalid query function name. Expecting \"query\"")
	}
	var ListID string // Entities
	var err error

	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting name of the person to query")
	}

	ListID = args[0]

	// Get the state from the ledger
	ListIDval, err := stub.GetState(ListID)
	if err != nil {
		jsonResp := "{\"Error\":\"Failed to get state for " + ListID + "\"}"
		return nil, errors.New(jsonResp)
	}

	if ListIDval == nil {
		jsonResp := "{\"Error\":\"Nil amount for " + ListID + "\"}"
		return nil, errors.New(jsonResp)
	}

	jsonResp := "{\"Name\":\"" + ListID + "\",\"Amount\":\"" + string(ListIDval) + "\"}"
	fmt.Printf("Query Response:%s\n", jsonResp)
	return ListIDval, nil
}

func main() {
	err := shim.Start(new(SimpleChaincode))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}
}
