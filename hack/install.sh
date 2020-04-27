# Copyright (c) Huawei Technologies Co., Ltd. 2019. All rights reserved.
# lxcfs-tools is licensed under the Mulan PSL v2.
# You can use this software according to the terms and conditions of the Mulan PSL v2.
# You may obtain a copy of Mulan PSL v2 at:
#     http://license.coscl.org.cn/MulanPSL
# THIS SOFTWARE IS PROVIDED ON AN "AS IS" BASIS, WITHOUT WARRANTIES OF ANY KIND, EITHER EXPRESS OR
# IMPLIED, INCLUDING BUT NOT LIMITED TO NON-INFRINGEMENT, MERCHANTABILITY OR FIT FOR A PARTICULAR
# PURPOSE.
# See the Mulan PSL v2 for more details.
# Description: install script
# Author: zhangsong
# Create: 2019-01-18

#!/bin/bash 


LXCFS_TOOLS_DIR=/usr/local/bin

echo "lxcfs-tools will be installed to $LXCFS_TOOLS_DIR"

install -m 0755 -p ../build/lxcfs-tools ${LXCFS_TOOLS_DIR}

echo "lxcfs-tools  install  done"
