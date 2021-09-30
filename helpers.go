package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"strings"
)

func isConfirm(text string) bool {
	return text != "new"
}

// Returns true if url is public or false when private
// Since the public view run only with the tls proxy
// we get the url as parameter and check if it contains localhost.run
func isPublic(url string) bool {
	return strings.Contains(url, "localhost.run")
}

// Readln returns a single line (without the ending \n)
// from the input buffered reader.
// An error is returned if there is an error with the
// buffered reader.
func Readln(r *bufio.Reader) (string, error) {
	var (
		isPrefix bool  = true
		err      error = nil
		line, ln []byte
	)
	for isPrefix && err == nil {
		line, isPrefix, err = r.ReadLine()
		ln = append(ln, line...)
	}
	return string(ln), err
}

// Returns either the public IP (api true)
// or the local IP (API false)
func getpIP(api bool) string {
	if api {
		req, err := http.Get("http://ip-api.com/json/")
		if err != nil {
			return err.Error()
		}
		defer req.Body.Close()

		body, err := ioutil.ReadAll(req.Body)
		if err != nil {
			return err.Error()
		}

		var ip IP
		json.Unmarshal(body, &ip)

		return ip.Query
	} else {
		conn, err := net.Dial("udp", "8.8.8.8:80")
		if err != nil {
			fmt.Println(err)
		}

		defer conn.Close()
		localAddr := conn.LocalAddr().String()
		idx := strings.LastIndex(localAddr, ":")

		return localAddr[0:idx]
	}
}

// Returns files from given folder
// files = map fileName => fileSize
func getFiles(folder string) map[string]int64 {
	files, err := ioutil.ReadDir(folder)
	if err != nil {
		log.Printf("%v", err)
	}
	// We map file name => file size
	m := make(map[string]int64)
	for _, f := range files {
		file, _ := os.Open(folder + pathSeparator + f.Name())
		stat, err := file.Stat()
		if err != nil {
			log.Printf("%v", err)
		}

		m[f.Name()] = stat.Size()
	}

	return m
}

// Returns true|false if filename is file type
// if it is file type filename = [a-z].ext.aes
// if it is message type filename = [a-z].aes
func checkFilename(filename string) bool {
	removeAES := filename[:len(filename)-4]
	ext := removeAES[len(removeAES)-3:]
	var isFile bool
	switch ext {
	// Audio file
	case "aif", "cda", "mid", "mp3", "mpa", "ogg", "wav", "wma", "wpl":
		isFile = true
	// Image file
	case "jpg", "gif", "png", "bmp", "ico", "psd", "svg", "tif":
		isFile = true
	// Application file
	case "pdf", "xls", "doc", "txt":
		isFile = true
	// Video file
	case "mp4", "avi", "mpg", "mkv", "mov", "3g2", "3gp", "flv", "m4v", "wmv":
		isFile = true
	default:
		isFile = false
	}

	return isFile
}
