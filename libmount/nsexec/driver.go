// Copyright (c) Huawei Technologies Co., Ltd. 2019. All rights reserved.
// lxcfs-tools is licensed under the Mulan PSL v1.
// You can use this software according to the terms and conditions of the Mulan PSL v1.
// You may obtain a copy of Mulan PSL v1 at:
//     http://license.coscl.org.cn/MulanPSL
// THIS SOFTWARE IS PROVIDED ON AN "AS IS" BASIS, WITHOUT WARRANTIES OF ANY KIND, EITHER EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO NON-INFRINGEMENT, MERCHANTABILITY OR FIT FOR A PARTICULAR
// PURPOSE.
// See the Mulan PSL v1 for more details.
// Description: ns exec driver
// Author: zhangsong
// Create: 2019-05-31

package nsexec

import (
	"fmt"
	"os"
)

var (
	// NSExecDriver is the nsexec driver which will use "c+go" to enter namespace
	NSExecDriver = "nsexec"
	// DefaultNSDriver is the default namespace driver name.
	DefaultNSDriver = NSExecDriver
)

// NsDriver is the namespace driver interface
type NsDriver interface {
	Mount(pid string, mount *Mount) error
	Umount(pid string, umount *Umount) error
}

// NewNsDriver creates the namespace driver by name
func NewNsDriver(name string) (NsDriver, error) {
	switch name {
	case NSExecDriver:
		return NewNSExecDriver(), nil
	}
	return nil, fmt.Errorf("Ns device driver (%s) not supported", name)
}

// NewDefaultNsDriver creates the default namespace driver
func NewDefaultNsDriver() NsDriver {
	drv, err := NewNsDriver(DefaultNSDriver)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
	return drv
}
