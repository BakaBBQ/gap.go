/*
  The MIT License (MIT)

  Copyright (c) 2015 Charles Liu

  Permission is hereby granted, free of charge, to any person obtaining a copy
  of this software and associated documentation files (the "Software"), to deal
  in the Software without restriction, including without limitation the rights
  to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
  copies of the Software, and to permit persons to whom the Software is
  furnished to do so, subject to the following conditions:

  The above copyright notice and this permission notice shall be included in
  all copies or substantial portions of the Software.

  THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
  IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
  FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
  AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
  LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
  OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
  THE SOFTWARE.
*/

package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"math/rand"
	"time"
)

type pptpSecret struct {
	client   string
	password string
}

const pptpConfigPath string = "/etc/ppp/chap-secrets"

const randPasswordLength int = 6


// I followed https://gobyexample.com/reading-files

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func parsesecrets() []pptpSecret {
	dat, err := ioutil.ReadFile(pptpConfigPath)
	check(err)

	lines := strings.Split(string(dat), "\n")

	var secrets []pptpSecret

	for _, oneLine := range lines {
		if i := strings.Index(oneLine, "#"); i != 0 {
			tokens := tokenizeLine(oneLine)

			if len(tokens) <= 1 {

			} else {
				newSecret := pptpSecret{tokens[0], tokens[2]}
				secrets = append(secrets, newSecret)
			}
		}

	}
	return secrets
}

func dumpsecrets(secrets *[]pptpSecret) {
	var buffer bytes.Buffer

	f, err := os.Create(pptpConfigPath)
	check(err)

	buffer.WriteString(serializesecrets(secrets))

	w := bufio.NewWriter(f)

	n, err := w.WriteString(buffer.String())
	check(err)
	fmt.Printf("wrote %d bytes \n", n)
	w.Flush()
	defer f.Close()
}

func serializesecrets(s *[]pptpSecret) string {
	var buffer bytes.Buffer

	for _, secret := range *s {
		buffer.WriteString(fmt.Sprintf("%s pptpd %s * \n", secret.client, secret.password))
	}
	return buffer.String()
}

func tokenizeLine(l string) []string {
	wordsWithSpaces := strings.Split(l, " ")
	return wordsWithSpaces
}

func genSecret(name string) pptpSecret{
	var letters = []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789!@#$^&*()")
	
	rand.Seed(time.Now().UTC().UnixNano())
	b := make([]rune, randPasswordLength)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	
	return pptpSecret{name, string(b)}
}
