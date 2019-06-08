package main

import (
	"fmt"

	"github.com/smiffy2/ictBridge"
	"time"
)

func main() {

	tag := "BRIDGE9TESTTONYA99999999999"
	address := "TEST9ADDRESSTONYSTKHG999999999999999999999999999999999999999999999999999999999999"
	msg := ictBridge.IctBridgeMessage{Address:address,Tag:tag,Message:"Hello Tony Again, one last message",Publish:true}

	client := ictBridge.CreateIctBridgeClient("35.204.45.126","7331")
	client1 := ictBridge.CreateIctBridgeClient("35.204.80.128","7331")

	client.SubmitMessage(msg)
	time.Sleep(5 * time.Second)
	//retTrans, err := client1.QueryByTag(tag)
	retTrans, err := client1.QueryByAddress(address)
        if(err != nil) {
                panic(err)
        }
        if(retTrans != nil) {
                for _,v := range retTrans {
                        fmt.Println(v.Message)
                }
        } else {
		fmt.Println("Nothing found")
	}

}

