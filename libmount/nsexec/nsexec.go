// Copyright (c) Huawei Technologies Co., Ltd. 2019. All rights reserved.
// iSulad-lxcfs-toolkit is licensed under the Mulan PSL v1.
// You can use this software according to the terms and conditions of the Mulan PSL v1.
// You may obtain a copy of Mulan PSL v1 at:
//     http://license.coscl.org.cn/MulanPSL
// THIS SOFTWARE IS PROVIDED ON AN "AS IS" BASIS, WITHOUT WARRANTIES OF ANY KIND, EITHER EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO NON-INFRINGEMENT, MERCHANTABILITY OR FIT FOR A PARTICULAR
// PURPOSE.
// See the Mulan PSL v1 for more details.
// Description: ns exec in container namespace
// Author: zhangsong
// Create: 2019-05-31

package nsexec

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"syscall"

	"github.com/opencontainers/runc/libcontainer"
	"github.com/vishvananda/netlink/nl"
)

const (
	// MountMsg is the mount msg
	MountMsg = 1
	// UmountMsg is the umount  msg
	UmountMsg = 2
	// InitPipe is a parent and child process env name, used to pass the init pipe number to child process
	InitPipe = "_LIBCONTAINER_INITPIPE"
	// WorkType is a parent and child process env name, used to pass the work type to child process
	WorkType = "_LXCFS_TOOLS_WORKTYPE"
	// NsEnterReexecName is the reexec name, see reexec package
	NsEnterReexecName = "nsenter-init"
)

// Mount is mount argument
type Mount struct {
    Rootfs    string
	SrcPaths  []string
	DestPaths []string
}

// Umount is umount argument
type Umount struct {
	Paths []string
}

type nsexecDriver struct {
}

type pid struct {
	Pid int `json:"Pid"`
}

// ErrMsg is error msg
type ErrMsg struct {
	Error string
}

// NewNSExecDriver creates the nsexecDriver
func NewNSExecDriver() NsDriver {
	return &nsexecDriver{}
}

// NewPipe creates a pair of  new socket pipe.
func newPipe() (parent, child *os.File, err error) {
	fds, err := syscall.Socketpair(syscall.AF_LOCAL, syscall.SOCK_STREAM|syscall.SOCK_CLOEXEC, 0)
	if err != nil {
		return nil, nil, err
	}
	return os.NewFile(uintptr(fds[1]), "parent"), os.NewFile(uintptr(fds[0]), "child"), nil
}

// WriteJSON writes the provided struct v to w using standard json marshaling
func WriteJSON(w io.Writer, v interface{}) error {
	data, err := json.Marshal(v)
	if err != nil {
		return err
	}
	_, err = w.Write(data)
	return err
}

func (ns *nsexecDriver) exec(nsPaths string, worktype int, data interface{}) error {
	parent, child, err := newPipe()
	if err != nil {
		return err
	}
	cmd := &exec.Cmd{
		Path:       "/proc/self/exe",
		Args:       []string{NsEnterReexecName},
		ExtraFiles: []*os.File{child},
		Env: []string{fmt.Sprintf("%s=3", InitPipe),
			fmt.Sprintf("%s=%d", WorkType, worktype)},
		Stdout: os.Stdout,
		Stderr: os.Stderr,
	}

	if err := cmd.Start(); err != nil {
		return err
	}
	r := nl.NewNetlinkRequest(int(libcontainer.InitMsg), 0)

	r.AddData(&libcontainer.Bytemsg{
		Type:  libcontainer.NsPathsAttr,
		Value: []byte(nsPaths),
	})

	// send nspath to child process through _LXCFS_TOOLS_INITPIPE, to join container ns.
	if _, err := io.Copy(parent, bytes.NewReader(r.Serialize())); err != nil {
		return err
	}
	// send the config to child
	if err := WriteJSON(parent, data); err != nil {
		return err
	}

	// wait for command
	if err := cmd.Wait(); err != nil {
		return err
	}

	decoder := json.NewDecoder(parent)
	var pid *pid
	if err := decoder.Decode(&pid); err != nil {
		fmt.Fprintf(os.Stderr, "fail to decode pid:%v, but it may not affect later process", err)
	}

	// read error message
	var msg ErrMsg
	if err := decoder.Decode(&msg); err != nil {
		return err
	}

	if msg.Error != "" {
		return fmt.Errorf("%s", msg.Error)
	}
	return nil
}

func (ns *nsexecDriver) Mount(pid string, mount *Mount) error {
	namespaces := []string{"mnt"}
	nsPaths := buildNSString(pid, namespaces)

	return ns.exec(nsPaths, MountMsg, mount)
}

func (ns *nsexecDriver) Umount(pid string, umount *Umount) error {
	namespaces := []string{"mnt"}
	nsPaths := buildNSString(pid, namespaces)

	return ns.exec(nsPaths, UmountMsg, umount)
}

func buildNSString(pid string, namespaces []string) string {
	var nsPaths string
	for _, ns := range namespaces {
		if nsPaths != "" {
			nsPaths += ","
		}
		nsPaths += fmt.Sprintf("%s:/proc/%s/ns/%s", ns, pid, ns)
	}
	return nsPaths
}
