// Copyright (c) Huawei Technologies Co., Ltd. 2019. All rights reserved.
// iSulad-lxcfs-toolkit is licensed under the Mulan PSL v1.
// You can use this software according to the terms and conditions of the Mulan PSL v1.
// You may obtain a copy of Mulan PSL v1 at:
//     http://license.coscl.org.cn/MulanPSL
// THIS SOFTWARE IS PROVIDED ON AN "AS IS" BASIS, WITHOUT WARRANTIES OF ANY KIND, EITHER EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO NON-INFRINGEMENT, MERCHANTABILITY OR FIT FOR A PARTICULAR
// PURPOSE.
// See the Mulan PSL v1 for more details.
// Description: mount/umount in container namespace
// Author: zhangsong
// Create: 2019-05-31

package libmount

import (
	"encoding/json"
	"fmt"
	"isulad-lxcfs-toolkit/libmount/nsexec"
	"os"
	"strconv"
	"strings"
	"syscall"

	"github.com/docker/docker/pkg/reexec"
)

var (
    lxcfsPath = "/var/lib/lxc/lxcfs/cgroup"
)

func init() {
	reexec.Register(nsexec.NsEnterReexecName, WorkInContainer)
}

func setupPipe(name string) (*os.File, error) {
	v := os.Getenv(name)

	fd, err := strconv.Atoi(v)
	if err != nil {
		return nil, fmt.Errorf("unable to convert %s=%s to int", name, v)
	}
	return os.NewFile(uintptr(fd), "pipe"), nil
}

func setupWorkType(name string) (int, error) {
	v := os.Getenv(name)

	worktype, err := strconv.Atoi(v)
	if err != nil {
		return -1, fmt.Errorf("unable to convert %s=%s to int", name, v)
	}
	return worktype, nil
}

// WorkInContainer will handle command in new namespace(container).
func WorkInContainer() {
	var err error
	var worktype int
	var pipe *os.File
	pipe, err = setupPipe(nsexec.InitPipe)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v", err)
		return
	}

	// when pipe setup, should always send back the errors.
	defer func() {
		var msg nsexec.ErrMsg
		if err != nil {
			msg.Error = fmt.Sprintf("%s", err.Error())
		}
		if err := nsexec.WriteJSON(pipe, msg); err != nil {
			fmt.Fprintf(os.Stderr, "%v\n", err)
		}
	}()

	worktype, err = setupWorkType(nsexec.WorkType)
	if err != nil {
		return
	}

	// handle work here:
	switch worktype {
	case nsexec.MountMsg:
		err = doMount(pipe)
	case nsexec.UmountMsg:
		err = doUmount(pipe)
	default:
		err = fmt.Errorf("unkown worktype=(%d)", worktype)
	}
	// do not need to check err here because we have check in defer
	return
}

func doMount(pipe *os.File) error {
	var mount nsexec.Mount
	if err := json.NewDecoder(pipe).Decode(&mount); err != nil {
		return err
	}

    // remount lxcfs cgroup path readonly
    if err := syscall.Mount(mount.Rootfs+lxcfsPath, mount.Rootfs+lxcfsPath, "none", syscall.MS_BIND, ""); err != nil {
        return err
    }
    if err := syscall.Mount(mount.Rootfs+lxcfsPath, mount.Rootfs+lxcfsPath, "none", syscall.MS_BIND|syscall.MS_REMOUNT|syscall.MS_RDONLY, ""); err != nil {
        return err
    }
	for i := 0; i < len(mount.SrcPaths) && i < len(mount.DestPaths); i++ {
		if err := syscall.Mount(mount.SrcPaths[i], mount.DestPaths[i], "none", syscall.MS_BIND, ""); err != nil {
			return err
		}
	}
	return nil
}

func doUmount(pipe *os.File) error {
	var umount nsexec.Umount
	if err := json.NewDecoder(pipe).Decode(&umount); err != nil {
		return err
	}
	for i := 0; i < len(umount.Paths); i++ {
		if err := syscall.Unmount(umount.Paths[i], syscall.MNT_DETACH); err != nil {
			if !strings.Contains(err.Error(), "invalid argument") {
				return err
			}
		}
	}
    if err := syscall.Unmount(lxcfsPath, 0); err != nil {
        if !strings.Contains(err.Error(), "invalid argument") {
            return err
        }
    }
	return nil
}
