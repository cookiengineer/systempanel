package nmconnection

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type Connection struct {
	Connection    ConnectionSection
	Wifi          WifiSection
	WifiSecurity  WifiSecuritySection
	IPv4          IPv4Section
	IPv6          IPv6Section
	Path          string
}

type ConnectionSection struct {
	ID            string
	UUID          string
	Type          string
	Autoconnect   bool
	InterfaceName string
	Timestamp     string
}

type WifiSection struct {
	Mode            string
	SSID            string
	Hidden          bool
	Band            string
	BSSID           string
	Channel         int
	ClonedMAC       string
	MTU             int
}

type WifiSecuritySection struct {
	KeyMgmt string
	PSK     string
	WEPKey0 string
	WEPKey1 string
	WEPKey2 string
	WEPKey3 string
	WEPTxKeyidx string
	AuthAlg  string
}

type IPv4Section struct {
	Method      string
	Address1    string
	DNS         string
	DNSOptions  string
	DNSSearch   string
	Gateway     string
	NeverDefault bool
	RouteMetric int
	DHCPClientID string
}

type IPv6Section struct {
	Method       string
	Address1     string
	DNS          string
	DNSSearch    string
	Gateway      string
	NeverDefault  bool
	RouteMetric  int
	AddrGenMode  string
	IP6Privacy   string
}

func Parse(path string) (*Connection, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	c := &Connection{Path: path}
	c.Connection.Autoconnect = true
	c.IPv4.Method = "auto"
	c.IPv6.Method = "auto"
	c.Wifi.Mode = "infrastructure"

	scanner := bufio.NewScanner(f)
	section := ""
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		if strings.HasPrefix(line, "[") && strings.HasSuffix(line, "]") {
			section = line[1 : len(line)-1]
			continue
		}
		eq := strings.Index(line, "=")
		if eq < 0 {
			continue
		}
		key := strings.TrimSpace(line[:eq])
		val := strings.TrimSpace(line[eq+1:])
		switch section {
		case "connection":
			c.parseConnection(key, val)
		case "wifi":
			c.parseWifi(key, val)
		case "wifi-security":
			c.parseWifiSecurity(key, val)
		case "ipv4":
			c.parseIPv4(key, val)
		case "ipv6":
			c.parseIPv6(key, val)
		}
	}

	if c.Connection.ID == "" && c.Wifi.SSID != "" {
		c.Connection.ID = c.Wifi.SSID
	}
	return c, scanner.Err()
}

func (c *Connection) parseConnection(key, val string) {
	switch key {
	case "id": c.Connection.ID = val
	case "uuid": c.Connection.UUID = val
	case "type": c.Connection.Type = val
	case "autoconnect": c.Connection.Autoconnect = val == "true"
	case "interface-name": c.Connection.InterfaceName = val
	case "timestamp": c.Connection.Timestamp = val
	}
}

func (c *Connection) parseWifi(key, val string) {
	switch key {
	case "mode": c.Wifi.Mode = val
	case "ssid": c.Wifi.SSID = val
	case "hidden": c.Wifi.Hidden = val == "true"
	case "band": c.Wifi.Band = val
	case "bssid": c.Wifi.BSSID = val
	case "channel": fmt.Sscanf(val, "%d", &c.Wifi.Channel)
	case "cloned-mac-address": c.Wifi.ClonedMAC = val
	case "mtu": fmt.Sscanf(val, "%d", &c.Wifi.MTU)
	}
}

func (c *Connection) parseWifiSecurity(key, val string) {
	switch key {
	case "key-mgmt": c.WifiSecurity.KeyMgmt = val
	case "psk": c.WifiSecurity.PSK = val
	case "wep-key0": c.WifiSecurity.WEPKey0 = val
	case "wep-key1": c.WifiSecurity.WEPKey1 = val
	case "wep-key2": c.WifiSecurity.WEPKey2 = val
	case "wep-key3": c.WifiSecurity.WEPKey3 = val
	case "wep-tx-keyidx": c.WifiSecurity.WEPTxKeyidx = val
	case "auth-alg": c.WifiSecurity.AuthAlg = val
	}
}

func (c *Connection) parseIPv4(key, val string) {
	switch key {
	case "method": c.IPv4.Method = val
	case "address1": c.IPv4.Address1 = val
	case "dns": c.IPv4.DNS = val
	case "dns-options": c.IPv4.DNSOptions = val
	case "dns-search": c.IPv4.DNSSearch = val
	case "gateway": c.IPv4.Gateway = val
	case "never-default": c.IPv4.NeverDefault = val == "true"
	case "route-metric": fmt.Sscanf(val, "%d", &c.IPv4.RouteMetric)
	case "dhcp-client-id": c.IPv4.DHCPClientID = val
	}
}

func (c *Connection) parseIPv6(key, val string) {
	switch key {
	case "method": c.IPv6.Method = val
	case "address1": c.IPv6.Address1 = val
	case "dns": c.IPv6.DNS = val
	case "dns-search": c.IPv6.DNSSearch = val
	case "gateway": c.IPv6.Gateway = val
	case "never-default": c.IPv6.NeverDefault = val == "true"
	case "route-metric": fmt.Sscanf(val, "%d", &c.IPv6.RouteMetric)
	case "addr-gen-mode": c.IPv6.AddrGenMode = val
	case "ip6-privacy": c.IPv6.IP6Privacy = val
	}
}

