package bencode

import (
	"encoding/hex"
	"reflect"
	"testing"
)

func TestDecodeDictionary(t *testing.T) {
	dict := make(map[string]interface{})
	dict["str"] = "string"
	dict["int"] = int64(1234)
	dict["list"] = []interface{}{
		"string",
		int64(1234),
		[]interface{}{"string", int64(1234)},
	}
	dict["dict"] = map[string]interface{}{
		"str": "string",
		"int": int64(1234),
		"list": []interface{}{
			"string",
			int64(1234),
			[]interface{}{"string", int64(1234)},
		},
	}

	res, err := Decode([]byte("d4:dictd3:inti1234e4:listl6:stringi1234el6:stringi1234eee3:str6:stringe3:inti1234e4:listl6:stringi1234el6:stringi1234eee3:str6:stringe"))
	switch {
	case err != nil:
		t.Errorf("while decoding dictionary: %v", err)
	case !reflect.DeepEqual(res, dict):
		t.Errorf("unexpected dictionary decoding: %v", res)
	}
}

func TestDecodeExampleTorrent(t *testing.T) {
	hexTorrent := "64373a636f6d6d656e7431353a6578616d706c6520636f6d6d656e743130" +
		"3a6372656174656420627931383a71426974746f7272656e742076342e30" +
		"2e3331333a6372656174696f6e2064617465693135383638303432383365" +
		"343a696e666f64363a6c656e6774686937383565343a6e616d65363a7465" +
		"73742e7631323a7069656365206c656e67746869313633383465363a7069" +
		"6563657332303a5fdf83383fbb7aad048ae7ae23c71f9c15f7802165383a" +
		"75726c2d6c697374303a65"

	data, err := hex.DecodeString(hexTorrent)
	if err != nil {
		t.Errorf("unable to hex decode torrent file: %v", err)
	}

	dict := make(map[string]interface{})
	dict["comment"] = "example comment"
	dict["created by"] = "qBittorrent v4.0.3"
	dict["creation date"] = int64(1586804283)
	dict["info"] = map[string]interface{}{
		"length":       int64(785),
		"name":         "test.v",
		"piece length": int64(16384),
		"pieces":       string([]byte{0x5f, 0xdf, 0x83, 0x38, 0x3f, 0xbb, 0x7a, 0xad, 0x04, 0x8a, 0xe7, 0xae, 0x23, 0xc7, 0x1f, 0x9c, 0x15, 0xf7, 0x80, 0x21}),
	}
	dict["url-list"] = ""

	res, err := Decode(data)
	switch {
	case err != nil:
		t.Errorf("while decoding torrent: %v", err)
	case !reflect.DeepEqual(res, dict):
		t.Errorf("unexpected torrent decoding: %v\n%v", res, dict)
	}
}
