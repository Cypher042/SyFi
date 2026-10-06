package main

import (
	"os/exec"
	"runtime"
	"strings"
)

func IsTargetSSID() (bool, error) {
	if runtime.GOOS != "windows" {
		return false, nil
	}
	ssid, err := CurrentSSID()
	if err != nil {
		return false, err
	}
	return strings.EqualFold(ssid, TargetSSID), nil
}

func CurrentSSID() (string, error) {
	cmd := exec.Command("netsh", "wlan", "show", "interfaces")
	hideCommandWindow(cmd)
	out, err := cmd.Output()
	if err != nil {
		return "", err
	}
	text := string(out)
	for _, line := range strings.Split(text, "\n") {
		if strings.Contains(line, "SSID") && strings.Contains(line, ":") {
			parts := strings.SplitN(line, ":", 2)
			if len(parts) == 2 {
				return strings.TrimSpace(parts[1]), nil
			}
		}
	}
	return "", nil
}
