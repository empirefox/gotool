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
)

type Gorm struct {
	Dialect string
	Url     string
	MaxIdle int
	MaxOpen int
}

func init() {
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
	return "0.0.0.0:" + GetEnv("PORT", "9999")
}

// TODO test with custom domain
func GetSubDomain() string {
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
	return "127.0.0.1:" + GetEnv("PORT", "9999")
}

func GetGorm() Gorm {
	switch PaasVendor {
	case HEROKU:
		return Gorm{
			Dialect: "postgres",
			Url:     os.Getenv("DATABASE_URL"),
			MaxIdle: 20,
			MaxOpen: 20,
		}
	case CLOUD_CONTROL:
		return Gorm{
			Dialect: "postgres",
			Url:     os.Getenv("ELEPHANTSQL_URL"),
			MaxIdle: 5,
			MaxOpen: 5,
		}
	case OPENSHIFT:
		// Vendor default value is 100
		// Can be set by OPENSHIFT_POSTGRESQL_MAX_CONNECTIONS
		// We only use 20 for now
		return Gorm{
			Dialect: "postgres",
			Url:     os.Getenv("OPENSHIFT_POSTGRESQL_DB_URL") + "/" + os.Getenv("OPENSHIFT_APP_NAME"),
			MaxIdle: 20,
			MaxOpen: 20,
		}
	}
	// must be test mode
	return Gorm{
		Dialect: "postgres",
		Url:     os.Getenv("DB_URL"),
		MaxIdle: 5,
		MaxOpen: 5,
	}
}
