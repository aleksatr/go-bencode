package bencode

import (
	"bytes"
	"fmt"
	"testing"
)

func TestEncodeNil(t *testing.T) {
	_, err := Encode(nil)
	if err != ErrUnsupportedType {
		t.Errorf("unexpected err type: %v", err)
	}
}

func TestEncodeEmptyString(t *testing.T) {
	check := func(res []byte, err error) {
		switch {
		case err != nil:
			t.Errorf("while encoding empty string: %v", err)
		case len(res) != 2:
			t.Errorf("unexpected result length")
		case res[0] != 48 || res[1] != 58:
			t.Errorf("empty string should be encoded to 0:")
		}
	}

	res, err := Encode("")
	check(res, err)

	var b []byte
	res, err = Encode(b)
	check(res, err)
}

func TestEncodeString(t *testing.T) {
	check := func(res []byte, err error) {
		switch {
		case err != nil:
			t.Errorf("while encoding string: %v", err)
		case string(res) != "25:this is an example string":
			t.Errorf("unexpected string encoding: %s", string(res))
		}
	}

	testStr := "this is an example string"
	res, err := Encode(testStr)
	check(res, err)

	// check the byte slice variant
	res, err = Encode([]byte(testStr))
	check(res, err)
}

func TestEncodeInt(t *testing.T) {
	res, err := Encode(1234)
	switch {
	case err != nil:
		t.Errorf("while encoding int: %v", err)
	case string(res) != "i1234e":
		t.Errorf("unexpected int encoding: %s", string(res))
	}
}

func TestEncodeList(t *testing.T) {
	list := []interface{}{
		"string",
		1234,
		[]interface{}{"string", 1234},
	}
	res, err := Encode(list)
	switch {
	case err != nil:
		t.Errorf("while encoding list: %v", err)
	case string(res) != "l6:stringi1234el6:stringi1234eee":
		t.Errorf("unexpected list encoding: %s", string(res))
	}
}

func TestEncodeDictionary(t *testing.T) {
	subDict := make(map[string]interface{})
	subDict["str"] = "string"
	subDict["int"] = 1234
	subDict["list"] = []interface{}{
		"string",
		1234,
		[]interface{}{"string", 1234},
	}

	dict := make(map[string]interface{})
	dict["str"] = "string"
	dict["int"] = 1234
	dict["list"] = []interface{}{
		"string",
		1234,
		[]interface{}{"string", 1234},
	}
	dict["dict"] = subDict

	res, err := Encode(dict)
	switch {
	case err != nil:
		t.Errorf("while encoding dictionary: %v", err)
	case string(res) != "d4:dictd3:inti1234e4:listl6:stringi1234el6:stringi1234eee3:str6:stringe3:inti1234e4:listl6:stringi1234el6:stringi1234eee3:str6:stringe":
		t.Errorf("unexpected dictionary encoding: %s", string(res))
	}

}

func BenchmarkEncodeString(b *testing.B) {
	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		err := encodeString(&bytes.Buffer{}, "this is an example string")
		if err != nil {
			b.Error(err)
			return
		}
	}
}

func BenchmarkAlternativeEncodeString(b *testing.B) {
	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		err := alternativeEncodeString(&bytes.Buffer{}, "this is an example string")
		if err != nil {
			b.Error(err)
			return
		}
	}
}

func alternativeEncodeString(buf *bytes.Buffer, s string) error {
	_, err := fmt.Fprintf(buf, "%d:%s", len(s), s)
	return err
}
