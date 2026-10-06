//go:build windows

package main

import (
	"errors"
	"fmt"
	"os"
	"time"

	"golang.org/x/sys/windows"
	"golang.org/x/sys/windows/svc"
	"golang.org/x/sys/windows/svc/mgr"
)

type ServiceInfo struct {
	Installed bool   `json:"installed"`
	Running   bool   `json:"running"`
	Status    string `json:"status"`
}

func serviceInfo() (ServiceInfo, error) {
	manager, err := mgr.Connect()
	if err != nil {
		return ServiceInfo{}, fmt.Errorf("connect to service manager: %w", err)
	}
	defer manager.Disconnect()

	service, err := manager.OpenService(ServiceName)
	if err != nil {
		if errors.Is(err, windows.ERROR_SERVICE_DOES_NOT_EXIST) {
			return ServiceInfo{Status: "Not installed"}, nil
		}
		return ServiceInfo{}, fmt.Errorf("query service: %w", err)
	}
	defer service.Close()

	status, err := service.Query()
	if err != nil {
		return ServiceInfo{}, fmt.Errorf("query service status: %w", err)
	}
	return ServiceInfo{
		Installed: true,
		Running:   status.State == svc.Running,
		Status:    serviceState(status.State),
	}, nil
}

func serviceState(state svc.State) string {
	switch state {
	case svc.Running:
		return "Running"
	case svc.Stopped:
		return "Stopped"
	case svc.StartPending:
		return "Starting"
	case svc.StopPending:
		return "Stopping"
	default:
		return "Unknown"
	}
}

func installService() error {
	executable, err := os.Executable()
	if err != nil {
		return fmt.Errorf("find SyFi executable: %w", err)
	}
	if err := os.MkdirAll(appLogDirectory(), 0o755); err != nil {
		return fmt.Errorf("prepare shared application data: %w", err)
	}

	manager, err := mgr.Connect()
	if err != nil {
		return fmt.Errorf("connect to service manager: %w", err)
	}
	defer manager.Disconnect()

	service, err := manager.CreateService(ServiceName, executable, mgr.Config{
		DisplayName: DisplayName,
		Description: ServiceDescription,
		StartType:   mgr.StartAutomatic,
	}, "--service")
	if err != nil {
		if !errors.Is(err, windows.ERROR_SERVICE_EXISTS) {
			return fmt.Errorf("install service: %w", err)
		}

		service, err = manager.OpenService(ServiceName)
		if err != nil {
			return fmt.Errorf("open existing service: %w", err)
		}

		config, err := service.Config()
		if err != nil {
			service.Close()
			return fmt.Errorf("read existing service configuration: %w", err)
		}
		config.StartType = mgr.StartAutomatic
		if err := service.UpdateConfig(config); err != nil {
			service.Close()
			return fmt.Errorf("enable service at startup: %w", err)
		}
	}
	defer service.Close()

	status, err := service.Query()
	if err != nil {
		return fmt.Errorf("query service before start: %w", err)
	}
	if status.State == svc.Stopped {
		if err := service.Start(); err != nil {
			return fmt.Errorf("start service: %w", err)
		}
	}
	for deadline := time.Now().Add(10 * time.Second); time.Now().Before(deadline); {
		status, err = service.Query()
		if err != nil {
			return fmt.Errorf("query started service: %w", err)
		}
		if status.State == svc.Running {
			return nil
		}
		if status.State == svc.Stopped {
			return errors.New("service stopped during startup")
		}
		time.Sleep(100 * time.Millisecond)
	}
	return errors.New("service did not start in time")
}

func uninstallService() error {
	manager, err := mgr.Connect()
	if err != nil {
		return fmt.Errorf("connect to service manager: %w", err)
	}
	defer manager.Disconnect()

	service, err := manager.OpenService(ServiceName)
	if err != nil {
		if errors.Is(err, windows.ERROR_SERVICE_DOES_NOT_EXIST) {
			return nil
		}
		return fmt.Errorf("open service: %w", err)
	}
	defer service.Close()

	status, err := service.Query()
	if err != nil {
		return fmt.Errorf("query service status: %w", err)
	}
	if status.State != svc.Stopped {
		if _, err := service.Control(svc.Stop); err != nil {
			return fmt.Errorf("stop service: %w", err)
		}
		for deadline := time.Now().Add(5 * time.Second); time.Now().Before(deadline); {
			status, err = service.Query()
			if err != nil {
				return fmt.Errorf("wait for service to stop: %w", err)
			}
			if status.State == svc.Stopped {
				break
			}
			time.Sleep(100 * time.Millisecond)
		}
		if status.State != svc.Stopped {
			return errors.New("service did not stop in time")
		}
	}
	if err := service.Delete(); err != nil {
		return fmt.Errorf("uninstall service: %w", err)
	}
	return nil
}
