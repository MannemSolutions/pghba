package hba

import (
	"bytes"
	"fmt"
	"math/bits"
	"net"
	"strconv"
	"strings"
)

type AddressType int

const (
	AddressTypeUnknown AddressType = iota
	AddressTypeIpV4
	AddressTypeIpV6
	AddressTypeHostName
	AddressTypeDomain
	AddressTypeAll
	AddressTypeSameHost
	AddressTypeSameNet
)

func (at AddressType) Weight() int {
	if at == AddressTypeUnknown {
		return int(AddressTypeSameNet) - int(AddressTypeUnknown) + 1
	} else {
		return int(AddressTypeSameNet) - int(AddressTypeUnknown)
	}
}

func aTypeFromIP(ip net.IP) AddressType {
	s := ip.String()
	if s == "<nil>" {
		return AddressTypeUnknown
	}
	if strings.Contains(s, ":") {
		return AddressTypeIpV6
	}
	return AddressTypeIpV4
}

func aTypeFromStr(addr string) AddressType {
	switch addr {
	case "all":
		return AddressTypeAll
	case "samehost":
		return AddressTypeSameHost
	case "samenet":
		return AddressTypeSameNet
	}
	if strings.HasPrefix(addr, ".") {
		return AddressTypeDomain
	}
	return AddressTypeHostName
}

type Address struct {
	ip    net.IP
	ipNet *net.IPNet
	str   string
	aType AddressType
}

type Addresses []Address

func NewAddress(addr string) (a Address, err error) {
	a.str = addr
	if strings.Contains(addr, "/") {
		// This is a CIDR (like "192.0.2.0/24" or "2001:db8::/32")
		_, ipNet, err := net.ParseCIDR(addr)
		if err != nil {
			return Address{}, fmt.Errorf("address %s seems like a CIDR, but isn't", addr)
		}
		a.ipNet = ipNet
		a.aType = aTypeFromIP(ipNet.IP)
		return a, nil
	}
	ip := net.ParseIP(addr)
	if ip != nil {
		// Let's store as IP. Maybe netmask follows later
		a.ip = ip
		a.aType = aTypeFromIP(ip)
	} else {
		a.aType = aTypeFromStr(addr)
	}
	return a, nil
}

func sizeFromNet(m net.IPMask) (uint, error) {
	// net.IPMask.Size() seems somewhat broken as in:
	// If the mask is not in the canonical form--ones followed by zeros--then Size returns 0, 0.
	// this function will return the number on leading zeros even when not in canonical form

	// First lets give it a shot for canonical form:
	if lz, _ := m.Size(); lz > 0 {
		return uint(lz), nil
	}
	// That didn't work, which is what this function is actually for
	s := m.String()
	var zeroes uint
	// Note that ipv6 has no netmask. Only a prefix (thank god)...
	for _, part := range strings.Split(s, ".") {
		b, err := strconv.Atoi(part)
		if err != nil {
			// Seems like this is not an IP after all. Weird...
			return 0, fmt.Errorf("this should not happen, but seems like netmask %s is not ipv4 (not int)", s)
		}
		if b > 255 || b < 0 {
			// Seems like this is not an IP after all. Weird...
			return 0, fmt.Errorf("this should not happen, but seems like netmask %s is not ipv4 (not byte)", s)
		}
		zeroes += uint(8 - bits.OnesCount(uint(b)))
	}
	return zeroes, nil
}

func (a Address) NetworkSize() (uint, error) {
	switch a.aType {
	case AddressTypeIpV4, AddressTypeIpV6:
		if a.ipNet.IP.IsUnspecified() {
			if a.aType == AddressTypeIpV4 {
				return 32, nil
			} else {
				return 128, nil
			}
		}
		return sizeFromNet(a.ipNet.Mask)
	default:
		return 0, fmt.Errorf("cannot get network size from address other then ipv4/ipv6")
	}
}

