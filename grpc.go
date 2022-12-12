package cosmo

import (
	"google.golang.org/grpc"
)

type Querier struct {
	*grpc.ClientConn
}

func NewQuerier(addr string) *Querier {
	conn, err := OpenConn(addr)
	if err != nil {
		panic(err)
	}
	return &Querier{conn}
}

func OpenConn(addr string) (*grpc.ClientConn, error) {
	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}
	return conn, nil
}
