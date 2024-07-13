package internal

import (
	"bytes"
	"log"
	"net"
	"time"
)

const updMaxPayloadSize = 65507

func forwardAndListenForResp(
	targetAddr net.UDPAddr,
	data []byte,
	respConn *net.UDPConn,
	respAddr net.UDPAddr,
	waitResponseDuration time.Duration,
) {
	listenAddr := &net.UDPAddr{}

	targetConn, err := net.ListenUDP("udp", listenAddr)

	if err != nil {
		log.Printf("response listener could not be started: %v\n", err)
		return
	}

	defer targetConn.Close()

	log.Printf("response listener started on: %v\n", targetConn.LocalAddr())

	_, err = targetConn.WriteToUDP(data, &targetAddr)

	if err != nil {
		log.Printf("response listener could not forward data: %v\n", err)
		return
	}

	var buf [updMaxPayloadSize]byte

	targetConn.SetReadDeadline(time.Now().Add(waitResponseDuration))

	for {
		n, addr, err := targetConn.ReadFromUDP(buf[:])

		if err, ok := err.(net.Error); ok && err.Timeout() {
			log.Printf("response listener timed out")
			break
		} else if err != nil {
			log.Printf("response listener read error: %v\n", err)
			continue
		}

		log.Printf("response listener recevied %v bytes form %v\n", n, addr)

		_, err = respConn.WriteToUDP(buf[0:n], &respAddr)

		if err != nil {
			log.Printf("response listener could not forward response: %v\n", err)
		}
	}
}

func RunServer(
	listenAddr *net.UDPAddr,
	targetAddr *net.UDPAddr,
	waitResponseDuration time.Duration,
) error {
	listenConn, err := net.ListenUDP("udp", listenAddr)

	if err != nil {
		return err
	}

	defer listenConn.Close()

	log.Printf("server started on %v, target address is %v", listenAddr, targetAddr)

	var buf [updMaxPayloadSize]byte

	for {
		n, addr, err := listenConn.ReadFromUDP(buf[:])

		if err != nil {
			log.Println("server udp read error: ", err)
			continue
		}

		log.Printf("server received %v bytes form %v", n, addr)

		go forwardAndListenForResp(*targetAddr, bytes.Clone(buf[0:n]), listenConn, *addr, waitResponseDuration)
	}
}
