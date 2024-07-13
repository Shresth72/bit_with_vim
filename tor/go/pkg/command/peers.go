package command

import "fmt"

type Peer struct {
	Address string
}

func parsePeers(peers string) []Peer {
	var peerList []Peer
	for i := 0; i < len(peers); i += 6 {
		ip := fmt.Sprintf("%d.%d.%d.%d", peers[i], peers[i+1], peers[i+2], peers[i+3])
		port := int(peers[i+4])<<8 + int(peers[i+5])
		peerList = append(peerList, Peer{
			Address: fmt.Sprintf("%s:%d", ip, port),
		})
	}

	return peerList
}

func getPeers(meta Meta) (string, TrackerGetRequest, error) {
	if meta.Announce == "" || meta.Info.Pieces == "" || meta.Info.Length == 0 {
		return "", TrackerGetRequest{}, fmt.Errorf("missing required fields")
	}

	req := TrackerGetRequest{
		InfoHash:   meta.Info.Pieces,
		PeerId:     "id420",
		Port:       6969,
		Uploaded:   0,
		Downloaded: 0,
		Left:       meta.Info.Length,
		Compact:    true,
	}

	return meta.Announce, req, nil
}
