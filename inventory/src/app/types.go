package app

type EnabledAppEntry struct {
	Users []string `yaml:"users" json:"users"`
}
type EnabledApps map[string]EnabledAppEntry

type CommonConfig struct {
	SSHPort           int    `yaml:"ssh_port" json:"ssh_port"`
	SSHKeyPath        string `yaml:"ssh_key_path" json:"ssh_key_path"`
	StartScript       string `yaml:"start_script" json:"start_script"`
	StopScript        string `yaml:"stop_script" json:"stop_script"`
	HealthCheckScript string `yaml:"health_check_script" json:"health_check_script"`
}

type HostEntry struct {
	Hostname  string `yaml:"hostname" json:"hostname"`
	IPAddress string `yaml:"ip_address" json:"ip_address"`
	Username  string `yaml:"username" json:"username"`
	// Only specified if different from common SSHPort
	SSHPort *int `yaml:"ssh_port,omitempty" json:"ssh_port,omitempty"`
}

type AppEnvData struct {
	Common CommonConfig `yaml:"common" json:"common"`
	Hosts  []HostEntry  `yaml:"hosts" json:"hosts"`
}

type AppEnvListResponse struct {
	Environments []EnvNameType `json:"environments"`
}

type AppEnvDataResponse struct {
	AppEnvData *AppEnvData
	Access     []string `json:"access,omitempty"`
}
