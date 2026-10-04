//go:build windows

package main

func saveCredentials(username, password string) error {
	return saveSimpleCredentials(username, password)
}

func loadCredentials() (string, string, error) {
	return loadSimpleCredentials()
}
