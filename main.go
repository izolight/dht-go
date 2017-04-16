package main

import (
	"bytes"
	"crypto/rand"
	"encoding/binary"
	"fmt"
	"github.com/marksamman/bencode"
	"net"
	"os"
)

type DHTResponse map[string]interface{}

type NodeInfo struct {
	id   string
	addr *net.UDPAddr
}

func CheckError(err error) {
	if err != nil {
		fmt.Println("Error: ", err)
		os.Exit(0)
	}
}

func randomString(size int) string {
	buf := make([]byte, size)
	_, err := rand.Read(buf)
	CheckError(err)
	return string(buf)
}

func findNodesQuery(id string) []byte {
	q := make(map[string]interface{})
	q["y"] = "q"
	q["q"] = "find_node"
	q["t"] = randomString(2)
	a := make(map[string]interface{})
	a["id"] = id
	a["target"] = randomString(20)
	q["a"] = a

	fmt.Printf("Sending: TX: %x\t Target: %x\n", q["t"], a["target"])

	return bencode.Encode(q)
}

//func (n NodeAddr) String() string {
//	ip := net.IPv4(n[0], n[1], n[2], n[3])
//	port := binary.BigEndian.Uint16([]byte(n[4:]))
//	return fmt.Sprintf("%v:%d", ip, port)
//}

func parseIP(addr string) (*net.UDPAddr, error) {
	ip := net.IPv4(addr[0], addr[1], addr[2], addr[3])
	port := binary.BigEndian.Uint16([]byte(addr[4:]))

	return net.ResolveUDPAddr("udp", fmt.Sprintf("%v:%d", ip, port))
}

func (d DHTResponse) String() string {
	nodeAddr, _ := parseIP(d["ip"].(string))
	tx := d["t"].(string)
	r := d["r"].(map[string]interface{})
	id := r["id"].(string)
	n := r["nodes"].(string)
	nodes := []NodeInfo{}
	for i := 0; i < len(n); {
		id := n[i : i+20]
		addr, _ := parseIP(n[i+20 : i+26])
		node := NodeInfo{id, addr}
		fmt.Printf("%s %x\n", addr, id)
		i = i + 26
		nodes = append(nodes, node)
	}

	//return fmt.Sprintf("Receiving: TX: %x ID: %x\nNodes: %v", tx, id, nodes)
	return fmt.Sprintf("Receiving: %s TX: %x ID: %x\n", nodeAddr, tx, id)
}

func main() {
	ServerAddr, err := net.ResolveUDPAddr("udp", ":12343")
	CheckError(err)

	id := randomString(20)
	//secret := randomString(4)
	buf := make([]byte, 65536)
	fmt.Printf("Started on: %v with id: %x\n", ServerAddr, id)

	bootstrapNode, err := net.ResolveUDPAddr("udp", "router.bittorrent.com:6881")
	CheckError(err)

	conn, err := net.DialUDP("udp", ServerAddr, bootstrapNode)
	CheckError(err)
	defer conn.Close()
	conn.Write(findNodesQuery(id))

	for {
		n, err := conn.Read(buf)
		CheckError(err)
		r := bytes.NewReader(buf[0:n])

		t, err := bencode.Decode(r)
		CheckError(err)

		res := DHTResponse(t)

		fmt.Printf("%v\n", res)
	}
}
