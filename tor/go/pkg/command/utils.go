package command

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"

	"github.com/Shresth72/tor/pkg/decode"
)

func EncodeHashToString(data string) string {
	hash := sha1.New()
	hash.Write([]byte(data))
	hashSum := hash.Sum(nil)
	return hex.EncodeToString(hashSum)
}

func DecodeHashFromString(encodedStr string) ([]byte, error) {
	decoded, err := hex.DecodeString(encodedStr)
	if err != nil {
		return nil, err
	}

	return decoded, nil
}

func getMeta(arg string) (Meta, error) {
	data, err := os.ReadFile(arg)
	if err != nil {
		return Meta{}, fmt.Errorf("read file: %w", err)
	}

	d, _, err := decode.DecodeDict(string(data), 0)
	if err != nil {
		return Meta{}, fmt.Errorf("decode dict: %w", err)
	}

	return mapToMeta(d.(map[string]interface{}))
}

func mapToMeta(data map[string]interface{}) (Meta, error) {
	var meta Meta

	if announce, ok := data["announce"].(string); ok {
		meta.Announce = announce
	}

	if infoMap, ok := data["info"].(map[string]interface{}); ok {
		var info MetaInfo

		if name, ok := infoMap["name"].(string); ok {
			info.Name = name
		}

		if length, ok := infoMap["length"].(int); ok {
			info.Length = int64(length)
		} else if length, ok := infoMap["length"].(int64); ok {
			info.Length = length
		}

		if pieceLength, ok := infoMap["piece length"].(int); ok {
			info.PieceLength = int64(pieceLength)
		} else if pieceLength, ok := infoMap["piece length"].(int64); ok {
			info.PieceLength = pieceLength
		}

		if pieces, ok := infoMap["pieces"].(string); ok {
			info.Pieces = EncodeHashToString(pieces)
		}

		meta.Info = info
	} else {
		return meta, fmt.Errorf("info section not found")
	}

	return meta, nil
}

func sendTrackerRequest(uri string, req TrackerGetRequest) (*http.Response, error) {
	params := url.Values{}
	params.Add("info_hash", req.InfoHash)
	params.Add("peer_id", req.PeerId)
	params.Add("port", fmt.Sprintf("%d", req.Port))
	params.Add("uploaded", fmt.Sprintf("%d", req.Uploaded))
	params.Add("downloaded", fmt.Sprintf("%d", req.Downloaded))
	params.Add("left", fmt.Sprintf("%d", req.Left))
	params.Add("compact", fmt.Sprintf("%t", req.Compact))

	trackerUrl := fmt.Sprintf("%s?%s", uri, params.Encode())

	res, err := http.Get(trackerUrl)
	if err != nil {
		return nil, fmt.Errorf("send GET request: %w", err)
	}

	return res, nil
}

func getTrackerResponse(argument string) (TrackerResponse, error) {
	meta, err := getMeta(argument)
	if err != nil {
		return TrackerResponse{}, err
	}

	url, req, err := getPeers(meta)
	if err != nil {
		return TrackerResponse{}, err
	}

	res, err := sendTrackerRequest(url, req)
	if err != nil {
		return TrackerResponse{}, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return TrackerResponse{}, fmt.Errorf("read response body: %w", err)
	}

	d, _, err := decode.DecodeBencode(string(body), 0)
	if err != nil {
		return TrackerResponse{}, fmt.Errorf("decode bencode: %w", err)
	}

	trackerResponse := d.(map[string]interface{})
	var interval int
	if trackInterval, ok := trackerResponse["interval"].(int); ok {
		interval = trackInterval
	} else {
		return TrackerResponse{}, fmt.Errorf("expected interval in trackerResponse")
	}

	var peers string
	if trackerPeers, ok := trackerResponse["peers"].(string); ok {
		peers = trackerPeers
	} else {
		return TrackerResponse{}, fmt.Errorf("expected peers in trackerResponse")
	}

	return TrackerResponse{
		Interval: interval,
		Peers:    peers,
	}, nil
}
