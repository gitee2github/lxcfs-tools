// Copyright (c) Huawei Technologies Co., Ltd. 2019. All rights reserved.
// iSulad-lxcfs-toolkit is licensed under the Mulan PSL v1.
// You can use this software according to the terms and conditions of the Mulan PSL v1.
// You may obtain a copy of Mulan PSL v1 at:
//     http://license.coscl.org.cn/MulanPSL
// THIS SOFTWARE IS PROVIDED ON AN "AS IS" BASIS, WITHOUT WARRANTIES OF ANY KIND, EITHER EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO NON-INFRINGEMENT, MERCHANTABILITY OR FIT FOR A PARTICULAR
// PURPOSE.
// See the Mulan PSL v1 for more details.
// Description: remount command
// Author: zhangsong
// Create: 2019-01-18

// go base main package
package main

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"isulad-lxcfs-toolkit/libmount"
	"os"
	"os/exec"
	"strings"
	"sync"
	"syscall"
	"time"

	isulad_lxcfs_log "github.com/sirupsen/logrus"
	"github.com/urfave/cli"
)

var recmountContainer = cli.Command{
	Name:        "remount",
	Usage:       "remount lxcfs to a certain container",
	ArgsUsage:   `[options] <container-id>`,
	Description: `You can remount lxcfs to a running container by container id.`,

	Flags: []cli.Flag{
		cli.BoolFlag{
			Name:  "all, a",
			Usage: "remount lxcfs to all running container",
		},
		cli.StringFlag{
			Name:  "container-id, c",
			Value: " ",
			Usage: "remount a certain container by container id",
		},
	},
	Action: func(context *cli.Context) {
		if context.NArg() > 0 {
			onfailf("%s: requires none args", context.Command.Name)
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
			isulad_lxcfs_log.Info("remount to all containers")
			if err := remountAll(initMountns, initUserns); err != nil {
				onfail(err)
			}

			return
		}

		containerid := context.String("container-id")
		containerid = strings.Replace(containerid, "\n", "", -1)
		if len(containerid) == 1 {
			onfailf("%s: You must choose at least one container", context.Command.Name)
		}

		if err := remountToContainer(initMountns, initUserns, containerid, "", false); err != nil {
			onfail(err)
		}

	},
}

var checklxcfs = cli.Command{
	Name:        "check-lxcfs",
	Usage:       "check whether lxcfs is running",
	ArgsUsage:   ` `,
	Description: `do lxcfs check.`,
	Flags:       []cli.Flag{},

	Action: func(context *cli.Context) {
		if context.NArg() > 0 {
			onfailf("%s: don't need any args", context.Command.Name)
		}
		if err := waitForLxcfs(); err != nil {
			onfail(err)
		} else {
			isulad_lxcfs_log.Info("lxcfs is running")
		}
	},
}

var execprestart = cli.Command{
	Name:        "prestart",
	Usage:       "bind mount before lxcfs start ",
	ArgsUsage:   ` `,
	Description: `do perstart.`,
	Flags:       []cli.Flag{},
	Action: func(context *cli.Context) {
		if context.NArg() > 0 {
			onfailf("%s: don't need any args", context.Command.Name)
		}
		if err := doprestart(); err != nil {
			onfail(err)
		} else {
			isulad_lxcfs_log.Info("prestart done")
		}
	},
}

func doprestart() error {
	isulad_lxcfs_log.Info("do prestart")

	if err := syscall.Unmount("/var/lib/lxc/lxcfs", syscall.MNT_DETACH); err == nil {
		isulad_lxcfs_log.Warning("releaseMountpoint: umount /var/lib/lxc/lxcfs")
	}

	if err := syscall.Unmount("/var/lib/lxc", syscall.MNT_DETACH); err == nil {
		isulad_lxcfs_log.Warning("releaseMountpoint: umount /var/lib/lxc")
	}

	prestartparm1 := []string{
		"-B",
		"/var/lib/lxc",
		"/var/lib/lxc",
	}

	if _, err := execCommond("mount", prestartparm1); err != nil {
		onfail(err)
	}
	prestartparm2 := []string{
		"--make-shared",
		"/var/lib/lxc",
	}
	if _, err := execCommond("mount", prestartparm2); err != nil {
		onfail(err)
	}
	return nil
}

