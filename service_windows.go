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
		case <-time.After(MonitorInterval):
			if ok, err := IsTargetSSID(); err == nil && ok {
				if !CheckInternetConnectivity() && DetectCaptivePortal() {
					user, pass, err := LoadCredentials()
					if err == nil && user != "" && pass != "" {
						if err := Login(user, pass); err == nil {
							logger.Info("Authentication successful")
						} else {
							logger.Warn("Authentication failed")
						}
					} else {
						logger.Warn("No credentials available for portal login")
					}
				}
			}
		}
	}
}
