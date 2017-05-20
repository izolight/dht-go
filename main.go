package main

import (
	"bytes"
	"fmt"
	"gitea.izolight.xyz/gabor/dht-go/dht"
	"gitea.izolight.xyz/gabor/dht-go/util"
	"github.com/marksamman/bencode"
	"github.com/op/go-logging"
	"net"
	"os"
	"time"
)

var log = logging.MustGetLogger("dht")
var format = logging.MustStringFormatter(
	`%{color}%{time:15:04:05.000} %{shortfunc} ▶ %{level:.4s} %{id:03x}%{color:reset} %{message}`,
)

type DHTResponse map[string]interface{}

type Node struct {
	id         string
	addr       *net.UDPAddr
	lastActive time.Time
}

type Request struct {
	infoHash []byte
	addr     *net.UDPAddr
}

var requests = make(chan Request)

func (d DHTResponse) String() string {
	nodeAddr, _ := util.ParseIP(d["ip"].(string))
	tx := d["t"].(string)
	r := d["r"].(map[string]interface{})
	id := r["id"].(string)
	n := r["nodes"].(string)
	nodes := []string{}
	for i := 0; i < len(n); {
		id := n[i : i+20]
		addr, _ := util.ParseIP(n[i+20 : i+26])
		node := fmt.Sprintf("%s %x\n", addr, id)
		i = i + 26
		nodes = append(nodes, node)
	}

	return fmt.Sprintf("Receiving: %s TX: %x ID: %x\nNodes: %s", nodeAddr, tx, id, nodes)
}

func (d DHTResponse) Nodes(routingTable chan<- Node) {
	r := d["r"].(map[string]interface{})
	n := r["nodes"].(string)
	for i := 0; i < len(n); {
		id := n[i : i+20]
		addr, err := util.ParseIP(n[i+20 : i+26])
		if err != nil {
			continue
		}
		node := Node{id: id, addr: addr, lastActive: time.Now()}
		routingTable <- node
		i = i + 26
	}
}

func initialize() (Node, []byte) {
	serverAddr, err := net.ResolveUDPAddr("udp", ":12343")
	util.CheckError(err)
	id := util.RandomString(20)
	log.Debug(fmt.Sprintf("Initialized with on %v with ID %x", serverAddr, id))

	return Node{id, serverAddr, time.Now()}, make([]byte, 65536)
}

func bootstrap(routingTable chan<- Node) {
	addr1, _ := net.ResolveUDPAddr("udp", "router.bittorrent.com:6881")
	//n2, _ := net.ResolveUDPAddr("udp", "dht.transmissionbt.com:6881")
	node1 := Node{addr: addr1, lastActive: time.Now()}
	routingTable <- node1
}

func main() {
	backend := logging.NewLogBackend(os.Stderr, "", 0)
	backendFormatter := logging.NewBackendFormatter(backend, format)
	logging.SetBackend(backendFormatter)

	self, buf := initialize()
	//secret := randomString(4)
	routingTable := make(chan Node, 2)
	go bootstrap(routingTable)

	for n := range routingTable {
		c, err := net.DialUDP("udp", self.addr, n.addr)
		util.CheckError(err)
		q := dht.FindNodes(self.id)
		log.Debug(fmt.Sprintf("Querying %v with %s Tx: %x Target: %x",
			n.addr,
			q["q"],
			q["t"],
			q["a"].(map[string]interface{})["target"]))
		c.Write(bencode.Encode(q))

		n, err := c.Read(buf)
		util.CheckError(err)
		payload := bytes.NewReader(buf[0:n])
		r, err := bencode.Decode(payload)
		util.CheckError(err)

		res := DHTResponse(r)
		log.Debug(res)
		res.Nodes(routingTable)
		c.Close()
	}
}
