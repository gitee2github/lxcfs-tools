// Copyright (c) Huawei Technologies Co., Ltd. 2021-2022. All rights reserved.
// rubik licensed under the Mulan PSL v2.
// You can use this software according to the terms and conditions of the Mulan PSL v2.
// You may obtain a copy of Mulan PSL v2 at:
//     http://license.coscl.org.cn/MulanPSL2
// THIS SOFTWARE IS PROVIDED ON AN "AS IS" BASIS, WITHOUT WARRANTIES OF ANY KIND, EITHER EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO NON-INFRINGEMENT, MERCHANTABILITY OR FIT FOR A PARTICULAR
// PURPOSE.
// See the Mulan PSL v2 for more details.
// Author: Jiaqi Yang
// Date: 2022-10-15
// Description: This file is used for remountcmd.go

// go base main package
package main

import (
	"os"
	"testing"
)

const (
	oldLxcPath = "/var/lib/lxc"
	newLxcPath = "/var/lib/lxc.bak"
)

func pathExists(path string) bool {
	if _, err := os.Stat(path); err != nil {
		return false
	}
	return true
}

func tryMovelxcfs() bool {
	if pathExists(oldLxcPath) {
		if err := rename(oldLxcPath, newLxcPath); err != nil {
			return false
		}
	}
	return true
}

func tryRecoverlxcfs() {
	if pathExists(newLxcPath) {
		if err := rename(newLxcPath, oldLxcPath); err != nil {
			return
		}
	}
}

// TestWaitForLxcfs tests waitForLxcfs
func TestWaitForLxcfs(t *testing.T) {
	if !tryMovelxcfs() {
		return
	}
	defer tryRecoverlxcfs()
	if err := waitForLxcfs(); err == nil {
		t.Errorf("doprestart() should fail but runs successfully")
	}
}

// TestRemountAll tests remountAll
func TestRemountAll(t *testing.T) {
	if !tryMoveiSulad() {
		return
	}
	defer recoveriSulad()

	tests := []struct {
		name, initMountns, initUserns string
		wantErr                       bool
	}{
		{
			name:        "TC1-absent of isulad",
			initMountns: "1",
			initUserns:  "5",
			wantErr:     true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := remountAll(tt.initMountns, tt.initUserns)
			if (err != nil) != tt.wantErr {
				t.Errorf("remountAll() = %v, want %v", err, tt.wantErr)
			}
		})
	}
}

// TestRemountToContainer tests remountToContainer
func TestRemountToContainer(t *testing.T) {
	if !tryMovelxcfs() {
		return
	}
	defer tryRecoverlxcfs()
	tests := []struct {
		name, initMountns, initUserns, containerid, pid string
		isAll                                           bool
		wantErr                                         bool
	}{
		{
			name:        "TC1-remountToContainer",
			initMountns: "1",
			initUserns:  "5",
			isAll:       true,
			wantErr:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := remountToContainer(tt.initMountns, tt.initUserns, tt.containerid, tt.pid, tt.isAll)
			if (err != nil) != tt.wantErr {
				t.Errorf("umountForContainer() = %v, want %v", err, tt.wantErr)
			}
		})
	}
}

// TestIsContainerExsit tests isContainerExsit
func TestIsContainerExsit(t *testing.T) {
	tests := []struct {
		name, containerid string
		wantErr           bool
	}{
		{
			name:    "TC1-remountToContainer",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := isContainerExsit(tt.containerid)
			if (err != nil) != tt.wantErr {
				t.Errorf("isContainerExsit() = %v, want %v", err, tt.wantErr)
			}
		})
	}
}
