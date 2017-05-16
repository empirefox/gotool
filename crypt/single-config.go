package crypt

import (
	"reflect"

	"github.com/caarlos0/env"
	"github.com/mcuadros/go-defaults"

	"gopkg.in/go-playground/validator.v9"
)

type EnvLoadable interface {
	GetEnvPtrs() []interface{}
}

type Validable interface {
	Validate(v interface{}) error
}

type ConfigOptions struct {
	XpsBootConfig       string `env:"XPS_BOOT_CONFIG" default:"xps-config.json"`
	XpsBootConfigDecode string `env:"XPS_BOOT_CONFIG_DECODE"` // support json ymal toml json5
	Password            string `env:"XPS_PASSWORD"`           // overwrite xps.Password if set
	XpsFile             string `env:"XPS_TARBALL"`            // overwrite xps.XpsFile if set
	ConfigFile          string `env:"XPS_APP_CONFIG"`         // overwrite xps.ConfigFile if set
	EquipTag            string `env:"XPS_TAG"`                // overwrite xps.EquipTag if set
}

func LoadConfig(config interface{}, opts *ConfigOptions) (err error) {
	if opts == nil {
		opts = new(ConfigOptions)
		if err = env.Parse(opts); err != nil {
			return err
		}
	}

	err = LoadXps(config, opts)
	if err != nil {
		return err
	}

	// overwrite with env
	if envloader, ok := config.(EnvLoadable); ok {
		for _, s := range envloader.GetEnvPtrs() {
			err = env.Parse(s)
			if err != nil {
				return err
			}
		}
	}

	defaults.SetDefaults(config)

	var v func(v interface{}) error
	if validable, ok := config.(Validable); ok {
		v = validable.Validate
	} else {
		v = defaultValidate()
	}

	if err = v(config); err != nil {
		return err
	}

	return nil
}

func LoadXps(config interface{}, opts *ConfigOptions) error {
	defaults.SetDefaults(opts)

	filetype := DetectFileType(opts.XpsBootConfigDecode, opts.XpsBootConfig)
	xps, err := NewXps(opts.XpsBootConfig, filetype)
	if err != nil {
		return err
	}

	if opts.XpsFile != "" {
		xps.XpsFile = opts.XpsFile
	}
	if opts.ConfigFile != "" {
		xps.ConfigFile = opts.ConfigFile
	}
	if opts.EquipTag != "" {
		xps.EquipTag = opts.EquipTag
	}

	files, err := xps.DecryptXhexFile(opts.Password)
	if err != nil {
		return err
	}

	filetype = DetectFileType("", xps.ConfigFile)
	if err = UnmarshalFormat(files[xps.ConfigFile], config, filetype); err != nil {
		return err
	}
	NewFiles(files, xps.EquipTag).Equip(config)
	return nil
}

func defaultValidate() func(v interface{}) error {
	validate := validator.New()
	validateRequireField := func(fl validator.FieldLevel) bool {
		field := fl.Field()
		if validate.Var(field.Interface(), "required") != nil {
			return true
		}

		rqField, rqKind, ok := fl.GetStructFieldOK()
		if !ok {
			return false
		}

		dep := rqField.Interface()
		if rqKind == reflect.Slice && rqField.Len() == 0 {
			return false
		}

		if validate.Var(dep, "required") != nil {
			return false
		}
		return true
	}
	validate.RegisterValidation("dep", validateRequireField)
	validate.RegisterAlias("zap_level", "len=0|eq=debug|eq=info|eq=warn|eq=error|eq=dpanic|eq=panic|eq=fatal")
	return validate.Struct
}
