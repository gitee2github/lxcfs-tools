// Copyright (c) Huawei Technologies Co., Ltd. 2019. All rights reserved.
// iSulad-lxcfs-toolkit is licensed under the Mulan PSL v1.
// You can use this software according to the terms and conditions of the Mulan PSL v1.
// You may obtain a copy of Mulan PSL v1 at:
//     http://license.coscl.org.cn/MulanPSL
// THIS SOFTWARE IS PROVIDED ON AN "AS IS" BASIS, WITHOUT WARRANTIES OF ANY KIND, EITHER EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO NON-INFRINGEMENT, MERCHANTABILITY OR FIT FOR A PARTICULAR
// PURPOSE.
// See the Mulan PSL v1 for more details.
// Description: umount command
// Author: zhangsong
// Create: 2019-01-18

// go base main package
package main

import (
	"fmt"
	"io/ioutil"
	"isulad-lxcfs-toolkit/libmount"
	"os"
	"strings"
	"sync"
	"time"

	isulad_lxcfs_log "github.com/sirupsen/logrus"
	"github.com/urfave/cli"
)

var umountContainer = cli.Command{
	Name:        "umount",
	Usage:       "umount lxcfs for a certain container",
	ArgsUsage:   `[options] <container-id>`,
	Description: `You can umount lxcfs for a running container by container id.`,

	Flags: []cli.Flag{
		cli.BoolFlag{
			Name:  "all, a",
			Usage: "umount lxcfs for all running container",
		},
		cli.StringFlag{
			Name:  "container-id, c",
			Value: " ",
			Usage: "umount a certain container by container id",
		},
	},
	Action: func(context *cli.Context) {
		if context.NArg() > 0 {
			onfailf("%s:requires at least one argument", context.Command.Name)
		}

		if err := waitForLxcfs(); err != nil {
			onfail(err)
		}

		initMountns, err := os.Readlink("/proc/1/ns/mnt")
		if err != nil {
			onfail(fmt.Errorf("read init mount namespace fail: %v", err))
		}
		initUserns, err := os.Readlink("/proc/1/ns/user")
		if err != nil {
			onfail(fmt.Errorf("read init user namespace fail: %v", err))
		}
		if context.Bool("all") {
			isulad_lxcfs_log.Info("umount for all containers")
			if err := umountAll(initMountns, initUserns); err != nil {
				onfail(err)
			}

			return
		}

		containerid := context.String("container-id")
		containerid = strings.Replace(containerid, "\n", "", -1)
		if len(containerid) == 1 {
			onfailf("%s: You must choose at least one container", context.Command.Name)
		}

		if err := umountForContainer(initMountns, initUserns, containerid, "", false); err != nil {
			onfail(err)
		}

	},
}

func umountAll(initMountns, initUserns string) error {
	isulad_lxcfs_log.Info("begin umount All runing container...")
	out, err := execCommond("isula", []string{"ps", "--format", "{{.ID}} {{.Pid}}"})
	if err != nil {
		return err
	}

	var wg sync.WaitGroup
	for _, value := range out {
		containerslice := strings.Fields(value)
		if len(containerslice) < 2 {
			continue
		}
		wg.Add(1)
		go func() {
			defer wg.Done()
			res := make(chan struct{}, 1)
			go func() {
				if err := umountForContainer(initMountns, initUserns, containerslice[0], containerslice[1], true); err != nil {
					isulad_lxcfs_log.Errorf("umount lxcfs dir from container(%s) failed: %v", containerslice[0], err)
				}
				res <- struct{}{}
			}()
			select {
			case <-res:
			case <-time.After(30 * time.Second): // 30s timeout
				isulad_lxcfs_log.Errorf("umount lxcfs dir from container(%s) timeout", containerslice[0])
			}
		}()

	}
	wg.Wait()
	isulad_lxcfs_log.Info(" umount All done...")
	return nil
}

func umountForContainer(initMountns, initUserns, containerid string, pid string, isAll bool) error {
	if isAll == false {
		var err error
		pid, err = isContainerExsit(containerid)
		if err != nil {
			onfail(err)
		}
	}

	isulad_lxcfs_log.Infof("begin umount container,container id: %s, pid: %s", containerid, pid)

	lxcfssubpath, err := ioutil.ReadDir("/var/lib/lxc/lxcfs/proc")
	if err != nil {
		return fmt.Errorf("Parse lxcfs dir failed: %v", err)
	}

	mountns, err := os.Readlink("/proc/" + pid + "/ns/mnt")
	if err != nil {
		return fmt.Errorf("read container mount namespace fail: %v", err)
	}
	if initMountns == mountns {
		return fmt.Errorf("container pid changed: container mount namespace is same as init namespace")
	}

	var valuePaths []string
	for _, value := range lxcfssubpath {
		valuePaths = append(valuePaths, fmt.Sprintf("/proc/%s", value.Name()))
	}

	if err := libmount.NsExecUmount(pid, valuePaths); err != nil {
		isulad_lxcfs_log.Errorf("unmount %v for container %s error: %v", valuePaths, containerid, err)
		return err
	}
	return nil
}
