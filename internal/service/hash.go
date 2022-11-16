package service

import (
	"crypto/aes"
	"encoding/hex"
	"fmt"
	"os"
)

func Hash(s string) string {
	key := os.Getenv("KEY")
	c, err := aes.NewCipher([]byte(key))
	if err != nil {
		fmt.Println(err)
	}
	msgByte := make([]byte, len(s))
	c.Encrypt(msgByte, []byte(s))
	return hex.EncodeToString(msgByte)
}

func DeHash(s string) string {
	txt, _ := hex.DecodeString(s)
	key := os.Getenv("KEY")
	c, err := aes.NewCipher([]byte(key))
	if err != nil {
		fmt.Println(err)
	}
	msgByte := make([]byte, len(txt))
	c.Decrypt(msgByte, []byte(txt))
	msg := string(msgByte[:])
	return msg
}
