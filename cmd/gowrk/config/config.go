package config

import (
	"errors"
	"flag"
	"io"
	"io/ioutil"
	"net/url"
	"os"
	"strings"
)

// ErrEmptyURL is returned if an empty URL is provided as a command line argument.
var ErrEmptyURL = errors.New("the URL cannot be empty")

//go:generate mockgen -destination=../../../test/mocks/config_mocks.go -package=mocks -source=config.go flagSet
type flagSet interface {
	SetOutput(io.Writer)
	String(string, string, string) *string
	Bool(string, bool, string) *bool
	Uint(string, uint, string) *uint
	Parse([]string) error
	PrintDefaults()
}

// Config represents command line configuration provided by the user.
// An empty config is ready to use in the program, however it can rarely be useful.
type Config struct {
	URL       url.URL
	Rate      uint
	IsVerbose bool
	PoolSize  uint32
	MaxConns  uint32
}

// New returns new Config, populated with value from command line flags.
// Second returned value contains help text on how to use the command line flags.
// If the flags could not be read or are malformed, an error is returned.
func New() (Config, string, error) {
	return NewFromFlagSet(
		flag.NewFlagSet("", flag.ContinueOnError),
		os.Args[1:],
	)
}

// NewFromFlagSet returns new Config initialised from provided flag set and arguments.
func NewFromFlagSet(flags flagSet, args []string) (Config, string, error) {
	var config Config

	flags.SetOutput(ioutil.Discard)

	urlFlag := flags.String("url", "", "the url which to send requests to")
	verboseFlag := flags.Bool("verbose", false, "enable verbose output")
	rateFlag := flags.Uint("rate", 10, "the rate at which requests are sent")
	poolSizeFlag := flags.Uint("poolsize", 1, "the amount of http clients in the pool")
	maxConnsFlag := flags.Uint("maxconns", 1, "the amount of idle connections per client")

	if err := flags.Parse(args); err != nil {
		return config, getUsage(flags), err
	}

	url, err := parseURL(*urlFlag)
	if err != nil {
		return config, getUsage(flags), err
	}

	config.Rate = *rateFlag
	config.IsVerbose = *verboseFlag
	config.PoolSize = uint32(*poolSizeFlag)
	config.MaxConns = uint32(*maxConnsFlag)
	config.URL = *url

	return config, "", nil
}

func getUsage(flags flagSet) string {
	var builder strings.Builder

	flags.SetOutput(&builder)
	flags.PrintDefaults()

	return builder.String()
}

func parseURL(urlString string) (*url.URL, error) {
	// url.Parse does not return an error if an empty string is passed,
	// so the URL string has to be validated manually.
	if urlString == "" {
		return nil, ErrEmptyURL
	}

	return url.Parse(urlString)
}
