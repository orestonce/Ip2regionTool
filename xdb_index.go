package Ip2regionTool

import (
	"fmt"
	"strings"
)

type IndexPolicy int

const (
	VectorIndexPolicy IndexPolicy = 1
	BTreeIndexPolicy  IndexPolicy = 2
)

func IndexPolicyFromString(str string) (IndexPolicy, error) {
	switch strings.ToLower(str) {
	case "vector":
		return VectorIndexPolicy, nil
	case "btree":
		return BTreeIndexPolicy, nil
	default:
		return VectorIndexPolicy, fmt.Errorf("invalid policy '%s'", str)
	}
}
