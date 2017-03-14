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

//民政局结算

import (
	"encoding/json"
	"errors"
	"fmt"
	//"strconv"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

// SimpleChaincode example simple Chaincode implementation
type SimpleChaincode struct {
}

type draftInfoStruct struct {
	Sum         string
	From        string
	FromAccount string
	To          string
	ToAccount   string
	Deadline    string
	SerialNum   string
	Status      string
}

//传入一个参数，[0]是操作人ID
func (t *SimpleChaincode) Init(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {

	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. ")
	}

	return nil, nil
}

func (t *SimpleChaincode) Invoke(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	if function == "create" {
		return t.create(stub, args)
	} else if function == "transfer" {
		return t.transfer(stub, args)
	}

	return nil, errors.New("no such a method on this chaincode")
}

//create传入3个参数：数字汇票ID，数字汇票信息，操作人ID
func (t *SimpleChaincode) create(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {

	if len(args) != 3 {
		return nil, errors.New("Incorrect number of arguments.Expecting 3 ")
	}
	draftID := args[0]
	draftInfo := args[1]
	var err error

	// 添加时需要判断ID对应的值是否已经存在，防止重复添加
	TempHashval, err := stub.GetState(draftID)
	if TempHashval != nil {
		return nil, errors.New("This ID already exists")
	}

	err = stub.PutState(draftID, []byte(draftInfo))
	if err != nil {
		return nil, err
	}
	return nil, nil
}

//delete传入3个参数：数字汇票ID，银行流水号，操作人ID
func (t *SimpleChaincode) transfer(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {

	if len(args) != 3 {
		return nil, errors.New("Incorrect number of arguments. Expecting 3")
	}
	draftID := args[0]
	SerialNum := args[1]
	var err error

	var draftInfo draftInfoStruct

	draftInfoTemp, errs := stub.GetState(draftID)
	if errs != nil {
		return nil, errors.New("The ID is not existed")
	}
	if draftInfoTemp == nil {
		return nil, errors.New("Entity not found")
	}

	//将byte的结果转换成struct
	err = json.Unmarshal(draftInfoTemp, &draftInfo)
	if err != nil {
		fmt.Println("error:", err)
	}

	draftInfo.SerialNum = SerialNum
	draftInfo.Status = "done"

	n, err := json.Marshal(draftInfo)
	if err != nil {
		return nil, errors.New("Can not translate struct to byte")
	}

	// Write the state to the ledger
	err = stub.PutState(draftID, []byte(n))
	if err != nil {
		return nil, errors.New("Can not put the new userInfo to the ledger")
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
