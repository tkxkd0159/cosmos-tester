package tm

import "net/http"

type Querier struct {
	*http.Client
	URL string
}

func NewQuerier(host string) *Querier {
	return &Querier{&http.Client{}, host}
}