func (c *Connection) Serialize() string {
	var b strings.Builder

	b.WriteString("[connection]\n")
	writeKey(&b, "id", c.Connection.ID)
	writeKey(&b, "uuid", c.Connection.UUID)
	writeKey(&b, "type", c.Connection.Type)
	writeKeyBool(&b, "autoconnect", c.Connection.Autoconnect)
	writeKey(&b, "interface-name", c.Connection.InterfaceName)
	b.WriteString("\n")

	if c.Connection.Type == "wifi" {
		b.WriteString("[wifi]\n")
		writeKey(&b, "mode", c.Wifi.Mode)
		writeKey(&b, "ssid", c.Wifi.SSID)
		writeKeyBool(&b, "hidden", c.Wifi.Hidden)
		writeKey(&b, "band", c.Wifi.Band)
		writeKey(&b, "bssid", c.Wifi.BSSID)
		if c.Wifi.Channel > 0 {
			writeKeyInt(&b, "channel", c.Wifi.Channel)
		}
		writeKey(&b, "cloned-mac-address", c.Wifi.ClonedMAC)
		if c.Wifi.MTU > 0 {
			writeKeyInt(&b, "mtu", c.Wifi.MTU)
		}
		b.WriteString("\n")

		if c.WifiSecurity.KeyMgmt != "" {
			b.WriteString("[wifi-security]\n")
			writeKey(&b, "key-mgmt", c.WifiSecurity.KeyMgmt)
			writeKey(&b, "psk", c.WifiSecurity.PSK)
			writeKey(&b, "wep-key0", c.WifiSecurity.WEPKey0)
			writeKey(&b, "wep-key1", c.WifiSecurity.WEPKey1)
			writeKey(&b, "wep-key2", c.WifiSecurity.WEPKey2)
			writeKey(&b, "wep-key3", c.WifiSecurity.WEPKey3)
			writeKey(&b, "wep-tx-keyidx", c.WifiSecurity.WEPTxKeyidx)
			writeKey(&b, "auth-alg", c.WifiSecurity.AuthAlg)
			b.WriteString("\n")
		}
	}

	b.WriteString("[ipv4]\n")
	writeKey(&b, "method", c.IPv4.Method)
	writeKey(&b, "address1", c.IPv4.Address1)
	writeKey(&b, "dns", c.IPv4.DNS)
	writeKey(&b, "dns-search", c.IPv4.DNSSearch)
	writeKey(&b, "gateway", c.IPv4.Gateway)
	writeKeyBool(&b, "never-default", c.IPv4.NeverDefault)
	if c.IPv4.RouteMetric > 0 {
		writeKeyInt(&b, "route-metric", c.IPv4.RouteMetric)
	}
	b.WriteString("\n")

	b.WriteString("[ipv6]\n")
	writeKey(&b, "method", c.IPv6.Method)
	writeKey(&b, "address1", c.IPv6.Address1)
	writeKey(&b, "dns", c.IPv6.DNS)
	writeKey(&b, "dns-search", c.IPv6.DNSSearch)
	writeKey(&b, "gateway", c.IPv6.Gateway)
	writeKeyBool(&b, "never-default", c.IPv6.NeverDefault)
	if c.IPv6.RouteMetric > 0 {
		writeKeyInt(&b, "route-metric", c.IPv6.RouteMetric)
	}
	writeKey(&b, "addr-gen-mode", c.IPv6.AddrGenMode)

	return b.String()
}

func NewWiFiConnection(ssid string) *Connection {
	return &Connection{
		Connection: ConnectionSection{
			ID:           ssid,
			Type:         "wifi",
			Autoconnect:  true,
		},
		Wifi: WifiSection{
			Mode: "infrastructure",
			SSID: ssid,
		},
		WifiSecurity: WifiSecuritySection{
			KeyMgmt: "wpa-psk",
		},
		IPv4: IPv4Section{
			Method: "auto",
		},
		IPv6: IPv6Section{
			Method: "auto",
		},
	}
}

func NewEthernetConnection(name string) *Connection {
	return &Connection{
		Connection: ConnectionSection{
			ID:          name,
			Type:        "802-3-ethernet",
			Autoconnect: true,
		},
		IPv4: IPv4Section{
			Method: "auto",
		},
		IPv6: IPv6Section{
			Method: "auto",
		},
	}
}

func writeKey(b *strings.Builder, key, val string) {
	if val != "" {
		b.WriteString(key)
		b.WriteString("=")
		b.WriteString(val)
		b.WriteString("\n")
	}
}

func writeKeyBool(b *strings.Builder, key string, val bool) {
	if val {
		b.WriteString(key + "=true\n")
	} else {
		b.WriteString(key + "=false\n")
	}
}

func writeKeyInt(b *strings.Builder, key string, val int) {
	b.WriteString(key)
	b.WriteString("=")
	b.WriteString(fmt.Sprintf("%d", val))
	b.WriteString("\n")
}
