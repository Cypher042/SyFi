package main

const (
	TargetSSID           = "ISM-Campus-Wifi"
	ServiceName          = "CollegeWiFiAutoLogin"
	DisplayName          = "College Wi‑Fi Auto Login"
	ServiceDescription   = "Automatically authenticates to the college Wi‑Fi captive portal."
	LogFileName          = "collegewifi.log"
	ConnectivityCheckURL = "https://detectportal.firefox.com/success.txt"
	PortalBaseURL        = "https://netaccess.iitism.ac.in:6082"
	PortalLoginURL       = "https://netaccess.iitism.ac.in:6082/?src="
	PortalCheckURL       = "https://netaccess.iitism.ac.in:6082"
	MonitorInterval      = 20 * 1000000000 // 20 seconds
	RetryBackoffBase     = 10 * 1000000000 // 10 seconds
	RetryBackoffMax      = 60 * 1000000000 // 60 seconds
)

func appLogDirectory() string {
	return userLocalAppDataDir()
}
