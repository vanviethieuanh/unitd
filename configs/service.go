package configs

type ServiceBlock struct {
	Type          string   `hcl:"type,optional" systemd:"Type"`
	ExecStart     string   `hcl:"exec_start,optional" systemd:"ExecStart"`
	ExecStartPre  []string `hcl:"exec_start_pre,optional" systemd:"ExecStartPre"`
	ExecStartPost []string `hcl:"exec_start_post,optional" systemd:"ExecStartPost"`
	ExecStop      []string `hcl:"exec_stop,optional" systemd:"ExecStop"`
	ExecStopPost  []string `hcl:"exec_stop_post,optional" systemd:"ExecStopPost"`
	ExecReload    []string `hcl:"exec_reload,optional" systemd:"ExecReload"`

	Restart                  string   `hcl:"restart,optional" systemd:"Restart"`
	RestartSec               string   `hcl:"restart_sec,optional" systemd:"RestartSec"`
	SuccessExitStatus        []string `hcl:"success_exit_status,optional" systemd:"SuccessExitStatus"`
	RestartPreventExitStatus []string `hcl:"restart_prevent_exit_status,optional" systemd:"RestartPreventExitStatus"`
	RestartForceExitStatus   []string `hcl:"restart_force_exit_status,optional" systemd:"RestartForceExitStatus"`

	TimeoutStartSec string `hcl:"timeout_start_sec,optional" systemd:"TimeoutStartSec"`
	TimeoutStopSec  string `hcl:"timeout_stop_sec,optional" systemd:"TimeoutStopSec"`
	WatchdogSec     string `hcl:"watchdog_sec,optional" systemd:"WatchdogSec"`

	RemainAfterExit bool `hcl:"remain_after_exit,optional" systemd:"RemainAfterExit"`
	GuessMainPID    bool `hcl:"guess_main_pid,optional" systemd:"GuessMainPID"`

	User         string            `hcl:"user,optional" systemd:"User"`
	Group        string            `hcl:"group,optional" systemd:"Group"`
	Environments map[string]string `hcl:"environments,optional" systemd:"Environment"`

	StandardOutput string `hcl:"standard_output,optional" systemd:"StandardOutput"`
	StandardError  string `hcl:"standard_error,optional" systemd:"StandardError"`
}

type Service struct {
	Name string `hcl:"name,label"`

	Unit    UnitBlock    `hcl:"unit,block"`
	Service ServiceBlock `hcl:"service,block"`
	Install InstallBlock `hcl:"install,block"`
}

func (s *Service) Encode() (*SystemdUnit, error) {
	unitSection, err := EncodeSystemdSection(s.Unit)
	if err != nil {
		panic(err)
	}

	serviceSection, err := EncodeSystemdSection(s.Service)
	if err != nil {
		panic(err)
	}

	installSection, err := EncodeSystemdSection(s.Install)
	if err != nil {
		panic(err)
	}

	b := NewSystemdUnitBuilder().
		AddEntries("Unit", unitSection...).
		AddEntries("Service", serviceSection...).
		AddEntries("Install", installSection...)

	return b.Build(s.Name + ".service"), nil
}
