package flags

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	goflags "github.com/jessevdk/go-flags"
)

// CommandlineArgs defines the command line argument groups.
type CommandlineArgs struct {
	Server ServerArgs `group:"Server" namespace:"server"`

	TLS TLSArgs `group:"Tls" namespace:"tls"`

	Stats StatsArgs `group:"Stats" namespace:"stats"`

	DoomHammer DoomHammerArgs `group:"DoomHammer" namespace:"doom-hammer"`

	ShopperAuth ShopperAuthArgs `group:"ShopperAuth" namespace:"shopper-auth-client"`

	ShopperAPI ShopperAPIArgs `group:"ShopperApi" namespace:"shopper-api-client"`

	ShopperAccountAPI ShopperAccountAPIArgs `group:"ShopperAccountApi" namespace:"shopper-account-api-client"`

	Verify VerifyArgs `group:"Verify" namespace:"verify-client"`

	ApigeeAPI ApigeeAPIArgs `group:"ApigeeApi" namespace:"apigee-client"`

	ForterAPI ForterAPIArgs `group:"ForterAPI" namespace:"forter-api"`

	MfaBypass MfaBypassArgs `group:"MfaBypass" namespace:"mfa-bypass"`

	Encryption EncryptionArgs `group:"Encryption" namespace:"encryption"`

	APM APMArgs `group:"APM" namespace:"apm"`

	ShopperToken ShopperTokenArgs `group:"ShopperToken" namespace:"shopper-token"`

	Domain DomainArg `group:"Domain" namespace:"domain"`

	HealthCheck HealthCheck `group:"Healthcheck" namespace:"healthcheck"`
}

type DomainArg struct {
	AppName string `long:"app-name" required:"true" description:"App name of application"`
	Suffix  string `long:"suffix" required:"true" description:"Domain suffix for application"`
}

type EncryptionArgs struct {
	KeyPath string `long:"key-path" required:"true" description:"Encryption key file path"`
	IvPath  string `long:"iv-path" required:"true" description:"Encryption initialization vector file path"`
}

// ServerArgs defines the server arguments.
type ServerArgs struct {
	HTTPSPort   int    `long:"https-port" required:"true" description:"Port for the server to listen on HTTPS"`
	Environment string `long:"environment" required:"true" description:"Environment type" choice:"prod" choice:"dev" choice:"ci" choice:"int" choice:"stage" choice:"perf" choice:"int-canary" choice:"prod-canary"`
}

// TLSArgs defines the TLS arguments.
type TLSArgs struct {
	CertPath string `long:"cert-path" required:"true" description:"Path to the path to the certificate used for SSL"`
	KeyPath  string `long:"key-path" required:"true" description:"Path to the private key for the certificate"`
}

// StatsArgs defines the Stats arguments.
type StatsArgs struct {
	Prefix string `long:"prefix"     required:"true" description:"Prefix of metrics key"`
	Host   string `long:"host"     required:"true" description:"Host of stats server"`
	Port   int    `long:"port"     required:"true" description:"Port of stats server"`
}

// DoomHammerArgs defines the arguments for the shopper api client.
type DoomHammerArgs struct {
	BaseURL string `long:"base-url" required:"true" description:"Doom Hammer API url (i.e.: https://account-doomhammer.p3.r53.nordstrom.net)"`
}

// APMArgs defines the arguments for the APM Client
type APMArgs struct {
	KeyPath string `long:"key-path" required:"true" description:"Path to the client key for APM"`
}

type ShopperTokenArgs struct {
	KeyPath string `long:"key-path" required:"true" description:"Path to the key for generating Shopper Token"`
}

// ShopperAuthArgs defines the arguments for the shopper auth client.
type ShopperAuthArgs struct {
	BaseURL         string `long:"base-url" required:"true" description:"ShopperAuth API url (i.e.: https://ci-auth-mapi.dev.nordstrom.com)"`
	Auth            string `long:"auth" required:"true" description:"ShopperAuth API Auth"`
	AcctAuth        string `long:"acct-auth" required:"true" description:"The ShopperAuth API Acct Auth"`
	VerifyTableName string `long:"verify-table-name" required:"true" description:"dynamo table with verify info"`
}

// ForterAPIArgs defines the arguments for the forter client.
type ForterAPIArgs struct {
	BaseURL            string `long:"base-url" required:"true" description:"Forter API url (i.e.: https://api.forter-secure.com)"`
	SiteID             string `long:"site-id" required:"true" description:"Forter API site id"`
	Key                string `long:"key" required:"true" description:"Forter key/username"`
	Version            string `long:"version" required:"true" description:"Forter API version"`
	MaxIdleConnections string `long:"maxIdle-connections" required:"true" description:"Forter Max Idle Connections"`
}

// ShopperAPIArgs defines the arguments for the shopper api client.
type ShopperAPIArgs struct {
	BaseURL string `long:"base-url" required:"true" description:"ShopperApi API url (i.e.: https://inapi.nonprod.ecomsecure.aws.cloud.nordstrom.net)"`
}

// ShopperAccountAPIArgs defines the arguments for the shopper account api client.
type ShopperAccountAPIArgs struct {
	BaseURL string `long:"base-url" required:"true" description:"ShopperAccountApi API url (i.e.: https://inapi.nonprod.ecomsecure.aws.cloud.nordstrom.net)"`
}

// VerifyArgs defines the arguments for the verify client.
type VerifyArgs struct {
	BaseURL string `long:"base-url" required:"true" description:"Verify API url"`
	IamRole string `long:"iam-role" required:"true"`
}

// ApigeeAPIArgs defines the arguments for the apigee client.
type ApigeeAPIArgs struct {
	BaseURL        string `long:"base-url" required:"true" description:"Apigee API url (i.e.: https://example.apigee.com)"`
	BaseSignoutURL string `long:"base-signout-url" required:"true" description:"Apigee API url (i.e.: https://example.apigee.com)"`
	SignOutAllURL  string `long:"signout-all-url" required:"true" description:"Apigee signout all API url (i.e.: https://example.apigee.com)"`
	APIKey         string `long:"apikey" required:"true" description:"Apigee API key"`
	RevokeKey      string `long:"revoke-apikey" required:"true" description:"Apigee revoke API key"`
}

// MfaBypassArgs defines the arguments for the MFA Bypass client.
type MfaBypassArgs struct {
	PublicKeyURL string `long:"public-key-url" required:"true" description:"MFA Bypass Public Key URL"`
}

type HealthCheck struct {
	FileName string `long:"file-name" required:"true" description:"Filename to use for healthcheck"`
}

// Parse parses command line arguments.
func Parse() CommandlineArgs {
	args := CommandlineArgs{}

	parser := goflags.NewParser(&args, goflags.Default)
	parser.NamespaceDelimiter = "-"
	var err error
	if len(os.Args) == 2 {
		filePath := os.Args[1]
		bArr, err := ioutil.ReadFile(filePath)
		if err != nil {
			fmt.Println("Error reading file: " + filePath)
			fmt.Println(err)
			os.Exit(1)
		}

		lines := strings.Split(string(bArr), "\n")
		_, err = parser.ParseArgs(lines)
	} else {
		_, err = parser.Parse()
	}
	if err != nil {
		fmt.Println("Could not parse args", err.Error())
		os.Exit(1)
	}

	return args
}
