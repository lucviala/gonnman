package connman

import (
	"fmt"

	"github.com/godbus/dbus"
)

type IPv4Config struct {
	Method  string `json:"method,omitempty"`
	Address string `json:"address,omitempty"`
	Netmask string `json:"netmask,omitempty"`
	Gateway string `json:"gateway,omitempty"`
}

type IPv6Config struct {
	Method       string `json:"method,omitempty"`
	Address      string `json:"address,omitempty"`
	PrefixLength uint8  `json:"prefix_length"`
	Gateway      string `json:"gateway,omitempty"`
	Privacy      string `json:"privacy,omitempty"`
}

type EthConfig struct {
	Method    string `json:"method,omitempty"`
	Interface string `json:"interface,omitempty"`
	Address   string `json:"address,omitempty"`
	MTU       uint16 `json:"mtu,omitempty"`
}

type ProxyConfig struct {
	Method   string   `json:"method,omitempty"`
	URL      string   `json:"url,omitempty"`
	Servers  []string `json:"servers,omitempty"`
	Excludes []string `json:"excludes,omitempty"`
}

type Provider struct {
	Host   string `json:"host,omitempty"`
	Domain string `json:"domain,omitempty"`
	Name   string `json:"name,omitempty"`
	Type   string `json:"type,omitempty"`
}

type Service struct {
	Path        dbus.ObjectPath `json:"path,omitempty"`
	Name        string          `json:"name,omitempty"`
	Type        string          `json:"type,omitempty"`
	State       string          `json:"state,omitempty"`
	Error       string          `json:"error,omitempty"`
	Security    []string        `json:"security,omitempty"`
	Strength    uint8           `json:"strength,omitempty"`
	Favorite    bool            `json:"favorite,omitempty"`
	AutoConnect bool            `json:"autoconnect,omitempty"`
	Immutable   bool            `json:"immutable,omitempty"`
	Roaming     bool            `json:"roaming,omitempty"`

	Ethernet           EthConfig   `json:"ethernet,omitempty"`
	IPv4               IPv4Config  `json:"ipv4,omitempty"`
	IPv4Configuration  IPv4Config  `json:"ipv4_configuration,omitempty"`
	IPv6               IPv6Config  `json:"ipv6,omitempty"`
	IPv6Configuration  IPv6Config  `json:"ipv6_configuration,omitempty"`
	Proxy              ProxyConfig `json:"proxy,omitempty"`
	ProxyConfiguration ProxyConfig `json:"proxy_configuration,omitempty"`
	Provider           Provider    `json:"provider,omitempty"`

	Domains                  []string `json:"domains,omitempty"`
	DomainsConfiguration     []string `json:"domains_configuration,omitempty"`
	Nameservers              []string `json:"nameservers,omitempty"`
	NameserversConfiguration []string `json:"nameservers_configuration,omitempty"`
	Timeservers              []string `json:"timeservers,omitempty"`
	TimeserversConfiguration []string `json:"timeservers_configuration,omitempty"`
}

func (s *Service) Connect(psk string) error {
	db, err := DBusService(s.Path)
	if err != nil {
		return err
	}

	secure := false
	for _, s := range s.Security {
		if s == "psk" || s == "wep" {
			secure = true
			break
		}
	}

	if !secure {
		_, err = db.Call("Connect")
		return err
	}

	ag := NewAgent(psk)
	if ag == nil {
		return fmt.Errorf("Could not spawn a new agent\n")
	}

	if err := RegisterAgent(ag); err != nil {
		return err
	}
	defer func() {
		UnregisterAgent(ag)
		ag.Destroy()
	}()

	_, err = db.Call("Connect")
	return err
}

func (s *Service) Disconnect() error {
	db, err := DBusService(s.Path)
	if err != nil {
		return err
	}

	_, err = db.Call("Disconnect")
	return err
}

func (s *Service) ApplyIP() error {
	db, err := DBusService(s.Path)
	if err != nil {
		return err
	}

	arg, err := structToDict(s.IPv4Configuration)
	if err != nil {
		return err
	}

	return db.Set("IPv4.Configuration", arg)
}

func (s *Service) ApplyDNS() error {
	db, err := DBusService(s.Path)
	if err != nil {
		return err
	}

	return db.Set("Nameservers.Configuration", s.NameserversConfiguration)
}
