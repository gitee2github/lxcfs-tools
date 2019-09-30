// Copyright (c) Huawei Technologies Co., Ltd. 2019. All rights reserved.
// iSulad-lxcfs-toolkit is licensed under the Mulan PSL v1.
// You can use this software according to the terms and conditions of the Mulan PSL v1.
// You may obtain a copy of Mulan PSL v1 at:
//     http://license.coscl.org.cn/MulanPSL
// THIS SOFTWARE IS PROVIDED ON AN "AS IS" BASIS, WITHOUT WARRANTIES OF ANY KIND, EITHER EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO NON-INFRINGEMENT, MERCHANTABILITY OR FIT FOR A PARTICULAR
// PURPOSE.
// See the Mulan PSL v1 for more details.
// Description: prestart hook
// Author: zhangsong
// Create: 2019-01-18

// go base main package
package main

import (
	"flag"
	"fmt"
	"isula.org/isulad-lxcfs-toolkit/hooks/lxcfs-hook/utils"
	"os"

	"github.com/docker/docker/pkg/reexec"
	_ "github.com/opencontainers/runc/libcontainer/nsenter"
	"github.com/sirupsen/logrus"
)

var (
	syslogTag = "lxcfs-hook"
	nsPID     = 1
	rootfs    = ""
)

func setupLog(logfile string) {
	logrus.SetLevel(logrus.DebugLevel)
	logrus.SetOutput(os.Stdout)
	if logfile != "" {
		f, err := os.OpenFile(logfile, os.O_CREATE|os.O_WRONLY|os.O_APPEND|os.O_SYNC, 0600)
		if err != nil {
			return
		}
		logrus.SetOutput(f)
		return
	}

	if err := utils.SetupSyslog("unix:///dev/log", syslogTag); err != nil {
		// failed to setup Syslog
		fmt.Fprintf(os.Stdout, "%v", err)
	}
}

func main() {
	if reexec.Init() {
		return
	}
	flLogfile := flag.String("log", "", "set output log file")
	flag.Parse()
	setupLog(*flLogfile)
	state, err := utils.ParseHookState(os.Stdin)
	if err != nil {
		logrus.Errorf("Parse Hook State Failed: %v", err)
		return
	}

	if state.Pid <= 0 {
		logrus.Errorf("Can't get correct pid of container:%d", state.Bundle)
	}
	logrus.Infof("PID:%d", state.Pid)
	logrus.Infof("Root:%s", state.Root)
	nsPID = state.Pid
	rootfs = state.Root

	if err := prestartMountHook(nsPID, rootfs); err != nil {
		logrus.Errorf("Can't mount lxcfs to certain container,%v", err)
	}

	return
}
