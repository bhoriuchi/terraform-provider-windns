package windns

import (
	"time"

	"github.com/bhoriuchi/go-winrmkrb5"
	"github.com/masterzen/winrm"
)

const (
	defaultTTL     = 300
	defaultTimeout = 60
	winrmHTTP      = 5985
	winrmHTTPS     = 5986
)

// Client a windns client
type Client struct {
	o *Options
	c *winrm.Client
}

// Record record
type Record struct {
	Type string      `json:"type"`
	Name string      `json:"name"`
	Zone string      `json:"zone"`
	Data interface{} `json:"data"`
	TTL  int         `json:"ttl"`
}

// Response a response object
type Response struct {
	Code    int       `json:"code"`
	Detail  string    `json:"detail"`
	Records []*Record `json:"records"`
}

// Options client options
type Options struct {
	DnsServer       string
	ProxyHost       string
	Username        string
	Password        string
	Domain          string
	SecureTransport bool
	SkipSSLVerify   bool
	Port            int
	KDCServers      []string
	KRB5Conf        string
	TimeoutSeconds  int
}

// NewClient creates a new client
func NewClient(opts *Options) (*Client, error) {
	var err error
	client := &Client{o: opts}

	port := opts.Port
	if port == 0 {
		port = winrmHTTP
		if opts.SecureTransport {
			port = winrmHTTPS
		}
	}

	host := opts.DnsServer
	if opts.ProxyHost != "" {
		host = opts.ProxyHost
	}

	timeout := time.Duration(opts.TimeoutSeconds)
	if timeout == 0 || timeout < 0 {
		timeout = defaultTimeout
	}

	kdcs := opts.KDCServers
	if len(opts.KDCServers) == 0 {
		kdcs = []string{opts.DnsServer}
	}

	endpoint := winrm.NewEndpoint(
		host,
		port,
		opts.SecureTransport,
		opts.SkipSSLVerify,
		nil,
		nil,
		nil,
		timeout*time.Second,
	)

	winrm.DefaultParameters.TransportDecorator = func() winrm.Transporter {
		return &winrmkrb5.Transport{
			Username: opts.Username,
			Password: opts.Password,
			Domain:   opts.Domain,
			KDC:      kdcs,
			KRB5Conf: opts.KRB5Conf,
			Timeout:  timeout * time.Second,
		}
	}

	if client.c, err = winrm.NewClientWithParameters(
		endpoint,
		"",
		"",
		winrm.DefaultParameters,
	); err != nil {
		return nil, err
	}

	return client, nil
}
