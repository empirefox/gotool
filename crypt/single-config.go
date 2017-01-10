package crypt

import (
	"encoding/json"
	"fmt"

	"github.com/mcuadros/go-defaults"
)

type ConfigOptions struct {
	XpsConfigFileBase string `default:"xps-config"` // xps-config[-dev].json out of xps file
	Mode              string // dev prod
	Password          string
	XpsFile           string
	ConfigFile        string // in xps file
	EquipTag          string
}

func LoadSingleConfig(config interface{}) error {
	return LoadSingleConfigWithOptions(config, new(ConfigOptions))
}

func LoadSingleConfigWithOptions(config interface{}, opts *ConfigOptions) error {
	defaults.SetDefaults(opts)
	xpsConfigFile := opts.XpsConfigFileBase + ".json"
	if opts.Mode != "" {
		xpsConfigFile = fmt.Sprintf("%s-%s.json", opts.XpsConfigFileBase, opts.Mode)
	}

	xps, err := NewXps(xpsConfigFile)
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

	if err = json.Unmarshal(files[xps.ConfigFile], config); err != nil {
		return err
	}
	NewFiles(files, xps.EquipTag).Equip(config)
	defaults.SetDefaults(config)
	return nil
}
