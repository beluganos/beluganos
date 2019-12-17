// -*- coding: utf-8 -*-

// Copyright (C) 2019 Nippon Telegraph and Telephone Corporation.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or
// implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package fflib

import (
	"fmt"
	"io"
	"os/exec"
	"strings"

	"golang.org/x/sys/unix"
)

func ExecAndOutput(w io.Writer, cmd string, args ...string) error {
	c := exec.Command(cmd, args...)
	return ExecCmdAndOutput(w, c)
}

func ExecCmdAndOutput(w io.Writer, c *exec.Cmd) error {
	if out, err := c.CombinedOutput(); err != nil {
		fmt.Fprintf(w, "error: %s %s\n", strings.Join(c.Args, " "), err)
		return err
	} else {
		fmt.Fprintf(w, "%s\n", out)
		return nil
	}
}

func ExecAndWait(cmds ...string) error {
	binary, err := exec.LookPath(cmds[0])
	if err != nil {
		return nil
	}

	return unix.Exec(binary, cmds, unix.Environ())
}
