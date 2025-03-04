package inventory

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
)
func Checksum(text string) string {
	hash := sha256.New()
	hash.Write([]byte(text))
	return hex.EncodeToString(hash.Sum(nil))
}

func JsonStringfy(v interface{}, indent bool) string {
	if indent {
		raw, _ := json.MarshalIndent(v, "", "  ")
		return string(raw)
	}
	raw, _ := json.Marshal(v)
	return string(raw)
}
func JsonParse(raw string, v interface{}) error {
	return json.Unmarshal([]byte(raw), v)
}