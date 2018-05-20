package AEDAclient

import (
	"bufio"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"strings"
	"time"

	"github.com/op/go-logging"

	"github.com/torlenor/AbyleEDA/AEDAcrypt"
	"github.com/torlenor/AbyleEDA/eventmessage"
)

var log = logging.MustGetLogger("AEDAlogger")

var rcvOK = []byte("0")
var rcvFAIL = []byte("1")

func checkError(err error) {
	if err != nil {
		log.Error("Error: ", err)
	}
}

// The UDPClient struct holds the udp client of AEDAclient
type UDPClient struct {
	DebugMode   bool
	isConnected bool
	isStarted   bool

	Conn    *net.UDPConn
	SrvAddr *net.UDPAddr

	packetQueue chan UDPPacket

	ResQueue chan ServerMessage

	ccfg AEDAcrypt.CryptCfg
}

func greetServer(Conn *net.UDPConn) error {
	// Greet the server
	fmt.Println("Sending a hi to the server...")
	p := make([]byte, 2048)
	msg := string("1001")
	buf := []byte(msg)
	_, err := Conn.Write(buf)
	checkError(err)
	if err != nil {
		log.Panic("Can't greet server!")
		return err
	}

	// TODO: Do this with timeout in response
	// TODO: Maybe resend the greeting then
	// Wait for the response
	fmt.Println("Waiting for the response...")
	_, err = bufio.NewReader(Conn).Read(p)
	checkError(err)
	if err != nil {
		log.Panic("Can't connect to server!")
		return err
	}
	fmt.Println("Got a repsonse from the server! Yey! Client ready!")

	return nil
}

func authenticateClient(Conn *net.UDPConn) (AEDAcrypt.CryptCfg, error) {
	// TODO: Do an authentication
	nonce, _ := hex.DecodeString("bb8ef84243d2ee95a41c6c57")

	ccfg := AEDAcrypt.CryptCfg{Key: []byte("AES256Key-32Characters1234567890"),
		Nonce: nonce}
	return ccfg, nil
}

// ConnectUDPClient function connects a client to a server
func ConnectUDPClient(srvAddr *net.UDPAddr) (*UDPClient, error) {
	client := &UDPClient{}

	// Define the server address and port
	Conn, err := net.DialUDP("udp", nil, srvAddr)
	checkError(err)
	if err != nil {
		log.Panic("Can't dial server!")
		return nil, err
	}

	err = greetServer(Conn)
	if err != nil {
		return nil, err
	}

	ccfg, err := authenticateClient(Conn)
	if err != nil {
		return nil, err
	}

	client.Conn = Conn
	client.SrvAddr = srvAddr
	client.ccfg = ccfg

	return client, nil
}

// DisconnectUDPClient disconnects a client from its server
func DisconnectUDPClient(client *UDPClient) {
	client.Conn.Close()
}

// SendMessageToServer sends a message to the connected server
func SendMessageToServer(client *UDPClient, event eventmessage.EventMessage) {
	// Set the current timestamp
	event.Timestamp = time.Now().UnixNano()

	msg, err := json.Marshal(event)
	checkError(err)

	// Encrypt the message
	encmsg := AEDAcrypt.Encrypter(msg, client.ccfg)

	buf := []byte(encmsg)
	msgmd5 := AEDAcrypt.GetMD5HashFromString(string(buf))

	succ := false
	cnt := 0
	for succ == false && cnt < 3 {
		tstartsend := time.Now()
		_, err := client.Conn.Write(buf)
		if err != nil {
			log.Error(msg, err, "... quitting ... ")
			os.Exit(1)
		}

		c1 := make(chan bool, 1)
		c2 := make(chan bool, 1)

		var md5fromsrv string

		go func() {
			p := make([]byte, 1024)
			deadline := time.Now().Add(10 * time.Second)
			client.Conn.SetReadDeadline(deadline)
			n, _, err := client.Conn.ReadFromUDP(p)
			checkError(err)
			if err != nil {
				c1 <- false
				return
			}

			md5fromsrv = string(p[0:n])
			if strings.Compare(md5fromsrv, msgmd5) == 0 {
				c2 <- true
			} else {
				c2 <- false
			}
		}()

		select {
		case res := <-c1:
			duration := time.Now().Sub(tstartsend)
			if !res {
				log.Error("Problem receiving answer from server (", msgmd5, "), (", duration.String(), ")")
				if cnt < 3 {
					log.Error("Trying again...")
					time.Sleep(time.Second * 1)
				}
				cnt++
			}

		case res := <-c2:
			duration := time.Now().Sub(tstartsend)
			if res {
				log.Info("Server said package was OK (", msgmd5, "), (", duration.String(), ")")
				succ = true
			} else {
				log.Error("Server said package was NOT OK, (", duration.String(), ")")
				log.Error("MD5 package sent:	", msgmd5)
				log.Error("MD5 received:		", md5fromsrv)
				if cnt < 3 {
					log.Error("Trying again...")
					time.Sleep(time.Second * 1)
				}
				cnt++
			}

		case <-time.After(time.Second * 30):
			log.Error("Received no answer from server!")
		}
	}
}
