//Package to help send data to the Ict Bridge.ixi interface
package ictBridge

import (
	"net"
	"fmt"
	"github.com/golang/protobuf/proto"
	. "github.com/smiffy2/ictBridge/proto"
	. "github.com/iotaledger/iota.go/converter"
	"github.com/iotaledger/iota.go/address"
	"github.com/iotaledger/iota.go/trinary"
	"bytes"
	"encoding/binary"
	"time"
	"math/rand"
)

//CreateIctBridgeClient: Creates a IctBridgeClient using given IPAddress and port
func CreateIctBridgeClient(ipaddress string, port string) *IctBridgeClient {

	conn,err  := net.Dial("tcp", ipaddress + ":" + port)
	seed := generateRandomSeed()
	if(err != nil) {
		return nil
	}
	return &IctBridgeClient{Conn:conn,seed:seed}
 
}

//struct for IctBridgeClient
type IctBridgeClient struct {
	Conn net.Conn
	seed string
	currentIndex uint64
}

//struct for IctBridgeMessage
type IctBridgeMessage struct {
	Address string
	Tag string
	Message string
	Publish bool
	Value int
}

func createTransactionFromMsg(msg IctBridgeMessage) (TransactionBuilder,error) {

	var transaction TransactionBuilder
	var err error
	if(msg.Address == "") {
		trytesAddress,err := address.GenerateAddresses(generateRandomSeed(),1,1,2)
		if(err != nil) {
			return TransactionBuilder{},err
		}
		transaction.Address = trytesAddress[0]
	} else {
		transaction.Address = msg.Address
	}
	if(msg.Message != "") {
		msg.Message,err = ASCIIToTrytes(msg.Message)
		if(err != nil) {
			return TransactionBuilder{},err
		}
		transaction.SignatureFragments = trinary.Pad(msg.Message,2187)
	}

	if(msg.Value > 0) {
		transaction.Value = IntToBytes(msg.Value)
	}	

	transaction.Tag = msg.Tag

	if(msg.Publish) {
		transaction.IsBundleHead = true
		transaction.IsBundleTail = true
	}

	return transaction,nil
}

//SendMessage: Sends WrapperMessage to Bridge.ixi
func (ict IctBridgeClient) sendMessage (mess WrapperMessage) error {

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
func (ict *IctBridgeClient) QueryByAddress(address string) ([]IctBridgeMessage,error) {

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
	
	trans :=  reply.GetFindTransactionsByAddressResponse().Transaction
        if(trans == nil) {
                return nil,nil
        }

        var messages []IctBridgeMessage

        for _,v := range trans {
                mess := IctBridgeMessage{}
                mess.Message,err = TrytesToASCII(v.SignatureFragments + "9")
                if(err != nil) {
                        fmt.Println(err)
                        mess.Message = ""
                }
                mess.Tag = v.Tag
                mess.Address = v.Address
                messages = append(messages,mess)
        }

        return messages,nil

}

//QueryByTag: Send a requst to Bridge.ixi to request for data on a given tag, wait for response
func (ict *IctBridgeClient) QueryByTag(tag string) ([]IctBridgeMessage,error) {

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

	trans :=  reply.GetFindTransactionsByTagResponse().Transaction
	if(trans == nil) {
		return nil,nil
	}

	var messages []IctBridgeMessage

	for _,v := range trans {
		mess := IctBridgeMessage{}
		mess.Message,err = TrytesToASCII(v.SignatureFragments + "9")
		if(err != nil) {
			fmt.Println(err)
			mess.Message = ""
		}
		mess.Tag = v.Tag
		mess.Address = v.Address
		messages = append(messages,mess)		
	}

	return messages,nil
}

//SubmitTransaction: Send a transaction to Ict
func (ict *IctBridgeClient) SubmitMessage(mess IctBridgeMessage) error {

	transaction,err := createTransactionFromMsg(mess)
	if(err != nil) {
		return err
	}
	wrapperNewTransaction := WrapperMessage {
                MessageType: WrapperMessage_SUBMIT_TRANSACTION_BUILDER_REQUEST,
                Msg:  &WrapperMessage_SubmitTransactionBuilderRequest{
                                &SubmitTransactionBuilderRequest{
                                        TransactionBuilder:&transaction,
                                },
                },
        }

	return ict.sendMessage(wrapperNewTransaction)
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

func generateRandomSeed() (string) {

	validCharacters := "ABCDEFGHIJKLMNOPQRSTUVWXYZ9"
	rand.Seed(time.Now().UnixNano())
	var seed  string
	for i:=0; i<81; i++ {
		num := rand.Intn(26)
		seed = seed + validCharacters[num:num+1]
	}
	return seed
}
