// Copyright (c) Huawei Technologies Co., Ltd. 2019. All rights reserved.
// lxcfs-tools is licensed under the Mulan PSL v2.
// You can use this software according to the terms and conditions of the Mulan PSL v2.
// You may obtain a copy of Mulan PSL v2 at:
//     http://license.coscl.org.cn/MulanPSL2
// THIS SOFTWARE IS PROVIDED ON AN "AS IS" BASIS, WITHOUT WARRANTIES OF ANY KIND, EITHER EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO NON-INFRINGEMENT, MERCHANTABILITY OR FIT FOR A PARTICULAR
// PURPOSE.
// See the Mulan PSL v2 for more details.
// Description: main funtion
// Author: zhangsong
// Create: 2019-01-18

// go base main package
package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/docker/docker/pkg/reexec"
	_ "github.com/opencontainers/runc/libcontainer/nsenter"
	lxcfs_log "github.com/sirupsen/logrus"
	"github.com/urfave/cli"
)

const (
	usage     = `Toolkit for reconnet to a running isulad using lxcfs`
	syslogTag = "lxcfs-tools"
)

var version = "0.1"

func onfail(err error) {
	lxcfs_log.Error(err)
	fmt.Fprint(os.Stderr, err)
	os.Exit(1)
}

func onfailf(t string, v ...interface{}) {
	onfail(fmt.Errorf(t, v...))
}

func runToolkit() {
	app := cli.NewApp()
	app.Name = "lxcfs-tools"
	app.Usage = usage

	v := []string{
		version,
	}
	app.Version = strings.Join(v, "\n")

	app.Commands = []cli.Command{
		recmountContainer,
		umountContainer,
		checklxcfs,
		execprestart,
	}
	if err := app.Run(os.Args); err != nil {
		onfail(err)
	}

}

func main() {
	if reexec.Init() {
		return
	}
	runToolkit()

}
