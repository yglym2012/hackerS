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
	//"strings"
	//"strconv"
	//"encoding/json"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

// SimpleChaincode example simple Chaincode implementation
type SimpleChaincode struct {
}

//用户账户信息结构体
type userInfoStruct struct {
	GOLD string 	//黄金
	CNY string 	//人民币
	BTC string 	//比特币
}

//部署，没有传入参数
func (t *SimpleChaincode) Init(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	if len(args) != 0 {
		return nil, errors.New("This is the init step. Incorrect number of arguments. Expecting 0")
	}
	return nil, nil
}

//有3个函数：addUser(用户注册)recharge(充值)exchange（清算）
func (t *SimpleChaincode) Invoke(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	if function == "addUser" {
		return t.addUser(stub, args)
	}else if function == "recharge" {
		return t.recharge(stub, args)
	}else if function == "exchange" {
		return t.exchange(stub, args)
	}

	return nil, errors.New("This is invoke step. no such a method on this chaincode")
}

//addUser 用户注册 
//参数有两个 1.用户ID 2.用户账户信息（json格式 3种资产GOLD CNY BTC）
func (t *SimpleChaincode) addUser(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var userID string	//用户ID
	var userInfo string	//用户账户信息

	var err error

	if len(args) != 2 {
		return nil, errors.New("This is addUser function. Incorrect number of arguments. Expecting 2")
	}

	// Initialize the chaincode
	userID = args[0]
	userInfo = args[1]

	// Write the state to the ledger
	err = stub.PutState(userID, []byte(userInfo))
	if err != nil {
		return nil, err
	}

	return nil, nil
}

//recharge 充值
//参数有3个 用户ID 充值资产类型（GOLD/CNY/BTC其中一种） 充值数量
func (t *SimpleChaincode) recharge(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var userID string 	//用户ID
	var assetType string	//资产种类
	var rechargeNum string	//充值数量

	var userInfo userInfoStruct 	//用户账户信息结构体
	var tmpuserInfo userInfoStruct 	//用户账户信息临时结构体
	var GOLD string 	//黄金
	var CNY string 	//人民币
	var BTC string 	//比特币
	var userInfoByte []byte 	//接收用户账户信息（byte格式）
	var newUserInfoByte []byte 	//接收更新后的用户账户信息结果（byte格式）

	var err error

	if len(args) != 3 {
		return nil, errors.New("This is recharge function. Incorrect number of arguments. Expecting 3")
	}

	// Initialize the chaincode
	userID = args[0]
	assetType = args[1]
	rechargeNum = args[2]

	//取出用户账户信息
	userInfoByte, err = stub.GetState(userID)
	if err != nil {
		return nil, errors.New("Failed to get state")
	}
	if userInfoByte == nil {
		return nil, errors.New("Entity not found")
	}
	//将byte的结果转换成struct
	err = json.Unmarshal(userInfoByte, &userInfo)  
	if err != nil {  
    	fmt.Println("error:", err)  
	}

	//将充值的数量写入相应资产的余额中
	//！！！！！！！！！！有机会换一个简单的写法
	if strings.EqualFold(assetType,GOLD) {
		userInfo.GOLD = rechargeNum
	} else if strings.EqualFold(assetType,CNY) {
		userInfo.CNY = rechargeNum
	} else if strings.EqualFold(assetType,BTC) {
		userInfo.BTC = rechargeNum
	} else {
		return nil, errors.New("No such assetType")
	}

	//用户账户信息变更完毕，将用户账户信息重新存进区块链中
	//将struct类型转换成bytes[]
	newUserInfoByte, err = json.Marshal(tmpDraftInfo)  
	if err != nil {  
		return nil, errors.New("Can not translate struct to byte")
	}

	// Write the state to the ledger
	err = stub.PutState(userID, []byte(newUserInfoByte))
	if err != nil {
		return nil, errors.New("Can not put the new userInfo to the ledger")
	}

	return nil, nil
}

