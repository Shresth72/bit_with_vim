package command

import (
	"encoding/json"
	"fmt"
	"net"
	"time"

	"github.com/Shresth72/tor/pkg/decode"
)

type MetaInfo struct {
	Name        string `bencode:"name"`
	Pieces      string `bencode:"pieces"`
	Length      int64  `bencode:"length"`
	PieceLength int64  `bencode:"piece length"`
}

type Meta struct {
	Announce string   `bencode:"announce"`
	Info     MetaInfo `bencode:"info"`
}

type TrackerGetRequest struct {
	InfoHash   string `json:"info_hash"` // URL encoded
	PeerId     string `json:"peer_id"`
	Port       int16  `json:"port"`
	Uploaded   uint   `json:"uploaded"`
	Downloaded uint   `json:"downloaded"`
	Left       int64  `json:"left"`
	Compact    bool   `json:"compact"`
}

type TrackerResponse struct {
	Interval int `json:"interval"`
	// peer(6 bytes) - ipAddr(4bytes) , port(2bytes)
	// Is string, but the contents are binary data
	Peers string `json:"peers"`
}

func ExecuteCommand(command, argument string) ([]byte, error) {
	switch command {
	case "decode":
		decoded, _, err := decode.DecodeBencode(argument, 0)
		if err != nil {
			return nil, fmt.Errorf("decode bencode: %w", err)
		}

		jsonOutput, err := json.Marshal(decoded)
		if err != nil {
			return nil, fmt.Errorf("encode to json: %w", err)
		}
		return jsonOutput, nil

	case "info":
		meta, err := getMeta(argument)
		if err != nil {
			return nil, err
		}

		fmt.Println("Tracker URL:", meta.Announce)
		fmt.Println("Length:", meta.Info.Length)
		fmt.Printf("Info Hash: %x\n", meta.Info.Pieces)
		fmt.Println("Piece Length:", meta.Info.PieceLength)

		return json.Marshal(meta)

	case "peers":
		trackerRes, err := getTrackerResponse(argument)
		if err != nil {
			return nil, err
		}
		return json.Marshal(trackerRes)

	case "handshake":
		trackerRes, err := getTrackerResponse(argument)
		if err != nil {
			return nil, err
		}

		peerList := parsePeers(trackerRes.Peers)
		if len(peerList) == 0 {
			return nil, fmt.Errorf("no peers found")
		}

		conn, err := net.DialTimeout("tcp", peerList[0].Address, 10*time.Second)
		if err != nil {
			return nil, fmt.Errorf("dial: %w", err)
		}
		defer conn.Close()

		err = sendHandshake(conn, trackerRes.Peers)
		if err != nil {
			return nil, fmt.Errorf("send handshake: %w", err)
		}

		response, err := readHandshake(conn)
		if err != nil {
			return nil, fmt.Errorf("read handshake: %w", err)
		}

		return json.Marshal(response)

	case "download":
		trackerRes, err := getTrackerResponse(argument)
		if err != nil {
			return nil, err
		}

		peerList := parsePeers(trackerRes.Peers)
		if len(peerList) == 0 {
			return nil, fmt.Errorf("no peers found")
		}

		for _, peer := range peerList {
			conn, err := net.DialTimeout("tcp", peer.Address, 10*time.Second)
			if err != nil {
				continue
			}
			defer conn.Close()

			err = sendHandshake(conn, trackerRes.Peers)
			if err != nil {
				continue
			}

			_, err = readHandshake(conn)
			if err != nil {
				return nil, err
			}

		}
		return nil, nil

	default:
		return nil, fmt.Errorf("unknown command: %s", command)
	}
}
