module lxcfs-tools

replace (
	github.com/docker/docker => github.com/docker/engine v0.0.0-20181106193140-f5749085e9cb
	golang.org/x/sys => github.com/golang/sys v0.0.0-20190813064441-fde4db37ae7a
	gopkg.in/yaml.v2 => github.com/go-yaml/yaml v2.1.0+incompatible
)

require (
	github.com/containerd/console v0.0.0-20181022165439-0650fd9eeb50 // indirect
	github.com/coreos/go-systemd v0.0.0-20190719114852-fd7a80b32e1f // indirect
	github.com/coreos/pkg v0.0.0-20180928190104-399ea9e2e55f // indirect
	github.com/cyphar/filepath-securejoin v0.2.2 // indirect
	github.com/docker/docker v0.0.0-00010101000000-000000000000
	github.com/docker/go-units v0.4.0 // indirect
	github.com/godbus/dbus v4.1.0+incompatible // indirect
	github.com/golang/protobuf v1.3.2 // indirect
	github.com/mrunalp/fileutils v0.0.0-20171103030105-7d4729fb3618 // indirect
	github.com/opencontainers/runc v1.0.0-rc6
	github.com/opencontainers/runtime-spec v1.0.1
	github.com/opencontainers/selinux v1.3.0 // indirect
	github.com/pkg/errors v0.8.1 // indirect
	github.com/seccomp/libseccomp-golang v0.9.1 // indirect
	github.com/sirupsen/logrus v1.4.2
	github.com/syndtr/gocapability v0.0.0-20180916011248-d98352740cb2 // indirect
	github.com/urfave/cli v1.21.0
	github.com/vishvananda/netlink v1.0.0
	github.com/vishvananda/netns v0.0.0-20190625233234-7109fa855b0f // indirect
)
