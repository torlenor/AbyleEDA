package AEDAserver

import (
	"errors"
	"net"
	"strconv"

	"github.com/op/go-logging"
	"github.com/torlenor/AbyleEDA/AEDAcrypt"
)

// Simple OK/NOTOK for the client
var rcvOK = []byte("0")
var rcvFAIL = []byte("1")

var maxQueue int = 12

var log = logging.MustGetLogger("AEDAlogger")

type ServerWriter interface {
	Write(addr *net.UDPAddr, msg []byte)
}

type Server interface {
	Start()
	Close()
	ServerWriter
}

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
		srv.Write(addr, rcvFAIL)
		srv.Stats.Pktserrcnt++
		return
	}

	srv.ResQueue <- ClientMessage{Addr: addr, Msg: msg}
	srv.Write(addr, []byte(msgmd5))
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

// Start starts receiving UDP packages for the server
func (srv *UDPServer) Start() error {
	if srv.isStarted == false {
		buf := make([]byte, 64*1024) // theoretical max of UDP package

		srv.isStarted = true
		go srv.startWorker()

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

// Close closes the server and no more packages are read from UDP
func (srv *UDPServer) Close() {
	srv.Conn.Close()
}

func (srv *UDPServer) Write(addr *net.UDPAddr, msg []byte) {
	_, err := srv.Conn.WriteToUDP(msg, addr)
	checkError(err)
	srv.Stats.Pktssentcnt++
}

func (srv *UDPServer) startWorker() {
	for {
		select {
		case pkt := <-srv.packetQueue:
			go parseUDPMessage(srv, pkt.Addr, pkt.Buf)
		}
	}
}

func checkError(err error) {
	if err != nil {
		log.Error("Error: ", err)
	}
}

func addNewClient(srv *UDPServer, addr *net.UDPAddr) {
	srv.Clients = appendClient(srv.Clients, *addr)
	log.Info("New client (", addr, ") connected ... greeting it!")
	srv.Write(addr, []byte("From server: Hello I got your mesage "))
}

// CreateUDPServer returns a new UDP server with the provided parameters
func CreateUDPServer(port int, ccfg AEDAcrypt.CryptCfg) (*UDPServer, error) {
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
	srv.ccfg = ccfg

	return srv, nil
}
