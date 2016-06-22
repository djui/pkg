package nif

import "net"

// IsUsable filters network interface that are usable, with some definition of
// usable.
func IsUsable(nif net.Interface) bool {
	// Filter out interfaces that have no hardware address.
	if len(nif.HardwareAddr) == 0 {
		return false
	}
	// Filter out loopback interfaces.
	if nif.Flags&net.FlagLoopback == net.FlagLoopback {
		return false
	}
	// Filter out Point-to-Point interfaces.
	if nif.Flags&net.FlagPointToPoint == net.FlagPointToPoint {
		return false
	}
	// Filter out interfaces that are not up.
	if nif.Flags&net.FlagUp != net.FlagUp {
		return false
	}

	// If the interface made it that far, assume it's of interest.
	return true
}
