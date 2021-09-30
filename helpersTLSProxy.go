// This helpers is used to start|stop the service tlsproxy
// (establish a TLS termination proxy between your public IP and your localhost)
// When the service start it generates a hidden file $HOME/.ssh/.tunneled:
//
// The process id is:9016
// 68487f8e8af111.localhost.run tunneled with tls termination, https://68487f8e8af111.localhost.run

package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

// Returns client homeDir
func HomeDir() string {
	home, err := os.UserHomeDir()
	if err != nil {
		log.Printf("==> Can't retrieve home directory: %v", err)
	}
	return home
}

// This function is used by stopTLS
// returns the pid of tlsproxy from .tunneled
func getPID() int {
	h := HomeDir()
	p := h + pathSeparator + ".ssh" + pathSeparator + ".tunneled"

	fileIO, err := os.OpenFile(p, os.O_RDONLY, 0400)
	if err != nil {
		log.Println(err)
	}
	defer fileIO.Close()
	rawBytes, err := ioutil.ReadAll(fileIO)
	if err != nil {
		log.Println(err)
	}
	lines := strings.Split(string(rawBytes), "\n")
	first := string(lines[0])
	pid, _ := strconv.Atoi(first[18:])
	return pid
}

// Stop tlsproxy service by killing the pid
func stopTLS() {
	pid := getPID()
	p, _ := os.FindProcess(pid)
	p.Kill()
	log.Println("Process killed")
}

// Start the service tlsproxy
func startTLS() {
	tlsproxy := Executable{
		Name:   "tlsproxy",
		Params: []string{},
	}

	go func() { tlsproxy.Run() }()
}

// This function is used by getURL
// to parse a given line and
// return what is before the given delimiter
func before(l, d string) string {
	pos := strings.Index(l, d)
	if pos == -1 {
		return ""
	}
	return l[0:pos]
}

// Returns the url generated from .tunneled
// Used to display URL in home page
// We don't check $HOME/.shh/.tunneled
// since it is required by the app and part of the setup
func getURL() string {
	h := HomeDir()
	p := h + pathSeparator + ".ssh" + pathSeparator + ".tunneled"

	fileIO, err := os.OpenFile(p, os.O_RDONLY, 0400)
	if err != nil {
		log.Println(err)
		return "noFile"
	}
	defer fileIO.Close()
	rawBytes, err := ioutil.ReadAll(fileIO)
	if err != nil {
		log.Println(err)
		return "noDataYet"
	}

	lines := strings.Split(string(rawBytes), "\n")
	second := string(lines[1])
	url := before(second, ".localhost.run tunneled with tls termination")
	return url
}

var client = http.Client{
	Transport: nil,
	CheckRedirect: func(req *http.Request, via []*http.Request) error {
		return nil
	},
	Jar:     nil,
	Timeout: 2 * time.Second,
}

// Returns true|false if link is alive
func checkLink(link string) (int, error) {
	url := "https://" + link + ".localhost.run/"
	req, err := http.NewRequest("HEAD", url, nil)
	if err != nil {
		return 0, err
	}
	resp, err := client.Do(req)
	if err != nil {
		return 0, err
	}
	resp.Body.Close()
	return resp.StatusCode, nil
}
