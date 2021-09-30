package main

import (
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

// Read uploaded files from the form post
// Encrypt each file content to pathEncryptedFiles
func UploadFiles(w http.ResponseWriter, r *http.Request, f, passphrase string) {
	isFile := true // we separate encrypted messages and encrypted files
	// collect files from form
	formdata := r.MultipartForm
	files := formdata.File[f]
	// for each file uploaded
	for i := range files {
		filename := files[i].Filename
		// open uploaded file
		file, err := files[i].Open()
		if err != nil {
			log.Println(w, err)
			return
		}
		defer file.Close()
		// create file to be encrypted
		out, err := os.Create(pathDecrypted + pathSeparator + filename)
		if err != nil {
			log.Printf("%s\n", "Unable to create the file for writing. Check your write access privilege")
			return
		}
		defer out.Close()
		// copy uploaded file content to new created file
		_, err = io.Copy(out, file) // file not files[i] !
		if err != nil {
			log.Println(w, err)
			return
		}
		// read new created file
		content, err := ioutil.ReadFile(pathDecrypted + pathSeparator + filename)
		if err != nil {
			log.Println(err.Error())
		}
		// encrypt new created file
		encryptionClient := newClient()
		keyHash := encryptionClient.GenerateMD5String(passphrase)
		AES := encryptionClient.EncryptAES(content, []byte(keyHash))
		// save it
		save := encryptionClient.saveToFile(AES, filename, isFile)
		if save != nil {
			log.Println(save)
		}
	}

}
