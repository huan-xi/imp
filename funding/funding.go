package main

import (
	"fmt"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/protos/peer"
)

/**
peer chaincode install -p github.com/chaincode/funding -n funding -v 1.0.0
peer chaincode instantiate -C mychannel -n funding -v 1.0.0 -c '{"Args":[]}'
peer chaincode upgrade -C mychannel  -o orderer.example.com:7050 -n funding -v 1.0.2 -c '{"Args":[]}'
peer chaincode query -n mycc -c '{"Args":["query","{\"selector\": {}}"]}' -C mychannel



众筹链代码
用户
众筹项目
注册用户
登入
添加
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
		return t.get(stub, args)
	} else if function == "set" {
		return t.set(stub, args)
	} else if function == "delete" {
		return t.delete(stub, args)
	} else if function == "update" {
		return t.update(stub, args)
	}

	return shim.Error("no this invoke function:" + function)
}
func main() {
	err := shim.Start(new(FundingChaincode))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}
}

func (t *FundingChaincode) set(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments. Expecting 2 key,value")
	}

	key := args[0]
	value := args[1]
	valAsbytes, err := stub.GetState(key)
	if valAsbytes != nil {
		return shim.Error("already exists this key")
	}
	err = stub.PutState(key, []byte(value))
	if err != nil {
		return shim.Error("put state error")
	}
	return shim.Success([]byte("ok"))
}

func (t *FundingChaincode) get(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1 key")
	}
	key := args[0]
	valAsbytes, err := stub.GetState(key) //get the milk from chaincode state
	if err != nil {
		jsonResp := "get statue error"
		return shim.Error(jsonResp)
	} else if valAsbytes == nil {
		jsonResp := "not exist this id"
		return shim.Error(jsonResp)
	}
	return shim.Success(valAsbytes)
}

func (t *FundingChaincode) delete(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1 key")
	}
	key := args[0]
	err := stub.DelState(key) //get the milk from chaincode state
	if err != nil {
		jsonResp := "Failed to get state for " + key + "\"}"
		return shim.Error(jsonResp)
	}
	return shim.Success([]byte("ok"))
}

func (t *FundingChaincode) update(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments. Expecting 2 key,value")
	}
	key := args[0]
	value := args[1]
	valAsbytes, err := stub.GetState(key) //get the milk from chaincode state
	if err != nil {
		jsonResp := "get statue error"
		return shim.Error(jsonResp)
	} else if valAsbytes == nil {
		jsonResp := "not exist this id"
		return shim.Error(jsonResp)
	}
	err = stub.PutState(key, []byte(value))
	if err != nil {
		return shim.Error("put state error")
	}
	return shim.Success([]byte("ok"))
}
