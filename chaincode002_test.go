package main

import (
	"fmt"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"testing"
)

func checkInit(t *testing.T, stub *shim.MockStub, args [][]byte) {
	res := stub.MockInit("1", args)
	if res.Status != shim.OK {
		fmt.Println("Init failed", string(res.Message))
		t.FailNow()
	}
}

func checkQuery(t *testing.T, stub *shim.MockStub, name string) {
	res := stub.MockInvoke("1", [][]byte{[]byte("get"), []byte(name)})
	if res.Status != shim.OK {
		fmt.Println("Query", name, "failed", string(res.Message))
		t.FailNow()
	}
	if res.Payload == nil {
		fmt.Println("Query", name, "failed to get value")
		t.FailNow()
	}

	fmt.Println("Query value", name, "was ", string(res.Payload))

}

func checkInvoke(t *testing.T, stub *shim.MockStub, args [][]byte) {
	res := stub.MockInvoke("1", args)
	if res.Status != shim.OK {
		fmt.Println("Invoke", args, "failed", string(res.Message))
		t.FailNow()
	}
}
func Test_Helloworld(t *testing.T) {
	hello := new(ImpChaincode)
	stub := shim.NewMockStub("hello", hello)
	//checkInvoke(t, stub, [][]byte{[]byte("initOrg"),[]byte("data")})
	//checkInvoke(t, stub, [][]byte{[]byte("updateOrg"),[]byte("data")})
	checkInvoke(t, stub, [][]byte{[]byte("initMilkPowder"), []byte("id1"), []byte("10000"), []byte("fasd"), []byte("fdsf")})
	//checkInvoke(t, stub, [][]byte{[]byte("processMilk"), []byte("milk1"), []byte("datafsdsdf")})
	//checkQuery(t, stub, "queryMilksByOwner")
	//checkInvoke(t, stub, [][]byte{[]byte("readMilk"), []byte("milk1")})
	//checkInvoke(t, stub, [][]byte{[]byte("queryMilksByOwner"), []byte("bob")})
}
