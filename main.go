//go:build windows

package main

import (
	"fmt"
	"os"

	"golang.org/x/sys/windows/svc"
)

func main() {
	if shouldRunService(os.Args, nil) {
		if err := svc.Run(ServiceName, &serviceRunner{}); err != nil {
			fmt.Fprintf(os.Stderr, "service error: %v\n", err)
			os.Exit(1)
		}
		return
	}

	isService, err := svc.IsWindowsService()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to determine session type: %v\n", err)
		os.Exit(1)
	}
	if isService {
		if err := svc.Run(ServiceName, &serviceRunner{}); err != nil {
			fmt.Fprintf(os.Stderr, "service error: %v\n", err)
			os.Exit(1)
		}
		return
	}

	runWailsGUI()
}
