package types

import (
	"context"

	"gopkg.in/yaml.v3"
)

type Command struct {
	Name  *string `yaml:"name"`
	Async bool    `yaml:"async"`
	Run   string  `yaml:"run"`
}

type Scripts struct {
	MacosArm   *Command `yaml:"macos_arm"`
	MacosIntel *Command `yaml:"macos_intel"`
	Arch       *Command `yaml:"arch"`
	Windows    *Command `yaml:"windows"`
}

type ConfigType string

const (
	ConfigTypeLink   = "link"
	ConfigTypeScript = "script"
)

type Config interface {
	GetConfigType() ConfigType
	GetBeforeScript() bool
}

type ConfigLink struct {
	Type         ConfigType `yaml:"type"`
	Source       string     `yaml:"source"`
	Target       string     `yaml:"target"`
	BeforeScript bool       `yaml:"before_script"`
}

func (c ConfigLink) GetConfigType() ConfigType { return c.Type }
func (c ConfigLink) GetBeforeScript() bool     { return c.BeforeScript }

type ConfigScript struct {
	Type         ConfigType `yaml:"type"`
	Platform     string     `yaml:"platform"`
	Script       string     `yaml:"script"`
	BeforeScript bool       `yaml:"before_script"`
}

func (c ConfigScript) GetConfigType() ConfigType { return c.Type }
func (c ConfigScript) GetBeforeScript() bool     { return c.BeforeScript }

type ExecutionStatus int

const (
	ExecutionStatusWaiting ExecutionStatus = iota
	ExecutionStatusInProgress
	ExecutionStatusCompleted
	ExecutionStatusFailed
	ExecutionStatusSkiped
	ExecutionStatusInstalled
	ExecutionStatusConfiguring
)

type Channels struct {
	globalLog chan string
	toolLog   chan string
	status    chan ExecutionStatus
}

type InstallScriptDefinition struct {
	Name        string
	Description string
	Stage       int
	Config      []Config
	Scripts     Scripts
	ToolPath    string
	channels    Channels
}

func (c *InstallScriptDefinition) UnmarshalYAML(value *yaml.Node) error {
	var tmpRoot struct {
		Name        string      `yaml:"name"`
		Description string      `yaml:"description"`
		Config      []yaml.Node `yaml:"config"`
		Stage       int         `yaml:"stage"`
		Scripts     Scripts     `yaml:"scripts"`
	}
	if err := value.Decode(&tmpRoot); err != nil {
		return err
	}

	configs := make([]Config, 0)
	for _, config := range tmpRoot.Config {

		var tmpConfig struct {
			Type ConfigType `yaml:"type"`
		}
		if err := config.Decode(&tmpConfig); err != nil {
			return err
		}

		switch tmpConfig.Type {
		case ConfigTypeLink:
			var configLink ConfigLink
			if err := config.Decode(&configLink); err != nil {
				return err
			}
			configs = append(configs, configLink)
		case ConfigTypeScript:
			var configScript ConfigScript
			if err := config.Decode(&configScript); err != nil {
				return err
			}
			configs = append(configs, configScript)
		}
	}

	c.Name = tmpRoot.Name
	c.Description = tmpRoot.Description
	c.Stage = tmpRoot.Stage
	c.Config = configs
	c.Scripts = tmpRoot.Scripts

	return nil
}

func (s *InstallScriptDefinition) InitChannels(ctx context.Context) {
	s.channels.globalLog = make(chan string, 100000)
	s.channels.toolLog = make(chan string, 100000)
	s.channels.status = make(chan ExecutionStatus, 1)
	go func() {
		<-ctx.Done()
		close(s.channels.globalLog)
		close(s.channels.toolLog)
		close(s.channels.status)
	}()
}

func (s *InstallScriptDefinition) GetChannels() (chan string, chan ExecutionStatus) {
	return s.channels.toolLog, s.channels.status
}

type ExecutionGrouping map[int][]*string
