package AEDAclient

import (
	"errors"
	"net"

	"github.com/torlenor/AbyleEDA/AEDAcrypt"
)

type ServerMessage struct {
	Addr *net.UDPAddr
	Msg  []byte
}

type UDPPacket struct {
	Addr *net.UDPAddr
	Buf  []byte
}

// Start starts receiving UDP packages comming from server to client
func (client *UDPClient) startReceiveChannel() error {
	if client.isStarted == false {
		buf := make([]byte, 64*1024) // theoretical max of UDP package

		client.isStarted = true
		go client.startWorker()

		for {
			n, addr, err := client.Conn.ReadFromUDP(buf)
			if err != nil {
				log.Error(err)
			}

			rcvmsg := make([]byte, len(buf[0:n]))
			copy(rcvmsg, buf[:])

			select {
			case client.packetQueue <- UDPPacket{Addr: addr, Buf: rcvmsg}:
			default:
				log.Error("Client packetQueue channel full. Discarding message from server!")
			}
		}
	}

	return errors.New("Server already running")
}

func (client *UDPClient) startWorker() {
	for {
		select {
		case pkt := <-client.packetQueue:
			go parseUDPMessage(client, pkt.Addr, pkt.Buf)
		}
	}
}

func parseUDPMessage(client *UDPClient, addr *net.UDPAddr, buf []byte) {
	msgmd5 := AEDAcrypt.GetMD5HashFromString(string(buf))
	msg, err := AEDAcrypt.Decrypter(buf, client.ccfg)

	if err != nil {
		client.Write(addr, rcvFAIL)
		return
	}

	select {
	case client.ResQueue <- ServerMessage{Addr: addr, Msg: msg}:
		client.Write(addr, []byte(msgmd5))
	default:
		log.Error("Client packetQueue channel full. Discarding message from server!")
		client.Write(addr, rcvFAIL)
	}
}

func (client *UDPClient) Write(addr *net.UDPAddr, msg []byte) {
	_, err := client.Conn.WriteToUDP(msg, addr)
	checkError(err)
}
