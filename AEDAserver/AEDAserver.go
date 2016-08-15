package AEDAserver

import (
	"encoding/hex"
	"errors"
	"net"
	"strconv"

	"github.com/op/go-logging"
	"github.com/torlenor/AbyleEDA/AEDAcrypt"
)

// Simple OK/NOTOK for the client
var rcvOK = []byte("0")
var rcvFAIL = []byte("1")

var log = logging.MustGetLogger("AEDAlogger")

func parseUDPMessage(srv *UDPServer, addr *net.UDPAddr, buf []byte) {
	srv.Stats.Pktsrecvcnt++

	if string(buf) == "1001" {
		addNewClient(srv, addr)
		log.Info("Current clients:", srv.Clients)
		return
	}

	msgmd5 := AEDAcrypt.GetMD5HashFromString(string(buf))
	msg, err := AEDAcrypt.Decrypter(buf, srv.ccfg)

	if err != nil {
		SendUDPmsg(*srv, addr, rcvFAIL)
		srv.Stats.Pktserrcnt++
		return
	}

	srv.ResQueue <- ClientMessage{Addr: addr, Msg: msg}
	SendUDPmsg(*srv, addr, []byte(msgmd5))
}

func appendClient(slice []net.UDPAddr, addr net.UDPAddr) []net.UDPAddr {
	n := len(slice)
	if n == cap(slice) {
		newSlice := make([]net.UDPAddr, len(slice), 2*len(slice)+1)
		copy(newSlice, slice)
		slice = newSlice
	}
	slice = slice[0 : n+1]
	slice[n] = addr
	return slice
}

type SrvStats struct {
	Pktssentcnt int
	Pktsrecvcnt int
	Pktserrcnt  int
}

type ClientMessage struct {
	Addr *net.UDPAddr
	Msg  []byte
}

type UDPPacket struct {
	Addr *net.UDPAddr
	Buf  []byte
}

type UDPServer struct {
	DebugMode bool
	isStarted bool

	Conn    *net.UDPConn
	Addr    *net.UDPAddr
	Clients []net.UDPAddr
	Stats   SrvStats

	packetQueue chan UDPPacket

	ResQueue chan ClientMessage

	ccfg AEDAcrypt.CryptCfg
}

var maxQueue int = 12

// TODO REPLACE WITH SOMETHING USEFUL AND CLIENT BASED
func getCryptKey() AEDAcrypt.CryptCfg {
	// TODO: Do an authentication
	nonce, _ := hex.DecodeString("bb8ef84243d2ee95a41c6c57")

	ccfg := AEDAcrypt.CryptCfg{Key: []byte("AES256Key-32Characters1234567890"),
		Nonce: nonce}
	return ccfg
}

func init() {
	// do nothing yet
}

func checkError(err error) {
	if err != nil {
		log.Error("Error: ", err)
	}
}

func SendUDPmsg(srv UDPServer, addr *net.UDPAddr, msg []byte) {
	_, err := srv.Conn.WriteToUDP(msg, addr)
	checkError(err)
	srv.Stats.Pktssentcnt++
}

func addNewClient(srv *UDPServer, addr *net.UDPAddr) {
	srv.Clients = appendClient(srv.Clients, *addr)
	log.Info("New client (", addr, ") connected ... greeting it!")
	SendUDPmsg(*srv, addr, []byte("From server: Hello I got your mesage "))
}

func CreateUDPServer(port int) (*UDPServer, error) {
	srv := &UDPServer{}

	var srvPort string = strconv.Itoa(port)
	ServerAddr, err := net.ResolveUDPAddr("udp", ":"+srvPort)
	checkError(err)

	ServerConn, err := net.ListenUDP("udp", ServerAddr)
	if err != nil {
		return nil, err
	}

	srv.Conn = ServerConn
	srv.Addr = ServerAddr
	srv.packetQueue = make(chan UDPPacket, maxQueue)
	srv.ResQueue = make(chan ClientMessage, maxQueue)
	srv.ccfg = getCryptKey()

	return srv, nil
}

func startWorker(srv *UDPServer) {
	for {
		select {
		case pkt := <-srv.packetQueue:
			go parseUDPMessage(srv, pkt.Addr, pkt.Buf)
		}
	}
}

func Start(srv *UDPServer) error {
	if srv.isStarted == false {
		buf := make([]byte, 64*1024) // until finding a better way, assume max of 64k packages

		srv.isStarted = true
		go startWorker(srv)

		for {
			n, addr, err := srv.Conn.ReadFromUDP(buf)
			if err != nil {
				log.Error(err)
			}

			rcvmsg := make([]byte, len(buf[0:n]))
			copy(rcvmsg, buf[:])

			srv.packetQueue <- UDPPacket{Addr: addr, Buf: rcvmsg}
		}
	}

	return errors.New("Server already running")
}
