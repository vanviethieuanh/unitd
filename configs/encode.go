package configs

import "strings"

// Encode converts a Service to a systemd .service unit file string.
// Sections are written in standard order: [Unit], [Service], [Install].
func (s *Service) Encode() (string, error) {
	var b strings.Builder

	sections := []struct {
		name string
		data any
	}{
		{"Unit", s.Unit},
		{"Service", s.Service},
		{"Install", s.Install},
	}

	for _, sec := range sections {
		entries, err := EncodeSystemdSection(sec.data)
		if err != nil {
			return "", err
		}
		if len(entries) > 0 {
			writeSection(&b, sec.name, entries)
		}
	}

	return strings.TrimSpace(b.String()) + "\n", nil
}
