# SyFi

SyFi is a Windows desktop GUI for automatically logging in to the ISM Campus
Wi-Fi captive portal. It can also install a Windows service that periodically
checks the connection and performs the login when required.

SyFi is supported on **Windows only** and is intended to be used through its
desktop GUI. There is no supported Linux, macOS, or command-line mode.

## Requirements

- Windows 10 or later
- Go 1.25 or later
- A network connection to the target Wi-Fi
- Microsoft Edge WebView2 Runtime

The Windows service requires administrator privileges to install, uninstall,
or start it.

## Build the Windows executable

Open PowerShell in the repository directory:

```powershell
go mod download
go build -tags desktop,production `
  -ldflags "-w -s -H windowsgui" `
  -o syfi-gui.exe .
```

The executable is created at:

```text
.\syfi-gui.exe
```

The repository includes the SyFi icon as `assets\syfi.ico`. The accompanying
`syfi.syso` Windows resource file embeds that icon automatically whenever Go
builds the package, so no additional icon arguments are needed.

The `-H windowsgui` linker option prevents a console window from appearing
when the GUI starts. The `-w -s` options reduce the executable size.

To build without stripping debug information:

```powershell
go build -tags desktop,production -o syfi-gui.exe .
```

## Build for a specific Windows architecture

For 64-bit Windows:

```powershell
$env:GOOS = "windows"
$env:GOARCH = "amd64"
go build -tags desktop,production `
  -ldflags "-w -s -H windowsgui" `
  -o syfi-gui.exe .
```

For Windows on ARM64:

```powershell
$env:GOOS = "windows"
$env:GOARCH = "arm64"
go build -tags desktop,production `
  -ldflags "-w -s -H windowsgui" `
  -o syfi-gui-arm64.exe .
```

The generated executables use the SyFi icon in Windows Explorer, the taskbar,
and the application window.

Unset the variables after cross-compiling if you will build again in the same
PowerShell session:

```powershell
Remove-Item Env:GOOS
Remove-Item Env:GOARCH
```

## Run SyFi

Start the GUI normally:

```powershell
.\syfi-gui.exe
```

Enter the portal username and password, then select **Save Credentials**.
SyFi displays the current SSID, target-network status, internet status, and
recent service logs in the GUI.

## Install the Windows service

1. Start PowerShell or the executable as **Administrator**.
2. Run `.\syfi-gui.exe`.
3. Select **Install Service** in the GUI.
4. Leave the executable at the same path after installation.

The service is configured to start automatically with Windows. It checks the
target Wi-Fi connection every 20 seconds and attempts a captive-portal login
when the internet is unavailable and a portal is detected.

To remove the service:

1. Run the GUI as **Administrator**.
2. Select **Uninstall Service**.

The GUI waits for the service to stop before removing it.

## Data and logs

Shared application data is stored in:

```text
C:\ProgramData\CollegeWiFiAutoLogin
```

This directory contains:

- `credentials.json` — saved portal credentials
- `syfi.log` — Windows service logs

Keep the executable at the same path after installing the service. Windows
stores that executable path in the service configuration.

## Test the project

Run the Go tests with:

```powershell
go test ./...
```

Check that the Windows build compiles:

```powershell
$env:GOOS = "windows"
$env:GOARCH = "amd64"
go test ./...
go build -tags desktop,production `
  -ldflags "-w -s -H windowsgui" `
  -o syfi-gui.exe .
Remove-Item Env:GOOS
Remove-Item Env:GOARCH
```

## Clean local build output

Build artifacts are ignored by Git. To remove the local GUI executable:

```powershell
Remove-Item .\syfi-gui.exe -ErrorAction SilentlyContinue
```

Do not delete `C:\ProgramData\CollegeWiFiAutoLogin` unless you also want to
remove the saved credentials and service logs.