//exchange 结算
//参数有6个 卖方用户ID 买方用户ID 出售资产类型 支付资产类型 出售资产数量 支付资产数量
func (t *SimpleChaincode) exchange(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var sellerID string 	//卖方用户ID
	var payerID string 	//买方用户ID
	var sellAssetType string	//出售资产类型
	var payAssetType string	//支付资产类型
	var sellAssetNum string //出售资产数量
	var payAssetNum string //支付资产数量

	var sellerInfoByte []byte 	//卖方用户账户信息（从ledger中读出）
	var payerInfoByte []byte 	//买方用户账户信息（从ledger中读出）
	var sellerInfo userInfoStruct 	//卖方用户账户信息
	var payerInfo userInfoStruct 	//卖方用户账户信息

	var err error

	if len(args) != 6 {
		return nil, errors.New("This is exchange function. Incorrect number of arguments. Expecting 6")
	}

	// Initialize the chaincode
	sellerID = args[0]
	payerID = args[1]
	sellAssetType = args[2]
	payAssetType = args[3]
	sellAssetNum = args[4]
	payAssetNum = args[5]

	//取出卖方用户的账户信息
	sellerInfoByte, err = stub.GetState(sellerID)
	if err != nil {
		return nil, errors.New("Failed to get state")
	}
	if sellerInfoByte == nil {
		return nil, errors.New("Entity not found")
	}
	//将byte的结果转换成struct
	err = json.Unmarshal(sellerInfoByte, &sellerInfo)  
	if err != nil {  
    	fmt.Println("error:", err)  
	}

	//取出买方用户的账户信息
	payerInfoByte, err = stub.GetState(payerID)
	if err != nil {
		return nil, errors.New("Failed to get state")
	}
	if payerInfoByte == nil {
		return nil, errors.New("Entity not found")
	}
	//将byte的结果转换成struct
	err = json.Unmarshal(payerInfoByte, &payerInfo)  
	if err != nil {  
    	fmt.Println("error:", err)  
	}

	//卖方资产转移
	if strings.EqualFold(sellAssetType,GOLD) {
		sellerInfo.GOLD,payerInfo.GOLD = typeExchangeAndCount(sellerInfo.GOLD,payerInfo.GOLD,sellAssetNum)
	} else if strings.EqualFold(sellAssetType,CNY) {
		sellerInfo.CNY,payerInfo.CNY = typeExchangeAndCount(sellerInfo.CNY,payerInfo.CNY,sellAssetNum)
	} else if strings.EqualFold(sellAssetType,BTC) {
		sellerInfo.BTC,payerInfo.BTC = typeExchangeAndCount(sellerInfo.BTC,payerInfo.BTC,sellAssetNum)
	} else {
		return nil, errors.New("No such assetType")
	}

	//买方资产转移
	if strings.EqualFold(payAssetType,GOLD) {
		payerInfo.GOLD,sellerInfo.GOLD = typeExchangeAndCount(payerInfo.GOLD,sellerInfo.GOLD,payAssetNum)
	} else if strings.EqualFold(payAssetType,CNY) {
		payerInfo.CNY,sellerInfo.CNY = typeExchangeAndCount(payerInfo.CNY,sellerInfo.CNY,payAssetNum)
	} else if strings.EqualFold(payAssetType,BTC) {
		payerInfo.BTC,sellerInfo.BTC = typeExchangeAndCount(payerInfo.BTC,sellerInfo.BTC,payAssetNum)
	} else {
		return nil, errors.New("No such assetType")
	}

	//用户账户信息变更完毕，将用户账户信息重新存进区块链中
	a, err := json.Marshal(sellerInfo)  
	if err != nil {  
	}

	// Write the state to the ledger
	err = stub.PutState(sellerID, []byte(a))
	if err != nil {
		return nil, err
	}

	b, err := json.Marshal(payerInfo)  
	if err != nil {  
	}

	// Write the state to the ledger
	err = stub.PutState(payerID, []byte(b))
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

func typeExchangeAndCount(sellerInit string , payerInit string , variationalValue string ) (string , string){
	str1 := sellerInit
	str2 := payerInit
	str3 := variationalValue

	num1,_ := strconv.ParseFloat(str1,32)
	num2,_ := strconv.ParseFloat(str2,32)
	x,_ := strconv.ParseFloat(str3,32)

	finalSeller := num1 + x
	finalPayer := num2 - x
	

	finalSellerStr := strconv.FormatFloat(finalSeller, 'f', 2, 32)
	finalPayerStr := strconv.FormatFloat(finalPayer, 'f', 2, 32)

	fmt.Println(finalSellerStr)
	fmt.Println(finalPayerStr)

	return finalSellerStr , finalPayerStr

}

func main() {
	err := shim.Start(new(SimpleChaincode))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}
}
