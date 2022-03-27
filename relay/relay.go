package relay

import (
	"fmt"
	"net"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	cconfig "github.com/notedit/relay-server/relay/config"
	"github.com/pion/turn/v2"
	"github.com/pkg/errors"
)

type RelayServer struct {
	turnServer *turn.Server
	config     *cconfig.Config
	relayGen   turn.RelayAddressGenerator
}

func NewRelayServer(configFile string) (*RelayServer, error) {
	config, err := cconfig.LoadConfig(configFile)
	if err != nil {
		return nil, err
	}

	server := &RelayServer{
		config: config,
		relayGen: &turn.RelayAddressGeneratorStatic{
			RelayAddress: net.ParseIP(config.Server.PublicIP),
			Address:      "0.0.0.0",
		},
	}

	serverConfig := turn.ServerConfig{
		Realm: config.Server.Realm,
	}

	tcpListener, err := net.Listen("tcp4", "0.0.0.0:"+strconv.Itoa(config.Server.TCPPort))
	if err != nil {
		return nil, err
	}

	listenerConfig := turn.ListenerConfig{
		Listener:              tcpListener,
		RelayAddressGenerator: server.relayGen,
	}
	serverConfig.ListenerConfigs = append(serverConfig.ListenerConfigs, listenerConfig)

	udpListener, err := net.ListenPacket("udp4", "0.0.0.0:"+strconv.Itoa(config.Server.UDPPort))
	if err != nil {
		return nil, errors.Wrap(err, "could not listen on TURN UDP port")
	}

	packetConfig := turn.PacketConnConfig{
		PacketConn:            udpListener,
		RelayAddressGenerator: server.relayGen,
	}
	serverConfig.PacketConnConfigs = append(serverConfig.PacketConnConfigs, packetConfig)

	serverConfig.AuthHandler = server.authHandler
	server.turnServer, err = turn.NewServer(serverConfig)

	if err != nil {
		return nil, errors.Wrap(err, "init turn server error")
	}

	fmt.Println("Relay server started")
	fmt.Println("Public IP:", config.Server.PublicIP)
	fmt.Println("Realm:", config.Server.Realm)
	fmt.Println("Password:", config.Server.Password)
	fmt.Println("TCP Port:", config.Server.TCPPort)
	fmt.Println("UDP Port:", config.Server.UDPPort)

	return server, nil
}

func (s *RelayServer) authHandler(username string, realm string, srcAddr net.Addr) ([]byte, bool) {

	fmt.Println("authHandler: ", username, realm, srcAddr.String())

	return turn.GenerateAuthKey(username, s.config.Server.Realm, s.config.Server.Password), true
}

func (s *RelayServer) Run() {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	<-sigs

	if err := s.turnServer.Close(); err != nil {
		fmt.Println(err)
	}
}
