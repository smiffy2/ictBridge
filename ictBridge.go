//Package to help send data to the Ict Bridge.ixi interface
package ictBridge

import (
	"net"
	"github.com/golang/protobuf/proto"
	. "github.com/smiffy2/ictBridge/proto"
	"bytes"
	"encoding/binary"
)

//CreateIctBridgeClient: Creates a IctBridgeClient using given IPAddress and port
func CreateIctBridgeClient(ipaddress string, port string) *IctBridgeClient {

	conn,err  := net.Dial("tcp", ipaddress + ":" + port)
	if(err != nil) {
		return nil
	}
	return &IctBridgeClient{Conn:conn}
 
}

//Class for IctBridgeClient
type IctBridgeClient struct {
	Conn net.Conn
}

//SendMessage: Sends WrapperMessage to Bridge.ixi
func (ict IctBridgeClient) SendMessage (mess WrapperMessage) error {

	err := processMessage(ict.Conn, mess)
	if(err != nil) {
		return err
	}
	return nil
}

//SendQuery: Sends a WrapperMessage to Bridge.ixi and waits for a rpely.
func (ict IctBridgeClient) SendQuery (mess WrapperMessage) (WrapperMessage, error) {

	err := processMessage(ict.Conn, mess)
	if(err != nil) {
		return WrapperMessage{},err
	}
	return ict.GetMessage()
}

//GetMessagë: Wait for message from Ict 
func (ict *IctBridgeClient) GetMessage() (WrapperMessage,error) {

	reply := WrapperMessage{}
	buf, err := readBytes(ict.Conn,4)
	if(err != nil) {
		return reply,err 
	}
	buf ,err = readBytes(ict.Conn,BytesToInt(buf))
	if(err != nil) {
		return reply,err
	}

	err = proto.Unmarshal(buf,&reply)
	if(err != nil) {
		return reply,err
	}

	return reply,nil	
}

//QueryByAddress: Send a request to Bridge.ixi to ask for information on an Address, wait for resposnse
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
	
	if(reply.GetMsg() == nil) {
		return nil,nil
	}
	return  reply.GetFindTransactionsByTagResponse().Transaction,nil
}

//QueryByTag: Send a requst to Bridge.ixi to request for data on a given tag, wait for response
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
	if(reply.GetMsg() == nil) {
		return nil,nil
	}

	return reply.GetFindTransactionsByTagResponse().Transaction, nil
}

//SubmitTransaction: Send a transaction to Ict
func (ict *IctBridgeClient) SubmitTransaction(transaction TransactionBuilder) error {

	wrapperNewTransaction := WrapperMessage {
                MessageType: WrapperMessage_SUBMIT_TRANSACTION_BUILDER_REQUEST,
                Msg:  &WrapperMessage_SubmitTransactionBuilderRequest{
                                &SubmitTransactionBuilderRequest{
                                        TransactionBuilder:&transaction,
                                },
                },
        }

	return ict.SendMessage(wrapperNewTransaction)
}

func processMessage(conn net.Conn, msg WrapperMessage) error {

	data, err := proto.Marshal(&msg)
	if(err != nil) {
		return err
	}
	_, err = conn.Write(IntToBytes(len(data)))
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

//BytesToInt: Helper method to convert bytes to int.
func BytesToInt(in []byte) int {

        var m uint32
        err := binary.Read(bytes.NewBuffer(in), binary.BigEndian, &m)
        if(err != nil) {
                return 0
        }
        return int(m)
}

//IntToBytes: Helper method to get byte representation of an int.
func IntToBytes(in int) []byte {

        buff := new(bytes.Buffer)
        binary.Write(buff, binary.BigEndian, uint32(in))
        return buff.Bytes()
}

