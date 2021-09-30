package main

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
)

// Return users from dbUsers
// users = map userName => isOnline
func getUsers() map[string]bool {
	m := make(map[string]bool)

	f, err := os.Open(dbUsers)
	if err != nil {
		fmt.Printf("error opening file: %v\n", err)
		return m
	}

	r := bufio.NewReader(f)
	s, e := Readln(r)

	for e == nil {
		l := strings.Split(string(s), "/") // one line of users.txt
		u := l[0]                          // user from users.txt
		m[u] = isOnline(u)
		s, e = Readln(r)
	}
	return m
}

// Return true when given user/password match local file
// which defined our allowed users
func checkUser(user, password string) bool {
	f, err := os.Open(dbUsers)
	if err != nil {
		fmt.Printf("error opening file: %v\n", err)
		os.Exit(1)
	}
	r := bufio.NewReader(f)
	s, e := Readln(r)
	for e == nil {
		l := strings.Split(string(s), "/") // one line of users.txt
		u := l[0]                          // user from users.txt
		p := l[1]                          // password from users.txt
		if u == user {
			if p == password {
				return true
			}
		}
		s, e = Readln(r)
	}

	return false
}

// Create a file login.txt in tmp
// when login.txt is found in tmp
// the user will be tag online
func writeOnline(l string) error {
	_, err := os.Create(pathTMP + pathSeparator + l + ".txt")
	if err != nil {
		fmt.Printf("error opening file: %v\n", err)
		return err
	}
	return nil
}

// Take login as parameter and return error
// Simply delete the txt from tmp
func removeOnline(l string) error {
	err := os.Remove(pathTMP + pathSeparator + l + ".txt")
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

// Check if login is online
// Simply check if txt exist in tmp
func isOnline(l string) bool {
	_, err := os.Stat("tmp" + pathSeparator + l + ".txt")
	if err == nil {
		return true
	}
	if errors.Is(err, os.ErrNotExist) {
		return false
	}
	return false
}
