package scan

import (
	"fmt"
	"net"
	"time"
)

func SendIcmpPing(ip string) error {
	conn, err := net.Dial("ip4:icmp", ip)
	if err != nil {
		return err
	}
	defer conn.Close()

	icmp := []byte{8, 0, 0, 0, 0, 0, 0, 0}
	checksum := calculateChecksum(icmp)
	icmp[2], icmp[3] = byte(checksum>>8), byte(checksum&0xff)

	if _, err := conn.Write(icmp); err != nil {
		return err
	}

	conn.SetReadDeadline(time.Now().Add(2 * time.Second))
	reply := make([]byte, 1024)
	if _, err := conn.Read(reply); err != nil {
		return err
	}

	if len(reply) < 1 || reply[0] != 0 {
		return fmt.Errorf("unexpected ICMP response from %s", ip)
	}
	return nil
}

func calculateChecksum(data []byte) uint16 {
	sum := 0
	for i := 0; i < len(data)-1; i += 2 {
		sum += int(data[i])<<8 | int(data[i+1])
	}
	if len(data)%2 == 1 {
		sum += int(data[len(data)-1]) << 8
	}
	sum = (sum & 0xffff) + (sum >> 16)
	return ^uint16(sum)
}
