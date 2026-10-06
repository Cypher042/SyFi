package main

import (
	"errors"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func userLocalAppDataDir() string {
	if programData := os.Getenv("ProgramData"); programData != "" {
		return filepath.Join(programData, "CollegeWiFiAutoLogin")
	}
	if os.Getenv("LOCALAPPDATA") != "" {
		return filepath.Join(os.Getenv("LOCALAPPDATA"), "CollegeWiFiAutoLogin")
	}
	if os.Getenv("APPDATA") != "" {
		return filepath.Join(os.Getenv("APPDATA"), "CollegeWiFiAutoLogin")
	}
	return filepath.Join(".", "localdata")
}

func CheckInternetConnectivity() bool {
	client := &http.Client{Timeout: 8 * time.Second}
	resp, err := client.Get(ConnectivityCheckURL)
	if err != nil {
		return false
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 200 && resp.StatusCode < 400 {
		return true
	}
	return false
}

func DetectCaptivePortal() bool {
	client := &http.Client{
		Timeout: 8 * time.Second,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			if len(via) >= 5 {
				return errors.New("too many redirects")
			}
			return nil
		},
	}

	resp, err := client.Get(PortalCheckURL)
	if err != nil {
		return false
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 && resp.StatusCode < 400 {
		return true
	}

	body, err := io.ReadAll(io.LimitReader(resp.Body, 2048))
	if err != nil {
		return false
	}
	text := strings.ToLower(string(body))
	for _, token := range []string{"login", "username", "password", "portal", "authenticate"} {
		if strings.Contains(text, token) {
			return true
		}
	}
	return false
}
