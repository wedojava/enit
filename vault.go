package enit

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
	"sync"

	"github.com/wedojava/enit/encrypt"
)

func File(encodingKey, filepath string) *Vault {
	return &Vault{
		encodingKey: encodingKey,
		filepath:    filepath,
		keyValues:   make(map[string]string),
	}
}

type Vault struct {
	encodingKey string
	filepath    string
	mutex       sync.Mutex
	keyValues   map[string]string
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
	decryptedJSON, err := encrypt.Decrypt(v.encodingKey, sb.String())
	if err != nil {
		return err
	}
	// fmt.Println(decryptedJSON) // TODO: why the print is wierd
	r := strings.NewReader(decryptedJSON)
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
	encryptJSON, err := encrypt.Encrypt(v.encodingKey, sb.String())
	if err != nil {
		return err
	}
	f, err := os.OpenFile(v.filepath, os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		v.keyValues = make(map[string]string)
		return nil
	}
	defer f.Close()
	_, err = fmt.Fprint(f, encryptJSON)
	if err != nil {
		return err
	}
	return nil
}

func (v *Vault) Get(key string) (string, error) {
	err := v.loadKeyValues()
	if err != nil {
		return "", err
	}
	value, ok := v.keyValues[key]
	if !ok {
		return "", errors.New("secret: no value for that key")
	}
	return value, nil
}

func (v *Vault) Set(key, value string) error {
	v.mutex.Lock()
	defer v.mutex.Unlock()
	err := v.loadKeyValues()
	if err != nil {
		return err
	}
	v.keyValues[key] = value
	err = v.saveKeyValues()
	return err
}
