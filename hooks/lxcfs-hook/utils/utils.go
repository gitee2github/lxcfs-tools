// Copyright (c) Huawei Technologies Co., Ltd. 2019. All rights reserved.
// lxcfs-tools is licensed under the Mulan PSL v1.
// You can use this software according to the terms and conditions of the Mulan PSL v1.
// You may obtain a copy of Mulan PSL v1 at:
//     http://license.coscl.org.cn/MulanPSL
// THIS SOFTWARE IS PROVIDED ON AN "AS IS" BASIS, WITHOUT WARRANTIES OF ANY KIND, EITHER EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO NON-INFRINGEMENT, MERCHANTABILITY OR FIT FOR A PARTICULAR
// PURPOSE.
// See the Mulan PSL v1 for more details.
// Description: utils
// Author: zhangsong
// Create: 2019-01-18

package utils

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log/syslog"
	"strings"

	"github.com/opencontainers/runtime-spec/specs-go"
	"github.com/sirupsen/logrus"
	logrus_syslog "github.com/sirupsen/logrus/hooks/syslog"
)

const (
	syslogUDPPrefix         = "udp://"
	syslogTCPPrefix         = "tcp://"
	syslogUnixSock          = "unix://"
	syslogDefaultUDPService = "localhost:541"
	syslogDefaultTCPService = "localhost:541"
)

// SyslogService syslog service
type SyslogService struct {
	Type string
	Addr string
}

// ParseSyslogService parses syslog service from input string
func ParseSyslogService(service string) (*SyslogService, error) {
	serviceType := ""
	serviceAddr := ""
	// UDP syslog service
	if strings.HasPrefix(service, syslogUDPPrefix) {
		serviceType = "udp"
		serviceAddr := service[len(syslogUDPPrefix):]
		if serviceAddr == "" {
			serviceAddr = syslogDefaultUDPService
		}
	} else if strings.HasPrefix(service, syslogTCPPrefix) {
		serviceType = "tcp"
		serviceAddr = service[len(syslogTCPPrefix):]
		if serviceAddr == "" {
			serviceAddr = syslogDefaultTCPService
		}
	} else if strings.HasPrefix(service, syslogUnixSock) {
		// syslog package will use empty string as network,
		// and syslog will lookup the unix socket on host, we do not care.
		serviceType = ""
		serviceAddr = service[len(syslogUnixSock):]
	} else {
		return nil, fmt.Errorf("Unspported syslog network: %s", service)
	}

	serv := &SyslogService{
		Type: serviceType,
		Addr: serviceAddr,
	}
	return serv, nil
}

// SetupSyslog setup syslog service
func SetupSyslog(service, tag string) error {
	serv, err := ParseSyslogService(service)
	if err != nil {
		return err
	}

	hook, err := logrus_syslog.NewSyslogHook(serv.Type, serv.Addr, syslog.LOG_INFO, tag)
	if err != nil {
		return fmt.Errorf("Unable to connect to syslog daemon")
	}
	logrus.AddHook(hook)
	return nil
}

// HookState is container hook state
type HookState struct {
	specs.State
	Root string `json:"root"`
}

// CompatHookState is container hook stat compat with old version
type CompatHookState struct {
	specs.State
	Bundle string `json:"bundlePath"`
}

// ParseHookState parses container hook state from json
func ParseHookState(reader io.Reader) (*HookState, error) {
	// We expect configs.HookState as a json string in <stdin>
	stateBuf, err := ioutil.ReadAll(reader)
	if err != nil {
		return nil, err
	}
	var state HookState
	if err = json.Unmarshal(stateBuf, &state); err != nil {
		return nil, err
	}

	var compatStat CompatHookState
	if state.Bundle == "" {
		if err = json.Unmarshal(stateBuf, &compatStat); err != nil {
			return nil, err
		}
		if compatStat.Bundle == "" {
			return nil, fmt.Errorf("unmarshal hook state failed %s", stateBuf)
		}
		state.Bundle = compatStat.Bundle
	}
	return &state, nil
}
