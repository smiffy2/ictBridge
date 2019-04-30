package main

import (
        "fmt"
        "github.com/smiffy2/ictBridge"
        . "github.com/smiffy2/ictBridge/proto"
)

func main () {

	client := ictBridge.CreateIctBridgeClient("35.204.80.128","7331")

        address := "TEST9ADDRESSTONYSTJH9999999999999999999999999999999999999999999999999999999999999"
	tag := "BRIDGE9TESTTONY999999999999"
        transaction := TransactionBuilder { Address:address,Tag:tag,Value:ictBridge.IntToBytes(500)}

        client.SubmitTransaction(transaction)

	//trans,err := client.QueryByAddress(address)
	//if(err != nil) {
	//	panic(err)
	//}
	//fmt.Printf("Transacion value = %v",trans[0].Value)

	trans, err := client.QueryByTag(tag)
	if(err != nil) {
		panic(err)
	}
	if(trans != nil) {
        	for _,v := range trans {
                	fmt.Println(v.Address)
        	}
	}
}

