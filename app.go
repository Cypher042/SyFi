package main

import (
	"embed"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
)

//go:embed frontend
var frontendAssets embed.FS

type App struct{}

type Credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func shouldRunService(args []string, env map[string]string) bool {
	for _, arg := range args {
		if strings.EqualFold(arg, "--service") || strings.EqualFold(arg, "-service") {
			return true
		}
	}

	lookup := func(key string) string {
		if env == nil {
			return os.Getenv(key)
		}
		return env[key]
	}

	value := strings.TrimSpace(strings.ToLower(lookup("SYFI_SERVICE")))
	return value == "1" || value == "true" || value == "yes"
}

func (a *App) LoadCredentials() (Credentials, error) {
	username, password, err := LoadCredentials()
	if err != nil {
		return Credentials{}, err
	}
	return Credentials{Username: username, Password: password}, nil
}

func (a *App) Status() (string, error) {
	ssid, err := CurrentSSID()
	if err != nil {
		return "", fmt.Errorf("SSID detection failed: %w", err)
	}

	isTarget, err := IsTargetSSID()
	if err != nil {
		return "", fmt.Errorf("target SSID validation failed: %w", err)
	}

	return fmt.Sprintf(
		"SSID: %s\nTarget Wi-Fi: %t\nInternet connected: %t\nCaptive portal detected: %t",
		ssid,
		isTarget,
		CheckInternetConnectivity(),
		DetectCaptivePortal(),
	), nil
}

func (a *App) ServiceInfo() (ServiceInfo, error) {
	return serviceInfo()
}

func (a *App) Logs() (string, error) {
	return ReadRecentLogs()
}

func (a *App) InstallService() (ServiceInfo, error) {
	if err := installService(); err != nil {
		return ServiceInfo{}, err
	}
	return serviceInfo()
}

func (a *App) UninstallService() (ServiceInfo, error) {
	if err := uninstallService(); err != nil {
		return ServiceInfo{}, err
	}
	return serviceInfo()
}

func (a *App) SaveCredentials(username, password string) error {
	if err := SaveCredentials(strings.TrimSpace(username), strings.TrimSpace(password)); err != nil {
		return fmt.Errorf("save failed: %w", err)
	}
	return nil
}

func (a *App) Login(username, password string) error {
	username = strings.TrimSpace(username)
	password = strings.TrimSpace(password)
	if username == "" || password == "" {
		return errors.New("portal username and password are required")
	}
	if err := SaveCredentials(username, password); err != nil {
		return fmt.Errorf("failed to save credentials: %w", err)
	}

	ssid, err := CurrentSSID()
	if err != nil {
		return fmt.Errorf("unable to detect SSID: %w", err)
	}

	isTarget, err := IsTargetSSID()
	if err != nil {
		return fmt.Errorf("SSID validation failed: %w", err)
	}
	if !isTarget {
		return fmt.Errorf("not connected to target Wi-Fi; current SSID: %s", ssid)
	}
	if CheckInternetConnectivity() {
		return errors.New("internet is already available; no login needed")
	}
	if !DetectCaptivePortal() {
		return errors.New("no captive portal was detected")
	}
	if err := Login(username, password); err != nil {
		return fmt.Errorf("portal login failed: %w", err)
	}
	return nil
}

func runWailsGUI() {
	err := wails.Run(&options.App{
		Title:  "SyFi - Captive Portal",
		Width:  920,
		Height: 620,
		AssetServer: &assetserver.Options{
			Assets: frontendAssets,
		},
		Bind: []interface{}{
			&App{},
		},
	})
	if err != nil {
		fmt.Printf("GUI error: %v\n", err)
	}
}
