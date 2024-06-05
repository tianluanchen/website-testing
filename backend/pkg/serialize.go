package pkg

import (
	"bytes"
	"encoding/json"
)

func Serialize(data any, indent ...string) ([]byte, error) {
	buffer := &bytes.Buffer{}
	encoder := json.NewEncoder(buffer)
	encoder.SetEscapeHTML(false)
	if len(indent) > 0 {
		encoder.SetIndent("", indent[0])
	} else {
		encoder.SetIndent("", "")
	}
	if err := encoder.Encode(data); err != nil {
		return nil, err
	}
	bs := buffer.Bytes()
	// remove last character: \n
	return bs[:len(bs)-1], nil
}

func Deserialize(data []byte, v any) error {
	return json.Unmarshal(data, v)
}
