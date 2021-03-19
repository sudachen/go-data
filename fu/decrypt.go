package fu

import (
	"bytes"
	"compress/gzip"
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/binary"
	"io"
	"sudachen.xyz/pkg/go-forge/errors"
)

func Decrypt(password string, data []byte) ([]byte, error) {
	key := sha512.Sum512([]byte(password))
	ds := bytes.NewReader(data)
	ivx := make([]byte, aes.BlockSize*2)
	if _, err := io.ReadFull(ds, ivx[:aes.BlockSize]); err != nil {
		return nil, err
	}
	sha := [20]byte{}
	if _, err := io.ReadFull(ds, sha[:]); err != nil {
		return nil, err
	}
	copy(ivx[aes.BlockSize:], key[sha512.Size-aes.BlockSize:])
	iv := sha256.Sum256(ivx)
	block, err := aes.NewCipher(key[:aes.BlockSize])
	if err != nil {
		return nil, err
	}
	bf := bytes.Buffer{}
	rd := &cipher.StreamReader{S: cipher.NewOFB(block, iv[:aes.BlockSize]), R: ds}
	var ln uint32
	if err := binary.Read(rd, binary.LittleEndian, &ln); err != nil {
		return nil, err
	}
	gr, err := gzip.NewReader(rd)
	if err != nil {
		return nil, err
	}
	if _, err := io.Copy(&bf, gr); err != nil {
		return nil, err
	}
	data = bf.Bytes()
	if len(data) != int(ln) || sha1.Sum(data) != sha {
		return nil, errors.New("encrypted data corrupted")
	}
	return data, nil
}
