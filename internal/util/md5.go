package util

import (
	"crypto/md5"
	"encoding/hex"
	"io"
	"log"
	"os"
)

// MD5CheckSum get md5 checksum given file path
func MD5CheckSum(path string) (string, error) {
	h := md5.New()
	f, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	if _, err := io.Copy(h, f); err != nil {
		return "", nil
	}
	return hex.EncodeToString(h.Sum(nil)), nil
}
