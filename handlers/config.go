package handlers

import (
	"flag"
	"fmt"
)

type BaseConfig struct {
	Address string
	Port    int
	Folder  string
}

type RecorderConfig struct {
	BaseConfig
	CreateNewFolder bool
	Prevent304      bool
	LogAsFocused    bool
	PrintOnly       bool
}

type SeekerConfig struct {
	BaseConfig
}

type ScenarioRecorderConfig struct {
	BaseConfig
	CreateNewFolder bool
	Prevent304      bool
}

type ScenarioSeekerConfig struct {
	BaseConfig
}

// BaseConfig
func (cfg *BaseConfig) AddressAndPort() string {
	return fmt.Sprintf("%s:%d", cfg.Address, cfg.Port)
}

func (cfg *BaseConfig) String() string {
	return fmt.Sprintf("Address=%s\nFolder=%s", cfg.AddressAndPort(), cfg.Folder)
}

func (cfg *BaseConfig) InitFlags(flags *flag.FlagSet, defaultFolder string) {
	flags.StringVar(&cfg.Address, "address", "127.0.0.1", "The address on which to run the proxy server")
	flags.IntVar(&cfg.Port, "port", 8123, "The port on which to run the proxy server")
	flags.StringVar(&cfg.Folder, "folder", defaultFolder, "The root data folder")
}

// RecorderConfig
func NewRecorderConfig(flags *flag.FlagSet, args []string, defaultFolder string) *RecorderConfig {
	cfg := &RecorderConfig{}
	cfg.InitFlags(flags, defaultFolder)

	if err := flags.Parse(args); err != nil {
		return nil
	}

	return cfg
}

func (cfg *RecorderConfig) String() string {
	return fmt.Sprintf("%s\nCreateNewFolder=%t\nPrevent304=%t\nLogAsFocused=%t\nPrintOnly=%t", cfg.BaseConfig.String(), cfg.CreateNewFolder, cfg.Prevent304, cfg.LogAsFocused, cfg.PrintOnly)
}

func (cfg *RecorderConfig) InitFlags(flags *flag.FlagSet, defaultFolder string) {
	cfg.BaseConfig.InitFlags(flags, defaultFolder)

	flags.BoolVar(&cfg.CreateNewFolder, "new_folder", true, "Create new folder inside root for log")
	flags.BoolVar(&cfg.Prevent304, "prevent_304", true, "Prevent 304 http answer")
	flags.BoolVar(&cfg.LogAsFocused, "focused", false, "Log all requests as focused")
	flags.BoolVar(&cfg.PrintOnly, "print_only", false, "Only print requests, no logs")
}

// SeekerConfig
func NewSeekerConfig(flags *flag.FlagSet, args []string, defaultFolder string) *SeekerConfig {
	cfg := &SeekerConfig{}
	cfg.InitFlags(flags, defaultFolder)

	if err := flags.Parse(args); err != nil {
		return nil
	}

	return cfg
}

func (cfg *SeekerConfig) String() string {
	return cfg.BaseConfig.String()
}

func (cfg *SeekerConfig) InitFlags(flags *flag.FlagSet, defaultFolder string) {
	cfg.BaseConfig.InitFlags(flags, defaultFolder)
}

// ScenarioRecorderConfig
func NewScenarioRecorderConfig(flags *flag.FlagSet, args []string, defaultFolder string) *ScenarioRecorderConfig {
	cfg := &ScenarioRecorderConfig{}
	cfg.InitFlags(flags, defaultFolder)

	if err := flags.Parse(args); err != nil {
		return nil
	}

	return cfg
}

func (cfg *ScenarioRecorderConfig) String() string {
	return fmt.Sprintf("%s\nCreateNewFolder=%t\nPrevent304=%t", cfg.BaseConfig.String(), cfg.CreateNewFolder, cfg.Prevent304)
}

func (cfg *ScenarioRecorderConfig) InitFlags(flags *flag.FlagSet, defaultFolder string) {
	cfg.BaseConfig.InitFlags(flags, defaultFolder)

	flags.BoolVar(&cfg.CreateNewFolder, "new_folder", false, "Create new folder inside root for log")
	flags.BoolVar(&cfg.Prevent304, "prevent_304", true, "Prevent 304 http answer")
}

// ScenarioSeekerConfig
func NewScenarioSeekerConfig(flags *flag.FlagSet, args []string, defaultFolder string) *ScenarioSeekerConfig {
	cfg := &ScenarioSeekerConfig{}
	cfg.InitFlags(flags, defaultFolder)

	if err := flags.Parse(args); err != nil {
		return nil
	}

	return cfg
}

func (cfg *ScenarioSeekerConfig) String() string {
	return cfg.BaseConfig.String()
}

func (cfg *ScenarioSeekerConfig) InitFlags(flags *flag.FlagSet, defaultFolder string) {
	cfg.BaseConfig.InitFlags(flags, defaultFolder)
}
