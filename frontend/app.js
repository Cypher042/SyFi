const username = document.getElementById("username");
const password = document.getElementById("password");
const status = document.getElementById("status");
const logs = document.getElementById("logs");
const message = document.getElementById("message");
const serviceStatus = document.getElementById("service-status");
const serviceAction = document.getElementById("service-action");
const themeToggle = document.getElementById("theme-toggle");
const refreshButton = document.getElementById("refresh");
let statusCheckRunning = false;
let serviceCheckRunning = false;
let logCheckRunning = false;

function setTheme(theme) {
  const dark = theme === "dark";
  document.documentElement.classList.toggle("dark", dark);
  themeToggle.setAttribute("aria-label", dark ? "Switch to light mode" : "Switch to dark mode");
  themeToggle.innerHTML = dark
    ? '<svg viewBox="0 0 24 24" aria-hidden="true"><path d="M21 12.8A8.5 8.5 0 1 1 11.2 3 6.5 6.5 0 0 0 21 12.8Z"/></svg>'
    : '<svg viewBox="0 0 24 24" aria-hidden="true"><circle cx="12" cy="12" r="4"/><path d="M12 2v2M12 20v2M4.93 4.93l1.42 1.42M17.65 17.65l1.42 1.42M2 12h2M20 12h2M4.93 19.07l1.42-1.42M17.65 6.35l1.42-1.42"/></svg>';
  localStorage.setItem("syfi-theme", theme);
}

function showMessage(text) {
  message.textContent = text;
}

function serviceErrorMessage(error) {
  const detail = String(error);
  if (detail.toLowerCase().includes("access is denied")) {
    return "Please run SyFi as administrator.";
  }
  return detail;
}

async function refreshDashboard() {
  refreshButton.disabled = true;
  refreshButton.classList.add("is-refreshing");
  await Promise.all([refreshStatus(), refreshService(), refreshLogs()]);
  refreshButton.classList.remove("is-refreshing");
  refreshButton.disabled = false;
}

function setStatusValue(elementId, value, state) {
  const element = document.getElementById(elementId);
  element.textContent = `• ${value}`;
  element.className = state;
}

function updateStatusDashboard(statusText) {
  status.textContent = statusText;
  const values = Object.fromEntries(
    statusText.split("\n").map((line) => {
      const separator = line.indexOf(": ");
      return separator < 0
        ? [line, ""]
        : [line.slice(0, separator), line.slice(separator + 2)];
    }),
  );

  document.getElementById("network-name").textContent = values.SSID || "Unknown";
  const target = values["Target Wi-Fi"] === "true";
  const internet = values["Internet connected"] === "true";
  const portal = values["Captive portal detected"] === "true";
  setStatusValue("target-status", target ? "Matched" : "Not matched", target ? "ok" : "warning");
  setStatusValue("internet-status", internet ? "Connected" : "Disconnected", internet ? "ok" : "warning");
  setStatusValue("portal-status", portal ? "Detected" : "Not detected", portal ? "warning" : "ok");

  const online = internet;
  document.getElementById("connection-title").textContent = online ? "You're online" : "Connection needs attention";
  document.getElementById("connection-subtitle").textContent = online
    ? "Portal login handled automatically"
    : "Checking portal connection";
}

async function refreshStatus() {
  if (statusCheckRunning) {
    return;
  }
  statusCheckRunning = true;
  try {
    updateStatusDashboard(await window.go.main.App.Status());
  } catch (error) {
    showMessage(`Status check failed: ${error}`);
  } finally {
    statusCheckRunning = false;
  }
}

async function loadCredentials() {
  try {
    const credentials = await window.go.main.App.LoadCredentials();
    username.value = credentials.username || "";
    password.value = credentials.password || "";
  } catch (error) {
    showMessage(`Credential load failed: ${error}`);
  }
}

async function refreshService() {
  if (serviceCheckRunning) {
    return;
  }
  serviceCheckRunning = true;
  try {
    const service = await window.go.main.App.ServiceInfo();
    serviceStatus.textContent = service.status;
    serviceStatus.style.color = service.running ? "#45c99d" : "#e2a43b";
    serviceAction.textContent = service.installed ? "Uninstall" : "Install";
  } catch (error) {
    serviceStatus.textContent = "Unavailable";
    showMessage(`Service status failed: ${serviceErrorMessage(error)}`);
  } finally {
    serviceCheckRunning = false;
  }
}

async function refreshLogs() {
  if (logCheckRunning) {
    return;
  }
  logCheckRunning = true;
  try {
    logs.textContent = await window.go.main.App.Logs();
  } catch (error) {
    logs.textContent = `Log read failed: ${error}`;
  } finally {
    logCheckRunning = false;
  }
}

async function waitForService() {
  for (let attempt = 0; attempt < 20; attempt += 1) {
    await refreshService();
    if (serviceStatus.textContent === "Running") {
      return;
    }
    await new Promise((resolve) => setTimeout(resolve, 500));
  }
}

serviceAction.addEventListener("click", async () => {
  serviceAction.disabled = true;
  try {
    const service = await window.go.main.App.ServiceInfo();
    const updated = service.installed
      ? await window.go.main.App.UninstallService()
      : await window.go.main.App.InstallService();
    serviceStatus.textContent = updated.status;
    serviceAction.textContent = updated.installed ? "Uninstall" : "Install";
    if (updated.installed && updated.status !== "Running") {
      await waitForService();
    }
  } catch (error) {
    showMessage(`Service action failed: ${serviceErrorMessage(error)}`);
  } finally {
    serviceAction.disabled = false;
  }
});

document.getElementById("save").addEventListener("click", async () => {
  try {
    await window.go.main.App.SaveCredentials(username.value, password.value);
    showMessage("Credentials saved successfully.");
  } catch (error) {
    showMessage(error.toString());
  }
});

document.getElementById("toggle-password").addEventListener("click", (event) => {
  const showing = password.type === "text";
  password.type = showing ? "password" : "text";
  event.currentTarget.setAttribute("aria-label", showing ? "Show password" : "Hide password");
});

themeToggle.addEventListener("click", () => {
  setTheme(document.documentElement.classList.contains("dark") ? "light" : "dark");
});

refreshButton.addEventListener("click", refreshDashboard);

setTheme(localStorage.getItem("syfi-theme") || "light");
loadCredentials().then(refreshDashboard);
