// Copyright (c) Huawei Technologies Co., Ltd. 2019. All rights reserved.
// lxcfs-tools is licensed under the Mulan PSL v1.
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
	"lxcfs-tools/libmount/nsexec"
)

// NsExecMount exec mount in container namespace
func NsExecMount(pid string, rootfs string, srcPaths []string, destPaths []string) error {
	driver := nsexec.NewDefaultNsDriver()
	mount := &nsexec.Mount{
        Rootfs: rootfs,
    }
	for i := 0; i < len(srcPaths) && i < len(destPaths); i++ {
		mount.SrcPaths = append(mount.SrcPaths, srcPaths[i])
		mount.DestPaths = append(mount.DestPaths, destPaths[i])
	}
	return driver.Mount(pid, mount)
}

// NsExecUmount exec umount in container namespace
func NsExecUmount(pid string, paths []string) error {
	driver := nsexec.NewDefaultNsDriver()
	umount := &nsexec.Umount{}
	for i := 0; i < len(paths); i++ {
		umount.Paths = append(umount.Paths, paths[i])
	}
	return driver.Umount(pid, umount)
}
