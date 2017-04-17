package dht

func makeQuery(t, q string, a map[string]interface{}) map[string]interface{} {
	return map[string]interface{}{
		"y": "q",
		"t": t,
		"q": q,
		"a": a,
	}
}

func makeResponse(t string, r map[string]interface{}) map[string]interface{} {
	return map[string]interface{}{
		"y": "r",
		"t": t,
		"r": r,
	}
}

func makeError(t string, ec uint8, em string) map[string]interface{} {
	return map[string]interface{}{
		"y": "e",
		"t": t,
		"e": []interface{}{
			ec,
			em,
		},
	}
}

func makePingQuery(t, id string) map[string]interface{} {
	return makeQuery(t, "ping", map[string]interface{}{
		"id": id,
	})
}

func makePingResponse(t, id string) map[string]interface{} {
	return makeResponse(t, map[string]interface{}{
		"id": id,
	})
}

func makeFindNodeQuery(t, id, target string) map[string]interface{} {
	return makeQuery(t, "find_node", map[string]interface{}{
		"id":     id,
		"target": target,
	})
}

func makeFindNodeResponse(t, id, nodes string) map[string]interface{} {
	return makeResponse(t, map[string]interface{}{
		"id":    id,
		"nodes": nodes,
	})
}

func makeGetPeersQuery(t, id, infoHash string) map[string]interface{} {
	return makeQuery(t, "get_peers", map[string]interface{}{
		"id":        id,
		"info_hash": infoHash,
	})
}

func makeGetPeersResponsePeers(t, id, token, peers string) map[string]interface{} {
	return makeResponse(t, map[string]interface{}{
		"id":     id,
		"token":  token,
		"values": []string{peers},
	})
}

func makeGetPeersResponseNodes(t, id, token, nodes string) map[string]interface{} {
	return makeResponse(t, map[string]interface{}{
		"id":    id,
		"token": token,
		"nodes": nodes,
	})
}

func makeAnnouncePeerQuery(t, id, infoHash, token string, port int) map[string]interface{} {
	return makeQuery(t, "announce_peers", map[string]interface{}{
		"id":        id,
		"token":     token,
		"info_hash": infoHash,
		"port":      port,
	})
}

func makeAnnouncePeerResponse(t, id string) map[string]interface{} {
	return makePingResponse(t, id)
}
