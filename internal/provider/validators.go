package provider

import (
	"fmt"
	"strings"
)

// Credits
// https://github.com/hashicorp/terraform-provider-dns/blob/main/internal/provider/validators.go

// IsFqdn checks if a domain name is fully qualified.
func IsFqdn(s string) bool {
	s2 := strings.TrimSuffix(s, ".")
	if s == s2 {
		return false
	}

	i := strings.LastIndexFunc(s2, func(r rune) bool {
		return r != '\\'
	})

	// Test whether we have an even number of escape sequences before
	// the dot or none.
	return (len(s2)-i)%2 != 0
}

func validateZone(v interface{}, k string) (ws []string, errors []error) {
	value := v.(string)
	if strings.TrimSpace(value) != value {
		errors = append(errors, fmt.Errorf("DNS zone name %q must not contain whitespace: %q", k, value))
	}
	if !IsFqdn(value) {
		errors = append(errors, fmt.Errorf("DNS zone name %q must be fully qualified: %q", k, value))
	}
	return
}

func validateName(v interface{}, k string) (ws []string, errors []error) {
	value := v.(string)
	if strings.TrimSpace(value) != value || len(value) == 0 {
		errors = append(errors, fmt.Errorf("DNS record name %q must not contain whitespace or be empty: %q", k, value))
	}
	if IsFqdn(value) {
		errors = append(errors, fmt.Errorf("DNS record name %q must not be fully qualified: %q", k, value))
	}
	return
}
