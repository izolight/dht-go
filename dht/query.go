package dht

import (
	"fmt"
	"gitea.izolight.xyz/gabor/dht-go/util"
	"github.com/marksamman/bencode"
)

func FindNodes(id string) []byte {
	q := make(map[string]interface{})
	q["y"] = "q"
	q["q"] = "find_node"
	q["t"] = util.RandomString(2)
	a := make(map[string]interface{})
	a["id"] = id
	a["target"] = util.RandomString(20)
	q["a"] = a

	fmt.Printf("Sending: TX: %x\t Target: %x\n", q["t"], a["target"])

	return bencode.Encode(q)
}

func ping(id string) []byte {
	q := make(map[string]interface{})
	q["y"] = "q"
	q["q"] = "ping"
	q["t"] = util.RandomString(2)
	a := make(map[string]interface{})
	a["id"] = id
	q["a"] = a

	return bencode.Encode(q)
}

func getPeers(id string, infoHash string) []byte {
	q := make(map[string]interface{})
	q["y"] = "q"
	q["q"] = "get_peers"
	q["t"] = util.RandomString(2)
	a := make(map[string]interface{})
	a["id"] = id
	a["info_hash"] = infoHash
	q["a"] = a

	return bencode.Encode(q)
}

func announcePeer(id string, infoHash string, port uint16, token string) []byte {
	q := make(map[string]interface{})
	q["y"] = "q"
	q["q"] = "announce_peer"
	q["t"] = util.RandomString(2)
	a := make(map[string]interface{})
	a["id"] = id
	a["info_hash"] = infoHash
	a["port"] = port
	a["token"] = token
	q["a"] = a

	return bencode.Encode(q)
}
