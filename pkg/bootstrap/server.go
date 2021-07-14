package bootstrap

import "github.com/zirain/mcpoverxds/pkg/config"

type Server struct {
	discoveryAddress string
	configController *config.Controller
}

func NewServer(serverArgs *ServerArgs) (*Server, error) {
	return &Server{
		discoveryAddress: serverArgs.DiscoveryAddress,
		configController: config.NewController(serverArgs.DiscoveryAddress),
	}, nil
}

func (s *Server) Start(stop <-chan struct{}) error {
	s.configController.RegisterEventHandler()
	s.configController.Run(stop)
	return nil
}
