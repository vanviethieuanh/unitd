package configs

import (
	"fmt"
	"strings"
)

// sanitizeHCLIdent converts a systemd unit name stem to a valid HCL identifier.
// Hyphens become underscores, "@" becomes "_at_", leading "-" becomes "root".
func sanitizeHCLIdent(name string) string {
	dot := strings.LastIndex(name, ".")
	if dot < 0 {
		return name
	}
	stem := name[:dot]

	if stem == "-" {
		return "root"
	}

	stem = strings.ReplaceAll(stem, "-", "_")
	stem = strings.ReplaceAll(stem, "@", "_at_")

	return stem
}

// TemplateUnitName builds the systemd template unit filename.
// ("worker", "queue", "service") → "worker-queue@.service"
// ("worker", "", "service")      → "worker@.service"
func TemplateUnitName(name, variant, unitType string) string {
	if variant != "" {
		return fmt.Sprintf("%s-%s@.%s", name, variant, unitType)
	}
	return fmt.Sprintf("%s@.%s", name, unitType)
}

// InstanceUnitName builds a systemd instance unit filename from a template name.
// ("worker-queue@.service", "q1") → "worker-queue@q1.service"
func InstanceUnitName(templateName, instance string) string {
	at := strings.Index(templateName, "@")
	if at < 0 {
		return templateName
	}
	dot := strings.LastIndex(templateName, ".")
	if dot <= at {
		return templateName
	}
	return templateName[:at+1] + instance + templateName[dot:]
}
