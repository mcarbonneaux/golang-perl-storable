package storable

import (
	"fmt"
	"strings"
	"reflect"
)


type StorableError struct {
	message string
	path []reflect.Value
}

func (e *StorableError) Error() string {
	var res []string
	for _, value := range e.path {
		res = append(res, value.String())
	}

	return fmt.Sprintf("%s // %s", e.message, strings.Join(res, " ⟶ "))
}
