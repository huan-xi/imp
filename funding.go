package main

import (
	"fmt"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/protos/peer"
)

/**
众筹链代码
用户
众筹项目
*/

type FundingChaincode struct {
}

//初始化方法
func (t *FundingChaincode) Init(stub shim.ChaincodeStubInterface) peer.Response {
	return shim.Success(nil)
}

//方法调用入口
func (t *FundingChaincode) Invoke(stub shim.ChaincodeStubInterface) peer.Response {
	function, args := stub.GetFunctionAndParameters()
	if function == "get" {
		return t.initOrg(stub, args)
	} else if function == "set" {
		return t.updateOrg(stub, args)
	} else if function == "delete" {
		return t.queryOrg(stub, args)
	}

	return shim.Error("no this invoke function:" + function)
}
func main() {
	err := shim.Start(new(FundingChaincode))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}
}
