package support

import (
	"crypto/sha256"

	"github.com/btcsuite/btcutil/base58"
	"golang.org/x/crypto/ripemd160"
)

func Hasher(code string) string {
	if code == "" {
		return ""
	}
	h2 := sha256.New()
	h2.Write([]byte(code))
	digest2 := h2.Sum(nil)
	h3 := ripemd160.New()
	h3.Write(digest2)
	digest3 := h3.Sum(nil)
	return base58.Encode(digest3)
}
func Hasher2(code string) string {
	if code == "" {
		return ""
	}
	h2 := sha256.New()
	h2.Write([]byte(code))
	digest2 := h2.Sum(nil)
	return base58.Encode(digest2)
}