# ictBridge

To download : go get github.com/smiffy2/ictBridge

Example of how to use

package main

import (
        "fmt"
        "github.com/smiffy2/ictBridge"
        . "github.com/smiffy2/ictBridge/proto"
)

func main () {

        client := ictBridge.CreateIctBridgeClient("35.204.80.128","7331")

        address := "TEST9ADDRESSTONYSDE99999999999999999999999999999999999999999999999999999999999999"
        wrapperNewTransaction := WrapperMessage {
                MessageType: WrapperMessage_SUBMIT_TRANSACTION_BUILDER_REQUEST,
                Msg:  &WrapperMessage_SubmitTransactionBuilderRequest{
                                &SubmitTransactionBuilderRequest{
                                        TransactionBuilder:&TransactionBuilder {
                                                Address:address,
                                                Tag:"BRIDGE9TESTTONY999999999999",
                                        },
                                },
                },
        }

        client.SendMessage(wrapperNewTransaction)

        wrapperTagQuery := WrapperMessage{
                MessageType: WrapperMessage_FIND_TRANSACTIONS_BY_TAG_REQUEST,
                Msg: &WrapperMessage_FindTransactionsByTagRequest{
                        &FindTransactionsByTagRequest{Tag:"BRIDGE9TESTTONY999999999999"}},
        }

        reply,err := client.SendQuery(wrapperTagQuery)
        if(err != nil) {
                panic(err)
        }
        for _,v := range reply.GetFindTransactionsByTagResponse().Transaction {
                fmt.Println(v.Address)
        }
}

