package main

import (
	"bytes"
	"io/ioutil"
	"os/exec"
	"strings"
)

func execCommand(commands ...string) error {
	log.Debugf("execute %s", strings.Join(commands, " "))
	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}
	cmd := exec.Command(commands[0], commands[1:]...)
	cmd.Stdout = stdout
	cmd.Stderr = stderr
	err := cmd.Run()
	bout, _ := ioutil.ReadAll(stdout)
	berr, _ := ioutil.ReadAll(stderr)
	if err != nil {
		log.Debugf("stdout: %s", string(bout))
		log.Debugf("stderr: %s", string(berr))
		return err
	}
	return nil
}

type ipsecTemporaryError struct {
	error
}

func (e *ipsecTemporaryError) Temporary() bool {
	err, ok := e.error.(*exec.ExitError)
	if !ok {
		return false
	}
	if err.ProcessState.ExitCode() == 24 {
		return true
	}
	return false
}

func ipsecCommand(options ...string) error {
	cmd := []string{"ipsec"}
	cmd = append(cmd, options...)
	err := execCommand(cmd...)
	if err != nil {
		return &ipsecTemporaryError{err}
	}
	return nil
}

