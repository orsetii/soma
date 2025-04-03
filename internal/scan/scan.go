package scan

import (
	"soma/internal/device"
	"sync"
)

type DiscoveryMethods struct {
	Icmp    bool
	Arp     bool
	Tcp     bool
	TcpPort int
}

// Scan a given subnet for devices
func Scan(subnet string, methods DiscoveryMethods) ([]device.DeviceInfo, error) {
	ips, err := ExpandCIDR(subnet)
	if err != nil {
		return nil, err
	}

	pw, tracker := initializeProgress(len(ips))
	defer (*pw).Stop()

	ipChan := make(chan string, len(ips))
	resultChan := make(chan device.DeviceInfo, len(ips))

	var wg sync.WaitGroup
	startWorkers(ipChan, resultChan, &wg, methods, tracker)

	// Send IPs to workers
	for _, ip := range ips {
		ipChan <- ip.String()
	}
	close(ipChan)

	// Collect results and update progress
	var devices []device.DeviceInfo
	go func() {
		for result := range resultChan {
			devices = append(devices, result)
			tracker.Increment(1)
		}
	}()

	wg.Wait()
	close(resultChan)

	return devices, nil
}
