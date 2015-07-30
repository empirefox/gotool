package paas

import (
	"flag"
	"fmt"
	"os"
	"strings"
)

const (
	HEROKU        = "heroku"
	OPENSHIFT     = "openshift"
	CLOUD_CONTROL = "cloudControl"
)

var (
	PaasVendor string
	IsDevMode  bool
	BindAddr   string
	SubDomain  string
	PortInTest string
)

type Gorm struct {
	Dialect string
	Url     string
	MaxIdle int
	MaxOpen int
}

func init() {
	PortInTest = GetEnv("PORT", "9999")
	PaasVendor = GetPaasVendor()
	IsDevMode = CheckIsDev()
	BindAddr = GetBindAddr()
	SubDomain = GetSubDomain()
	if IsDevMode {
		flag.Set("stderrthreshold", "INFO")
	}
}

func IsSystemMode() bool {
	switch strings.ToLower(os.Getenv("SYSTEM_MODE_ON")) {
	case "true", "yes", "1", "on", "ok":
		return true
	}
	return false
}

func GetPaasVendor() string {
	if os.Getenv("DYNO") != "" {
		return HEROKU
	}
	if os.Getenv("OPENSHIFT_APP_NAME") != "" {
		return OPENSHIFT
	}
	switch os.Getenv("PAAS_VENDOR") {
	case CLOUD_CONTROL:
		return CLOUD_CONTROL
	}
	return ""
}

func CheckIsDev() bool {
	return PaasVendor == ""
}

func GetEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value != "" {
		return value
	}
	return defaultValue
}

// access from client
func GetWsPorts() (ws, wss string) {
	switch PaasVendor {
	case HEROKU, CLOUD_CONTROL:
		return "80", "443"
	case OPENSHIFT:
		return "8000", "8443"
	}
	// must be test mode
	return PortInTest, PortInTest
}

func GetBindAddr() string {
	// all copy from the official examples
	switch PaasVendor {
	case HEROKU:
		return fmt.Sprintf(":%v", os.Getenv("PORT"))
	case CLOUD_CONTROL:
		return fmt.Sprintf(":%v", GetEnv("PORT", "8080"))
	case OPENSHIFT:
		return fmt.Sprintf("%s:%s", os.Getenv("HOST"), os.Getenv("PORT"))
	}
	// must be test mode
	return "0.0.0.0:" + PortInTest
}

// TODO test with custom domain
func GetSubDomain() string {
	if sd := os.Getenv("PUBLIC_DOMAIN"); sd != "" {
		return sd
	}
	switch PaasVendor {
	case HEROKU:
		// not supported
		return ""
	case CLOUD_CONTROL:
		// use websocket domain
		app := strings.Split(os.Getenv("DEP_NAME"), "/")[0] + "."
		domain := os.Getenv("DOMAIN")
		if domain == "cloudcontrolled.com" {
			domain = "cloudcontrolapp.com"
		}
		return app + domain
	case OPENSHIFT:
		// use default domain
		return os.Getenv("OPENSHIFT_APP_DNS")
	}
	// must be test mode
	return "127.0.0.1"
}

func GetGorm() Gorm {
	if url := os.Getenv("DB_URL"); url != "" {
		// useful in test
		return Gorm{
			Dialect: "postgres",
			Url:     url,
			MaxIdle: 5,
			MaxOpen: 5,
		}
	}
	switch PaasVendor {
	case HEROKU:
		return Gorm{
			Dialect: "postgres",
			Url:     os.Getenv("DATABASE_URL"),
			MaxIdle: 20,
			MaxOpen: 20,
		}
	case CLOUD_CONTROL:
		if url := os.Getenv("ELEPHANTSQL_URL"); url != "" {
			return Gorm{
				Dialect: "postgres",
				Url:     url,
				MaxIdle: 5,
				MaxOpen: 5,
			}
		}
		return Gorm{
			Dialect: "mysql",
			Url: fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8",
				os.Getenv("MYSQLS_USERNAME"),
				os.Getenv("MYSQLS_PASSWORD"),
				os.Getenv("MYSQLS_HOSTNAME"),
				os.Getenv("MYSQLS_PORT"),
				os.Getenv("MYSQLS_DATABASE")),
			MaxIdle: 2,
			MaxOpen: 2,
		}
	case OPENSHIFT:
		// Vendor default value is 100
		// Can be set by OPENSHIFT_POSTGRESQL_MAX_CONNECTIONS
		// We only use 20 for now
		return Gorm{
			Dialect: "postgres",
			Url: fmt.Sprintf("%s/%s?sslmode=disable",
				os.Getenv("OPENSHIFT_POSTGRESQL_DB_URL"),
				os.Getenv("OPENSHIFT_APP_NAME")),
			MaxIdle: 20,
			MaxOpen: 20,
		}
	}
	panic("db param must be set, like DB_URL")
}
