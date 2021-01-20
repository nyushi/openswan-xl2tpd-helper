package main

import (
	"fmt"
	"os"
	"text/template"
)

func xl2tpdConfFile(remote string) error {
	path := "/etc/xl2tpd/xl2tpd.conf"
	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("failed to create %s: %w", path, err)
	}
	t := template.Must(template.New("xl2tpd_conf").Parse(xl2tpdConfTemplate))
	if err := t.Execute(f, struct {
		Remote string
	}{
		Remote: remote,
	}); err != nil {
		return fmt.Errorf("failed to fill template %s: %w", path, err)
	}
	return nil
}

var xl2tpdConfTemplate = `[lac vpn-connection]
lns = {{ .Remote }}
ppp debug = yes
pppoptfile = /etc/ppp/options.l2tpd.client
length bit = yes
`
