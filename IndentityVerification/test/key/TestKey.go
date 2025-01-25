package main

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
)

func main() {
	// 生成 32 字节的随机密钥（256 位）
	key := make([]byte, 32)
	_, err := rand.Read(key)
	if err != nil {
		panic("Failed to generate random key")
	}

	// 将密钥编码为 Base64 字符串（方便存储和传输）
	keyBase64 := base64.StdEncoding.EncodeToString(key)
	fmt.Println("Generated Key:", keyBase64)
}
