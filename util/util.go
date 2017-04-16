package util

import (
	"crypto/rand"
	"encoding/binary"
	"fmt"
	"net"
	"os"
)

func CheckError(err error) {
	if err != nil {
		fmt.Println("Error: ", err)
		os.Exit(0)
	}
}

func RandomString(size int) string {
	buf := make([]byte, size)
	_, err := rand.Read(buf)
	CheckError(err)
	return string(buf)
}

func ParseIP(addr string) (*net.UDPAddr, error) {
	ip := net.IPv4(addr[0], addr[1], addr[2], addr[3])
	port := binary.BigEndian.Uint16([]byte(addr[4:]))

	return net.ResolveUDPAddr("udp", fmt.Sprintf("%v:%d", ip, port))
}
