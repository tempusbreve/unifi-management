// Package unifi provides tools for interacting with Unifi.
package unifi

import (
	"fmt"
	"time"
)

// Device describes a UniFi network client.
type Device struct {
	ID                  string `json:"_id,omitempty"`
	MAC                 string `json:"mac,omitempty"`
	SiteID              string `json:"site_id,omitempty"`
	OUI                 string `json:"oui,omitempty"`
	NetworkID           string `json:"network_id,omitempty"`
	IP                  string `json:"ip,omitempty"`
	FixedIP             string `json:"fixed_ip,omitempty"`
	Hostname            string `json:"hostname,omitempty"`
	UsergroupID         string `json:"usergroup_id,omitempty"`
	Name                string `json:"name,omitempty"`
	FirstSeen           int64  `json:"first_seen,omitempty"`
	LastSeen            int64  `json:"last_seen,omitempty"`
	DeviceIDOverride    int    `json:"dev_id_override,omitempty"`
	FingerprintOverride bool   `json:"fingerprint_override,omitempty"`
	Blocked             bool   `json:"blocked,omitempty"`
	IsGuest             bool   `json:"is_guest,omitempty"`
	IsWired             bool   `json:"is_wired,omitempty"`
	Noted               bool   `json:"noted,omitempty"`
	UseFixedIP          bool   `json:"use_fixedip,omitempty"`
}

// DisplayName is the friendly name for this device.
func (d Device) DisplayName() string {
	return firstNonEmpty(d.Name, d.Hostname, d.FixedIP, d.IP, d.OUI, d.MAC)
}

func (d Device) String() string {
	name := d.DisplayName()
	ip := firstNonEmpty(d.FixedIP, d.IP)
	last := time.Unix(d.LastSeen, 0)

	blocked := ""
	if d.Blocked {
		blocked = " blocked"
	}

	return fmt.Sprintf("%20s %-16s %-20s (%s)%s", d.MAC, ip, name, last.Format(time.RFC3339), blocked)
}

// Response encapsulates a UniFi http response.
type Response struct {
	Meta struct {
		RC string `json:"rc,omitempty"`
	}
	Data []Device `json:"data,omitempty"`
}

func firstNonEmpty(options ...string) string {
	for _, option := range options {
		if len(option) > 0 {
			return option
		}
	}

	return ""
}
