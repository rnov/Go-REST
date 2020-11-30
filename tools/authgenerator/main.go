package main

import (
	"crypto/sha256"
	b64 "encoding/base64"
	"fmt"
)

type basicAuth struct {
	user string
	pwd  string
	ba64 string
	hash string
}

// script that generates Basic Auth and prints the its values
func main() {
	baList := []basicAuth{
		{
			user: "username",
			pwd:  "password",
		},
		{
			user: "admin",
			pwd:  "A&BEtr*!n^51",
		},
		{
			user: "chef-1",
			pwd:  "I49Zq!0Pqc",
		},
	}

	for i, ba := range baList {
		aux := ba.user + ":" + ba.pwd
		ba64 := b64.StdEncoding.EncodeToString([]byte(aux))
		hash := hash(ba64)
		baList[i].ba64 = ba64
		baList[i].hash = hash
	}

	for _, ba := range baList {
		fmt.Println("username: " + ba.user + " password: " + ba.pwd + " base64: " + ba.ba64 + " hash: " + ba.hash)
	}
}

func hash(ba string) string {
	h := sha256.New()
	h.Write([]byte(ba))

	return fmt.Sprintf("%x", h.Sum(nil))
}
