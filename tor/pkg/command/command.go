package command

import (
	"encoding/json"
	"fmt"
	"io"

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
		meta, err := getMeta(argument)
		if err != nil {
			return nil, err
		}

		url, req, err := getPeers(meta)
		if err != nil {
			return nil, err
		}

		res, err := sendTrackerRequest(url, req)
		if err != nil {
			return nil, err
		}
		defer res.Body.Close()

		body, err := io.ReadAll(res.Body)
		if err != nil {
			return nil, fmt.Errorf("read response body: %w", err)
		}

		d, _, err := decode.DecodeBencode(string(body), 0)
		if err != nil {
			return nil, fmt.Errorf("decode bencode: %w", err)
		}

		trackerResponse := d.(map[string]interface{})
    var interval int
    if trackInterval, ok := trackerResponse["interval"].(int); ok {
      interval = trackInterval
    } else {
      return nil, fmt.Errorf("expected interval in trackerResponse")
    }

    var peers string
    if trackerPeers, ok := trackerResponse["peers"].(string); ok {
      peers = trackerPeers
    } else {
      return nil, fmt.Errorf("expected peers in trackerResponse")
    }

    tres := TrackerResponse{
      Interval: interval,
      Peers: peers,
    }
		return json.Marshal(tres)

	default:
		return nil, fmt.Errorf("unknown command: %s", command)
	}
}
