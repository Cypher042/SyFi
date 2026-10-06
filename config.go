package main

const (
	TargetSSID           = "ISM-Campus-Wi-Fi"
	ServiceName          = "SyFi"
	DisplayName          = "SyFi"
	ServiceDescription   = "Automatically authenticates to the college Wi‑Fi captive portal."
	LogFileName          = "syfi.log"
	ConnectivityCheckURL = "https://google.com/generate_204"
	PortalLoginURL       = "https://netaccess.iitism.ac.in:6082/?src="
	PortalCheckURL       = "https://netaccess.iitism.ac.in:6082"
	MonitorInterval      = 20 * 1000000000 // 20 seconds
)

func appLogDirectory() string {
	return userLocalAppDataDir()
}
