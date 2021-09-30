package main

import (
	"html/template"
	"log"
	"os"

	"github.com/joho/godotenv"
)

const (
	cookieName            = "anonymous"
	claimsIssuer          = "anonymous-issuer"
	cipherKey             = "ekuy@x}(~q3qZE;w"
	pathSeparator         = string(os.PathSeparator)
	pathTMP               = "tmp"
	pathEncryptedMessages = "messages-encrypted"
	pathEncryptedFiles    = "files-encrypted"
	pathToDecrypt         = "files-toDecrypt"
	pathToShare           = "assets" + pathSeparator + "files-toShare"
	pathDecrypted         = "files-decrypted"
	dbUsers               = "users.txt"
)

var (
	tmpl             = template.Must(template.ParseGlob("templates/*"))
	messages         = []Message{}
	decryptedMessage = ""
)

// Decrypt is used to populate Decrypt (private view).
type Decrypt struct {
	Confirm   bool
	DecFiles  map[string]int64
	Decrypted string
}

// Index is used to populate home (private view).
type Index struct {
	URL         string
	Status      int
	PublicIPAPI string
	PublicIPUDP string
	EncMessages map[string]int64
	EncFiles    map[string]int64
}

// IP is used by Index (when query public API)
type IP struct {
	Query string
}

// Chat is used to populate chat room
type Chat struct {
	Messages []Message
	Users    map[string]bool
	Login    string
	Files    map[string]int64
}

// Message is used by Chat
type Message struct {
	Content string
	User    string
	Time    string
}

func setLog() {
	log.SetPrefix("[INFO]: ")
	log.SetFlags(0) // remove file:line and timestamps from log liness
}

// This function take 2 parameters:
// 	expected Variable to be found in .env
// 	expected Default value if not found
// Check if .env exist
// Check if variable set inside .env
// Returns either .env value or default value
func getEnv(v, d string) string {
	setLog()
	// Get env var
	if err := godotenv.Load(); err != nil {
		log.Printf("%s\n", "error loading .env file, set default value to "+d)
		return d
	}

	// Defined default value If not defined in env
	val := os.Getenv(v)
	if val == "" {
		val = ""
		log.Printf("%s %s %s%s%s\n", "error loading", v, ", set default value to `", d, "`")
		return val
	}

	return val
}
