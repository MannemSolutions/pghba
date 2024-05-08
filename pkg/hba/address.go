package hba

import (
	"bytes"
	"fmt"
	"math/bits"
	"net"
	"strconv"
	"strings"
)

// AddressType helps discern between the different concepts of 'Address' as
// pg_hba.conf defines them.
// See https://www.postgresql.org/docs/current/auth-pg-hba-conf.html
type AddressType int

// Define all AddressTypes that pg_hba.conf specifies and the AddressTypeUnknown.
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

// Return the relative Weight of the AddressType for a rule by comparing it arithmatically against the lowest weight
// AddressType.
func (at AddressType) Weight() int {
	if at == AddressTypeUnknown {
		return int(AddressTypeSameNet) - int(AddressTypeUnknown) + 1
	} else {
		return int(at) - int(AddressTypeUnknown)
	}
}

// Determine and return the IP address family type (AddressTypeIPV4 or AddressTypeIPV6) or AddressTypeUknown when nil.
// TODO: I believe this should be replaced with the function from the native net/netIP package
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

// Determine and return the pg_hba AddressType from the string 'addr'.
// This function implicitly assumes that 'addr' is not an IP type address.
// TODO: This determination may be too naive for all real-world scenarios and additionally, could be redundant because:
// 'all' is the same as 0.0.0.0/0
// 'samehost' is the same as w.x.y.z/32
// 'samenet' is (roughly) the same as host/netmask for referenced network.
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

// Address holds a string with the 'raw' address, an IP or IP network derived
// from it (optionally) and the address type accordig to pg_hba standards.
type Address struct {
	ip    net.IP
	ipNet *net.IPNet
	str   string
	aType AddressType
}

// A slice of Addresses
type Addresses []Address

// NewAddress returns a pg_hba Address from a string 'addr'
// 1) regexes are not allowed for the address field so anything with a '/' should be an IP with a CIDR netmask
// 2) So parse it and if it returns anything other than 'nil', assume it is an IP
// 3) Else parse it as a label, host or domain.
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

// TODO I'm uncertain as to why you would want to accept subnetmasks that have leading zeros
// This function returns the length of the significant part of a dotted quad net mask.
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

// NetworkSize returns the length of the bitmask of Address 'a' or the host mask for the associated IP type if none is
// specified. Assumes 'a' is an IPv4 or IPv6 address. Returns an error if that assumption fails.
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

// Return the 'Weight' for the Address 'a' based on the netmask length.
// The longer the netmask, the more specific an address is, the higher the weight.
// TODO this is a different type of 'weight' than that for the AddressType. I believe a better term would be
// Specificity.
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

// Return a string representation of Address 'a'
// TODO unclear to me why you would want to have non-canonical addresses and their mask separately for our purposes,
// and in a different format to boot.
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

// The SetMask method on Address sets ipNet on 'a'
func (a *Address) SetMask(mask string) error {
	if a.aType != AddressTypeIpV4 && a.aType != AddressTypeIpV6 {
		// TODO shouldn't we simply reject anything but a valid address type instead of continuing if mask is unset?
		// (In other words: can this be applied to an Address that isn't an IP?)
		if mask == "" {
			return nil
		}
		return fmt.Errorf("cannot set mask on something other then ipv4 or ipv6 address")
	}
	if a.ip.IsUnspecified() {
		return fmt.Errorf("cannot apply mask %s to address that is not ip %s", mask, a.str)
	}
	var size = -1
	// If no mask is defined, use default value for full host mask for each IP type.
	if mask == "" {
		if a.aType == AddressTypeIpV4 {
			size = 32
		} else {
			size = 128
		}
		// If 'mask' was specified, make sure the mask size is smaller than the maximum allowed for each IP type.
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
			// TODO is this correct? Shouldn't the second value be 'maxSize'?
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
			IP:   a.ip,
			Mask: net.IPMask{parts[0], parts[1], parts[2], parts[3]},
		}
		return nil
	}
	return fmt.Errorf("mask %s is not a valid netmask", mask)
}

// Returns true if Address 'a' has an AddressTypeUnknown
func (a Address) Unset() bool {
	return a.aType == AddressTypeUnknown
}

// Return a clone of Address 'a'
func (a Address) Clone() Address {
	return Address{
		ip:    a.ip,
		ipNet: a.ipNet,
		str:   a.str,
		aType: a.aType,
	}
}

// Returns true if the 'other' address is contained within the 'a' Address specification
func (a Address) Contains(other Address) bool {
	switch a.aType {
	// if 'a' is AddressTypeUnkown, return true if 'other' is also AddressType unknown. Otherwise return false.
	case AddressTypeUnknown:
		return a.aType == other.aType
		// if both are the same IP type, compare logically based on netmask.
	case AddressTypeIpV4, AddressTypeIpV6:
		if a.aType == other.aType {
			return a.ipNet.Contains(other.ipNet.IP)
		}
	case AddressTypeSameHost:
		return a.aType == other.aType
	case AddressTypeSameNet:
		return other.aType == AddressTypeSameHost || other.aType == AddressTypeSameNet
	case AddressTypeHostName:
		// with delete command, hostname can be set to "" which means all hosts
		return a.str == "" || a.aType == other.aType && a.str == other.str
	case AddressTypeDomain:
		return (other.aType == AddressTypeDomain || other.aType == AddressTypeHostName) && strings.HasSuffix(other.str, a.str)
		// TODO isn't everything contained within 'all'?
	case AddressTypeAll:
		return a.aType == other.aType
	}
	return false
}

// Return the relative Weight of one addressess' type against another.
// TODO this naming is confusing to me. It is unclear that the address type is being compared here in terms of weight
// as meant in AddressTypeWeight.
func (a Address) Compare(other Address) int {
	if a.aType == AddressTypeUnknown || other.aType == AddressTypeUnknown {
		// For delete command. If it is an unknown address type, make it equal to all addresses
		return 0
	}
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
