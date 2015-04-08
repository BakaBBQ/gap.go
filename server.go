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
	"fmt"
	"github.com/eknkc/amber"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"os/signal"
	"time"
	"os/exec"
	"regexp"
)

type vpnSession struct {
	pptpSecret
	startTime time.Time
	permanent bool
}

var maxSessions int = 3

func main() {
	router := gin.Default()

	var vpnSessions []vpnSession
	
	reflectPermanentSessions(&vpnSessions)
	
	vpnSessionTicker := time.NewTicker(time.Millisecond * 5000)
	go func() {
		for t := range vpnSessionTicker.C {
			activeCleanVpns(&vpnSessions, &t)
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for sig := range c {
			fmt.Println("Cleaning up dead vpns", sig)
			cleanDeadSessions(&vpnSessions)
			os.Exit(1)
		}
	}()

	router.Static("/assets", "./assets")
	router.GET("/", func(c *gin.Context) {
		template, err := amber.CompileFile("templates/index.amber", amber.DefaultOptions)
		
		obj := gin.H{"currentUserAmount": getCurrentSessionAmount(&vpnSessions)}
		
		if err != nil {
			panic(err)
		}
		router.SetHTMLTemplate(template)
		c.HTML(http.StatusOK, "index.amber", obj)
	})
	
	router.GET("/sessionCnt", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"sessionCnt": getCurrentSessionAmount(&vpnSessions)})
	})
	
	
	userRegex, _ := regexp.Compile("^[a-z]{2,}-[a-z]{2,}$")
	router.GET("/newSession", func(c *gin.Context) {
		c.Request.ParseForm()
		
		secretWord := c.Request.Form.Get("secret")
		userName := c.Request.Form.Get("user")
		
		//fmt.Println("It matches?",userRegex.MatchString(userName))
		
		if secretWord != "42" {
			c.JSON(http.StatusOK, gin.H{"status": "passcode invalid"})
			return
		}
		
		if userName == "" {
			c.JSON(http.StatusOK, gin.H{"status": "username required"})
			return
		}
		
		
		if !userRegex.MatchString(userName) {
			c.JSON(http.StatusOK, gin.H{"status": "username format invalid"})
			return
		}
		
		
		if isUserAlreadyInSession(userName, &vpnSessions) {
			c.JSON(http.StatusOK, gin.H{"status": "user already in session"})
			return
		}
		
		
		
		newSession, success := addNewSession(userName, &vpnSessions)
		if success {
			clientName := (*newSession).pptpSecret.client
			passwordName := (*newSession).pptpSecret.password
			c.JSON(http.StatusOK, gin.H{"status": "ok", "username": clientName, "password": passwordName})
			return
		} else {
			c.JSON(http.StatusOK, gin.H{"status": "maximum session reached"})
			return
		}
		
		//c.JSON(http.StatusOK, gin.H{"sessionCnt": getCurrentSessionAmount(&vpnSessions)})
	})
	
	
	router.Run(":8080")
}

func reflectPermanentSessions(vpnSessions *[]vpnSession){
	rawSecrets := parsesecrets()
	for _, s := range rawSecrets {
		*vpnSessions = append((*vpnSessions), vpnSession{s, time.Now(), true})
	}
}

func activeCleanVpns(vpnSessions *[]vpnSession, t *time.Time) {
	timeNow := time.Now()
	for i, ele := range *vpnSessions {
		if (!ele.permanent) && timeNow.Sub(ele.startTime).Hours() >= 1.0 {
			*vpnSessions = append((*vpnSessions)[:i], (*vpnSessions)[i+1:]...)
			fmt.Println("Cleaned a vpn account at ", timeNow)
			manifestVpnSessions(vpnSessions)
		}
	}
	fmt.Println("Vpn sessions", *vpnSessions)
}

func addNewSession(username string, vpnSessions *[]vpnSession) (*vpnSession, bool){
	if getCurrentSessionAmount(vpnSessions) >= 3 {
		return nil, false
	} else {
		newSession := vpnSession{genSecret(username), time.Now(), false}
		*vpnSessions = append(*vpnSessions, newSession)
		manifestVpnSessions(vpnSessions)
		return &newSession, true
	}
	
	panic("You should never reach here")
}

func manifestVpnSessions(vpnSessions *[]vpnSession) {
	secrets := make([]pptpSecret, len(*vpnSessions))
	for i, s := range *vpnSessions {
		secrets[i] = s.pptpSecret
	}
	dumpsecrets(&secrets)
	
	out, err := exec.Command("sh","-c","service pptpd restart").Output()
	
	if err != nil {
		//panic(err)
	}
	fmt.Println(out)
}

func getCurrentSessionAmount(vpnSessions *[]vpnSession) int {
	cnt := 0
	for _, s := range *vpnSessions {
		if (! s.permanent) {
			cnt += 1
		}
	}
	return cnt
}

func isUserAlreadyInSession(username string, vpnSessions *[]vpnSession) bool {
	found := false
	
	for _, s := range *vpnSessions {
		if s.pptpSecret.client == username {
			found = true
		}
	}
	return found
}

func cleanDeadSessions(vpnSessions *[]vpnSession) {
	fmt.Println("Cleaning dead sessions")
	for i, ele := range *vpnSessions {
		if (!ele.permanent){
			*vpnSessions = append((*vpnSessions)[:i], (*vpnSessions)[i+1:]...)
			fmt.Println("Cleaned a dead vpn account")
			manifestVpnSessions(vpnSessions)
		}
	}
}
