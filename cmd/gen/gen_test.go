package main

import "testing"

func TestToPascalCase(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"timer", "Timer"},
		{"on-calendar", "OnCalendar"},
		{"OnActiveSec", "OnActiveSec"},
		{"AccuracySec", "AccuracySec"},
		{"on-boot-sec", "OnBootSec"},
	}

	for _, tt := range tests {
		result := toPascalCase(tt.input)
		if result != tt.expected {
			t.Errorf("toPascalCase(%q) = %q; expected %q", tt.input, result, tt.expected)
		}
	}
}

func TestToSnakeCase(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"OnActiveSec", "on_active_sec"},
		{"OnCalendar", "on_calendar"},
		{"AccuracySec", "accuracy_sec"},
		{"on-boot-sec", "on_boot_sec"},
	}

	for _, tt := range tests {
		result := toSnakeCase(tt.input)
		if result != tt.expected {
			t.Errorf("toSnakeCase(%q) = %q; expected %q", tt.input, result, tt.expected)
		}
	}
}
