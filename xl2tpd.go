package main

import (
	"fmt"
	"os"
)

func l2tpstart() error {
	path := "/var/run/xl2tpd/l2tp-control"
	payload := `c vpn-connection`
	if err := writeToPipe(path, payload); err != nil {
		return fmt.Errorf("failed write `%s` to %s: %w", path, payload, err)
	}
	return nil
}

func l2tpstop() error {
	path := "/var/run/xl2tpd/l2tp-control"
	payload := `d vpn-connection`
	if err := writeToPipe(path, payload); err != nil {
		return fmt.Errorf("failed write `%s` to %s: %w", path, payload, err)
	}
	return nil
}

func writeToPipe(path, payload string) error {
	f, err := os.OpenFile(path, os.O_WRONLY, 0600)
	if err != nil {
		return fmt.Errorf("open pipe file: %w", err)
	}
	defer f.Close()
	if _, err := f.Write([]byte(payload)); err != nil {
		return fmt.Errorf("failed to write: %w", err)
	}
	return nil
}
