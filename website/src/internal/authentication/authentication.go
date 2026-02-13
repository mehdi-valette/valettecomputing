package authentication

import (
	"crypto/rand"
	"errors"

	"valette.software/internal/config"
)

var sessions map[string]struct{}
var password string

var ErrWrongPassword = errors.New("wrong password")

func Init(config config.Configurator) {
	sessions = make(map[string]struct{})
	password = config.GetAdminPassword()
}

func CheckSession(sessionId string) bool {
	_, ok := sessions[sessionId]

	return ok
}

func Authenticate(pwd string) (string, error) {
	if pwd == password {
		sessionId := rand.Text()
		sessions[sessionId] = struct{}{}
		return sessionId, nil
	}

	return "", ErrWrongPassword
}

func Logout() {
	sessions = make(map[string]struct{})
}
