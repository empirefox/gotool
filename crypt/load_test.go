package crypt

import (
	"testing"

	"github.com/mcuadros/go-defaults"
	"github.com/stretchr/testify/assert"
)

type Server struct {
	Cert        []byte `json:"-"`
	Key         []byte `json:"-" xps:"key.pem"`
	CertIgnored string `xps:"cert.pem"` // only accept []byte with xps tag
	Port        uint   `env:"PORT" default:"443"`
}

type PayService struct {
	ClientID string
	Key      []byte `xps:"pay-key.pem"`
}

type Config struct {
	Server     Server
	PayService PayService
}

func (c *Config) GetEnvPtrs() []interface{} {
	return []interface{}{&c.Server}
}

func TestLoadProd(t *testing.T) {
	assert := assert.New(t)

	config := new(Config)
	err := LoadConfig(config, &ConfigOptions{Password: "yourpassword"})
	assert.Nil(err)

	assert.Empty(config.Server.Cert)
	assert.Empty(config.Server.CertIgnored)
	assert.Equal("prod-key", string(config.Server.Key))
	assert.Equal("real_client_id", config.PayService.ClientID)
	assert.Equal("prod-paykey", string(config.PayService.Key))
	defaults.SetDefaults(config)
	assert.Equal(uint(443), config.Server.Port)
}

func TestLoadDev(t *testing.T) {
	assert := assert.New(t)

	config := new(Config)
	err := LoadConfig(config, &ConfigOptions{XpsBootConfig: "xps-config-dev.json"})
	assert.Nil(err)

	assert.Empty(config.Server.Cert)
	assert.Empty(config.Server.CertIgnored)
	assert.Equal("dev-key", string(config.Server.Key))
	assert.Equal(uint(8443), config.Server.Port)
	assert.Equal("fake_client_id", config.PayService.ClientID)
	assert.Equal("dev-paykey", string(config.PayService.Key))
}
