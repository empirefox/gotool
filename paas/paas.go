package paas

import (
	"flag"
	"fmt"
	"os"
	"regexp"
	"strings"
)

const (
	HEROKU         = "heroku"
	OPENSHIFT      = "openshift"
	CLOUD_CONTROL  = "cloudControl"
	CLOUD_AND_HEAT = "cloudandheat"
	BLUEMIX        = "bluemix"
	DAOCLOUD       = "daocloud"
)

var (
	Vendor     string
	IsDevMode  bool
	BindAddr   string
	PortInTest string
	Gorm       GormParams
	Info       ApiInfo
)

type ApiInfo struct {
	IsDevMode  bool
	HttpDomain string
	WsDomain   string
	WssDomain  string
}

type info struct {
	ApiInfo
	Vendor     string
	BindAddr   string
	GormParams GormParams
}

type GormParams struct {
	Dialect string
	Url     string
	MaxIdle int
	MaxOpen int
}

func init() {
	PortInTest = GetEnv("PORT", "9999")
	i := GetPaasInfo()
	Vendor = i.Vendor
	BindAddr = i.BindAddr
	Gorm = i.GormParams
	Info = i.ApiInfo
	if i.Vendor == "" {
		IsDevMode = true
		Info.IsDevMode = true
		flag.Set("stderrthreshold", "INFO")
	}
}

func GetPaasInfo() info {
	switch {
	case os.Getenv("PAAS_VENDOR") == CLOUD_CONTROL:
		// also fit for dotcloud
		return getCloudControl()
	case os.Getenv("DYNO") != "":
		return getHeroku()
	case os.Getenv("OPENSHIFT_APP_NAME") != "":
		return getOpenshift()
	case strings.Contains(os.Getenv("VCAP_APPLICATION"), ".mybluemix.net"):
		return getBluemix()
	}

	switch os.Getenv("DOCKER_VENDOR") {
	case DAOCLOUD:
		return getDaocloud()
	}
	return getTest()
}

func getDaocloud() info {
	domain := os.Getenv("DEFAULT_DOMAIN")
	var gps GormParams
	if os.Getenv("POSTGRESQL_INSTANCE_NAME") != "" {
		gps = GormParams{
			Dialect: "postgres",
			Url: fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
				os.Getenv("POSTGRESQL_USERNAME"),
				os.Getenv("POSTGRESQL_PASSWORD"),
				os.Getenv("POSTGRESQL_PORT_5432_TCP_ADDR"),
				os.Getenv("POSTGRESQL_PORT_5432_TCP_PORT"),
				os.Getenv("POSTGRESQL_INSTANCE_NAME"),
			),
			MaxIdle: 10,
			MaxOpen: 10,
		}
	} else {
		gps = GormParams{
			Dialect: "mysql",
			Url: fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True",
				os.Getenv("MYSQL_USERNAME"),
				os.Getenv("MYSQL_PASSWORD"),
				os.Getenv("MYSQL_PORT_3306_TCP_ADDR"),
				os.Getenv("MYSQL_PORT_3306_TCP_PORT"),
				os.Getenv("MYSQL_INSTANCE_NAME"),
			),
			MaxIdle: 10,
			MaxOpen: 10,
		}
	}
	return info{
		Vendor:   DAOCLOUD,
		BindAddr: fmt.Sprintf(":%v", GetEnv("PORT", "8080")),
		ApiInfo: ApiInfo{
			HttpDomain: domain,
			WsDomain:   domain,
			WssDomain:  domain,
		},
		GormParams: gps,
	}
}

func getHeroku() info {
	domain := os.Getenv("DEFAULT_DOMAIN")
	return info{
		Vendor:   HEROKU,
		BindAddr: fmt.Sprintf(":%v", os.Getenv("PORT")),
		ApiInfo: ApiInfo{
			HttpDomain: domain,
			WsDomain:   domain,
			WssDomain:  domain,
		},
		GormParams: GormParams{
			Dialect: "postgres",
			Url:     os.Getenv("DATABASE_URL"),
			MaxIdle: 19,
			MaxOpen: 19,
		},
	}
}

func getOpenshift() info {
	return info{
		Vendor:   OPENSHIFT,
		BindAddr: fmt.Sprintf("%s:%s", os.Getenv("OPENSHIFT_GO_IP"), os.Getenv("OPENSHIFT_GO_PORT")),
		ApiInfo: ApiInfo{
			HttpDomain: os.Getenv("OPENSHIFT_APP_DNS"),
			WsDomain:   os.Getenv("OPENSHIFT_APP_DNS") + ":8000",
			WssDomain:  os.Getenv("OPENSHIFT_APP_DNS") + ":8443",
		},
		// Vendor default value is 100
		// Can be set by OPENSHIFT_POSTGRESQL_MAX_CONNECTIONS
		// We only use 20 for now
		GormParams: GormParams{
			Dialect: "postgres",
			Url: fmt.Sprintf("%s/%s?sslmode=disable",
				strings.TrimRight(os.Getenv("OPENSHIFT_POSTGRESQL_DB_URL"), "/"),
				os.Getenv("OPENSHIFT_APP_NAME")),
			MaxIdle: 20,
			MaxOpen: 20,
		},
	}
}

