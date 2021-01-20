package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/tempusbreve/unifi-management/consul"
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

	cmd := "list"
	if len(os.Args) > 1 {
		cmd = os.Args[1]
	}

	var matches []string
	if len(os.Args) > 2 {
		matches = os.Args[2:]
	}

	var ops []operation

	switch cmd {
	case "kv":
		kv, err := consul.NewKV()
		if err != nil {
			fmt.Fprintf(os.Stderr, "cannot connect to key-value store: %v\n", err)
			os.Exit(1)
		}

		toBlock, err := kv.Get("blocked")
		if err != nil {
			fmt.Fprintf(os.Stderr, "cannot retrieve blocked list: %v\n", err)
			os.Exit(1)
		}

		toUnBlock, err := kv.Get("unblocked")
		if err != nil {
			fmt.Fprintf(os.Stderr, "cannot retrieve unblocked list: %v\n", err)
			os.Exit(1)
		}

		kvOp := "list"
		if len(os.Args) > 2 {
			kvOp = os.Args[2]
		}

		switch kvOp {
		case "block":
			if err = kv.Put("blocked", matches[1:]); err != nil {
				fmt.Fprintf(os.Stderr, "cannot update blocked matches in key-value store: %v\n", err)
				os.Exit(1)
			}

			return
		case "unblock":
			if err = kv.Put("unblocked", matches[1:]); err != nil {
				fmt.Fprintf(os.Stderr, "cannot update unblocked matches in key-value store: %v\n", err)
				os.Exit(1)
			}

			return
		case "sync":
			ops = append(
				ops,
				operation{action: blockDevice(ses), filter: nonGreedy(toBlock)},
				operation{action: unBlockDevice(ses), filter: nonGreedy(toUnBlock)},
			)
		case "list":
			fallthrough
		default:
			fmt.Fprintf(os.Stdout, "KV Config:\n  block %v\n  unblock %v\n", toBlock, toUnBlock)

			return
		}
	case "block":
		ops = append(ops, operation{action: blockDevice(ses), filter: greedy(matches)})
	case "unblock":
		ops = append(ops, operation{action: unBlockDevice(ses), filter: greedy(matches)})
	case "list":
		fallthrough
	default:
		ops = append(ops, operation{action: displayDevice, filter: greedy(matches)})
	}

	if err := apply(ses, ops); err != nil {
		fmt.Fprintf(os.Stderr, "executing %s: %v\n", cmd, err)
		os.Exit(1)
	}
}

type operation struct {
	action func(unifi.Device) error
	filter func(unifi.Device) bool
}

func apply(s *unifi.Session, ops []operation) error {
	devices, err := s.ListDevices()
	if err != nil {
		return err
	}

	for _, op := range ops {
		for _, device := range devices {
			if op.filter(device) {
				op.action(device)
			}
		}
	}

	return nil
}

func unBlockDevice(s *unifi.Session) func(unifi.Device) error { return fnDevice(s.Unblock) }
func blockDevice(s *unifi.Session) func(unifi.Device) error   { return fnDevice(s.Block) }

func fnDevice(fn func(string) ([]unifi.Device, error)) func(unifi.Device) error {
	return func(d unifi.Device) error {
		res, err := fn(d.MAC)
		if err != nil {
			return err
		}

		fmt.Fprintf(os.Stdout, "%v:\n", d.DisplayName())
		for _, dev := range res {
			fmt.Fprintf(os.Stdout, "  %v\n", dev)
		}

		return nil
	}
}

func displayDevice(d unifi.Device) error {
	fmt.Fprintf(os.Stdout, "%v\n", d)
	return nil
}

func greedy(m []string) func(unifi.Device) bool {
	return func(d unifi.Device) bool {
		return matchHelper(true, d, m)
	}
}

func nonGreedy(m []string) func(unifi.Device) bool {
	return func(d unifi.Device) bool {
		return matchHelper(false, d, m)
	}
}

func matchHelper(greedy bool, device unifi.Device, matches []string) bool {
	if len(matches) == 0 {
		return greedy
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
