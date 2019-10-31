package config

import pluginsConfig "github.com/lyft/flyteplugins/go/tasks/config"

//go:generate pflags Config

// Per Plugin optional configuration
type Config struct {
	// TODO Add any configuration parameters here. For example, if region should be preset? or if there are some common parameters that should be passed
	// for every execution. If this is not the case, just delete the config folder
	MyConfigParam string `json:"my-cfg-param" pflag:",Some configuration that is used by the backend plugin "`
}

var (
	// Default values can be provided directly when registering the section
	section = pluginsConfig.MustRegisterSubSection("sagemaker", &Config{MyConfigParam: "default-value"})
)

func Get() *Config {
	return section.GetConfig().(*Config)
}
