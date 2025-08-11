package config

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

var (
	dirs = []string{".", "config"}
)

type Config struct {
	name      string
	dirs      []string
	bindEnvKV map[string][]string
	defaultKV map[string]interface{}
}

func NewConfigHandler(name string, dirs []string, bindEnvKV map[string][]string, defaultKV map[string]interface{}) *Config {
	return &Config{
		name:      name,
		dirs:      dirs,
		bindEnvKV: bindEnvKV,
		defaultKV: defaultKV,
	}
}

func NewConfigHandlerWithDefaults(name string) *Config {
	return &Config{
		name:      name,
		dirs:      dirs,
		bindEnvKV: map[string][]string{},
		defaultKV: map[string]interface{}{},
	}
}

func (c *Config) SetDirs(dirs []string) {
	c.dirs = dirs
}

func (c *Config) SetBindEnv(bindEnvKV map[string][]string) {
	c.bindEnvKV = bindEnvKV
}

func (c *Config) SetDefault(defaultKV map[string]interface{}) {
	c.defaultKV = defaultKV
}

func (c *Config) BindEnv(key string, envKeys ...string) {
	c.bindEnvKV[key] = envKeys
}

func (c *Config) DefaultValue(key string, value interface{}) {
	c.defaultKV[key] = value
}

func (c *Config) LoadConfig(config interface{}) error {
	if c.name == "" {
		return fmt.Errorf("config name cannot be empty")
	}

	viper.SetEnvPrefix("")
	viper.AutomaticEnv()
	replacer := strings.NewReplacer(".", "_") // 支持 app.port -> APP_PORT
	viper.SetEnvKeyReplacer(replacer)

	for k, v := range c.defaultKV {
		viper.SetDefault(k, v)
	}

	for k, vals := range c.bindEnvKV {
		bindArr := append([]string{k}, vals...)
		err := viper.BindEnv(bindArr...)
		if err != nil {
			return fmt.Errorf("failed to bind environment variable %s: %w", k, err)
		}
	}

	// Read config file
	viper.SetConfigName(c.name)
	viper.SetConfigType("yaml")
	for _, dir := range c.dirs {
		viper.AddConfigPath(dir)
	}
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return fmt.Errorf("failed to read config file: %w", err)
		}
	}

	if err := viper.Unmarshal(&config); err != nil {
		return fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return nil
}

var (
	GetBool                 = viper.GetBool
	GetString               = viper.GetString
	GetInt                  = viper.GetInt
	GetFloat64              = viper.GetFloat64
	GetStringSlice          = viper.GetStringSlice
	GetStringMap            = viper.GetStringMap
	GetStringMapString      = viper.GetStringMapString
	GetStringMapStringSlice = viper.GetStringMapStringSlice
	GetDuration             = viper.GetDuration
	GetTime                 = viper.GetTime
	GetIntSlice             = viper.GetIntSlice
	GetInt64                = viper.GetInt64
	GetInt32                = viper.GetInt32
	GetUint                 = viper.GetUint
	GetUint64               = viper.GetUint64
	GetUint32               = viper.GetUint32
	GetUint16               = viper.GetUint16
	GetUint8                = viper.GetUint8

	GetSizeInBytes = viper.GetSizeInBytes

	GetViper  = viper.GetViper
	Get       = viper.Get
	Unmarshal = viper.UnmarshalKey
)
