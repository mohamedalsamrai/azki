package azki

import (
	"net/http"
)
type Param struct {
	Key   string
	Value string
}
type Params []Param
type Handle func(http.ResponseWriter, *http.Request, Params)