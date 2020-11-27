package httpserver_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "gitlab.nordstrom.com/sentry/authorize/httpserver"
	"gitlab.nordstrom.com/sentry/authorize/httpserver/flags"
)

var _ = Describe("StatsConfigBuilder", func() {
	var (
		subject StatsConfigBuilder
	)

	BeforeEach(func() {
		subject = NewStatsConfigBuilder()
	})

	It("reads the correct values from the command line args", func() {
		config := subject.StatsConfigFromArgs(flags.CommandlineArgs{
			Stats: flags.StatsArgs{
				Host:   "127.0.0.1",
				Port:   9999,
				Prefix: "the-stats-prefix",
			},
		})

		Expect(config.Prefix).To(Equal("the-stats-prefix"))
		Expect(config.HostName).To(Equal("127.0.0.1"))
		Expect(config.Port).To(Equal(9999))
	})

	It("resolves the statsd host to an IP", func() {
		config := subject.StatsConfigFromArgs(flags.CommandlineArgs{
			Stats: flags.StatsArgs{
				Host: "localhost",

				Port:   9999,
				Prefix: "the-stats-prefix",
			},
		})

		Expect(config.HostName).To(Equal("127.0.0.1"))
	})

	It("uses the hostname if unable to resolve the statsd IP", func() {
		config := subject.StatsConfigFromArgs(flags.CommandlineArgs{
			Stats: flags.StatsArgs{
				Host: "localhost1234",

				Port:   9999,
				Prefix: "the-stats-prefix",
			},
		})

		Expect(config.HostName).To(Equal("localhost1234"))
	})
})
