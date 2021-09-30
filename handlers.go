package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

func unauthorized(w http.ResponseWriter, r *http.Request) {
	d := []string{}
	if err := tmpl.ExecuteTemplate(w, "unauthorized", d); err != nil {
		log.Printf("%v", err)
	}
}
func expired(w http.ResponseWriter, r *http.Request) {
	d := []string{}
	if err := tmpl.ExecuteTemplate(w, "expired", d); err != nil {
		log.Printf("%v", err)
	}
}
func logout(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	login := vars["login"]
	removeOnline(login)
	clearSession(w)
	time.Sleep(2 * time.Second)
	http.Redirect(w, r, "/", http.StatusMovedPermanently)
}
func home(w http.ResponseWriter, r *http.Request) {
	const (
		priv = "Localhost request incoming"
		pub  = "Public request incoming"
	)

	// POST
	if r.Method == "POST" {
		// Start tlsproxy
		if proxyon := r.FormValue("proxyon"); proxyon == "y" {
			startTLS()
			time.Sleep(4 * time.Second)
			http.Redirect(w, r, "/", http.StatusMovedPermanently)
		}
		//  Stop tlsproxy
		if proxyoff := r.FormValue("proxyoff"); proxyoff == "y" {
			stopTLS()
		}
		/*
			POST PUBLIC (Login)
		*/
		flogin := r.FormValue("login")
		fpassword := r.FormValue("password")
		if flogin != "" && fpassword != "" {
			auth := checkUser(flogin, fpassword)
			if auth {
				log.Printf("%s \n", "Authentication successful")
				setSession(flogin, w)
				err := writeOnline(flogin)
				if err != nil {
					log.Println(err)
				}
				log.Printf("%s %s %s\n", "Session successfully saved in", cookieName, "cookie")
				http.Redirect(w, r, "/chatSafe", http.StatusMovedPermanently)
			} else {
				log.Printf("%s \n", "Authentication failed")
			}
		}
	}

	// GET
	url := getURL()
	status, _ := checkLink(url)

	// Populate data
	d := &Index{
		URL:         url,
		Status:      status,
		PublicIPAPI: getpIP(true),
		PublicIPUDP: getpIP(false),
		EncMessages: getFiles(pathEncryptedMessages),
		EncFiles:    getFiles(pathEncryptedFiles),
	}

	// Check if incoming request is public
	request := isPublic(r.Host)
	if request {
		// Public
		log.Printf("%s \n", pub)
		if err := tmpl.ExecuteTemplate(w, "chatLogin", d); err != nil {
			log.Printf("%v", err)
		}

	} else {
		// Private
		log.Printf("%s \n", priv)
		if err := tmpl.ExecuteTemplate(w, "home", d); err != nil {
			log.Printf("%v", err)
		}
	}
}
func encrypt(w http.ResponseWriter, r *http.Request) {
	var d bool
	vars := mux.Vars(r)
	verb := vars["verb"]
	d = isConfirm(verb)

	if err := tmpl.ExecuteTemplate(w, "encrypt", d); err != nil {
		log.Printf("%v", err)
	}
}
func encryptPOST(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		// we separate encrypted messages and encrypted files
		var isFile bool

		// Collect plaintext
		passphrase := r.FormValue("passphrase")
		message := r.FormValue("message")
		filename := r.FormValue("filename")

		// Encrypt message
		if message != "" {
			isFile = false
			encryptionClient := newClient()
			keyHash := encryptionClient.GenerateMD5String(passphrase)
			AES := encryptionClient.EncryptAES([]byte(message), []byte(keyHash))
			save := encryptionClient.saveToFile(AES, filename, isFile)
			if save != nil {
				log.Println(save)
			}
		}

		// Collect files
		files := r.MultipartForm.File["files"]

		// Upload & encrypt files
		if len(files) > 0 {
			isFile = true
			UploadFiles(w, r, "files", passphrase)
		}

		log.Println("Message encrypted and saved successfully")
		http.Redirect(w, r, "/encrypt/successfully", http.StatusMovedPermanently)
	}
}
func decrypt(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	verb := vars["verb"]
	c := isConfirm(verb)
	dec := getFiles(pathToDecrypt)
	d := &Decrypt{
		Confirm:   c,
		DecFiles:  dec,
		Decrypted: decryptedMessage,
	}
	if err := tmpl.ExecuteTemplate(w, "decrypt", d); err != nil {
		log.Printf("%v", err)
	}
}
func decryptPOST(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		passphrase := r.FormValue("passphrase")
		filename := r.FormValue("filename")
		encryptionClient := newClient()
		encryptedContent := encryptionClient.readFile(filename)
		keyHash := encryptionClient.GenerateMD5String(passphrase)
		plainContent := encryptionClient.DecryptAES(encryptedContent, []byte(keyHash))

		encryptedContentIsFile := checkFilename(filename)

		if encryptedContentIsFile {
			saveIt := encryptionClient.decryptToFile(plainContent, filename)
			if saveIt != nil {
				log.Println(saveIt)
			}
			decryptedMessage = pathDecrypted + "/" + filename[:len(filename)-4]
		} else {
			decryptedMessage = string(plainContent)
		}

		http.Redirect(w, r, "/decrypt/successfully", http.StatusMovedPermanently)
		log.Println("File successfully decrypted")
	}
}

// WITH MIDDLEWARE
var chatSafe = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	claims, err := getSession(r)
	if err != nil {
		log.Println(err)
		return
	}

	users := getUsers()

	d := &Chat{
		Messages: messages,
		Users:    users,
		Login:    claims.Login,
		Files:    getFiles(pathToShare),
	}

	if err := tmpl.ExecuteTemplate(w, "chatSafe", d); err != nil {
		log.Printf("%v", err)
	}
})

var addMessage = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		login := r.FormValue("login")
		message := r.FormValue("message")
		if message != "" {
			now := time.Now()
			messages = append(messages, Message{Content: message, User: login, Time: now.Format("2006-01-02 15:04:05")})
		}

		http.Redirect(w, r, "/chatSafe", http.StatusSeeOther)
	}
})

var getMessages = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	//update content type
	w.Header().Set("Content-Type", "application/json")

	//specify HTTP status code
	w.WriteHeader(http.StatusOK)

	//convert struct to JSON
	jsonResponse, err := json.Marshal(messages)
	if err != nil {
		return
	}

	//update response
	w.Write(jsonResponse)
})
