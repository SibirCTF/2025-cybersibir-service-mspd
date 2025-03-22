package utils

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"net/http"
	"regexp"
)

func CheckAuthKey(auth_key string, username string, r *http.Request) (err error) {
	if r.Referer() != "https://www.gosuslugi.ru/" {
		return errors.New("bad request")
	}
	if len(auth_key) != 108 {
		return errors.New("auth key is too short")
	}
	ending := auth_key[len(auth_key)-32:]
	hash := md5.Sum([]byte(username))
	hashed_str := hex.EncodeToString(hash[:])
	if ending != hashed_str {
		return errors.New("failed to check auth key")
	}
	re := regexp.MustCompile(`^NXvP\d{3}[a-zA-Z]{2}-\d{6}#([a-zA-Z]{6})-(?:([A-Z][a-z]){3})-(\d{10})@(\d+)#\?{3}[a-zA-Z0-9]{32}$`)
	if re.MatchString(auth_key) {
		return nil
	}
	return errors.New("failed to check auth key")
}
