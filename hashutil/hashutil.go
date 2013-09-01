package hashutil

import (
	"hash"
	"crypto/sha256"
	"fmt"
)

func CalcHexHash(h hash.Hash, data []byte) string {
	h.Write(data)
	return fmt.Sprintf("%x", h.Sum(nil))
}

func CalcSha256Sum(data []byte) string {
	return CalcHexHash(sha256.New(), data)
}
