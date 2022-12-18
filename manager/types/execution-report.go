package types

type ExecutionReport struct {
	Name             string          `yaml:"name"`
	Status           ExecutionStatus `yaml:"status"`
	Version          *string         `yaml:"version"`
	InstallationPath *string         `yaml:"installation_path"`
}
