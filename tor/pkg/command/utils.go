package command

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
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

func getPeers(meta Meta) (string, TrackerGetRequest, error) {
  if meta.Announce == "" || meta.Info.Pieces == "" || meta.Info.Length == 0 {
    return "", TrackerGetRequest{}, fmt.Errorf("missing required fields")
  }

  req := TrackerGetRequest{
    InfoHash: meta.Info.Pieces,
    PeerId: "id420",
    Port: 6969,
    Uploaded: 0,
    Downloaded: 0,
    Left: meta.Info.Length,
    Compact: true,
  }

  return meta.Announce, req, nil
}