func (a Address) Weight() int {
	switch a.aType {
	case AddressTypeIpV4, AddressTypeIpV6:
		if a.ipNet.IP.IsUnspecified() {
			return 0
		}
		size, err := sizeFromNet(a.ipNet.Mask)
		if err != nil {
			log.Errorf("Could not get weight of address %s, sorting at the end of the file", a.ipNet.String())
			return -1
		}
		return int(size)
	case AddressTypeHostName, AddressTypeSameHost:
		return 0
	case AddressTypeAll:
		return 128
	case AddressTypeDomain, AddressTypeSameNet:
		// Assumption, but this corresponds to a /24 network
		return 32
	default:
		// This should not occur. Sort all the way at the bottom.
		return -1
	}
}

func (a Address) String() string {
	switch a.aType {
	case AddressTypeIpV4, AddressTypeIpV6:
		if a.ipNet.IP.IsUnspecified() {
			size, err := a.NetworkSize()
			if err != nil {
				return a.str
			}
			ip := a.ip.String()
			return fmt.Sprintf("%s/%d", ip, size)
		}
		if _, size := a.ipNet.Mask.Size(); size != 0 {
			// Mask is canonical. We can just return the String() representation of IPNet.
			return a.ipNet.String()
		}
		// Non-canonical mask. We should return as 'IP MASK'...
		return fmt.Sprintf("%s %s", a.ipNet.IP.String(), a.ipNet.Mask.String())
	default:
		return a.str
	}
}

func (a *Address) SetMask(mask string) error {
	if a.aType != AddressTypeIpV4 && a.aType != AddressTypeIpV6 {
		return fmt.Errorf("cannot set mask on something other then ipv4 or ipv6 mask")
	}
	if a.ip.IsUnspecified() {
		return fmt.Errorf("cannot apply mask %s to address that is not ip %s", mask, a.str)
	}
	var size = -1
	// use ParseIP to find ou the bytes
	if mask == "" {
		if a.aType == AddressTypeIpV4 {
			size = 32
		} else {
			size = 128
		}
	} else if i, err := strconv.Atoi(mask); err != nil {
		var maxSize int
		if a.aType == AddressTypeIpV4 {
			maxSize = 32
		} else {
			maxSize = 128
		}
		if size > maxSize {
			return fmt.Errorf("invalid prefix %s (too large)", mask)
		}
		size = i
	}
	if size >= 0 {
		a.ipNet = &net.IPNet{
			IP: a.ip,
			Mask: net.CIDRMask(size, size),
		}
		return nil
	}
	if a.aType == AddressTypeIpV4 && strings.Count(mask, ".") == 4 {
		var parts [4]byte
		for i, part := range strings.Split(mask, ".") {
			b, err := strconv.Atoi(part)
			if err != nil {
				return fmt.Errorf("invalid ipv4 part '%s' in (*Address).SetMask(string), %e", part, err)
			}
			if b < 0 || b > 255 {
				return fmt.Errorf("invalid ipv4 part '%s' in (*Address).SetMask(string), %e", part, err)
			}
			parts[i] = byte(b)
		}
		a.ipNet = &net.IPNet{
			IP: a.ip,
			Mask: net.IPMask{parts[0], parts[1], parts[2], parts[3]},
		}
		return nil
	}
	return fmt.Errorf("mask %s is not a valid netmask", mask)
}

func (a Address) Unset() bool {
	return a.aType == AddressTypeUnknown
}

func (a Address) Clone() Address {
	return Address{
		ip:    a.ip,
		ipNet: a.ipNet,
		str:   a.str,
		aType: a.aType,
	}
}
func (a Address) Compare(other Address) int {
	if a.aType != other.aType {
		return a.aType.Weight() - other.aType.Weight()
	}
	aWeight := a.Weight()
	oWeight := other.Weight()
	if aWeight != oWeight {
		return aWeight - oWeight
	}
	return bytes.Compare(a.ipNet.IP, other.ipNet.IP)
}
