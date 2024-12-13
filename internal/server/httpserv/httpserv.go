package httpserv

import (
	"errors"
	"net"
	"net/http"

	"github.com/go-chi/chi/v5"

	_ "net/http/pprof"

	"github.com/Ord1nI/metrix/internal/server"
	"github.com/Ord1nI/metrix/internal/server/httpserv/middlewares"
)


type HttpServer struct {
	*server.Server
	Serv        http.Server
	Middlewares chi.Middlewares
}

func Default() (*HttpServer, error) {
	serv, err := New()
	if err != nil {
		return nil, err;
	}
	serv.Add(middlewares.LoggerMW(serv.Logger))

	if serv.Config.PrivateKeyFile != "" {
		serv.Add(middlewares.Decrypt(serv.Logger,serv.Config.PrivateKeyFile))
	}

	if serv.Config.Key != "" {
		serv.Add(middlewares.SignMW(serv.Logger, []byte(serv.Config.Key)))
	}

	if serv.Config.TrustedSubnet != "" {
		serv.Add(middlewares.CheckSubnet(serv.Logger, net.ParseIP(serv.Config.TrustedSubnet)))
	}

	serv.Add(middlewares.CompressorMW(serv.Logger))

	err = serv.Init()
	if err != nil {
		serv.Logger.Errorln("Fail while starting server")
		return nil, err
	}

	return serv, nil
}

// New constructor for Server
// Also calls GetConf
// And adds HeadMW as first middleware
func New() (*HttpServer, error) {
	mServer, err := server.New()

	if err != nil {
		return nil, err
	}

	s := &HttpServer{
		Server: mServer,
	}

	s.Serv.Addr = s.Config.Address

	s.Add(middlewares.HeadMW(s.Logger))

	return s, nil
}

// Init method that apply middlewarees must be called befor start.
func (s *HttpServer) Init() error {
	s.InitRouter(s.Middlewares...)
	return nil
}

// Add method to add middlewares in server list must be call before Init.
func (s *HttpServer) Add(mw ...func(http.Handler) http.Handler) {
	s.Middlewares = append(s.Middlewares, mw...)
}


func (s *HttpServer) startServ(stop <-chan struct{}) {
	go s.Serv.ListenAndServe();

	<-stop
	s.Serv.Close()
	if err := s.Repo.Close(); err != nil {
		s.Logger.Error(err)
	}
}

// Run Method to start server
func (s *HttpServer) Run(stop <-chan struct{}) error {

	if s.Serv.Handler != nil {
		go s.startServ(stop)
		return nil;
	}

	s.Repo.Close()

	return errors.New("router not initialized")
}
