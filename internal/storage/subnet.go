package storage

import (
	"errors"
)

var (
	ErrSubnetAlreadyExists = errors.New("subnet already exists")
	ErrSubnetNotExists     = errors.New("subnet not exists")
)

type Subnet struct {
	Subnet string `json:"subnet"`
}

func NewSubnet(subnet string) Subnet {
	return Subnet{
		Subnet: subnet,
	}
}
