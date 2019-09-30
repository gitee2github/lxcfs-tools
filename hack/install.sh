# Copyright (c) Huawei Technologies Co., Ltd. 2019. All rights reserved.
# iSulad-lxcfs-toolkit is licensed under the Mulan PSL v1.
# You can use this software according to the terms and conditions of the Mulan PSL v1.
# You may obtain a copy of Mulan PSL v1 at:
#     http://license.coscl.org.cn/MulanPSL
# THIS SOFTWARE IS PROVIDED ON AN "AS IS" BASIS, WITHOUT WARRANTIES OF ANY KIND, EITHER EXPRESS OR
# IMPLIED, INCLUDING BUT NOT LIMITED TO NON-INFRINGEMENT, MERCHANTABILITY OR FIT FOR A PARTICULAR
# PURPOSE.
# See the Mulan PSL v1 for more details.
# Description: install script
# Author: zhangsong
# Create: 2019-01-18

#!/bin/bash 


ISULAD_LXCFS_TOOLKIT_DIR=/usr/local/bin

echo "isulad_lxcfs_toolkit will be installed to $ISULAD_LXCFS_TOOLKIT_DIR"

install -m 0755 -p ../build/isulad_lxcfs_toolkit ${ISULAD_LXCFS_TOOLKIT_DIR}

echo "isulad_lxcfs_toolkit  install  done"
