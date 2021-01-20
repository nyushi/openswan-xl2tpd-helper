package main

import (
	"fmt"
	"os"
	"text/template"
)

func ipsecSecretFile(server, psk string) error {
	path := "/etc/ipsec.secrets"
	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("failed to create %s: %w", path, err)
	}
	tmpl := `%any {{ .Server }} : PSK "{{ .PSK }}"`
	t := template.Must(template.New("psk").Parse(tmpl))
	if err := t.Execute(f, struct {
		Server string
		PSK    string
	}{
		Server: server,
		PSK:    psk,
	}); err != nil {
		return fmt.Errorf("failed to fill template %s: %w", path, err)
	}
	return nil
}
