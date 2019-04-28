package ictBridge

import (
	"net"
	"github.com/golang/protobuf/proto"
	. "github.com/smiffy2/ictBridge/proto"
	"bytes"
	"encoding/binary"
)

func CreateIctBridgeClient(ipaddress string, port string) *IctBridgeClient {

	conn,err  := net.Dial("tcp", ipaddress + ":" + port)
	if(err != nil) {
		return nil
	}
	return &IctBridgeClient{Conn:conn}
 
}

type IctBridgeClient struct {
	Conn net.Conn
}

func (ict IctBridgeClient) SendMessage (mess WrapperMessage) error {

	err := processMessage(ict.Conn, mess)
	if(err != nil) {
		return err
	}
	return nil
}

func (ict IctBridgeClient) SendQuery (mess WrapperMessage) (WrapperMessage, error) {

	err := processMessage(ict.Conn, mess)
	if(err != nil) {
		return WrapperMessage{},err
	}
	return ict.GetMessage()
}

func (ict *IctBridgeClient) GetMessage() (WrapperMessage,error) {

	reply := WrapperMessage{}
	buf, err := readBytes(ict.Conn,4)
	if(err != nil) {
		return reply,err 
	}
	
	buf ,err = readBytes(ict.Conn,bytesToInt(buf))
	if(err != nil) {
		return reply,err
	}

	err = proto.Unmarshal(buf,&reply)
	if(err != nil) {
		return reply,err
	}

	return reply,nil	
}

func (ict *IctBridgeClient) QueryByAddress(address string) ([]*Transaction,error) {

	wrapperQuery := WrapperMessage{
                MessageType: WrapperMessage_FIND_TRANSACTIONS_BY_ADDRESS_REQUEST,
                Msg: &WrapperMessage_FindTransactionsByAddressRequest{
                        &FindTransactionsByAddressRequest{Address:address}},
        }

	reply,err := ict.SendQuery(wrapperQuery)

	if(err != nil) {
		return nil, err
	}

	return  reply.GetFindTransactionsByAddressResponse().Transaction,nil
}

func (ict *IctBridgeClient) QueryByTag(tag string) ([]*Transaction,error) {

	wrapperQuery := WrapperMessage{
                MessageType: WrapperMessage_FIND_TRANSACTIONS_BY_TAG_REQUEST,
                Msg: &WrapperMessage_FindTransactionsByTagRequest{
                        &FindTransactionsByTagRequest{Tag:tag}},
        }

	reply,err := ict.SendQuery(wrapperQuery)

	if(err != nil) {
		return nil, err
	}

	return  reply.GetFindTransactionsByTagResponse().Transaction,nil
}

func (ict *IctBridgeClient) SubmitTransaction(transaction TransactionBuilder) error {

	wrapperNewTransaction := WrapperMessage {
                MessageType: WrapperMessage_SUBMIT_TRANSACTION_BUILDER_REQUEST,
                Msg:  &WrapperMessage_SubmitTransactionBuilderRequest{
                                &SubmitTransactionBuilderRequest{
                                        TransactionBuilder:&transaction,
                                },
                },
        }

	ict.SendMessage(wrapperNewTransaction)

	return nil

}

func processMessage(conn net.Conn, msg WrapperMessage) error {

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
 
