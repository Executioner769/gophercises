package secret

import (
	"encoding/json"
	"errors"
	"io"
	"os"
	"strings"
	"sync"

	"gopher.com/secret/encrypt"
)

type Vault struct {
	encodingKey string
	fileMutex   sync.Mutex
	filepath    string
	keyValues   map[string]string
}

func NewVaultFile(encodingKey, file string) *Vault {
	return &Vault{
		encodingKey: encodingKey,
		filepath:    file,
		keyValues:   make(map[string]string),
	}
}

func (v *Vault) loadKeyValues() error {
	f, err := os.Open(v.filepath)
	if err != nil {
		v.keyValues = make(map[string]string)
		return nil
	}
	defer f.Close()

	var sb strings.Builder
	_, err = io.Copy(&sb, f)
	if err != nil {
		return err
	}

	decryptedJson, err := encrypt.Decrypt(v.encodingKey, sb.String())
	if err != nil {
		return err
	}

	r := strings.NewReader(decryptedJson)
	dec := json.NewDecoder(r)
	err = dec.Decode(&v.keyValues)
	if err != nil {
		return err
	}
	return nil
}

func (v *Vault) saveKeyValues() error {
	var sb strings.Builder
	enc := json.NewEncoder(&sb)
	err := enc.Encode(v.keyValues)
	if err != nil {
		return err
	}
	encryptedJson, err := encrypt.Encrypt(v.encodingKey, sb.String())
	if err != nil {
		return err
	}
	f, err := os.OpenFile(v.filepath, os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = f.Write([]byte(encryptedJson))
	if err != nil {
		return err
	}
	return nil

}

func (v *Vault) GetForFile(key string) (string, error) {
	v.fileMutex.Lock()
	defer v.fileMutex.Unlock()
	err := v.loadKeyValues()
	if err != nil {
		return "", err
	}
	value, ok := v.keyValues[key]
	if !ok {
		return "", errors.New("secret: no value for the key")
	}
	return value, nil
}

func (v *Vault) SetFromFile(key, value string) error {
	v.fileMutex.Lock()
	defer v.fileMutex.Unlock()

	encryptedValue, err := encrypt.Encrypt(v.encodingKey, value)
	if err != nil {
		return err
	}

	err = v.loadKeyValues()
	if err != nil {
		return err
	}
	v.keyValues[key] = encryptedValue
	err = v.saveKeyValues()
	return err

}

func NewVault(encodingKey string) *Vault {
	return &Vault{
		encodingKey: encodingKey,
		filepath:    "",
		keyValues:   make(map[string]string),
	}
}

func (v *Vault) Get(key string) (string, error) {
	hex, ok := v.keyValues[key]
	if !ok {
		return "", errors.New("secret: no value for the key")
	}
	ret, err := encrypt.Decrypt(v.encodingKey, hex)
	if err != nil {
		return "", err
	}
	return ret, nil
}

func (v *Vault) Set(key, value string) error {

	encryptedValue, err := encrypt.Encrypt(v.encodingKey, value)
	if err != nil {
		return err
	}
	v.keyValues[key] = encryptedValue
	return nil
}
