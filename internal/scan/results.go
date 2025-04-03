package scan

import "soma/internal/device"

// collectResults gathers results from the result channel into a slice.
func collectResults(resultChan chan device.DeviceInfo) []device.DeviceInfo {
	var devices []device.DeviceInfo
	for result := range resultChan {
		devices = append(devices, result)
	}
	return devices
}
