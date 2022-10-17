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
// Date: 2022-10-14
// Description: This file is used for testing umountcmd

// go base main package
package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"testing"
)

var (
	oldPath, newPath string
)

func isuladExisted() bool {
	cmd := exec.Command("isula", "ps")
	if err := cmd.Run(); err != nil {
		return false
	}
	return true
}

func whereIsIsulad() string {
	cmd := exec.Command("whereis", "isula")
	out, err := cmd.CombinedOutput()
	if err != nil {
		return ""
	}
	path := strings.Trim(strings.Split(string(out), ":")[1], " ")
	path = strings.ReplaceAll(path, "\n", "")
	return path
}

func rename(oldPath, newPath string) error {
	if err := os.Rename(oldPath, newPath); err != nil {
		fmt.Println(err)
		return err
	}
	fmt.Printf("rename isula path from %s to %s\n", oldPath, newPath)
	return nil
}

func tryMoveiSulad() bool {
	oldPath = ""
	newPath = ""
	if isuladExisted() {
		fmt.Println("isuila existed ")
		oldPath = whereIsIsulad()
		if oldPath == "" {
			return false
		}
		newPath = oldPath + ".bak"
		if err := rename(oldPath, newPath); err != nil {
			return false
		}
	}
	return true
}

func recoveriSulad() {
	if oldPath == "" || newPath == "" {
		return
	}
	if err := rename(newPath, oldPath); err != nil {
		return
	}
}

// TestUmountAll tests umountAll
func TestUmountAll(t *testing.T) {
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
			err := umountAll(tt.initMountns, tt.initUserns)
			if (err != nil) != tt.wantErr {
				t.Errorf("umountAll() = %v, want %v", err, tt.wantErr)
			}
		})
	}
}

// TestUmountForContainer tests umountForContainer
func TestUmountForContainer(t *testing.T) {
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
			name:        "TC1-umountForContainer",
			initMountns: "1",
			initUserns:  "5",
			isAll:       true,
			wantErr:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := umountForContainer(tt.initMountns, tt.initUserns, tt.containerid, tt.pid, tt.isAll)
			if (err != nil) != tt.wantErr {
				t.Errorf("umountForContainer() = %v, want %v", err, tt.wantErr)
			}
		})
	}
}
