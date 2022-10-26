package keystore

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"io/ioutil"
	"math"

	"golang.org/x/crypto/scrypt"
)

const (
	saltLength        = 32
	keyLength         = 32
	defaultWorkFactor = 20
)

func encrypt(key, data []byte) ([]byte, error) {
	key, salt, err := deriveKey(key, nil, defaultWorkFactor)
	if err != nil {
		return nil, err
	}

	blockCipher, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(blockCipher)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err = rand.Read(nonce); err != nil {
		return nil, err
	}

	ciphertext := gcm.Seal(nonce, nonce, data, nil)

	ciphertext = append([]byte{defaultWorkFactor}, append(salt, ciphertext...)...)

	return ciphertext, nil
}

func decrypt(key, data []byte) ([]byte, error) {
	workFactor, data := data[0], data[1:]
	salt, data := data[:saltLength], data[saltLength:]

	key, _, err := deriveKey(key, salt, workFactor)
	if err != nil {
		return nil, err
	}

	blockCipher, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(blockCipher)
	if err != nil {
		return nil, err
	}

	nonce, ciphertext := data[:gcm.NonceSize()], data[gcm.NonceSize():]

	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, err
	}

	return plaintext, nil
}

func deriveKey(password, salt []byte, workFactor byte) ([]byte, []byte, error) {
	if salt == nil {
		salt = make([]byte, saltLength)
		if _, err := rand.Read(salt); err != nil {
			return nil, nil, err
		}
	}

	n := int(math.Pow(2, float64(workFactor)))
	key, err := scrypt.Key(password, salt, n, 8, 1, keyLength)
	if err != nil {
		return nil, nil, err
	}

	return key, salt, nil
}

// Read reads an AES encrypted file using a scrypt password
func Read(password, filename string) ([]byte, error) {

	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	decoded, err := base64.StdEncoding.DecodeString(string(data))
	if err != nil {
		return nil, err
	}

	plaintext, err := decrypt([]byte(password), decoded)
	if err != nil {
		return nil, err
	}

	return plaintext, nil
}

// Write writes an AES encrypted file using a scrypt password
func Write(password, filename string, data []byte) error {

	ciphertext, err := encrypt([]byte(password), data)
	if err != nil {
		return err
	}

	encoded := base64.StdEncoding.EncodeToString(ciphertext)

	err = ioutil.WriteFile(filename, []byte(encoded), 0644)
	if err != nil {
		return err
	}

	return nil
}
