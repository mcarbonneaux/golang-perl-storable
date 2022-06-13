package storable

import (
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
	var data int8
	bytes := freeze_perl("123")
	
	t.Logf("%q", bytes)
	
	err := Unmarshal(bytes, &data)
	if err != nil {
		t.Errorf("Unmarshal error: %s", err.Error())
	} else if data != 123 {
		t.Errorf("expected: 123, got: %d", data)
	}
}




