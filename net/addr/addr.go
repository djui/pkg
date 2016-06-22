package addr

import (
	"net"
	"strings"
)

// Partition takes a list of network addresses and partitions them into v4 and
// v6 addresses.
func Partition(addrs []net.Addr) ([]net.Addr, []net.Addr) {
	var v4Addrs []net.Addr
	var v6Addrs []net.Addr

	for _, addr := range addrs {
		if IsIPv4(addr) {
			v4Addrs = append(v4Addrs, addr)
		} else {
			v6Addrs = append(v6Addrs, addr)
		}
	}

	return v4Addrs, v6Addrs
}

// IsIPv4 takes a network address and checks if the address is an IPv4 address.
func IsIPv4(addr net.Addr) bool {
	return strings.Contains(addr.String(), ".")
}

// IPs takes a list of network addresses and returns its list of IPs.
func IPs(addrs []net.Addr) ([]net.IP, error) {
	var IPs []net.IP

	for _, addr := range addrs {
		IP, _, err := net.ParseCIDR(addr.String())
		if err != nil {
			return nil, err
		}
		IPs = append(IPs, IP)
	}

	return IPs, nil
}
