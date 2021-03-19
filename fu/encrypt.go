package fu

import (
	"bytes"
	"compress/gzip"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/binary"
	"io"
)

func Encrypt(password string, data []byte) ([]byte, error) {
	key := sha512.Sum512([]byte(password))
	ivx := make([]byte, aes.BlockSize*2)
	if _, err := rand.Read(ivx[:aes.BlockSize]); err != nil {
		return nil, err
	}
	copy(ivx[aes.BlockSize:], key[sha512.Size-aes.BlockSize:])
	iv := sha256.Sum256(ivx)
	block, err := aes.NewCipher(key[:aes.BlockSize])
	if err != nil {
		return nil, err
	}
	bf := bytes.Buffer{}
	bf.Write(ivx[:aes.BlockSize])
	sha := sha1.Sum(data)
	bf.Write(sha[:])
	wr := &cipher.StreamWriter{S: cipher.NewOFB(block, iv[:aes.BlockSize]), W: &bf}
	if err := binary.Write(wr, binary.LittleEndian, uint32(len(data))); err != nil {
		return nil, err
	}
	gw := gzip.NewWriter(wr)
	if _, err := io.Copy(gw, bytes.NewReader(data)); err != nil {
		return nil, err
	}
	_ = gw.Close()
	_ = wr.Close()
	return bf.Bytes(), nil
}
