package main

import (
	"encoding/json"
	"fmt"
	"github.com/hyperledger/fabric/core/chaincode/lib/cid"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/protos/peer"
	"strconv"
)

type productAsset struct {
	Count     int64  `json:"count"`
	Owner     string `json:"owner"`
	ProductId string `json:"product_id"`
	Time      int64  `json:"time"`
}

//送检
type toInspect struct {
	ProductId    string `json:"product_id"`
	Time         int64  `json:"time"` //提交时间
	ToInspection string `json:"to_inspection"`
}

/*
出厂奶粉
输入出厂产品,奶粉id,数量，每个重量，详细描述
判断奶粉是否充足
奶粉减，增加产品资产
*/

func (t *ImpChaincode) manufacture(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	if len(args) != 5 {
		return shim.Error("Incorrect number of arguments. Expecting 5,product, milkId,count,weight,desc")
	}
	productId := args[0]
	milkId := args[1]
	count := args[2]
	weight := args[3]
	desc := args[4]
	w, err := strconv.ParseInt(weight, 10, 64)
	if err != nil {
		return shim.Error("arg" + err.Error())
	}
	c, err := strconv.ParseInt(count, 10, 64)
	if err != nil {
		return shim.Error("arg" + err.Error())
	}
	mspId, err := cid.GetMSPID(stub)
	if err != nil {
		return shim.Error("get id failed" + err.Error())
	}

	//检查奶粉是否充足
	milkPowderAssetSendByte, err := stub.GetState(getPowderAssetId(mspId, milkId))
	if err != nil {
		return shim.Error("get state failed" + err.Error())
	} else if milkPowderAssetSendByte == nil {
		return shim.Error("this org has not this milk powder")
	}
	milkPowderSendAsset := milkPowderAsset{}
	_ = json.Unmarshal(milkPowderAssetSendByte, &milkPowderSendAsset)
	if milkPowderSendAsset.Weight < c*w {
		return shim.Error("not sufficient funds")
	}
	//建立资产
	time, _ := stub.GetTxTimestamp()
	p, err := stub.GetState(productId)
	if p != nil {
		return shim.Error("already exist this product")
	}
	//建立product
	product := product{productId, desc, w, mspId, time.Seconds, 0, inspectInfo{"", 0, ""}}
	productAsset := productAsset{c, mspId, productId, time.Seconds}
	//保存奶粉产品资产
	productByte, _ := json.Marshal(product)
	productAssetByte, _ := json.Marshal(productAsset)
	_ = stub.PutState(productId, productByte)
	_ = stub.PutState(getProductAssetId(mspId, productId), productAssetByte)
	return shim.Success(nil)
}

/*
将产品交给检测机构（送检）
输入产品编号，检测机构id,检测信息
检测奶粉是否存在
*/
func (t *ImpChaincode) toInspect(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments. Expecting 2,productId,OrgId")
	}
	productId := args[0]
	inspectionId := args[1]
	productByte, err := stub.GetState(productId)
	if err != nil {
		return shim.Error("get state failed" + err.Error())
	}
	if productByte != nil {
		return shim.Error("already exist this product")
	}
	time, _ := stub.GetTxTimestamp()
	toInspect := toInspect{productId, time.Seconds, inspectionId}
	product := product{}
	json.Unmarshal(productByte, product)
	//product.InspectInfo=inspectInfo
	product.Status = 1
	productByte, _ = json.Marshal(product)
	toInspectByte, _ := json.Marshal(toInspect)
	stub.PutState(productId, productByte)
	stub.PutState(getToInspectionId(productId, inspectionId), toInspectByte)
	return shim.Success(nil)
}

//检测奶粉
//奶粉id，详情
func (t *ImpChaincode) inspect(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments. Expecting 2,productId,desc")
	}
	productId := args[0]
	desc := args[1]
	productByte, err := stub.GetState(productId)
	mspId, err := cid.GetMSPID(stub)
	if err != nil {
		return shim.Error("get id failed" + err.Error())
	}
	if err != nil {
		return shim.Error("get state failed" + err.Error())
	}
	if productByte == nil {
		return shim.Error("not exist this product")
	}
	product := product{}
	json.Unmarshal(productByte, product)
	time, _ := stub.GetTxTimestamp()
	inspectInfo := inspectInfo{mspId, time.Seconds, desc}
	product.InspectInfo = inspectInfo
	product.Status = 2
	productByte, _ = json.Marshal(product)
	stub.PutState(productId, productByte)
	return shim.Success(nil)
}

//查询我的待检测检测
//没有输入
func (t *ImpChaincode) getMyInspect(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	mspId, err := cid.GetMSPID(stub)
	if err != nil {
		return shim.Error("get id failed" + err.Error())
	}
	sql := "{\"selector\": {\"to_inspection\": \"%s\"}}"
	queryString := fmt.Sprintf(sql, mspId)
	queryResults, err := getQueryResultForQueryString(stub, queryString)
	if err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success(queryResults)
}

//查询资产
//输入id 查询起产品所有资产 Owner
func (t *ImpChaincode) queryProductAsset(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	if len(args) < 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1，owner")
	}
	owner := args[0]
	sql := "{\"selector\": {\"owner\": \"%s\",\"count\": {\"$gt\": 0}}}"
	queryString := fmt.Sprintf(sql, owner)
	queryResults, err := getQueryResultForQueryString(stub, queryString)
	if err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success(queryResults)
}

//查询产品
func (t *ImpChaincode) queryProduct(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1 id ")
	}
	name := args[0]
	valAsbytes, err := stub.GetState(name) //get the milk from chaincode state
	if err != nil {
		jsonResp := "{\"Error\":\"Failed to get state for " + name + "\"}"
		return shim.Error(jsonResp)
	} else if valAsbytes == nil {
		jsonResp := "{\"Error\":\"product does not exist: " + name + "\"}"
		return shim.Error(jsonResp)
	}
	return shim.Success(valAsbytes)
}

func getToInspectionId(productId string, inspectiosnId string) string {
	return productId + inspectiosnId + "productAsset"
}
func getProductAssetId(mspid string, productAsset string) string {
	return mspid + productAsset + "productAsset"
}
