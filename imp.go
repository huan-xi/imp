/**
imp溯源系统链代码
实现功能
初始化组织信息（组织id，组织描述）
org1 添加奶粉信息（出厂编号，时间，出产量，出产奶牛信息）
奶粉转让，指定重量给厂家
org2 加工奶粉成产品（批次编号，数量，出厂时间，配料，检测信息（此时为空））
产品转让，指定数量给检测机构
org3获得产品，填写检测信息（营养成分，检测时间）
org3还给org2

CAFILE=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem
TLSRootCertFiles=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt
peer chaincode install -p github.com/chaincode/imp -n cc -v 0
peer chaincode instantiate -n cc -v 0 -c '{"Args":[]}' -C $CHANNEL_NAME --tls --cafile $CAFILE
peer chaincode query -n cc -c '{"Args":["readMilk", "milk1"]}' -C $CHANNEL_NAME
peer chaincode invoke -C mychannel -n cc  --cafile $CAFILE -c '{"Args":["queryOrg","Org1MSP"]}' --tls true
peer chaincode invoke -C mychannel -n cc  --cafile $CAFILE -c '{"Args":[query,a]}' --tls true
peer chaincode invoke -C mychannel -n cc  --cafile $CAFILE -c '{"Args":["initMilkPowder","1","4234","123123","432423"]}' --tls true
peer chaincode query -n cc -c '{"Args":["queryOrg"]}' -C $CHANNEL_NAME --cafile $CAFILE

--cafile $CAFILE --peerAddresses peer0.org1.example.com:9051
--tlsRootCertFiles $TLSRootCertFiles
--peerAddresses peer0.org2.example.com:9051 --tlsRootCertFiles /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/tls/ca.crt
peer chaincode invoke -C tracechannel -n cc -c '{"Args":["initMilkPowder","milk1","10000000","desc","cattleinfo"]}'
peer chaincode invoke -C tracechannel -n cc -c '{"Args":["transferMilkPowder","Org2MSP","milk1","10000000"]}'
资产：奶粉，产品
*/
package main

import (
	"fmt"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/protos/peer"
)

type ImpChaincode struct {
}

type orgInfo struct {
	MspId string `json:"msp_id"`
	Desc  string `json:"desc"` //描述信息
}

type milkPowder struct {
	Id         string `json:"id"` // 出厂编号
	Time       int64  `json:"time"`
	Weight     int64  `json:"weight"`
	Desc       string `json:"Desc"`
	CattleInfo string `json:"cattle_info"`
}

type product struct {
	Id   string `json:"id"` // 出厂编号
	Time int64  `json:"time"`
}
type inspectInfo struct {
	Time int64  `json:"time"` //检测时间
	Desc string `json:"Desc"` //检测信息营养成分等
}

//初始化方法
func (t *ImpChaincode) Init(stub shim.ChaincodeStubInterface) peer.Response {
	return shim.Success(nil)
}

//方法调用入口
func (t *ImpChaincode) Invoke(stub shim.ChaincodeStubInterface) peer.Response {
	function, args := stub.GetFunctionAndParameters()
	if function == "initOrg" {
		return t.initOrg(stub, args)
	} else if function == "updateOrg" {
		return t.updateOrg(stub, args)
	} else if function == "queryOrg" {
		return t.queryOrg(stub, args)
	} else if function == "initMilkPowder" {
		return t.initMilkPowder(stub, args)
	} else if function == "transferMilkPowder" {
		return t.transferMilkPowder(stub, args)
	} else if function == "queryMilkPowder" {
		return t.queryMilkPowder(stub, args)
	} else if function == "queryMilkPowderAsset" {
		return t.queryMilkPowderAsset(stub, args)
	}
	return shim.Error("no this invoke function:" + function)
}

func main() {
	err := shim.Start(new(ImpChaincode))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}
}
