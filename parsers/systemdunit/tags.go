package systemdunit

var LinuxCapabilities = []string{
	"CAP_AUDIT_CONTROL",
	"CAP_AUDIT_READ",
	"CAP_AUDIT_WRITE",
	"CAP_BLOCK_SUSPEND",
	"CAP_BPF",
	"CAP_CHECKPOINT_RESTORE",
	"CAP_CHOWN",
	"CAP_DAC_OVERRIDE",
	"CAP_DAC_READ_SEARCH",
	"CAP_FOWNER",
	"CAP_FSETID",
	"CAP_IPC_LOCK",
	"CAP_IPC_OWNER",
	"CAP_KILL",
	"CAP_LEASE",
	"CAP_LINUX_IMMUTABLE",
	"CAP_MAC_ADMIN",
	"CAP_MAC_OVERRIDE",
	"CAP_MKNOD",
	"CAP_NET_ADMIN",
	"CAP_NET_BIND_SERVICE",
	"CAP_NET_BROADCAST",
	"CAP_NET_RAW",
	"CAP_PERFMON",
	"CAP_SETFCAP",
	"CAP_SETGID",
	"CAP_SETPCAP",
	"CAP_SETUID",
	"CAP_SYS_ADMIN",
	"CAP_SYS_BOOT",
	"CAP_SYS_CHROOT",
	"CAP_SYS_MODULE",
	"CAP_SYS_NICE",
	"CAP_SYS_PACCT",
	"CAP_SYS_PTRACE",
	"CAP_SYS_RAWIO",
	"CAP_SYS_RESOURCE",
	"CAP_SYS_TIME",
	"CAP_SYS_TTY_CONFIG",
	"CAP_SYSLOG",
	"CAP_WAKE_ALARM",
}

var SyscallGroups = []string{
	"@privileged",
	"@filesystem",
	"@network-io",
	"@network-bind",
	"@process",
	"@ipc",
	"@mount",
	"@chown",
	"@sync",
	"@clock",
	"@module",
	"@reboot",
	"@swap",
	"@resources",
	"@raw-io",
	"@sandbox",
	"@setuid",
	"@signal",
	"@timer",
	"@keyring",
	"@debug",
	"@pkey",
	"@bpf",
	"@aio",
	"@membarrier",
	"@memlock",
	"@default",
	"@basic-io",
	"@io-event",
}

func GetTagOptions(propertyKey string) []string {
	switch propertyKey {
	case "CapabilityBoundingSet", "AmbientCapabilities":
		return LinuxCapabilities
	case "SystemCallFilter":
		return SyscallGroups
	}
	return nil
}
