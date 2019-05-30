# ictBridge

To download : go get github.com/smiffy2/ictBridge

Example of how to use
```
package main

import (
        "fmt"

        "github.com/smiffy2/ictBridge"
        "time"
)

func main() {

        tag := "BRIDGE9TESTTONYA99999999999"
        msg := ictBridge.IctBridgeMessage{Address:"TEST9ADDRESSTONYSTKHG9999999999999999999999999999999999
99999999999999999999999999",Tag:tag,Message:"Hello Tony Again, one last message",Publish:true}

        client := ictBridge.CreateIctBridgeClient("127.0.0.1","7331")

        client.SubmitMessage(msg)
     
	retTrans, err := client.QueryByTag(tag)
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

```
