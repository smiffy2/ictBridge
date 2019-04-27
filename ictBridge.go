package main

import (
	"net"
	"fmt"
	"github.com/golang/protobuf/proto"
	. "github.com/smiffy2/ictBridge/proto"
	"bytes"
	"encoding/binary"
)

func main() {

	conn, err := net.Dial("tcp", "35.204.45.126:7331")

	if(err != nil) {
		fmt.Println("Unable to connect\n")
		panic(err)
	}

	defer conn.Close()
	

	address := "TEST9ADDRESSTONYSDD99999999999999999999999999999999999999999999999999999999999999"
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

	sendMessage(conn, wrapperNewTransaction)
	
	wrapperQuery := WrapperMessage{
		MessageType: WrapperMessage_FIND_TRANSACTIONS_BY_ADDRESS_REQUEST,
		Msg: &WrapperMessage_FindTransactionsByAddressRequest{
			&FindTransactionsByAddressRequest{Address:address}},
	}

	sendMessage(conn, wrapperQuery)
	
	reply, err := getMessage(conn)
	
	fmt.Println(reply.GetFindTransactionsByAddressResponse().Transaction[0].Tag)

	wrapperTagQuery := WrapperMessage{
                MessageType: WrapperMessage_FIND_TRANSACTIONS_BY_TAG_REQUEST,
                Msg: &WrapperMessage_FindTransactionsByTagRequest{
                        &FindTransactionsByTagRequest{Tag:"BRIDGE9TESTTONY999999999999"}},
        }

	sendMessage(conn, wrapperTagQuery)
	reply, err = getMessage(conn)
	if(err != nil) {
		panic(err)
	}
	transaction := reply.GetFindTransactionsByTagResponse().Transaction
	for _,v := range transaction {
		fmt.Println(v.Address)
	}	
}

func sendMessage(conn net.Conn, msg WrapperMessage) error {

	data, err := proto.Marshal(&msg)
	if(err != nil) {
		return err
	}
	buff := new(bytes.Buffer)
	err = binary.Write(buff, binary.BigEndian, uint32(len(data)))

	_, err = conn.Write(buff.Bytes())
	if(err != nil) {
		return err
	}
	_, err = conn.Write(data)
	if(err != nil) {
		return err
	}

	return nil
}

func getMessage(conn net.Conn) (WrapperMessage,error) {

	reply := WrapperMessage{}
	buf, err := readBytes(conn,4)
	if(err != nil) {
		return reply,err 
	}
	
	buf ,err = readBytes(conn,bytesToInt(buf))
	if(err != nil) {
		return reply,err
	}

	err = proto.Unmarshal(buf,&reply)
	if(err != nil) {
		return reply,err
	}

	return reply,nil	
}

func readBytes(conn net.Conn, len int) ([]byte, error) {
	
	readBytes := 0
	var returnBytes []byte
	for readBytes < len {
		buff := make([]byte,len - readBytes)
		n, err := conn.Read(buff)
		if(err != nil) {
			return returnBytes,nil
		}
		readBytes = readBytes + n
		returnBytes = append(returnBytes,buff[:n]...)
	}
	return returnBytes, nil
}

func bytesToInt(in []byte) int {

	var m uint32 
	err := binary.Read(bytes.NewBuffer(in), binary.BigEndian, &m)
	if(err != nil) {
		return 0
	}
	return int(m)
}
 
