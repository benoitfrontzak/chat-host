package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/rand"
	"encoding/hex"
	"io"
	"io/ioutil"
	"log"
	mathRand "math/rand"
	"os"
	"time"
)

type Client struct{}

// Hashing string to Compatible Cipher Keys
// We’re going to be using a simple MD5 hash.
// It is insecure, but it doesn’t really matter since we won’t be storing the output.
// The function will take a passphrase or any string, hash it, then return the hash as a hexadecimal value.
// Remember, we just need keys that meet the length criteria that AES demands.

func newClient() (c *Client) {
	mathRand.Seed(time.Now().UnixNano())
	c = &Client{}
	return
}

func (c *Client) EncryptAES(plainData, secret []byte) (cipherData []byte) {
	block, _ := aes.NewCipher(secret)
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return
	}
	nonce := make([]byte, gcm.NonceSize())
	_, err = io.ReadFull(rand.Reader, nonce)
	if err != nil {
		return
	}

	cipherData = gcm.Seal(
		nonce,
		nonce,
		plainData,
		nil)

	return
}

func (c *Client) DecryptAES(cipherData, secret []byte) (plainData []byte) {
	block, err := aes.NewCipher(secret)
	if err != nil {
		return
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return
	}
	nonceSize := gcm.NonceSize()

	nonce, ciphertext := cipherData[:nonceSize], cipherData[nonceSize:]
	plainData, err = gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return
	}
	return
}

func (c *Client) GenerateMD5String(passphrase string) (result string) {
	hasher := md5.New()
	hasher.Write([]byte(passphrase))
	return hex.EncodeToString(hasher.Sum(nil))
}

func (c *Client) saveToFile(cipherData []byte, filename string, isFile bool) (err error) {
	var path string
	if isFile {
		path = pathEncryptedFiles
	} else {
		path = pathEncryptedMessages
	}
	ff, err := os.Create(path + pathSeparator + filename + ".aes")
	if err != nil {
		return err
	}
	defer ff.Close()
	ff.Write(cipherData)
	return nil
}
func (c *Client) decryptToFile(cipherData []byte, filename string) (err error) {
	f := filename[:len(filename)-4]
	ff, err := os.Create(pathDecrypted + pathSeparator + f)
	if err != nil {
		return err
	}
	defer ff.Close()
	ff.Write(cipherData)
	return nil
}
func (c *Client) readFile(filename string) (cipherData []byte) {
	content, err := ioutil.ReadFile(pathToDecrypt + pathSeparator + filename)
	if err != nil {
		log.Println(err.Error())
	}
	return content
}
