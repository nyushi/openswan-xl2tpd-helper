package main

import (
	"fmt"
	"os"
	"text/template"
)

func optionsL2TPDClientFile(user, pass string) error {
	path := "/etc/ppp/options.l2tpd.client"
	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("failed to create %s: %w", path, err)
	}
	t := template.Must(template.New("options_l2tpd_client").Parse(optionsL2TPDClientTemplate))
	if err := t.Execute(f, struct {
		User string
		Pass string
	}{
		User: user,
		Pass: pass,
	}); err != nil {
		return fmt.Errorf("failed to fill template %s: %w", path, err)
	}
	return nil
}

var optionsL2TPDClientTemplate = `
ipcp-accept-local
ipcp-accept-remote
refuse-eap
require-mschap-v2
noccp
noauth
idle 1800
mtu 1410
mru 1410
defaultroute
usepeerdns
debug
connect-delay 5000
name {{ .User }}
password {{ .Pass }}
`
