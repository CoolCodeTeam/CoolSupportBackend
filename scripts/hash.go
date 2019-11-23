package main

import (
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"io/ioutil"
	"os"
	"strings"
)

func main() {
	pwd, _ := ioutil.ReadAll(os.Stdin)
	passwordString := string(pwd)
	passwordStrings := strings.Split(passwordString, "\n")
	for _, password := range passwordStrings[:len(passwordStrings)-1] {
		hash, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
		fmt.Println(string(hash))
	}
}
