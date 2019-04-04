package main

/**
组织信息操作
*/
import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/hyperledger/fabric/core/chaincode/lib/cid"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/protos/peer"
)

func (t *ImpChaincode) initOrg(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1 desc")
	}
	mspid, err := cid.GetMSPID(stub)
	if err != nil {
		return shim.Error("get id failed" + err.Error())
	}
	org := &orgInfo{mspid, args[0]}
	orgJson, err := json.Marshal(org)
	if err != nil {
		return shim.Error(err.Error())
	}
	err = stub.PutState(mspid, orgJson)
	if err != nil {
		return shim.Error("put state error")
	}
	return shim.Success(nil)
}

func (t *ImpChaincode) updateOrg(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1 desc")
	}
	mspid, err := cid.GetMSPID(stub)
	if err != nil {
		return shim.Error("get id failed" + err.Error())
	}
	orgBytes, err := stub.GetState(mspid)
	if err != nil {
		return shim.Error("get state error" + err.Error())
	}
	if orgBytes == nil {
		return shim.Error("no this org:" + mspid)
	}
	//反序列化
	org := orgInfo{}
	err = json.Unmarshal(orgBytes, &org)
	if err != nil {
		return shim.Error("json unmarshal error" + err.Error())
	}
	org.Desc = args[0]
	orgBytes, err = json.Marshal(org)
	if err != nil {
		return shim.Error("json marsha :" + err.Error())
	}

	err = stub.PutState(mspid, orgBytes)
	if err != nil {
		return shim.Error("put state error" + err.Error())
	}
	return shim.Success(nil)
}

func getQueryResultForQueryString(stub shim.ChaincodeStubInterface, queryString string) ([]byte, error) {

	fmt.Printf("- getQueryResultForQueryString queryString:\n%s\n", queryString)
	resultsIterator, err := stub.GetQueryResult(queryString)
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()
	buffer, err := constructQueryResponseFromIterator(resultsIterator)
	if err != nil {
		return nil, err
	}
	fmt.Printf("- getQueryResultForQueryString queryResult:\n%s\n", buffer.String())
	return buffer.Bytes(), nil
}
func constructQueryResponseFromIterator(resultsIterator shim.StateQueryIteratorInterface) (*bytes.Buffer, error) {
	// buffer is a JSON array containing QueryResults
	var buffer bytes.Buffer
	buffer.WriteString("[")
	bArrayMemberAlreadyWritten := false
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}
		// Add a comma before array members, suppress it for the first array member
		if bArrayMemberAlreadyWritten == true {
			buffer.WriteString(",")
		}
		buffer.WriteString("{\"Key\":")
		buffer.WriteString("\"")
		buffer.WriteString(queryResponse.Key)
		buffer.WriteString("\"")

		buffer.WriteString(", \"Record\":")
		// Record is a JSON object, so we write as-is
		buffer.WriteString(string(queryResponse.Value))
		buffer.WriteString("}")
		bArrayMemberAlreadyWritten = true
	}
	buffer.WriteString("]")

	return &buffer, nil
}

func (t *ImpChaincode) queryOrg(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1 mspId")
	}
	mspid := args[0]
	orgByte, err := stub.GetState(mspid) //get the milk from chaincode state
	jsonResp := ""
	if err != nil {
		jsonResp = "{\"Error\":\"Failed to get state for " + mspid + "\"}"
		return shim.Error(jsonResp)
	} else if orgByte == nil {
		jsonResp = "{\"Error\":\"org does not exist: " + mspid + "\"}"
		return shim.Error(jsonResp)
	}
	return shim.Success(orgByte)
}
