//go:build windows

package main

import (
	"time"

	"golang.org/x/sys/windows/svc"
)

type serviceRunner struct{}

func (s *serviceRunner) Execute(args []string, requests <-chan svc.ChangeRequest, statuses chan<- svc.Status) (bool, uint32) {
	logger := NewLogger()
	defer logger.Close()
	logger.Info("Service started")
	statuses <- svc.Status{State: svc.Running, Accepts: svc.AcceptStop | svc.AcceptShutdown}

	ticker := time.NewTicker(MonitorInterval)
	defer ticker.Stop()

	checkConnection := func() {
		logger.Info("Connection check started")

		isTarget, err := IsTargetSSID()
		if err != nil {
			logger.Warn("Target SSID check failed: " + err.Error())
			return
		}
		if !isTarget {
			logger.Info("Not connected to target SSID: " + TargetSSID)
			return
		}
		logger.Info("Connected to target SSID: " + TargetSSID)

		if CheckInternetConnectivity() {
			logger.Info("Internet is available; no login needed")
			return
		}
		logger.Info("Internet is unavailable")

		if !DetectCaptivePortal() {
			logger.Info("No captive portal detected")
			return
		}
		logger.Info("Captive portal detected")

		user, pass, err := LoadCredentials()
		if err != nil {
			logger.Warn("Credential load failed: " + err.Error())
		} else if user != "" && pass != "" {
			if err := Login(user, pass); err == nil {
				logger.Info("Authentication successful")
			} else {
				logger.Warn("Authentication failed: " + err.Error())
			}
		} else {
			logger.Warn("No credentials available for portal login")
		}
	}

	go checkConnection()
	for {
		select {
		case req := <-requests:
			switch req.Cmd {
			case svc.Interrogate:
				statuses <- req.CurrentStatus
			case svc.Stop, svc.Shutdown:
				logger.Info("Service stopping")
				statuses <- svc.Status{State: svc.StopPending}
				return false, 0
			}
		case <-ticker.C:
			checkConnection()
		}
	}
}
