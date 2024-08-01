package main

import (
	"fmt"
	"net"
	"net/netip"
	"os"
	"sync"
	"time"

	"io"
	"log"
	"net/http"

	"golang.zx2c4.com/wireguard/conn"
	"golang.zx2c4.com/wireguard/device"
	"golang.zx2c4.com/wireguard/tun/netstack"
	gtcp "gvisor.dev/gvisor/pkg/tcpip/transport/tcp"
	gudp "gvisor.dev/gvisor/pkg/tcpip/transport/udp"

	"golang.zx2c4.com/wireguard/wgctrl/wgtypes"

	"wiretap/peer"
	"wiretap/transport/icmp"
	"wiretap/transport/tcp"
	"wiretap/transport/udp"

	"github.com/rainforestapp/rainforest-cli/rainforest"
	"github.com/urfave/cli"
)

func launchDirectConnect(c cliContext, rfApi *rainforest.Client) error {
	tunnelId := c.String("tunnel-id")
	if tunnelId == "" {
		log.Println("Starting Direct Connect Tunnel")
	} else {
		log.Println("Starting Direct Connect Tunnel for tunnel ID:", tunnelId)
	}

	// Generate a wireguard public/private keypair
	privateKey, err := wgtypes.GeneratePrivateKey()
	if err != nil {
		return cli.NewExitError(fmt.Errorf("failed to generate private key: %w", err), 1)
	}
	publicKey := privateKey.PublicKey()

	// Ask the rainforest api to set up our connection, returning the server details
	serverDetails, err := rfApi.SetupDirectConnectTunnel(tunnelId, publicKey.String())
	if err != nil {
		return cli.NewExitError(fmt.Errorf("failed to configure direct connect: %w", err), 1)
	}

	var (
		wg   sync.WaitGroup
		lock sync.Mutex
	)

	// Static configuration settings
	listenPort := 51820
	persistentKeepaliveInterval := 25
	mtu := 1420
	allowedIPs := []string{"198.18.18.1/32", "fd:16::1/128"}
	ipv4Addr, err := netip.ParseAddr("198.18.18.2")
	if err != nil {
		return cli.NewExitError(fmt.Errorf("failed to parse ipv4 address: %w", err), 1)
	}
	ipv6Addr, err := netip.ParseAddr("fd:16::2")
	if err != nil {
		return cli.NewExitError(fmt.Errorf("failed to parse ipv6 address: %w", err), 1)
	}

	wireguardConfigArgs := peer.ConfigArgs{
		PrivateKey: privateKey.String(),
		ListenPort: listenPort,
		Peers: []peer.PeerConfigArgs{
			{
				PublicKey:                   serverDetails.ServerPublicKey,
				Endpoint:                    fmt.Sprintf("%s:%d", serverDetails.ServerEndpoint, serverDetails.ServerPort),
				PersistentKeepaliveInterval: persistentKeepaliveInterval,
				AllowedIPs:                  allowedIPs,
			},
		},
		Addresses: []string{ipv4Addr.String() + "/32"},
	}

	wireguardConfig, err := peer.GetConfig(wireguardConfigArgs)
	if err != nil {
		return cli.NewExitError(fmt.Errorf("failed to make wireguard configuration: %w", err), 1)
	}

	// Create virtual wireguard interface with these addresses and MTU.
	tunnelAddrs := []netip.Addr{ipv4Addr, ipv6Addr}
	tunWireguard, tnetWireguard, err := netstack.CreateNetTUN(
		tunnelAddrs,
		[]netip.Addr{},
		mtu,
	)
	if err != nil {
		return cli.NewExitError(fmt.Errorf("failed to create wireguard TUN: %w", err), 1)
	}
	// Set the interface to promiscuous mode so forwarding works properly
	s := tnetWireguard.Stack()
	s.SetPromiscuousMode(1, true)

	// TCP Forwarding settings.
	catchTimeout := 5 * 1000
	connTimeout := 5 * 1000
	keepaliveIdle := 60
	keepaliveCount := 3
	keepaliveInterval := 60

	// TCP Forwarding mechanism.
	tcpConfig := tcp.Config{
		CatchTimeout:      time.Duration(catchTimeout) * time.Millisecond,
		ConnTimeout:       time.Duration(connTimeout) * time.Millisecond,
		KeepaliveIdle:     time.Duration(keepaliveIdle) * time.Second,
		KeepaliveInterval: time.Duration(keepaliveInterval) * time.Second,
		KeepaliveCount:    int(keepaliveCount),
		Tnet:              tnetWireguard,
		StackLock:         &lock,
	}
	tcpForwarder := gtcp.NewForwarder(s, 0, 65535, tcp.Handler(tcpConfig))
	s.SetTransportProtocolHandler(gtcp.ProtocolNumber, tcpForwarder.HandlePacket)

	// UDP Forwarding mechanism.
	udpConfig := udp.Config{
		Tnet:      tnetWireguard,
		StackLock: &lock,
	}
	s.SetTransportProtocolHandler(gudp.ProtocolNumber, udp.Handler(udpConfig))

	// Configure the wireguard device with the wireguard details
	var logger int = device.LogLevelError // Or device.LogLevelVerbose or device.LogLevelSilent
	devWireguard := device.NewDevice(tunWireguard, conn.NewDefaultBind(), device.NewLogger(logger, ""))
	// Configure wireguard.
	err = devWireguard.IpcSet(wireguardConfig.AsIPC())
	if err != nil {
		return cli.NewExitError(fmt.Errorf("failed to configure wireguard device: %w", err), 1)
	}
	err = devWireguard.Up()
	if err != nil {
		return cli.NewExitError(fmt.Errorf("failed to bring up device: %w", err), 1)
	}

	// Handlers that require long-running routines:
	// Start ICMP Handler.
	wg.Add(1)
	go func() {
		icmp.Handle(tnetWireguard, &lock)
		wg.Done()
	}()

	// Start API handler.
	wg.Add(1)
	go func() {
		apiHandler(tnetWireguard, ipv4Addr, uint16(80))
		wg.Done()
	}()
	log.Printf("Rainforest Direct Connect running (name: %s, PID: %d)! Press Ctrl+C to exit\n", serverDetails.Name, os.Getpid())

	wg.Wait()

	return nil
}

func apiHandler(tnet *netstack.Net, addr netip.Addr, port uint16) {
	// Stand up API server.
	listener, err := tnet.ListenTCP(&net.TCPAddr{IP: addr.AsSlice(), Port: int(port)})
	if err != nil {
		log.Panic(err)
	}

	http.HandleFunc("/ping", handlePing)

	err = http.Serve(listener, nil)
	if err != nil {
		log.Panic(err)
	}
}

// handlePing responds with pong message.
func handlePing(w http.ResponseWriter, r *http.Request) {
	log.Printf("(client %s) - API: %s", r.RemoteAddr, r.RequestURI)
	_, err := io.WriteString(w, "pong\n")
	if err != nil {
		log.Printf("API Error: %v", err)
	}
}
