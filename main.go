package main

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"errors"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
)

const (
	ShellToUse = "bash"
)

var (
	plainFile  string
	cipherFile string
	exeMode    string
)

var (
	// ErrInvalidBlockSize indicates hash blocksize <= 0.
	ErrInvalidBlockSize = errors.New("invalid blocksize")

	// ErrInvalidPKCS7Data indicates bad input to PKCS7 pad or unpad.
	ErrInvalidPKCS7Data = errors.New("invalid PKCS7 data (empty or not padded)")

	// ErrInvalidPKCS7Padding indicates PKCS7 unpad fails to bad input.
	ErrInvalidPKCS7Padding = errors.New("invalid padding on input")
)

var (
	key = []byte{
		243, 195, 232, 2, 7, 110, 92, 52, 168, 157, 188, 189, 160, 62, 202, 208, 156, 70, 197, 147, 49, 18, 108, 139, 238, 176, 113, 15, 222, 177, 124, 29,
	}
)

func executeContent(content []byte) error {
	cmd := exec.Command(ShellToUse)
	cmd.Stdin = bytes.NewReader(content)
	var buf bytes.Buffer
	cmd.Stdout = &buf

	if err := cmd.Start(); err != nil {
		return err
	}

	if err := cmd.Wait(); err != nil {
		return err
	}
	_, err := io.Copy(os.Stdout, &buf)
	return err
}

func readFile(filename string) []byte {
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}

	return content
}

func writeFile(filename string, content []byte) {
	err := ioutil.WriteFile(filename, content, 0644)
	if err != nil {
		panic(err)
	}
}

func encrypt(key []byte, plaintext []byte) []byte {
	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}
	plaintext, _ = pkcs7Pad(plaintext, block.BlockSize())
	ciphertext := make([]byte, aes.BlockSize+len(plaintext))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		panic(err)
	}
	bm := cipher.NewCBCEncrypter(block, iv)
	bm.CryptBlocks(ciphertext[aes.BlockSize:], plaintext)

	return ciphertext
}

func decrypt(key []byte, ciphertext []byte) []byte {
	// Create the AES cipher
	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}
	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]
	bm := cipher.NewCBCDecrypter(block, iv)
	bm.CryptBlocks(ciphertext, ciphertext)
	ciphertext, _ = pkcs7Unpad(ciphertext, aes.BlockSize)
	return ciphertext
}

func pkcs7Pad(b []byte, blocksize int) ([]byte, error) {
	if blocksize <= 0 {
		return nil, ErrInvalidBlockSize
	}
	if b == nil || len(b) == 0 {
		return nil, ErrInvalidPKCS7Data
	}
	n := blocksize - (len(b) % blocksize)
	pb := make([]byte, len(b)+n)
	copy(pb, b)
	copy(pb[len(b):], bytes.Repeat([]byte{byte(n)}, n))
	return pb, nil
}

func pkcs7Unpad(b []byte, blocksize int) ([]byte, error) {
	if blocksize <= 0 {
		return nil, ErrInvalidBlockSize
	}
	if b == nil || len(b) == 0 {
		return nil, ErrInvalidPKCS7Data
	}
	if len(b)%blocksize != 0 {
		return nil, ErrInvalidPKCS7Padding
	}
	c := b[len(b)-1]
	n := int(c)
	if n == 0 || n > len(b) {
		return nil, ErrInvalidPKCS7Padding
	}
	for i := 0; i < n; i++ {
		if b[len(b)-n+i] != c {
			return nil, ErrInvalidPKCS7Padding
		}
	}
	return b[:len(b)-n], nil
}

func argsParse() {
	flag.StringVar(&plainFile, "p", "plain.sh", "Plain text file (Default: plain.sh)")
	flag.StringVar(&cipherFile, "c", "data.dat", "Cipher text file (Default: data.dat)")
	flag.StringVar(&exeMode, "m", "run", "Execute mode[enc/run] {enc: 'encrypt script', run: 'run script'}")
	flag.Parse()
}

func main() {
	argsParse()

	if exeMode == "run" {
		cipherText := readFile(cipherFile)
		plainText := decrypt(key, cipherText)
		executeContent(plainText)
	} else if exeMode == "enc" {
		plainText := readFile(plainFile)
		encrypted := encrypt(key, plainText)
		writeFile(cipherFile, encrypted)
	} else {
		fmt.Printf("mode '%s' not support!!\n", exeMode)
	}

}
