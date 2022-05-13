package main

import (
	"errors"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
)

type CmdLineArgs struct {
	ipFrom string
	ipTo   string
	port   uint16
}

func main() {
	cmdLineArgs, err := parseCmdLineArgs(os.Args)

	if err != nil {
		fmt.Println("Usage: psc <ip-from> [<ip-to>] <port>")
		return
	}

	nextIp := cmdLineArgs.ipFrom
	for {
		if tryConnect(nextIp, cmdLineArgs.port) {
			fmt.Printf("AVAILABLE %s:%d\n", nextIp, cmdLineArgs.port)
		} else {
			fmt.Printf("unavailable %s:%d\n", nextIp, cmdLineArgs.port)
		}

		if nextIp == cmdLineArgs.ipTo {
			break
		}

		nextIp, err = getNextIp(nextIp)
		if err != nil {
			break
		}
	}
}

func parseCmdLineArgs(args []string) (CmdLineArgs, error) {
	if len(args) < 3 {
		return CmdLineArgs{}, errors.New("Invalid arguments count")
	}

	parseFn := func(ipFromArg string, ipToArg string, portArg string) (CmdLineArgs, error) {
		var ipFrom, ipTo string
		var port uint16
		var err error

		ipFrom, err = parseIp(ipFromArg)
		if err != nil {
			return CmdLineArgs{}, errors.New("Invalid ip-from value")
		}

		ipTo, err = parseIp(ipToArg)
		if err != nil {
			return CmdLineArgs{}, errors.New("Invalid ip-to value")
		}

		port64, err := strconv.ParseUint(portArg, 10, 16)
		if err != nil {
			return CmdLineArgs{}, errors.New("Invalid port value")
		}

		port = uint16(port64)

		return CmdLineArgs{ipFrom, ipTo, port}, nil
	}

	if len(args) == 3 {
		return parseFn(args[1], args[1], args[2])
	} else {
		return parseFn(args[1], args[2], args[3])
	}
}

func parseIp(s string) (string, error) {
	isValidOctetFn := func(octet int) bool {
		return octet >= 0 && octet <= 255
	}

	octets := strings.Split(s, ".")

	if len(octets) != 4 {
		return "", errors.New("Invalid octet count")
	}

	for _, octet := range octets {
		octetInt32, err := strconv.Atoi(octet)
		if err != nil {
			return "", err
		}

		if !isValidOctetFn(octetInt32) {
			return "", errors.New("Invalid octet value")
		}
	}

	return s, nil
}

func getNextIp(ip string) (string, error) {
	octets := strings.Split(ip, ".")

	for i := 3; i >= 0; i-- {
		octet := octets[i]
		octetInt, _ := strconv.Atoi(octet)
		octetInt++

		overflow := false
		if octetInt > 255 {
			octetInt -= 255
			overflow = true
		}

		octets[i] = strconv.Itoa(octetInt)

		if !overflow {
			break
		}

		if overflow && i == 0 {
			return "", errors.New("IP address overflow")
		}
	}

	return strings.Join(octets, "."), nil
}

func tryConnect(ip string, port uint16) bool {
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", ip, port))
	if err != nil {
		return false
	}

	conn.Close()
	return true
}
