#Global macro or variable
%define debug_package %{nil}

#Basic Information
Name:		isulad-lxcfs-toolkit
Version:	0.3
Release:	15
Summary:	toolkit for lxcfs to remount a running isulad
License:  Mulan PSL v1
URL:      https://gitee.com/src-openeuler/iSulad-lxcfs-toolkit
Source0:	%{name}.tar.gz
BuildRoot:      %{_tmppath}/%{name}-root

#Dependency
BuildRequires:	golang > 1.7
BuildRequires:  glibc-static
Requires: iSulad

%description
A toolkit for lxcfs to remount a running isulad when crashes recover

#Build sections
%prep
export RPM_BUILD_DIR=%_topdir/BUILD
export RPM_BUILD_SOURCE=%_topdir/SOURCES

cd $RPM_BUILD_DIR

mkdir -p $RPM_BUILD_DIR/src/isula.org/isulad-lxcfs-toolkit && cd $RPM_BUILD_DIR/src/isula.org/isulad-lxcfs-toolkit
gzip -dc $RPM_BUILD_SOURCE/%{name}.tar.gz | tar -xvvf -

%build
cd $RPM_BUILD_DIR/src/isula.org/isulad-lxcfs-toolkit
make

%install
HOOK_DIR=$RPM_BUILD_ROOT/var/lib/lcrd/hooks
ISULAD_LXCFS_TOOLKIT_DIR=$RPM_BUILD_ROOT/usr/local/bin

cd $RPM_BUILD_DIR/src/isula.org/isulad-lxcfs-toolkit
mkdir -p -m 0700 ${HOOK_DIR}
mkdir -p -m 0700 ${ISULAD_LXCFS_TOOLKIT_DIR}

install -m 0750 build/lxcfs-hook ${HOOK_DIR}
install -m 0750 build/isulad-lxcfs-toolkit ${ISULAD_LXCFS_TOOLKIT_DIR}

#Install and uninstall scripts
%pre

%preun

%post
GRAPH=`lcrc info | grep -Eo "iSulad Root Dir:.+" | grep -Eo "\/.*"` 
if [ "$GRAPH" == "" ]; then
    GRAPH="/var/lib/lcrd"
fi

if [[ ("$GRAPH" != "/var/lib/lcrd") ]]; then
    mkdir -p -m 0550 $GRAPH/hooks
    install -m 0550 -p /var/lib/lcrd/hooks/lxcfs-hook $GRAPH/hooks

    echo
    echo "=================== WARNING! ================================================"
    echo " 'iSulad Root Dir' is $GRAPH, move /var/lib/lcrd/hooks/lxcfs-hook to  $GRAPH/hooks"
    echo "============================================================================="
    echo
fi
HOOK_SPEC=${GRAPH}/hooks
HOOK_DIR=${GRAPH}/hooks
touch ${HOOK_SPEC}/hookspec.json
cat << EOF > ${HOOK_SPEC}/hookspec.json
{
        "prestart": [
        {
                "path": "${HOOK_DIR}/lxcfs-hook",
                "args": ["lxcfs-hook"],
                "env": []
        }
        ],
        "poststart":[],
        "poststop":[]
}

EOF
chmod 0640 ${HOOK_SPEC}/hookspec.json

%postun

#Files list
%files
%defattr(0550,root,root,0550)
/usr/local/bin/isulad-lxcfs-toolkit
%attr(0550,root,root) /var/lib/lcrd/hooks
%attr(0550,root,root) /var/lib/lcrd/hooks/lxcfs-hook

#Clean section
%clean 
rm -rfv %{buildroot}


%changelog
* Thu Feb 1 2018 Tanzhe <tanzhe@huawei.com>
- add require version
