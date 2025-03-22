package utils

import (
	"math/rand"
	"sibir2025/service/config"

	"github.com/gorilla/securecookie"
)

var SC *securecookie.SecureCookie

func init() {
	SC = securecookie.New([]byte(config.KeyDict[rand.Int()%len(config.KeyDict)]), []byte(config.KeyDict[rand.Int()%len(config.KeyDict)]))
}
