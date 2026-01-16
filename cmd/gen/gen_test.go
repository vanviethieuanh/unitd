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

func Test_gen(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		u         *Unit
		propTypes map[DirectiveIdentifier]string
		want      string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := gen(tt.u, tt.propTypes)
			// TODO: update the condition below to compare got with tt.want.
			if true {
				t.Errorf("gen() = %v, want %v", got, tt.want)
			}
		})
	}
}
