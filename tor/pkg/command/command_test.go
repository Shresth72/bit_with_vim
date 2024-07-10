package command

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExecuteCommandInfo(t *testing.T) {
	argument := "../../test.torrent"
	expectedLength := int64(92063)
	expectedTrackerURL := "url"
	expectedHash := "ec1fb802d679814f3a84c0215c26e7e17449024a"

	jsonOutput, err := ExecuteCommand("info", argument)
	if err != nil {
		t.Fatal(err)
	}

	var meta Meta
	err = json.Unmarshal(jsonOutput, &meta)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, expectedLength, meta.Info.Length)
	assert.Equal(t, expectedTrackerURL, meta.Announce)
	assert.Equal(t, expectedHash, meta.Info.Pieces)
}
