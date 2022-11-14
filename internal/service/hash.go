package service

import (
	"crypto/sha1"
	"crypto/sha512"
	"fmt"
	"strings"
)

func Hash(s string) string {
	salt := sha1.New()
	h := sha512.New384()
	p := fmt.Sprintf("%xо%sо%x", h.Sum([]byte(s)), s, salt.Sum([]byte(s)))
	return p
}

func DeHash(s string) string {
	p := strings.Split(s, "о")
	return p[1]
}
