package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/tempusbreve/unifi-management/unifi"
)

func main() {
	endpoint := os.Getenv("UNIFI_ENDPOINT")
	username := os.Getenv("UNIFI_USERNAME")
	password := os.Getenv("UNIFI_PASSWORD")

	ses, err := unifi.NewSession(
		unifi.Endpoint(endpoint),
		unifi.Credentials(username, password),
	)
	if err != nil {
		fmt.Fprintf(os.Stderr, "cannot create session: %v\n", err)
		os.Exit(1)
	}

	if _, err = ses.Login(); err != nil {
		fmt.Fprintf(os.Stderr, "cannot authenticate session: %v\n", err)
		os.Exit(1)
	}

	devices, err := ses.ListDevices()
	if err != nil {
		fmt.Fprintf(os.Stderr, "cannot retrieve devices: %v\n", err)
		os.Exit(1)
	}

	cmd := "list"
	if len(os.Args) > 1 {
		cmd = os.Args[1]
	}

	var matches []string
	if len(os.Args) > 2 {
		matches = os.Args[2:]
	}

	var action func(unifi.Device)

	switch cmd {
	case "block":
		action = blockDevice(ses)
	case "unblock":
		action = unBlockDevice(ses)
	case "list":
		action = displayDevice
	default:
		action = displayDevice
	}

	for _, device := range devices {
		if match(device, matches) {
			action(device)
		}
	}
}

func unBlockDevice(s *unifi.Session) func(unifi.Device) { return fnDevice(s.Unblock) }
func blockDevice(s *unifi.Session) func(unifi.Device)   { return fnDevice(s.Block) }

func fnDevice(fn func(string) ([]unifi.Device, error)) func(unifi.Device) {
	return func(d unifi.Device) {
		res, err := fn(d.MAC)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%v error; %v\n", d.DisplayName(), err)

			return
		}

		fmt.Fprintf(os.Stdout, "%v:\n", d.DisplayName())
		for _, dev := range res {
			fmt.Fprintf(os.Stdout, "  %v\n", dev)
		}
	}
}

func displayDevice(d unifi.Device) {
	fmt.Fprintf(os.Stdout, "%v\n", d)
}

func match(device unifi.Device, matches []string) bool {
	if len(matches) == 0 {
		return true
	}

	for _, match := range matches {
		if strings.Contains(device.Name, match) {
			return true
		}

		if strings.Contains(device.Hostname, match) {
			return true
		}
	}

	return false
}