func waitForLxcfs() error {
	count := 0
	maxCount := 100

	for count < maxCount {
		_, err := ioutil.ReadDir("/var/lib/lxc/lxcfs/proc")
		if err != nil {
			time.Sleep(time.Millisecond * 50) // sleep time 50 Millisecond
		} else {
			break
		}
		count++
	}

	if count == maxCount {
		err := fmt.Errorf("lxcfs is not ready")
		isulad_lxcfs_log.Errorf("%v", err)
		return err
	}

	return nil
}

func remountAll(initMountns, initUserns string) error {
	isulad_lxcfs_log.Info("begin remount All runing container...")
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
				if err := remountToContainer(initMountns, initUserns, containerslice[0], containerslice[1], true); err != nil {
					isulad_lxcfs_log.Errorf("remount lxcfs dir to container(%s) failed: %v", containerslice[0], err)
				}
				res <- struct{}{}
			}()
			select {
			case <-res:
			case <-time.After(30 * time.Second): // 30s timeout
				isulad_lxcfs_log.Errorf("remount lxcfs dir to container(%s) timeout", containerslice[0])
			}
		}()
	}
	wg.Wait()
	isulad_lxcfs_log.Info(" remount All done...")
	return nil
}

func remountToContainer(initMountns, initUserns, containerid string, pid string, isAll bool) error {
	if isAll == false {
		var err error
		pid, err = isContainerExsit(containerid)
		if err != nil {
			onfail(err)
		}
	}

	isulad_lxcfs_log.Infof("begin remount container,container id: %s, pid: %s", containerid, pid)

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
	var valueMountPaths []string
	for _, value := range lxcfssubpath {
		valuePaths = append(valuePaths, fmt.Sprintf("/proc/%s", value.Name()))
		valueMountPaths = append(valueMountPaths, fmt.Sprintf("/var/lib/lxc/lxcfs/proc/%s", value.Name()))
	}

	if err := libmount.NsExecUmount(pid, valuePaths); err != nil {
		isulad_lxcfs_log.Errorf("unmount %v for container error: %v", valuePaths, err)
	}

	if err := libmount.NsExecMount(pid, "", valueMountPaths, valuePaths); err != nil {
		isulad_lxcfs_log.Errorf("mount %v into container %s error: %v", valueMountPaths, containerid, err)
		return err
	}
	return nil
}

func isContainerExsit(containerid string) (string, error) {
	isulad_lxcfs_log.Info("begin isContainerExsit...")
	if containerid == "" {
		return "", fmt.Errorf("Containerid mustn't be empty")
	}

	out, err := execCommond("isula", []string{"ps", "--format", "{{.ID}} {{.Pid}}"})
	if err != nil {
		onfail(err)
	}

	for _, value := range out {
		containerslice := strings.Fields(value)
		if len(containerslice) < 2 {
			continue
		}
		if strings.Contains(containerslice[0], containerid) {
			return containerslice[1], nil
		}
	}

	return "", fmt.Errorf("No such Container %s in this node", containerid)
}

func execCommond(command string, params []string) ([]string, error) {
	cmd := exec.Command(command, params...)
	res := []string{
		" ",
	}
	isulad_lxcfs_log.Info("exec cmd :", cmd.Args)

	stdout, err := cmd.StdoutPipe()

	if err != nil {
		return res, err
	}

	if err := cmd.Start(); err != nil {
		return res, err
	}

	reader := bufio.NewReader(stdout)

	for {
		line, err2 := reader.ReadString('\n')
		if err2 == io.EOF {
			break
		} else if err2 != nil {
			onfail(err2)
		}
		res = append(res, line)
	}

	if err := cmd.Wait(); err != nil {
		return res, err
	}

	return res, nil
}
