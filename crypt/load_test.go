package crypt

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type Server struct {
	Cert        []byte
	Key         []byte `xps:"key.pem"`
	CertIgnored string `xps:"cert.pem"` // only accept []byte with xps tag
	Port        uint   `default:"443"`
}

type PayService struct {
	ClientID string
	Key      []byte `xps:"pay-key.pem"`
}

type Config struct {
	Server     Server
	PayService PayService
}

func TestLoadProd(t *testing.T) {
	assert := assert.New(t)

	config := new(Config)
	err := LoadSingleConfigWithOptions(config, &ConfigOptions{Password: "yourpassword"})
	assert.Nil(err)

	assert.Empty(config.Server.Cert)
	assert.Empty(config.Server.CertIgnored)
	assert.Equal("prod-key", string(config.Server.Key))
	assert.Equal(uint(443), config.Server.Port)
	assert.Equal("real_client_id", config.PayService.ClientID)
	assert.Equal("prod-paykey", string(config.PayService.Key))
}

func TestLoadDev(t *testing.T) {
	assert := assert.New(t)

	config := new(Config)
	err := LoadSingleConfigWithOptions(config, &ConfigOptions{Mode: "dev"})
	assert.Nil(err)

	assert.Empty(config.Server.Cert)
	assert.Empty(config.Server.CertIgnored)
	assert.Equal("dev-key", string(config.Server.Key))
	assert.Equal(uint(8443), config.Server.Port)
	assert.Equal("fake_client_id", config.PayService.ClientID)
	assert.Equal("dev-paykey", string(config.PayService.Key))
}
