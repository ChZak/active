package main

import (
	"flag"
	"fmt"
	"net"
	"strconv"
	"strings"
	"time"
)

func main() {
	udpPtr := flag.Bool("u", false, "UDP scan")
	tcpPtr := flag.Bool("t", false, "TCP scan")
	portRangePtr := flag.String("p", "", "Range of ports to scan")
	helpPtr := flag.Bool("help", false, "Show help message")

	flag.Parse()

	if *helpPtr {
		flag.Usage()
		return
	}

	host := flag.Arg(0)

	if *tcpPtr {
		ports, err := parsePortRange(*portRangePtr)
		if err != nil {
			fmt.Println("Invalid port range")
			return
		}

		for _, port := range ports {
			scanTCP(host, port)
		}
	}

	if *udpPtr {
		ports, err := parsePortRange(*portRangePtr)
		if err != nil {
			fmt.Println("Invalid port range")
			return
		}
		for _, port := range ports {
			scanUDP(host, port)
		}
	}
}

// Effectue un balayage TCP sur le port donné
func scanTCP(host string, port int) {
	address := net.JoinHostPort(host, strconv.Itoa(port))
	conn, err := net.DialTimeout("tcp", address, 1*time.Second)
	if err != nil {
		fmt.Printf("Port %d is closed\n", port)
		return
	}

	conn.Close()
	fmt.Printf("Port %d is open\n", port)
}

// scanUDP effectue un balayage UDP sur un port spécifique
func scanUDP(host string, port int) {
	address := net.JoinHostPort(host, strconv.Itoa(port))
	conn, err := net.Dial("udp", address)
	if err != nil {
		fmt.Printf("Port %d may be closed (error: %s)\n", port, err)
		return
	}
	defer conn.Close()

	// Envoi d'un paquet UDP vide
	_, err = conn.Write([]byte{})
	if err != nil {
		fmt.Printf("Port %d may be closed (error: %s)\n", port, err)
		return
	}

	// Écoute d'une réponse ou d'un timeout
	conn.SetReadDeadline(time.Now().Add(2 * time.Second))
	buffer := make([]byte, 1024)
	_, err = conn.Read(buffer)
	if err != nil {
		if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
			fmt.Printf("Port %d is open or filtered\n", port)
		} else {
			fmt.Printf("Port %d may be closed (error: %s)\n", port, err)
		}
		return
	}

	fmt.Printf("Received response from port %d, it might be open\n", port)
}

// Analyse une chaine de ports et retourne une tranche de ports
func parsePortRange(portRange string) ([]int, error) {
	var ports []int
	if strings.Contains(portRange, "-") {
		parts := strings.Split(portRange, "-")
		start, err := strconv.Atoi(parts[0])
		if err != nil {
			return nil, err
		}
		end, err := strconv.Atoi(parts[1])
		if err != nil {
			return nil, err
		}
		for i := start; i <= end; i++ {
			ports = append(ports, i)
		}
	} else {
		singlePort, err := strconv.Atoi(portRange)
		if err != nil {
			return nil, err
		}

		ports = append(ports, singlePort)
	}

	return ports, nil
}
