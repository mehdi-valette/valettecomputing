package main

import (
	"log"
	"net/smtp"
	"os"
	"strings"
)

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
	user     string
	password string
	host     string
	port     string
	from     string
	to       []string
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

func beginWith(line string, begin string) bool {
	return len(line) >= len(begin) && line[0:len(begin)] == begin
}

func ReadConfig(file string) {
	data, err := os.ReadFile(file)
	smtpData := SmtpData{}

	if err != nil {
		log.Fatal(err)
	}

	for line := range strings.Lines(string(data)) {
		if beginWith(line, "smtp_user=") {
			smtpData.user = strings.Trim(line[10:], "\n")
		} else if beginWith(line, "smtp_password=") {
			smtpData.password = strings.Trim(line[14:], "\n")
		} else if beginWith(line, "smtp_from=") {
			smtpData.from = strings.Trim(line[10:], "\n")
		} else if beginWith(line, "smtp_to=") {
			smtpData.to = []string{strings.Trim(line[8:], "\n")}
		} else if beginWith(line, "smtp_host=") {
			smtpData.host = strings.Trim(line[10:], "\n")
		} else if beginWith(line, "smtp_port=") {
			smtpData.port = strings.Trim(line[10:], "\n")
		}
	}

	config.setData(smtpData, smtp.PlainAuth("", smtpData.user, smtpData.password, smtpData.host))
}

func GetConfig() Configurator {
	return config
}
