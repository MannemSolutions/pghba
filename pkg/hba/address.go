package hba

import (
	"fmt"
	"net"
	"strings"
)

type Address struct {
	ip net.IP
	ipNet *net.IPNet
	str string
}

func NewAddress(addr string) (a Address, err error) {
	a.str = addr
	if strings.Contains(addr, "/") {
		// This is a CIDR (like "192.0.2.0/24" or "2001:db8::/32")
		_ , ipnet, err := net.ParseCIDR(addr)
		if err != nil {
			return Address{}, fmt.Errorf("address %s seems like a CIDR, but isn't", addr)
		}
		a.ipNet = ipnet
		return a, nil
	}
	ip := net.ParseIP(addr)
	if ip != nil {
		// Let's store as IP. Maybe netmask follows later
		a.ip = ip
	}
	return a, nil
}

func (a Address) String() string {
	if ! a.ipNet.IP.IsUnspecified() {
		if _, bits := a.ipNet.Mask.Size(); bits != 0 {
			// Mask is canonical. We can jus return the String() representation of IPNet.
			return a.ipNet.String()
		}
		// Weird mask. We should return as 'IP MASK'...
		return fmt.Sprintf("%s %s", a.ipNet.IP.String(), a.ipNet.Mask.String())
	}
		ip := a.ip.String()
		if ip == "<nil>" {
			return a.str
		}
		if strings.Contains(ip, ":") {
			return fmt.Sprintf("%s/%d", ip, 128)
		}
		return fmt.Sprintf("%s/%d", ip, 32)
}

func (a *Address) SetMask(mask string) error {
	// use ParseIP to find ou the bytes
	if ! a.ip.IsUnspecified() {
		return fmt.Errorf("cannot apply mask %s to address that is not ip %s", mask, a.str)
	}
	m := net.ParseIP(mask)
	if m == nil {
		return fmt.Errorf("%s is not a valid mask", mask)
	}
	// This is a valid mask representation
	mask = m.String()
	ip :=  a.ip.String()
	if (strings.Contains(mask, ":") && strings.Contains(ip, ":")) ||
		(strings.Contains(mask, ".") && strings.Contains(ip, ".")) {
		a.str = fmt.Sprintf("%s %s", ip, mask)
		a.ip = nil
		return nil
	}
	return fmt.Errorf("ip %s and mask %s are not same version (one is ipv4 and other is ipv6)", ip, mask)
}