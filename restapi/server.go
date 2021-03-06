package restapi

import (
	"crypto/tls"
	"fmt"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/go-openapi/swag"
	flags "github.com/jessevdk/go-flags"
	graceful "github.com/tylerb/graceful"

	"github.com/tompscanlan/packerd/restapi/operations"
)

//go:generate swagger generate server -t ../.. -A Packerd -f ./swagger.yml

// NewServer creates a new api packerd server but does not configure it
func NewServer(api *operations.PackerdAPI) *Server {
	s := new(Server)
	s.api = api
	return s
}

// ConfigureAPI configures the API and handlers. Needs to be called before Serve
func (s *Server) ConfigureAPI() {
	if s.api != nil {
		s.handler = configureAPI(s.api)
	}
}

// ConfigureFlags configures the additional flags defined by the handlers. Needs to be called before the parser.Parse
func (s *Server) ConfigureFlags() {
	if s.api != nil {
		configureFlags(s.api)
	}
}

// Server for the packerd API
type Server struct {
	Host        string `long:"host" description:"the IP to listen on" default:"localhost" env:"HOST"`
	Port        int    `long:"port" description:"the port to listen on for insecure connections, defaults to a random value" env:"PORT"`
	httpServerL net.Listener

	TLSHost           string         `long:"tls-host" description:"the IP to listen on for tls, when not specified it's the same as --host" env:"TLS_HOST"`
	TLSPort           int            `long:"tls-port" description:"the port to listen on for secure connections, defaults to a random value" env:"TLS_PORT"`
	TLSCertificate    flags.Filename `long:"tls-certificate" description:"the certificate to use for secure connections" required:"true" env:"TLS_CERTIFICATE"`
	TLSCertificateKey flags.Filename `long:"tls-key" description:"the private key to use for secure conections" required:"true" env:"TLS_PRIVATE_KEY"`
	httpsServerL      net.Listener

	api          *operations.PackerdAPI
	handler      http.Handler
	hasListeners bool
}

// Logf logs message either via defined user logger or via system one if no user logger is defined.
func (s *Server) Logf(f string, args ...interface{}) {
	if s.api != nil && s.api.Logger != nil {
		s.api.Logger(f, args...)
	} else {
		log.Printf(f, args...)
	}
}

// SetAPI configures the server with the specified API. Needs to be called before Serve
func (s *Server) SetAPI(api *operations.PackerdAPI) {
	if api == nil {
		s.api = nil
		s.handler = nil
		return
	}

	s.api = api
	s.api.Logger = log.Printf
	s.handler = configureAPI(api)
}

// Serve the api
func (s *Server) Serve() (err error) {
	if !s.hasListeners {
		if err := s.Listen(); err != nil {
			return err
		}
	}

	httpServer := &graceful.Server{Server: new(http.Server)}
	httpServer.Handler = s.handler

	s.Logf("Serving packerd at http://%s", s.httpServerL.Addr())
	go func(l net.Listener) {
		if err := httpServer.Serve(tcpKeepAliveListener{l.(*net.TCPListener)}); err != nil {
			log.Fatalln(err)
		}
	}(s.httpServerL)

	httpsServer := &graceful.Server{Server: new(http.Server)}
	httpsServer.Handler = s.handler
	httpsServer.TLSConfig = new(tls.Config)
	httpsServer.TLSConfig.NextProtos = []string{"http/1.1"}
	// https://www.owasp.org/index.php/Transport_Layer_Protection_Cheat_Sheet#Rule_-_Only_Support_Strong_Protocols
	httpsServer.TLSConfig.MinVersion = tls.VersionTLS11
	httpsServer.TLSConfig.Certificates = make([]tls.Certificate, 1)
	httpsServer.TLSConfig.Certificates[0], err = tls.LoadX509KeyPair(string(s.TLSCertificate), string(s.TLSCertificateKey))

	configureTLS(httpsServer.TLSConfig)

	if err != nil {
		return err
	}
	s.Logf("Serving packerd at https://%s", s.httpsServerL.Addr())

	wrapped := tls.NewListener(tcpKeepAliveListener{s.httpsServerL.(*net.TCPListener)}, httpsServer.TLSConfig)
	if err := httpsServer.Serve(wrapped); err != nil {
		return err
	}

	return nil
}

// Listen creates the listeners for the server
func (s *Server) Listen() error {
	if s.hasListeners { // already done this
		return nil
	}

	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", s.Host, s.Port))
	if err != nil {
		return err
	}

	h, p, err := swag.SplitHostPort(listener.Addr().String())
	if err != nil {
		return err
	}
	s.Host = h
	s.Port = p
	s.httpServerL = listener

	if s.TLSHost == "" {
		s.TLSHost = s.Host
	}
	tlsListener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", s.TLSHost, s.TLSPort))
	if err != nil {
		return err
	}

	sh, sp, err := swag.SplitHostPort(tlsListener.Addr().String())
	if err != nil {
		return err
	}
	s.TLSHost = sh
	s.TLSPort = sp
	s.httpsServerL = tlsListener
	s.hasListeners = true
	return nil
}

// Shutdown server and clean up resources
func (s *Server) Shutdown() error {
	s.api.ServerShutdown()
	return nil
}

// tcpKeepAliveListener is copied from the stdlib net/http package

// tcpKeepAliveListener sets TCP keep-alive timeouts on accepted
// connections. It's used by ListenAndServe and ListenAndServeTLS so
// dead TCP connections (e.g. closing laptop mid-download) eventually
// go away.
type tcpKeepAliveListener struct {
	*net.TCPListener
}

func (ln tcpKeepAliveListener) Accept() (c net.Conn, err error) {
	tc, err := ln.AcceptTCP()
	if err != nil {
		return
	}
	tc.SetKeepAlive(true)
	tc.SetKeepAlivePeriod(3 * time.Minute)
	return tc, nil
}
