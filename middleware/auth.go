// Copyright 2014 Manu Martinez-Almeida.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package middleware

import (
	"crypto/md5"
	"encoding/base64"
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
)

// AuthUserKey is the cookie name for user credential in basic auth.
const AuthUserKey = "Authorization230x"

// Accounts defines a key/value for user/pass list of authorized logins.
type Accounts map[string]string

type authPair struct {
	value string
	user  string
}

type authPairs []authPair

type authResult struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    string `json:"data"`
}

func (a authPairs) searchCredential(authValue string) (string, bool) {
	if authValue == "" {
		return "", false
	}
	for _, pair := range a {
		if pair.value == authValue {
			return pair.user, true
		}
	}
	return "", false
}

// CookieAuthForRealm returns a Basic HTTP Authorization middleware. It takes as arguments a map[string]string where
// the key is the user name and the value is the password
func CookieAuthForRealm(accounts Accounts) gin.HandlerFunc {
	pairs := processAccounts(accounts)
	return func(c *gin.Context) {
		var user string
		if strings.Contains(c.Request.URL.Path, "logout") {
			c.SetCookie("auth", "", -1, "/", c.GetHeader("Host"), false, true)
			c.JSON(200, authResult{Code: 200, Message: "ok"})
			return
		}
		// Search user in the slice of allowed credentials
		val, _ := c.Cookie("auth")
		foundCookie := false
		if val != "" {
			user, foundCookie = pairs.searchCredential(val)
		}
		if foundCookie {
			c.Set(AuthUserKey, user)
			if strings.Contains(c.Request.URL.Path, "login") {
				c.JSON(200, authResult{Code: 200, Message: "already logged", Data: user})
				c.Abort()
			}
			return
		}
		name := c.Query("name")
		password := c.Query("password")
		if name == "" || password == "" {
			c.JSON(401, authResult{Code: 401, Message: "not login"})
			c.Abort()
			return
		}
		encryption := authorizationCookie(name, password)
		user, foundPass := pairs.searchCredential(encryption)
		if foundPass {
			c.SetCookie("auth", encryption, 86400, "/", c.GetHeader("Host"), false, false)
			c.Set(AuthUserKey, user)
			if strings.Contains(c.Request.URL.Path, "login") {
				c.JSON(200, authResult{Code: 200, Message: "ok", Data: encryption})
			}
			return
		} else {
			c.JSON(401, authResult{Code: 401, Message: "name or password incorrect"})
			c.Abort()
			return
		}
		// The user credentials was found, set user's id to key AuthUserKey in this context, the user's id can be read later using
		// c.MustGet(gin.AuthUserKey).
	}
}

// CookieAuth returns a Basic HTTP Authorization middleware. It takes as argument a map[string]string where
// the key is the user name and the value is the password.
func CookieAuth(accounts Accounts) gin.HandlerFunc {
	return CookieAuthForRealm(accounts)
}

func processAccounts(accounts Accounts) authPairs {
	pairs := make(authPairs, 0, len(accounts))
	if len(accounts) == 0 {
		fmt.Println("Empty list of authorized credentials")
		return pairs
	}
	for user, password := range accounts {
		value := authorizationCookie(user, password)
		pairs = append(pairs, authPair{
			value: value,
			user:  user,
		})
	}
	return pairs
}

func authorizationCookie(user, password string) string {
	data := []byte(password)
	has := md5.Sum(data)
	base := user + ":" + fmt.Sprintf("%x", has)
	return base64.StdEncoding.EncodeToString([]byte(base))
}
