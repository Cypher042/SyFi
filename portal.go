package main

import (
	"errors"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"
)

func currentIPv4() string {
	interfaces, err := net.Interfaces()
	if err != nil {
		return "127.0.0.1"
	}
	for _, iface := range interfaces {
		if iface.Flags&net.FlagUp == 0 || iface.Flags&net.FlagLoopback != 0 {
			continue
		}
		addrs, err := iface.Addrs()
		if err != nil {
			continue
		}
		for _, addr := range addrs {
			ipNet, ok := addr.(*net.IPNet)
			if !ok || ipNet.IP == nil || ipNet.IP.To4() == nil {
				continue
			}
			return ipNet.IP.To4().String()
		}
	}
	return "127.0.0.1"
}

func portalLoginURL() string {
	return PortalLoginURL + url.QueryEscape(currentIPv4())
}

func Login(username, password string) error {
	if strings.TrimSpace(username) == "" || strings.TrimSpace(password) == "" {
		return errors.New("username and password are required")
	}

	client := &http.Client{
		Timeout: 15 * time.Second,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			if len(via) >= 10 {
				return errors.New("too many redirects")
			}
			return nil
		},
	}

	form := url.Values{}
	form.Set("username", username)
	form.Set("password", password)

	resp, err := client.PostForm(portalLoginURL(), form)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return fmt.Errorf("login request failed with status %s", resp.Status)
	}

	if !CheckInternetConnectivity() {
		return errors.New("authentication did not restore internet access")
	}
	return nil
}
