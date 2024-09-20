package config

import (
	"fmt"
	"net"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/caarlos0/env/v11"
	"github.com/ex0rcist/gophermart/internal/entities"
	"github.com/ex0rcist/gophermart/internal/utils"
	"github.com/spf13/pflag"
	"golang.org/x/sync/errgroup"
)

type DB struct {
	DSN string `env:"DATABASE_URI"`
}

type Server struct {
	Address string `env:"RUN_ADDRESS"`
	Timeout int
	Secret  entities.Secret `env:"APP_KEY"`
}

type Accrual struct {
	Address        string `env:"ACCRUAL_SYSTEM_ADDRESS"`
	Timeout        int
	RefillInterval time.Duration
}

type Config struct {
	DB      DB
	Server  Server
	Accrual Accrual
}

func Parse() (*Config, error) {
	var err error

	cfg := &Config{}

	// порядок парсинга настроек: дефолтные; ENV; flags
	fns := []func(*Config) (*Config, error){
		NewDefault, ConfigFromEnv, ConfigFromFlags,
	}

	for _, fn := range fns {
		cfg, err = fn(cfg)
		if err != nil {
			return nil, err
		}
	}

	if err = validateConfig(cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}

func NewDefault(_ *Config) (*Config, error) {
	config := &Config{
		DB: DB{},
		Server: Server{
			Address: "0.0.0.0:8282",
			Timeout: 60,
		},
		Accrual: Accrual{
			Address:        "0.0.0.0:8080",
			RefillInterval: 5 * time.Second,
			Timeout:        60,
		},
	}

	return config, nil
}

func ConfigFromFlags(config *Config) (*Config, error) {
	flags := pflag.NewFlagSet(os.Args[0], pflag.ContinueOnError)

	flags.StringVarP(&config.DB.DSN, "database", "d", config.DB.DSN, "PostgreSQL database DSN")
	flags.StringVarP(&config.Accrual.Address, "accrual-address", "r", config.Accrual.Address, "address:port for accrual service")
	flags.StringVarP(&config.Server.Address, "gophermart-address", "a", config.Server.Address, "address:port for HTTP API requests")
	flags.VarP(&config.Server.Secret, "secret", "k", "a key to sign data; will be generated automatically if empty")

	err := flags.Parse(os.Args[1:])
	if err != nil {
		return nil, err
	}

	if len(config.Server.Secret) <= 0 {
		config.Server.Secret = entities.Secret(utils.GenerateRandomString(16))
	}

	return config, err
}

func ConfigFromEnv(config *Config) (*Config, error) {
	if err := env.Parse(config); err != nil {
		return nil, err
	}

	return config, nil
}

func validateConfig(c *Config) error {
	g := &errgroup.Group{}
	g.Go(func() error { return validateAddr(c.Server.Address) })
	g.Go(func() error { return validateAddr(c.Accrual.Address) })
	g.Go(func() error { return validateDSN(c.DB.DSN) })
	return g.Wait()
}

func validateAddr(address string) error {
	if !strings.HasPrefix(address, "http://") && !strings.HasPrefix(address, "https://") {
		address = "http://" + address
	}

	parsedURL, err := url.Parse(address)
	if err != nil {
		return fmt.Errorf("error parsing URL: %w", err)
	}

	// проверяем наличие хоста
	host := parsedURL.Host
	if host == "" {
		return fmt.Errorf("error parsing URL: host is missing")
	}

	// проверяем хост и порт с помощью net.SplitHostPort
	hostOnly, port, err := net.SplitHostPort(host)
	if err != nil {
		return fmt.Errorf("error parsing URL: %w", err)
	}

	// проверяем корректность IP или доменного имени
	if net.ParseIP(hostOnly) == nil {
		if _, err := net.LookupHost(hostOnly); err != nil {
			return fmt.Errorf("error parsing ip: invalid host or IP: %w", err)
		}
	}

	// проверяем корректность порта (если он был указан)
	if port != "" {
		if _, err := net.LookupPort("tcp", port); err != nil {
			return fmt.Errorf("error parsing ip: invalid port: %w", err)
		}
	}

	return nil
}

func validateDSN(dsn string) error {
	if dsn == "" {
		return fmt.Errorf("empty DB DSN")
	}

	return nil
}
