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
	//"strconv"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

// SimpleChaincode example simple Chaincode implementation
type SimpleChaincode struct {
}
//传入一个参数，[0]是操作人ID
func (t *SimpleChaincode) Init(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	//var ID string    // operator's ID
	//var IDval string // ID of the contracts list
	//var err error

	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. ")
	}

	// Initialize the chaincode
	//ID = args[0]
	//IDval = args[1]

	
	//err = stub.PutState(ID, []byte(IDval))
	//if err != nil {
	//	return nil, err
	//}

	return nil, nil
}

//实现3个功能：add，update（输入参数3个：[0]是医疗救助对象的ID，[1]是hash信息，[2]是操作人ID）；delete（输入参数2个：[0]医疗救助对象ID，[1]操作人员ID）
func (t *SimpleChaincode) Invoke(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	if function == "delete" {
		return t.delete(stub, args)
	}else if function == "update"{
		return t.update(stub,args)
	}else if function == "add"{
		return t.add(stub,args)
	}

	return nil, errors.New("no such a method on this chaincode")
}
//update传入3个参数：医疗救助人员ID，信息Hash，操作人ID
func (t *SimpleChaincode) update(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {

	if len(args) != 3 {
		return nil, errors.New("Incorrect number of arguments.Expecting 3 ")
	}
	NListID := args[0] //需更新的医疗救助人员ID
	NListIDval := args[1]//需要更新的信息
	var err error

	// 判断更新的ID是否在链上
		ListIDvalTemp, errs := stub.GetState(NListID)
		if errs != nil {
			return nil, errors.New("list is not here")
		}
		if ListIDvalTemp == nil {
			return nil, errors.New("Entity not found")
		}
	//若存在，则进行更新
		err = stub.PutState(NListID, []byte(NListIDval))
		if err != nil {
			return nil, err
		}
		return nil, nil
}
//add传入3个参数：医疗救助人员ID，信息Hash，操作人ID
func (t *SimpleChaincode) add(stub shim.ChaincodeStubInterface, args []string) ([]byte, error){
		if len(args) != 3 {
			return nil, errors.New("Incorrect number of arguments.Expecting 3 ")
		}
		ListID := args[0]   //新增加的医疗救助人员ID
		ListIDval := args[1]//新增加的医疗救助人员信息
		var err error
		// 添加时需要判断ID对应的值是否已经存在，防止重复添加
        TempHashval, err:= stub.GetState(ListID)

        if TempHashval != nil {
        	return nil, errors.New("This ID already exists")
        }
		// Write the state back to the ledger
		err = stub.PutState(ListID, []byte(ListIDval))
		if err != nil {
			return nil, err
		}
		return nil, nil
}
//delete传入2个参数：医疗救助人员ID，操作人ID
func (t *SimpleChaincode) delete(stub shim.ChaincodeStubInterface, args []string) ([]byte, error){	
	
		if len(args) != 2 {
			return nil, errors.New("Incorrect number of arguments. Expecting 2")
		}
		OListID := args[0] //医疗救助人员ID
		//ListIDvalTemp, errs := stub.GetState(ListID)
		//if errs != nil {
		//	return nil, errors.New("The ID is not existed")
		//}
		//if ListIDvalTemp == nil {
		//	return nil, errors.New("Entity not found")
		//}
		err := stub.DelState(OListID)
		if err != nil {
			return nil, errors.New("Failed to delete state")
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
