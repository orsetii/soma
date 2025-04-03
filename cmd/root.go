package cmd

import (
	"fmt"
	"os"

	"soma/internal/display"
	"soma/internal/scan"

	"github.com/urfave/cli/v2"
)

type Flags struct {
	Subnet string
	Ping   bool
	TCP    int
}

func Execute() {
	app := &cli.App{
		Name:  "Soma",
		Usage: "Scan a given subnet for devices and display their information",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "subnet",
				Usage:    "Specify the subnet to scan (e.g., 192.168.1.0/24)",
				Required: true,
			},
			&cli.BoolFlag{
				Name:  "ping",
				Usage: "Use ICMP pings on IPs in the range",
			},
			&cli.BoolFlag{
				Name:  "arp",
				Usage: "Use ARP Requests on IPs in the range",
			},
			&cli.IntFlag{
				Name:  "tcp",
				Usage: "Use TCP SYN requests on a given port, on IPs in the range",
			},
		},
		Action: func(c *cli.Context) error {
			subnet := c.String("subnet")
			ping := c.Bool("ping")
			arp := c.Bool("arp")
			tcp := c.Int("tcp")

			discoveryMethods := scan.DiscoveryMethods{
				Icmp:    ping,
				Arp:     arp,
				Tcp:     tcp > 0,
				TcpPort: tcp,
			}
			if !discoveryMethods.Icmp &&
				!discoveryMethods.Arp &&
				!discoveryMethods.Tcp {
				return fmt.Errorf("at least one discovery method must be specified")
			}

			if discoveryMethods.Tcp && discoveryMethods.TcpPort == 0 {
				return fmt.Errorf("TCP port must be specified when TCP discovery is enabled")
			}

			results, err := scan.Scan(subnet, discoveryMethods)
			if err != nil {
				return err
			}

			display.DisplayResults(results)

			return nil
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
