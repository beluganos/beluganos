// -*- coding: utf-8 -*-

package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	"golang.org/x/sys/unix"
)

func execAndOutput(cmd string, args ...string) error {
	c := exec.Command(cmd, args...)
	out, err := c.CombinedOutput()
	fmt.Printf("%s\n", out)
	return err
}

func execAndWait(cmds ...string) error {
	binary, err := exec.LookPath(cmds[0])
	if err != nil {
		return nil
	}

	return unix.Exec(binary, cmds, unix.Environ())
}

func indexOf(s string, arr []string) int {
	for index, a := range arr {
		if s == a {
			return index
		}
	}

	return -1
}

func trimLine(s string) string {
	return strings.Trim(s, " \n")
}

func containerDevices(name string, excludeDevices []string) ([]string, error) {
	cmd := exec.Command("lxc", "profile", "device", "list", name)
	out, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	ifnames := []string{}
	devices := strings.Split(string(out), "\n")
	for _, device := range devices {
		device = trimLine(device)
		if len(device) == 0 {
			continue
		}
		if indexOf(device, excludeDevices) >= 0 {
			continue
		}
		ifnames = append(ifnames, device)
	}
	return ifnames, nil
}

func containerDeviceProperty(name string, device string, property string) (string, error) {
	cmd := exec.Command("lxc", "profile", "device", "get", name, device, property)
	out, err := cmd.Output()
	if err != nil {
		return "", err
	}

	prop := trimLine(string(out))
	if len(prop) == 0 {
		return "", fmt.Errorf("property '%s' not found in %s/%s.", property, name, device)
	}
	return trimLine(string(out)), nil
}

func createFile(path string, overwrite bool, f func(string)) (*os.File, error) {
	_, err := os.Stat(path)
	if !os.IsNotExist(err) && !overwrite {
		jst := time.FixedZone("Asia/Tokyo", 9*60*60)
		now := time.Now().UTC().In(jst).Format("20060102_150405")
		backupPath := fmt.Sprintf("%s_%s", path, now)
		if err := os.Rename(path, backupPath); err != nil {
			return nil, err
		}

		f(backupPath)
	}

	return os.Create(path)
}
