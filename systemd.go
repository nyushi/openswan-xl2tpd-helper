package main

import (
	"fmt"

	"github.com/coreos/go-systemd/dbus"
)

func systemdOperation(svc string, f func(*dbus.Conn, chan string) error) error {
	sdcon, err := dbus.NewSystemdConnection()
	if err != nil {
		return fmt.Errorf("failed to get systemd socket conn: %w", err)
	}
	c := make(chan string)
	if err := f(sdcon, c); err != nil {
		return err
	}
	result := <-c
	switch result {
	case "done":
		return nil
	default:
		return fmt.Errorf("failed to process operation: %s", result)
	}
}
func startSystemdService(svc string) error {
	err := systemdOperation(svc, func(sdcon *dbus.Conn, c chan string) error {
		_, err := sdcon.StartUnit(svc, "replace", c)
		return err
	})
	if err != nil {
		return fmt.Errorf("error at systemd %s start: %w", svc, err)
	}
	return nil
}

func stopSystemdService(svc string) error {
	err := systemdOperation(svc, func(sdcon *dbus.Conn, c chan string) error {
		_, err := sdcon.StartUnit(svc, "replace", c)
		return err
	})
	if err != nil {
		return fmt.Errorf("error at systemd %s start: %w", svc, err)
	}
	return nil
}
