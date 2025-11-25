package config

import (
	"errors"
	"log"
	"net/smtp"
	"os"
	"strings"
)

var errEqualSignMissing = errors.New("all config lines should have the form 'key=value'")
var errKeyEmpty = errors.New("the key must not be empty")
var errUnknown = errors.New("unknown error")

var config Configurator = &Config{}

type Configurator interface {
	GetSmtpAuth() smtp.Auth
	GetSmtp() SmtpData
	setData(data SmtpData, auth smtp.Auth)
}

type Config struct {
	SmtpAuth smtp.Auth
	SmtpData SmtpData
}

type SmtpData struct {
	User     string
	Password string
	Host     string
	Port     string
	From     string
	To       []string
}

var _ Configurator = &Config{}

func (c Config) GetSmtpAuth() smtp.Auth {
	return c.SmtpAuth
}

func (c Config) GetSmtp() SmtpData {
	return c.SmtpData
}

func (c *Config) setData(data SmtpData, auth smtp.Auth) {
	c.SmtpData = data
	c.SmtpAuth = auth
}

func getValue(line string) (string, string, error) {
	line = strings.Trim(line, "\n")

	if len(line) == 0 {
		return "", "", nil
	}

	equalIndex := strings.Index(line, "=")

	if equalIndex == -1 {
		return "", "", errEqualSignMissing
	}

	if equalIndex == 0 {
		return "", "", errKeyEmpty
	}

	key := line[0:equalIndex]

	if len(line) >= equalIndex && line[0:len(key)] == key {
		return key, line[equalIndex+1:], nil
	}

	return "", "", errUnknown
}

func ReadConfig(file string) {
	data, err := os.ReadFile(file)
	smtpData := SmtpData{}

	if err != nil {
		log.Fatal(err)
	}

	for line := range strings.Lines(string(data)) {
		key, value, err := getValue(line)

		if err != nil {
			log.Print(err)
		}

		switch key {
		case "smtp_from":
			smtpData.From = value
		case "smtp_host":
			smtpData.Host = value
		case "smtp_password":
			smtpData.Password = value
		case "smtp_port":
			smtpData.Port = value
		case "smtp_to":
			smtpData.To = []string{value}
		case "smtp_user":
			smtpData.User = value
		default:
			log.Printf("the key '%s' is unknown", key)
		}
	}

	config.setData(smtpData, smtp.PlainAuth("", smtpData.User, smtpData.Password, smtpData.Host))
}

func GetConfig() Configurator {
	return config
}
