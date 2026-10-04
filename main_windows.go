//go:build windows

package main

import (
	"fmt"
	"os"

	"golang.org/x/sys/windows/svc"
)

func main() {
	isService, err := svc.IsAnInteractiveSession()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to determine session type: %v\n", err)
		os.Exit(1)
	}
	if !isService {
		if err := svc.Run(ServiceName, &serviceRunner{}); err != nil {
			fmt.Fprintf(os.Stderr, "service error: %v\n", err)
			os.Exit(1)
		}
		return
	}
	runConsole()
}

func runConsole() {
	logger := NewLogger()
	defer logger.Close()
	logger.Info("Service started")
	fmt.Println("=== College Wi‑Fi Auto Login CLI ===")
	fmt.Printf("1) Target SSID: %s\n", TargetSSID)

	ssid, err := CurrentSSID()
	if err != nil {
		fmt.Printf("2) SSID detection failed: %v\n", err)
		logger.Warn(fmt.Sprintf("SSID detection failed: %v", err))
		return
	}
	fmt.Printf("2) Current SSID: %s\n", ssid)

	isTarget, err := IsTargetSSID()
	if err != nil {
		fmt.Printf("3) Target check failed: %v\n", err)
		logger.Warn(fmt.Sprintf("Target check failed: %v", err))
		return
	}
	if !isTarget {
		fmt.Println("3) Not connected to target SSID. Nothing to do.")
		logger.Info("Not connected to target SSID")
		return
	}
	fmt.Println("3) Connected to target SSID.")

	internetOK := CheckInternetConnectivity()
	fmt.Printf("4) Internet connectivity: %t\n", internetOK)
	if internetOK {
		fmt.Println("4) Internet available. No login needed.")
		logger.Info("Internet available")
		return
	}

	portalDetected := DetectCaptivePortal()
	fmt.Printf("5) Captive portal detected: %t\n", portalDetected)
	if !portalDetected {
		fmt.Println("5) Portal not detected. Skipping login.")
		logger.Warn("Portal not detected")
		return
	}

	username, password, err := LoadCredentials()
	if err != nil {
		fmt.Printf("6) Credential load failed: %v\n", err)
		logger.Warn(fmt.Sprintf("Credential load failed: %v", err))
		return
	}
	if username == "" || password == "" {
		fmt.Println("6) No credentials saved. Please save username/password first.")
		logger.Warn("No credentials loaded")
		return
	}
	fmt.Printf("6) Loaded username: %s\n", username)

	fmt.Println("7) Attempting portal login...")
	if err := Login(username, password); err != nil {
		fmt.Printf("7) Login failed: %v\n", err)
		logger.Warn(fmt.Sprintf("Login failed: %v", err))
		return
	}

	fmt.Println("8) Login attempt complete.")
	logger.Info("Authentication successful")
	fmt.Println("9) Final check: Internet should now be available.")
}
