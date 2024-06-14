package encrypt

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
)

func Encrypt(key, plaintext string) (string, error) {
	block, err := newCipher(key)
	if err != nil {
		return "", err
	}

	ciphertext := make([]byte, aes.BlockSize+len(plaintext))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return "", err
	}

	if len(ciphertext)%aes.BlockSize != 0 {
		fmt.Println(len(plaintext))
		return "", errors.New("ciphertext is not a multiple of the block size")
	}

	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(ciphertext[aes.BlockSize:], []byte(plaintext))

	// It's important to remember that ciphertexts must be authenticated
	// (i.e. by using crypto/hmac) as well as being encrypted in order to
	// be secure.

	return fmt.Sprintf("%x", ciphertext), nil
}

func Decrypt(key, cipherhex string) (string, error) {

	block, err := newCipher(key)
	if err != nil {
		return "", err
	}

	ciphertext, err := hex.DecodeString(cipherhex)
	if err != nil {
		return "", err
	}

	if len(ciphertext) < aes.BlockSize {
		return "", errors.New("encrypt: cipher too short")
	}
	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]

	// CBC mode always works in whole blocks.
	if len(ciphertext)%aes.BlockSize != 0 {
		return "", errors.New("ciphertext is not a multiple of the block size")
	}

	mode := cipher.NewCBCDecrypter(block, iv)

	// CryptBlocks can work in-place if the two arguments are the same.
	mode.CryptBlocks(ciphertext, ciphertext)

	return string(ciphertext), nil
}

func newCipher(key string) (cipher.Block, error) {
	hasher := md5.New()
	hasher.Write([]byte(key))
	cipherKey := hasher.Sum(nil)
	block, err := aes.NewCipher(cipherKey)
	return block, err
}
