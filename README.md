# ictBridge

To download : go get github.com/smiffy2/ictBridge

Example of how to use
```
package main

import (
        "fmt"
        "github.com/smiffy2/ictBridge"
        . "github.com/smiffy2/ictBridge/proto"
)

func main () {

        client := ictBridge.CreateIctBridgeClient("127.0.0.1","7331")

        address := "TEST9ADDRESSTONYSDJ99999999999999999999999999999999999999999999999999999999999999"
        tag := "BRIDGE9TESTTONY999999999999"
        transaction := TransactionBuilder { Address:address,Tag:tag}

        client.SubmitTransaction(transaction)

        trans,err := client.QueryByAddress(address)
        if(err != nil) {
                panic(err)
        }
        fmt.Printf("Tag for address %v is %v\n",trans[0].Address,trans[0].Tag)

        trans, err = client.QueryByTag(tag)
        for _,v := range trans {
                fmt.Println(v.Address)
        }
}
```
