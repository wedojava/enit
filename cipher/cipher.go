package cipher

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/rand"
	"errors"
	"fmt"
	"io"
)

func encryptStream(key string, iv []byte) (cipher.Stream, error) {
	block, err := newCipherBlock(key)
	if err != nil {
		return nil, err
	}

	return cipher.NewCFBEncrypter(block, iv), nil
}

// Encrypt will take in a key and plaintext and return a hex
// representation of the encrypted value.
// This code is based on the standard library examples at:
//   - https://golang.org/pkg/crypto/cipher/#example_NewCFBEncrypter
// func Encrypt(key, plaintext string) (string, error) {
//         // The IV needs to be unique, but not secure. Therefore it's common to
//         // include it at the beginning of the ciphertext.
//         ciphertext := make([]byte, aes.BlockSize+len(plaintext))
//         iv := ciphertext[:aes.BlockSize]
//         if _, err := io.ReadFull(rand.Reader, iv); err != nil {
//                 return "", err
//         }
//
//         stream, err := encryptStream(key, iv)
//         if err != nil {
//                 return "", err
//         }
//
//         stream.XORKeyStream(ciphertext[aes.BlockSize:], []byte(plaintext))
//
//         // It's important to remember that ciphertexts must be authenticated
//         // (i.e. by using crypto/hmac) as well as being encrypted in order to
//         // be secure.
//         return fmt.Sprintf("%x", ciphertext), nil
// }

// EncryptWriter will return a writer that will write encrypted data to
// the original writer.
func EncryptWriter(key string, w io.Writer) (*cipher.StreamWriter, error) {
	iv := make([]byte, aes.BlockSize)
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return nil, err
	}
	stream, err := encryptStream(key, iv)
	if err != nil {
		return nil, err
	}
	n, err := w.Write(iv)
	if n != len(iv) || err != nil {
		return nil, errors.New("encrypt: unable to write full iv to writer")
	}
	return &cipher.StreamWriter{S: stream, W: w}, nil
}

func decryptStream(key string, iv []byte) (cipher.Stream, error) {
	block, err := newCipherBlock(key)
	if err != nil {
		return nil, err
	}
	return cipher.NewCFBDecrypter(block, iv), nil
}

// Decrypt will take in a key and a cipherHex
// (hex representation of the ciphertext) and decrypt it.
// This code is based on the standard library examples at:
//   - https://golang.org/pkg/crypto/cipher/#example_NewCFBDecrypter
// func Decrypt(key, cipherHex string) (string, error) {
//         ciphertext, err := hex.DecodeString(cipherHex)
//         if err != nil {
//                 return "", err
//         }
//
//         // The IV needs to be unique, but not secure. Therefore it's common to
//         // include it at the beginning of the ciphertext.
//         if len(ciphertext) < aes.BlockSize {
//                 return "", errors.New("encrypt: ciphertext too short")
//         }
//         iv := ciphertext[:aes.BlockSize]
//         ciphertext = ciphertext[aes.BlockSize:]
//
//         stream, err := decryptStream(key, iv)
//         if err != nil {
//                 return "", err
//         }
//
//         // XORKeyStream can work in-place if the two arguments are the same.
//         stream.XORKeyStream(ciphertext, ciphertext)
//         return string(ciphertext), nil
// }

// DecryptReader will return a reader that will descrypt data from the
// provided reader and give the user a way to read that data as if it was
// not encrypted.
func DecryptReader(key string, r io.Reader) (*cipher.StreamReader, error) {
	iv := make([]byte, aes.BlockSize)
	n, err := r.Read(iv)
	if n < len(iv) || err != nil {
		return nil, errors.New("decrypt: unable to read the full iv")
	}
	stream, err := decryptStream(key, iv)
	if err != nil {
		return nil, err
	}
	return &cipher.StreamReader{S: stream, R: r}, nil
}

func newCipherBlock(key string) (cipher.Block, error) {
	h := md5.New()
	fmt.Fprint(h, key)
	cipherKey := h.Sum(nil)
	return aes.NewCipher(cipherKey)
}
