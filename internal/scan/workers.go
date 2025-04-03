package scan

import (
	"net"
	"os/exec"
	"soma/internal/device"
	"strings"
	"sync"

	"github.com/jedib0t/go-pretty/v6/progress"
)

// startWorkers initializes worker goroutines for scanning.
func startWorkers(ipChan chan string, resultChan chan device.DeviceInfo, wg *sync.WaitGroup, methods DiscoveryMethods, tracker *progress.Tracker) {
	numWorkers := 4
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for ip := range ipChan {
				result := scanIP(ip, methods)
				resultChan <- result
				tracker.Increment(1)
			}
		}()
	}
}

// scanIP performs the scanning for a single IP address.
func scanIP(ip string, methods DiscoveryMethods) device.DeviceInfo {
	current := device.DeviceInfo{
		IP:       ip,
		MAC:      "N/A",
		Hostname: "N/A",
	}

	// Get the hostname
	hostname, err := getHostname(ip)
	if err == nil {
		current.Hostname = hostname
	}

	// Get the MAC address
	mac, err := getMacAddress(ip)
	if err == nil && mac != "" {
		current.MAC = mac
	}

	if methods.Icmp {
		err := SendIcmpPing(ip)
		current.IcmpResponse.Requested = true
		current.IcmpResponse.Responded = (err == nil)
	}
	if methods.Arp {
		err := SendArpRequest(ip)
		current.ArpResponse.Requested = true
		current.ArpResponse.Responded = (err == nil)
	}
	if methods.Tcp {
		err := SendTcpSyn(ip, methods.TcpPort)
		current.TcpResponse.Requested = true
		current.TcpResponse.Responded = (err == nil)
	}

	return current
}

// getHostname performs a reverse DNS lookup to get the hostname for an IP.
func getHostname(ip string) (string, error) {
	names, err := net.LookupAddr(ip)
	if err != nil || len(names) == 0 {
		return "", err
	}
	return strings.TrimSuffix(names[0], "."), nil
}

// getMacAddress retrieves the MAC address for an IP using the `arp` command.
func getMacAddress(ip string) (string, error) {
	out, err := exec.Command("arp", "-a", ip).Output()
	if err != nil {
		return "", err
	}

	lines := strings.Split(string(out), "\n")
	for _, line := range lines {
		if strings.Contains(line, ip) {
			fields := strings.Fields(line)
			if len(fields) >= 2 {
				return fields[1], nil // MAC address is usually the second field
			}
		}
	}
	return "", nil
}
