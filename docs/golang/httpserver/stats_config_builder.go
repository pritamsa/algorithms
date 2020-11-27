package httpserver

import (
	"net"

	"github.com/aws/aws-sdk-go/aws/ec2metadata"
	"github.com/aws/aws-sdk-go/aws/session"
	"gitlab.nordstrom.com/sentry/authorize/httpserver/flags"
	"gitlab.nordstrom.com/sentry/gohttp/statsd_wrapper"
)

//go:generate counterfeiter . StatsConfigBuilder
type StatsConfigBuilder interface {
	StatsConfigFromArgs(flags.CommandlineArgs) statsd_wrapper.StatsdConfig
}

func NewStatsConfigBuilder() StatsConfigBuilder {
	return statsConfigBuilder{}
}

var sess = session.Must(session.NewSession())

type statsConfigBuilder struct{}

func (builder statsConfigBuilder) StatsConfigFromArgs(args flags.CommandlineArgs) statsd_wrapper.StatsdConfig {
	statsConfig := statsd_wrapper.StatsdConfig{
		Prefix:      args.Stats.Prefix,
		HostName:    args.Stats.Host,
		Port:        args.Stats.Port,
		Environment: args.Server.Environment,
	}

	var ipAddressString string
	var err error

	//set host ip tp hostip on non-local env
	if args.Stats.Host == "hostip" {
		ipAddressString, err = getHostIP()
	} else {
		var ipAddress *net.IPAddr
		ipAddress, err = net.ResolveIPAddr("ip", args.Stats.Host)
		ipAddressString = ipAddress.String()
	}

	if err == nil {
		statsConfig.HostName = ipAddressString
	}
	return statsConfig
}

func getHostIP() (string, error) {
	var err error
	var ipAddressString string
	var instDocument ec2metadata.EC2InstanceIdentityDocument
	svc := ec2metadata.New(sess)
	instDocument, err = svc.GetInstanceIdentityDocument()
	if err == nil {
		ipAddressString = instDocument.PrivateIP
	}
	return ipAddressString, err
}
