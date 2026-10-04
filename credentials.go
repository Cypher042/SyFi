package main

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"strings"
)

type savedCredentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func credentialsPath() string {
	base := appLogDirectory()
	_ = os.MkdirAll(base, 0o755)
	return filepath.Join(base, "credentials.json")
}

func SaveCredentials(username, password string) error {
	if strings.TrimSpace(username) == "" {
		return errors.New("username is required")
	}
	if strings.TrimSpace(password) == "" {
		return errors.New("password is required")
	}
	return saveCredentials(username, password)
}

func LoadCredentials() (string, string, error) {
	return loadCredentials()
}

func saveSimpleCredentials(username, password string) error {
	payload := savedCredentials{Username: username, Password: password}
	content, err := json.MarshalIndent(payload, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(credentialsPath(), content, 0o600)
}

func loadSimpleCredentials() (string, string, error) {
	content, err := os.ReadFile(credentialsPath())
	if err != nil {
		if os.IsNotExist(err) {
			return "", "", nil
		}
		return "", "", err
	}
	var payload savedCredentials
	if err := json.Unmarshal(content, &payload); err != nil {
		return "", "", err
	}
	if payload.Username == "" {
		return "", "", nil
	}
	return payload.Username, payload.Password, nil
}

