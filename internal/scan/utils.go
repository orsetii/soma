package scan

import (
	"fmt"
	"net"
)

// ExpandCIDR returns all usable IP addresses in a CIDR block.
// Skips network and broadcast addresses for IPv4.
func ExpandCIDR(cidr string) ([]net.IP, error) {
	ip, ipNet, err := net.ParseCIDR(cidr)
	if err != nil {
		return nil, fmt.Errorf("invalid CIDR: %w", err)
	}

	var ips []net.IP

	// Convert IP to a uint32 for easier math (IPv4 only for now)
	ip4 := ip.To4()
	if ip4 == nil {
		return nil, fmt.Errorf("only IPv4 is supported")
	}

	ipInt := ipToUint32(ip4)
	mask := ipToUint32(net.IP(ipNet.Mask))
	network := ipInt & mask
	broadcast := network | ^mask

	// Iterate from network + 1 to broadcast - 1 (usable IPs)
	for i := network + 1; i < broadcast; i++ {
		ips = append(ips, uint32ToIP(i))
	}

	return ips, nil
}

func ipToUint32(ip net.IP) uint32 {
	ip = ip.To4()
	return uint32(ip[0])<<24 | uint32(ip[1])<<16 | uint32(ip[2])<<8 | uint32(ip[3])
}

func uint32ToIP(n uint32) net.IP {
	return net.IPv4(
		byte(n>>24),
		byte(n>>16),
		byte(n>>8),
		byte(n),
	)
}
