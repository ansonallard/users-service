package env

import (
	"fmt"
	"os"
	"strings"
)

func GetPort() string {
	port := os.Getenv("PORT")
	if port == "" {
		port = "5000"
	}
	return port
}

func GetHostname() string {
	if IsDevMode() {
		return fmt.Sprintf("http://localhost:%s", GetPort())
	}
	hostname := os.Getenv("HOSTNAME")
	if !IsDevMode() && hostname == "" {
		panic("HOSTNAME required when running in production")
	}
	if !IsDevMode() && !strings.HasPrefix(hostname, "http://") && !strings.HasPrefix(hostname, "https://") {
		return fmt.Sprintf("https://%s", hostname)
	}
	return hostname
}

func IsDevMode() bool {
	dev := os.Getenv("IS_DEV")
	return strings.ToLower(dev) == "true"
}
