package main

import (
	"bufio"
	"github.com/charmbracelet/log"
	"os"
	"path/filepath"
	"strings"
)

type DirectiveIdentifier struct {
	Section string
	Key     string
}

type Directive struct {
	Idenfitier DirectiveIdentifier
	Parser     string
}

func parseGperfLine(line string) (*Directive, bool) {
	line = strings.TrimSpace(line)

	if line == "" || strings.HasPrefix(line, "#") || strings.HasPrefix(line, "/*") {
		return nil, false
	}

	fields := strings.Split(line, ",")
	if len(fields) < 2 {
		return nil, false
	}

	key := strings.TrimSpace(fields[0])
	parser := strings.TrimSpace(fields[1])

	parts := strings.SplitN(key, ".", 2)
	if len(parts) != 2 {
		return nil, false
	}

	return &Directive{
		Idenfitier: DirectiveIdentifier{
			Section: parts[0],
			Key:     parts[1],
		},
		Parser: parser,
	}, true
}

func parserToGoType(parser string) string {
	switch parser {

	case "config_parse_bool":
		return "bool"

	case "config_parse_int",
		"config_parse_long",
		"config_parse_unsigned",
		"config_parse_sec",
		"config_parse_sec_fix_0",
		"config_parse_job_timeout_sec",
		"config_parse_job_running_timeout_sec",
		"config_parse_service_timeout",
		"config_parse_service_timeout_abort",
		"config_parse_service_timeout_failure_mode",
		"config_parse_concurrency_max",
		"config_parse_iec_size",
		"config_parse_swap_priority",
		"config_parse_ip_tos":
		return "int"

	case "config_parse_string",
		"config_parse_path_spec",
		"config_parse_unit_condition_string",
		"config_parse_unit_string_printf",
		"config_parse_unit_path_printf",
		"config_parse_unit_path_strv_printf",
		"config_parse_trigger_unit",
		"config_parse_unit_deps",
		"config_parse_user_group_compat",
		"config_parse_reboot_parameter":
		return "string"

	case "config_parse_strv",
		"config_parse_service_sockets",
		"config_parse_exec",
		"config_parse_exec_preserve_mode",
		"config_parse_socket_service",
		"config_parse_unit_mounts_for",
		"config_parse_fdname",
		"config_parse_timer":
		return "[]string"

	case "config_parse_job_mode",
		"config_parse_job_mode_isolate",
		"config_parse_notify_access",
		"config_parse_service_type",
		"config_parse_service_restart",
		"config_parse_service_restart_mode",
		"config_parse_signal",
		"config_parse_socket_protocol",
		"config_parse_socket_bind",
		"config_parse_socket_bindtodevice",
		"config_parse_socket_defer_trigger",
		"config_parse_socket_timestamping",
		"config_parse_collect_mode",
		"config_parse_documentation",
		"config_parse_set_status",
		"config_parse_oom_policy",
		"config_parse_obsolete_unit_deps":
		return "string"

	default:
		return "string"
	}
}

func getDirectives() map[DirectiveIdentifier]string {
	fragPath := filepath.Join("systemd-man", "src", "core", "load-fragment-gperf.gperf.in")
	f, err := os.Open(fragPath)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	var result = make(map[DirectiveIdentifier]string)
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		d, ok := parseGperfLine(line)
		if !ok || d.Idenfitier.Section == "{{type}}" {
			continue
		}

		result[d.Idenfitier] = parserToGoType(d.Parser)
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	return result
}
