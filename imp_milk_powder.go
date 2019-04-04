package main

import (
	"encoding/json"
	"fmt"
	"github.com/hyperledger/fabric/core/chaincode/lib/cid"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/protos/peer"
	"strconv"
)

//奶粉资产
type milkPowderAsset struct {
	Weight       int64  `json:"weight"`
	Owner        string `json:"owner"`
	MilkPowderId string `json:"milk_powder_id"`
	Time         int64  `json:"time"`
}

//只有组织一有权限
//添加奶粉信息
func (t *ImpChaincode) initMilkPowder(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	//===== 输入参数检测
	if len(args) != 4 {
		return shim.Error("Incorrect number of arguments. Expecting 4 id ,weight(g),desc ,cattleInfo")
	}
	id := args[0]
	weight := args[1]
	desc := args[2]
	cattleInfo := args[3]
	time, _ := stub.GetTxTimestamp()
	w, err := strconv.ParseInt(weight, 10, 64)
	if err != nil {
		return shim.Error("weight arg error")
	}
	m, _ := stub.GetState(id)
	if m != nil {
		return shim.Error("this id  exist")
	}
	milkPowder := milkPowder{id, time.Seconds, w, desc, cattleInfo}
	milkPowderByte, _ := json.Marshal(milkPowder)
	//建立资产
	//mspid := "org1"
	mspid, err := cid.GetMSPID(stub)
	if err != nil {
		return shim.Error("get id failed" + err.Error())
	}
	milkPowderAssetByte, err := stub.GetState(mspid + id)
	if err != nil {
		return shim.Error("get state failed" + err.Error())
	}
	if milkPowderAssetByte == nil {
		//创建资产
		time, _ := stub.GetTxTimestamp()
		milkPowderAsset := milkPowderAsset{w, mspid, id, time.Seconds}
		milkPowderAssetByte, _ = json.Marshal(milkPowderAsset)
	} else {
		milkPowderAsset := milkPowderAsset{}
		_ = json.Unmarshal(milkPowderAssetByte, &milkPowderAsset)
		milkPowderAsset.Weight += w
		milkPowderAssetByte, _ = json.Marshal(milkPowderAsset)
	}
	_ = stub.PutState(id, milkPowderByte)
	_ = stub.PutState(mspid+id, milkPowderAssetByte)
	return shim.Success(nil)
}

/*
奶粉转账
发件人是否存在
检查发件人余额
发件人资产减少
收件人增加
交易结束
*/
func (t *ImpChaincode) transferMilkPowder(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	//===== 输入参数检测(收件人地址，奶粉id，重量)
	if len(args) != 3 {
		return shim.Error("Incorrect number of arguments. Expecting 3 recipients，milkPowderId,weight(g)")
	}
	recipients := args[0]
	milkPowderId := args[1]
	weight := args[2]
	w, err := strconv.ParseInt(weight, 10, 64)
	if err != nil {
		return shim.Error("weight arg error")
	}
	//建立资产
	//mspid := "org1"
	mspid, err := cid.GetMSPID(stub)
	if err != nil {
		return shim.Error("get id failed" + err.Error())
	}
	//检测收件人
	recipientsBytes, err := stub.GetState(recipients)
	if err != nil {
		return shim.Error("Failed to get recipients: " + err.Error())
	} else if recipientsBytes == nil {
		return shim.Error("This recipients not exists: " + recipients)
	}
	//检查发件人余额
	milkPowderAssetSendByte, err := stub.GetState(getAssetId(mspid, milkPowderId))
	if err != nil {
		return shim.Error("get state failed" + err.Error())
	} else if milkPowderAssetSendByte == nil {
		return shim.Error("this org has no milk powder asset")
	}
	milkPowderSendAsset := milkPowderAsset{}
	_ = json.Unmarshal(milkPowderAssetSendByte, &milkPowderSendAsset)
	if milkPowderSendAsset.Weight < w {
		return shim.Error("not sufficient funds")
	}
	//建立收件人资产
	milkPowderRecipientsAssetByte, err := stub.GetState(getAssetId(recipients, milkPowderId))
	time, _ := stub.GetTxTimestamp()
	milkPowderRecipientsAsset := milkPowderAsset{0, recipients, milkPowderId, time.Seconds}
	if err != nil {
		return shim.Error("get state failed" + err.Error())
	} else if milkPowderRecipientsAssetByte != nil {
		_ = json.Unmarshal(milkPowderRecipientsAssetByte, milkPowderRecipientsAsset)
	}
	//减少
	milkPowderSendAsset.Weight += w
	milkPowderRecipientsAsset.Weight += w
	//存
	milkPowderRecipientsAssetByte, _ = json.Marshal(milkPowderRecipientsAsset)
	milkPowderAssetSendByte, _ = json.Marshal(milkPowderSendAsset)
	_ = stub.PutState(recipients+milkPowderId, milkPowderRecipientsAssetByte)
	_ = stub.PutState(mspid+milkPowderId, milkPowderAssetSendByte)
	return shim.Success(nil)
}

//查询资产
//输入id 查询起所有资产 Owner
func (t *ImpChaincode) queryMilkPowderAsset(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	if len(args) < 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}
	owner := args[0]
	sql := "{\"selector\": {\"owner\": \"%s\",\"weight\": {\"$gt\": 0}}}"
	queryString := fmt.Sprintf(sql, owner)
	queryResults, err := getQueryResultForQueryString(stub, queryString)
	if err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success(queryResults)
}

//查询奶粉
func (t *ImpChaincode) queryMilkPowder(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1 id ")
	}
	name := args[0]
	valAsbytes, err := stub.GetState(name) //get the milk from chaincode state
	if err != nil {
		jsonResp := "{\"Error\":\"Failed to get state for " + name + "\"}"
		return shim.Error(jsonResp)
	} else if valAsbytes == nil {
		jsonResp := "{\"Error\":\"milk does not exist: " + name + "\"}"
		return shim.Error(jsonResp)
	}
	return shim.Success(valAsbytes)
}

func getAssetId(mspid string, milkPowderId string) string {
	return mspid + milkPowderId + "milkPowderAsset"
}
