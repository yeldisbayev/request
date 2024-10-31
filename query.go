package request

import (
	"fmt"
	"math/big"
)

type Query interface {
	bool | string |
		int8 | uint8 | int16 | uint16 | int32 | uint32 |
		int | uint | int64 | uint64 | float32 | float64 | big.Int | big.Float
}

func Queries[T Query](values ...T) []string {
	result := make([]string, len(values))

	for i, v := range values {
		result[i] = fmt.Sprintf("%v", v)
	}

	return result

}
