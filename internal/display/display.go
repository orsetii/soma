package display

import (
	"os"
	"soma/internal/device"
	"sort"

	"github.com/jedib0t/go-pretty/v6/table"
)

// Displays a DeviceInfo slice in a table format
func DisplayResults(results []device.DeviceInfo) {
	// Sort results by IP address
	sort.Slice(results, func(i, j int) bool {
		return results[i].IP < results[j].IP
	})

	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"IP Address", "Hostname", "MAC Address", "ICMP", "ARP", "TCP"})

	for _, device := range results {
		t.AppendRow(table.Row{device.IP, device.Hostname, device.MAC, getResponseSymbol(device.IcmpResponse), getResponseSymbol(device.ArpResponse), getResponseSymbol(device.TcpResponse)})
	}
	t.Render()
}

func getResponseSymbol(resp device.ResponseType) string {
	if resp.Requested && resp.Responded {
		return "✓"
	} else if resp.Requested && !resp.Responded {
		return "✗"
	} else {
		return "?"
	}
}