func getCloudControl() info {
	app := strings.Split(os.Getenv("DEP_NAME"), "/")[0] + "."
	// gorm
	g := GormParams{
		Dialect: "postgres",
		Url:     os.Getenv("ELEPHANTSQL_URL"),
		MaxIdle: 4,
		MaxOpen: 4,
	}
	if g.Url == "" {
		g = GormParams{
			Dialect: "mysql",
			Url: fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8",
				os.Getenv("MYSQLS_USERNAME"),
				os.Getenv("MYSQLS_PASSWORD"),
				os.Getenv("MYSQLS_HOSTNAME"),
				os.Getenv("MYSQLS_PORT"),
				os.Getenv("MYSQLS_DATABASE")),
			MaxIdle: 1,
			MaxOpen: 1,
		}
	}
	var apiInfo ApiInfo
	if os.Getenv("DOMAIN") == "dotcloudapp.com" {
		apiInfo = ApiInfo{
			HttpDomain: app + "dotcloudapp.com",
			WsDomain:   app + "dotcloudapp.com",
			WssDomain:  app + "dotcloudapp.com",
		}
	} else {
		apiInfo = ApiInfo{
			HttpDomain: app + "cloudcontrolled.com",
			WsDomain:   app + "cloudcontrolapp.com",
			WssDomain:  app + "cloudcontrolapp.com",
		}
	}
	return info{
		Vendor:   CLOUD_CONTROL,
		BindAddr: fmt.Sprintf(":%v", GetEnv("PORT", "8080")),
		ApiInfo:  apiInfo,
		// first chech Elephant
		GormParams: g,
	}
}

func getCloudandheat() info {
	app := strings.Split(os.Getenv("DEP_NAME"), "/")[0] + "."
	return info{
		Vendor:   CLOUD_AND_HEAT,
		BindAddr: fmt.Sprintf(":%v", GetEnv("PORT", "8080")),
		ApiInfo: ApiInfo{
			HttpDomain: app + "cnh-apps.com",
			WsDomain:   app + "cnh-faster-apps.com",
			WssDomain:  app + "cnh-faster-apps.com",
		},
		// first chech Elephant
		GormParams: GormParams{
			Dialect: "mysql",
			Url: fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8",
				os.Getenv("MYSQLS_USERNAME"),
				os.Getenv("MYSQLS_PASSWORD"),
				os.Getenv("MYSQLS_HOSTNAME"),
				os.Getenv("MYSQLS_PORT"),
				os.Getenv("MYSQLS_DATABASE")),
			MaxIdle: 9,
			MaxOpen: 9,
		},
	}
}

var domainRegexp = regexp.MustCompile(`\"application_uris\"\:\[\"([^\"\s]+)\"\]`)

func getBluemix() info {
	domain := domainRegexp.FindStringSubmatch(os.Getenv("VCAP_APPLICATION"))[1]
	return info{
		Vendor:   BLUEMIX,
		BindAddr: fmt.Sprintf("%v:%v", os.Getenv("VCAP_APP_HOST"), os.Getenv("VCAP_APP_PORT")),
		ApiInfo: ApiInfo{
			HttpDomain: domain,
			WsDomain:   domain,
			WssDomain:  domain,
		},
		GormParams: GormParams{
			Dialect: "postgres",
			Url:     os.Getenv("DATABASE_URL") + "?sslmode=disable",
			MaxIdle: 19,
			MaxOpen: 19,
		},
	}
}

func getTest() info {
	url := os.Getenv("DB_URL")
	return info{
		Vendor:   "",
		BindAddr: ":" + PortInTest,
		ApiInfo: ApiInfo{
			HttpDomain: "127.0.0.1:" + PortInTest,
			WsDomain:   "127.0.0.1:" + PortInTest,
			WssDomain:  "127.0.0.1:" + PortInTest,
		},
		GormParams: GormParams{
			Dialect: GetEnv("DB_DIALECT", "postgres"),
			Url:     url,
			MaxIdle: 10,
			MaxOpen: 10,
		},
	}
}

func IsSystemMode() bool {
	switch strings.ToLower(os.Getenv("SYSTEM_MODE_ON")) {
	case "true", "yes", "1", "on", "ok":
		return true
	}
	return false
}

func GetEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value != "" {
		return value
	}
	return defaultValue
}
