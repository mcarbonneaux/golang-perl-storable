package storable

import (
	"bytes"
	"fmt"
	"os/exec"
	"testing"
)


// var RECURSION_ARRAY = [123, -1.23, nil, '123', '1' * 1000, [5], {'x': 6}, 'Привет!']
// RECURSION_ARRAY = append(RECURSION_ARRAY, RECURSION_ARRAY)


func freeze_perl(code string) []byte {
	code = fmt.Sprintf("$x = %s; print Storable::freeze(ref $x? $x: \\$x)", code)
	cmd := exec.Command("perl", "-MStorable", "-e", code)
    stdout, err := cmd.Output()

    if err != nil {
        panic(err.Error())
    }
	
	return stdout
}

func TestByte(t *testing.T) {
	var data uint8
	storable := freeze_perl("123")
	
	err := Unmarshal(storable, &data)
	
	if err != nil {
		t.Errorf("Unmarshal error: %s", err.Error())
	} else if data != 123 {
		t.Errorf("Expected: 123, got: %v", data)
	}
}

func TestByteInterface(t *testing.T) {
	var data any
	storable := freeze_perl("123")
	
	err := Unmarshal(storable, &data)
	if err != nil {
		t.Errorf("Unmarshal error: %s", err.Error())
	} else if data.(int8) != 123 {
		t.Errorf("Expected: 123, got: %v", data)
	}
}

func TestInteger(t *testing.T) {
	var data int64
	storable := freeze_perl("1000")
	
	err := Unmarshal(storable, &data)
	if err != nil {
		t.Errorf("Unmarshal error: %s", err.Error())
	} else if data != 1000 {
		t.Errorf("Expected: 1000, got: %d", data)
	}
}

func TestDouble(t *testing.T) {
	var data float32
	storable := freeze_perl("-2.5")
	
	err := Unmarshal(storable, &data)
	if err != nil {
		t.Errorf("Unmarshal error: %s", err.Error())
	} else if data != -2.5 {
		t.Errorf("Expected: -2.5, got: %g", data)
	}
}

func TestScalarAsBytes(t *testing.T) {
	var data []byte
	storable := freeze_perl(`"abc"`)
	
	err := Unmarshal(storable, &data)
	if err != nil {
		t.Errorf("Unmarshal error: %s", err.Error())
	} else if !bytes.Equal(data, []byte("abc")) {
		t.Errorf("Expected: `abc`, got: %q", data)
	}
}

func TestScalarAsString(t *testing.T) {
	var data string
	storable := freeze_perl(`"abc"`)
	
	err := Unmarshal(storable, &data)
	if err != nil {
		t.Errorf("Unmarshal error: %s", err.Error())
	} else if data != "abc" {
		t.Errorf("Expected: `abc`, got: %q", data)
	}
}

func TestLScalar(t *testing.T) {
	var data []byte
	storable := freeze_perl(`"a" x 1000`)
	
	err := Unmarshal(storable, &data)
	if err != nil {
		t.Errorf("Unmarshal error: %s", err.Error())
	} else if !bytes.Equal(data, bytes.Repeat([]byte{'a'}, 1000)) {
		t.Errorf("Expected: `a` repeated 1000 times, got: %q", data)
	}
}
