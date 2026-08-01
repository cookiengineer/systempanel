package systemdunit

type PropertyType int

const (
	PropString      PropertyType = iota
	PropBoolean
	PropEnum
	PropNumber
	PropFilePath
	PropFolderPath
	PropMultiLine
	PropTagList
	PropTarget
)

type PropertyMeta struct {
	Key        string
	PropType   PropertyType
	EnumValues []string
	Min        float64
	Max        float64
	Step       float64
	MimeTypes  []string
}

var actionEnums = []string{"none", "reboot", "reboot-force", "reboot-immediate", "poweroff", "poweroff-force", "poweroff-immediate", "exit", "exit-force", "soft-reboot", "soft-reboot-force", "kexec", "kexec-force", "halt", "halt-force"}

var jobModeEnums = []string{"replace", "fail", "replace-irreversibly", "isolate", "flush", "ignore-dependencies", "ignore-requirements"}

var collectModeEnums = []string{"inactive", "inactive-or-failed"}

var serviceTypeEnums = []string{"simple", "exec", "forking", "oneshot", "dbus", "notify", "idle"}

var exitTypeEnums = []string{"main", "cgroup"}

var restartEnums = []string{"no", "always", "on-success", "on-failure", "on-abnormal", "on-watchdog", "on-abort"}

var notifyAccessEnums = []string{"none", "main", "exec", "all"}

var oomPolicyEnums = []string{"continue", "stop", "kill", "abort"}

var privateUsersEnums = []string{"no", "yes", "self", "identity"}

var protectSystemEnums = []string{"no", "yes", "full", "strict"}

var protectHomeEnums = []string{"no", "yes", "read-only", "tmpfs"}

var privateTmpEnums = []string{"no", "yes", "disconnected"}

var devicePolicyEnums = []string{"auto", "closed", "strict"}

var protectProcEnums = []string{"default", "invisible", "ptraceable"}

var procSubsetEnums = []string{"all", "pid"}

var UnitProperties = []PropertyMeta{
	{Key: "Description", PropType: PropString},
	{Key: "Documentation", PropType: PropString},
	{Key: "Wants", PropType: PropTarget},
	{Key: "Requires", PropType: PropTarget},
	{Key: "Requisite", PropType: PropTarget},
	{Key: "BindsTo", PropType: PropTarget},
	{Key: "PartOf", PropType: PropTarget},
	{Key: "Upholds", PropType: PropTarget},
	{Key: "Conflicts", PropType: PropTarget},
	{Key: "Before", PropType: PropTarget},
	{Key: "After", PropType: PropTarget},
	{Key: "OnFailure", PropType: PropTarget},
	{Key: "OnSuccess", PropType: PropTarget},
	{Key: "OnFailureJobMode", PropType: PropEnum, EnumValues: jobModeEnums},
	{Key: "IgnoreOnIsolate", PropType: PropBoolean},
	{Key: "StopWhenUnneeded", PropType: PropBoolean},
	{Key: "RefuseManualStart", PropType: PropBoolean},
	{Key: "RefuseManualStop", PropType: PropBoolean},
	{Key: "AllowIsolate", PropType: PropBoolean},
	{Key: "DefaultDependencies", PropType: PropBoolean},
	{Key: "CollectMode", PropType: PropEnum, EnumValues: collectModeEnums},
	{Key: "FailureAction", PropType: PropEnum, EnumValues: actionEnums},
	{Key: "SuccessAction", PropType: PropEnum, EnumValues: actionEnums},
	{Key: "FailureActionExitStatus", PropType: PropString},
	{Key: "SuccessActionExitStatus", PropType: PropString},
	{Key: "JobTimeoutSec", PropType: PropNumber, Min: 0, Max: 0, Step: 1},
	{Key: "JobRunningTimeoutSec", PropType: PropNumber, Min: 0, Max: 0, Step: 1},
	{Key: "JobTimeoutAction", PropType: PropEnum, EnumValues: actionEnums},
	{Key: "JobTimeoutRebootArgument", PropType: PropString},
	{Key: "StartLimitIntervalSec", PropType: PropNumber, Min: 0, Max: 0, Step: 1},
	{Key: "StartLimitBurst", PropType: PropNumber, Min: 0, Max: 0, Step: 1},
	{Key: "StartLimitAction", PropType: PropEnum, EnumValues: actionEnums},
	{Key: "RebootArgument", PropType: PropString},
	{Key: "SourcePath", PropType: PropFilePath},
	{Key: "ConditionArchitecture", PropType: PropString},
	{Key: "ConditionVirtualization", PropType: PropString},
	{Key: "ConditionHost", PropType: PropString},
	{Key: "ConditionKernelCommandLine", PropType: PropString},
	{Key: "ConditionKernelVersion", PropType: PropString},
	{Key: "ConditionSecurity", PropType: PropString},
	{Key: "ConditionCapability", PropType: PropString},
	{Key: "ConditionACPower", PropType: PropString},
	{Key: "ConditionNeedsUpdate", PropType: PropString},
	{Key: "ConditionFirstBoot", PropType: PropString},
	{Key: "ConditionFileNotEmpty", PropType: PropFilePath},
	{Key: "ConditionFileIsExecutable", PropType: PropFilePath},
	{Key: "ConditionUser", PropType: PropString},
	{Key: "ConditionGroup", PropType: PropString},
	{Key: "ConditionControlGroupController", PropType: PropString},
	{Key: "AssertArchitecture", PropType: PropString},
	{Key: "AssertVirtualization", PropType: PropString},
	{Key: "AssertHost", PropType: PropString},
	{Key: "AssertKernelCommandLine", PropType: PropString},
	{Key: "AssertKernelVersion", PropType: PropString},
	{Key: "AssertSecurity", PropType: PropString},
	{Key: "AssertCapability", PropType: PropString},
	{Key: "AssertACPower", PropType: PropString},
	{Key: "AssertNeedsUpdate", PropType: PropString},
	{Key: "AssertFirstBoot", PropType: PropString},
	{Key: "AssertFileNotEmpty", PropType: PropFilePath},
	{Key: "AssertFileIsExecutable", PropType: PropFilePath},
	{Key: "AssertUser", PropType: PropString},
	{Key: "AssertGroup", PropType: PropString},
	{Key: "AssertControlGroupController", PropType: PropString},
}

var ServiceProperties = []PropertyMeta{
	{Key: "Type", PropType: PropEnum, EnumValues: serviceTypeEnums},
	{Key: "ExitType", PropType: PropEnum, EnumValues: exitTypeEnums},
	{Key: "RemainAfterExit", PropType: PropBoolean},
	{Key: "GuessMainPID", PropType: PropBoolean},
	{Key: "PIDFile", PropType: PropFilePath},
	{Key: "BusName", PropType: PropString},
	{Key: "ExecStart", PropType: PropFilePath},
	{Key: "ExecStartPre", PropType: PropFilePath},
	{Key: "ExecStartPost", PropType: PropFilePath},
	{Key: "ExecCondition", PropType: PropFilePath},
	{Key: "ExecReload", PropType: PropFilePath},
	{Key: "ExecStop", PropType: PropFilePath},
	{Key: "ExecStopPost", PropType: PropFilePath},
	{Key: "RestartSec", PropType: PropNumber, Min: 0, Max: 0, Step: 1},
	{Key: "TimeoutStartSec", PropType: PropNumber, Min: 0, Max: 0, Step: 1},
	{Key: "TimeoutStopSec", PropType: PropNumber, Min: 0, Max: 0, Step: 1},
	{Key: "TimeoutAbortSec", PropType: PropNumber, Min: 0, Max: 0, Step: 1},
	{Key: "TimeoutSec", PropType: PropNumber, Min: 0, Max: 0, Step: 1},
	{Key: "RuntimeMaxSec", PropType: PropNumber, Min: 0, Max: 0, Step: 1},
	{Key: "RuntimeRandomizedExtraSec", PropType: PropNumber, Min: 0, Max: 0, Step: 1},
	{Key: "WatchdogSec", PropType: PropNumber, Min: 0, Max: 0, Step: 1},
	{Key: "Restart", PropType: PropEnum, EnumValues: restartEnums},
	{Key: "SuccessExitStatus", PropType: PropString},
	{Key: "RestartPreventExitStatus", PropType: PropString},
	{Key: "RestartForceExitStatus", PropType: PropString},
	{Key: "RootDirectoryStartOnly", PropType: PropBoolean},
	{Key: "NonBlocking", PropType: PropBoolean},
	{Key: "NotifyAccess", PropType: PropEnum, EnumValues: notifyAccessEnums},
	{Key: "Sockets", PropType: PropString},
	{Key: "FileDescriptorStoreMax", PropType: PropNumber, Min: 0, Max: 0, Step: 1},
	{Key: "FileDescriptorStorePreserve", PropType: PropString},
	{Key: "USBFunctionDescriptors", PropType: PropString},
	{Key: "USBFunctionStrings", PropType: PropString},
	{Key: "OOMPolicy", PropType: PropEnum, EnumValues: oomPolicyEnums},
	{Key: "OpenFile", PropType: PropString},
	{Key: "ReloadSignal", PropType: PropString},
	{Key: "Environment", PropType: PropMultiLine},
	{Key: "EnvironmentFile", PropType: PropFilePath},
	{Key: "PassEnvironment", PropType: PropString},
	{Key: "UnsetEnvironment", PropType: PropString},
	{Key: "WorkingDirectory", PropType: PropFolderPath},
	{Key: "RootDirectory", PropType: PropFolderPath},
	{Key: "RootImage", PropType: PropFilePath},
	{Key: "RootImageOptions", PropType: PropString},
	{Key: "RootHash", PropType: PropString},
	{Key: "RootHashSignature", PropType: PropFilePath},
	{Key: "RootVerity", PropType: PropString},
	{Key: "MountAPIVFS", PropType: PropBoolean},
	{Key: "ProtectProc", PropType: PropEnum, EnumValues: protectProcEnums},
	{Key: "ProcSubset", PropType: PropEnum, EnumValues: procSubsetEnums},
	{Key: "BindPaths", PropType: PropString},
	{Key: "BindReadOnlyPaths", PropType: PropString},
	{Key: "User", PropType: PropString},
	{Key: "Group", PropType: PropString},
	{Key: "DynamicUser", PropType: PropBoolean},
	{Key: "SupplementaryGroups", PropType: PropString},
	{Key: "PAMName", PropType: PropString},
	{Key: "CapabilityBoundingSet", PropType: PropTagList},
	{Key: "AmbientCapabilities", PropType: PropTagList},
	{Key: "NoNewPrivileges", PropType: PropBoolean},
	{Key: "SecureBits", PropType: PropString},
	{Key: "SELinuxContext", PropType: PropString},
	{Key: "AppArmorProfile", PropType: PropString},
	{Key: "SmackProcessLabel", PropType: PropString},
	{Key: "LimitCPU", PropType: PropNumber, Min: 0, Max: 0, Step: 1},
	{Key: "LimitFSIZE", PropType: PropNumber, Min: 0, Max: 0, Step: 1},
	{Key: "LimitDATA", PropType: PropNumber, Min: 0, Max: 0, Step: 1},
	{Key: "LimitSTACK", PropType: PropNumber, Min: 0, Max: 0, Step: 1},
	{Key: "LimitCORE", PropType: PropNumber, Min: 0, Max: 0, Step: 1},
	{Key: "LimitRSS", PropType: PropNumber, Min: 0, Max: 0, Step: 1},
	{Key: "LimitNOFILE", PropType: PropNumber, Min: 0, Max: 0, Step: 1},
	{Key: "LimitAS", PropType: PropNumber, Min: 0, Max: 0, Step: 1},
	{Key: "LimitNPROC", PropType: PropNumber, Min: 0, Max: 0, Step: 1},
	{Key: "LimitMEMLOCK", PropType: PropNumber, Min: 0, Max: 0, Step: 1},
	{Key: "LimitLOCKS", PropType: PropNumber, Min: 0, Max: 0, Step: 1},
	{Key: "LimitSIGPENDING", PropType: PropNumber, Min: 0, Max: 0, Step: 1},
	{Key: "LimitMSGQUEUE", PropType: PropNumber, Min: 0, Max: 0, Step: 1},
	{Key: "LimitNICE", PropType: PropNumber, Min: -20, Max: 19, Step: 1},
	{Key: "LimitRTPRIO", PropType: PropNumber, Min: 0, Max: 0, Step: 1},
	{Key: "LimitRTTIME", PropType: PropNumber, Min: 0, Max: 0, Step: 1},
	{Key: "ReadWritePaths", PropType: PropString},
	{Key: "ReadOnlyPaths", PropType: PropString},
	{Key: "InaccessiblePaths", PropType: PropString},
	{Key: "ExecPaths", PropType: PropString},
	{Key: "NoExecPaths", PropType: PropString},
	{Key: "ExecSearchPath", PropType: PropString},
	{Key: "PrivateTmp", PropType: PropEnum, EnumValues: privateTmpEnums},
	{Key: "PrivateDevices", PropType: PropBoolean},
	{Key: "PrivateNetwork", PropType: PropBoolean},
	{Key: "PrivateUsers", PropType: PropEnum, EnumValues: privateUsersEnums},
	{Key: "PrivateMounts", PropType: PropBoolean},
	{Key: "PrivateIPC", PropType: PropBoolean},
	{Key: "PrivatePIDs", PropType: PropBoolean},
	{Key: "ProtectHome", PropType: PropEnum, EnumValues: protectHomeEnums},
	{Key: "ProtectSystem", PropType: PropEnum, EnumValues: protectSystemEnums},
	{Key: "ProtectHostname", PropType: PropBoolean},
	{Key: "ProtectKernelTunables", PropType: PropBoolean},
	{Key: "ProtectKernelModules", PropType: PropBoolean},
	{Key: "ProtectKernelLogs", PropType: PropBoolean},
	{Key: "ProtectClock", PropType: PropBoolean},
	{Key: "ProtectControlGroups", PropType: PropBoolean},
	{Key: "RestrictAddressFamilies", PropType: PropString},
	{Key: "RestrictFileSystems", PropType: PropString},
	{Key: "RestrictNetworkInterfaces", PropType: PropString},
	{Key: "LockPersonality", PropType: PropBoolean},
	{Key: "MemoryDenyWriteExecute", PropType: PropBoolean},
	{Key: "RestrictRealtime", PropType: PropBoolean},
	{Key: "RestrictSUIDSGID", PropType: PropBoolean},
	{Key: "RestrictNamespaces", PropType: PropString},
	{Key: "RemoveIPC", PropType: PropBoolean},
	{Key: "SystemCallFilter", PropType: PropTagList},
	{Key: "SystemCallArchitectures", PropType: PropString},
	{Key: "SystemCallErrorNumber", PropType: PropNumber, Min: 0, Max: 4095, Step: 1},
	{Key: "SystemCallLog", PropType: PropString},
	{Key: "MemoryMax", PropType: PropString},
	{Key: "MemoryHigh", PropType: PropString},
	{Key: "MemorySwapMax", PropType: PropString},
	{Key: "TasksMax", PropType: PropNumber, Min: 0, Max: 0, Step: 1},
	{Key: "IOWeight", PropType: PropNumber, Min: 1, Max: 10000, Step: 1},
	{Key: "IODeviceWeight", PropType: PropString},
	{Key: "IOReadBandwidthMax", PropType: PropString},
	{Key: "IOWriteBandwidthMax", PropType: PropString},
	{Key: "IOReadIOPSMax", PropType: PropString},
	{Key: "IOWriteIOPSMax", PropType: PropString},
	{Key: "CPUWeight", PropType: PropNumber, Min: 1, Max: 10000, Step: 1},
	{Key: "CPUShares", PropType: PropNumber, Min: 0, Max: 0, Step: 1},
	{Key: "CPUQuota", PropType: PropString},
	{Key: "Nice", PropType: PropNumber, Min: -20, Max: 19, Step: 1},
	{Key: "OOMScoreAdjust", PropType: PropNumber, Min: -1000, Max: 1000, Step: 1},
	{Key: "DevicePolicy", PropType: PropEnum, EnumValues: devicePolicyEnums},
	{Key: "DeviceAllow", PropType: PropString},
	{Key: "IPAccounting", PropType: PropBoolean},
	{Key: "CPUAccounting", PropType: PropBoolean},
	{Key: "MemoryAccounting", PropType: PropBoolean},
	{Key: "TasksAccounting", PropType: PropBoolean},
	{Key: "IOAccounting", PropType: PropBoolean},
	{Key: "TTYReset", PropType: PropBoolean},
	{Key: "TTYVHangup", PropType: PropBoolean},
	{Key: "TTYVTDisallocate", PropType: PropBoolean},
	{Key: "TemporaryFileSystem", PropType: PropString},
}

var SocketProperties = []PropertyMeta{
	{Key: "ListenStream", PropType: PropString},
	{Key: "ListenDatagram", PropType: PropString},
	{Key: "ListenSequentialPacket", PropType: PropString},
	{Key: "ListenFIFO", PropType: PropString},
	{Key: "ListenSpecial", PropType: PropString},
	{Key: "ListenNetlink", PropType: PropString},
	{Key: "ListenMessageQueue", PropType: PropString},
	{Key: "ListenUSBFunction", PropType: PropString},
	{Key: "SocketProtocol", PropType: PropString},
	{Key: "BindIPv6Only", PropType: PropString},
	{Key: "Backlog", PropType: PropNumber, Min: 0, Max: 0, Step: 1},
	{Key: "BindToDevice", PropType: PropString},
	{Key: "SocketUser", PropType: PropString},
	{Key: "SocketGroup", PropType: PropString},
	{Key: "SocketMode", PropType: PropNumber, Min: 0, Max: 511, Step: 1},
	{Key: "DirectoryMode", PropType: PropNumber, Min: 0, Max: 511, Step: 1},
	{Key: "Accept", PropType: PropBoolean},
	{Key: "Writable", PropType: PropBoolean},
	{Key: "FlushPending", PropType: PropBoolean},
	{Key: "MaxConnections", PropType: PropNumber, Min: 0, Max: 0, Step: 1},
	{Key: "MaxConnectionsPerSource", PropType: PropNumber, Min: 0, Max: 0, Step: 1},
	{Key: "KeepAlive", PropType: PropBoolean},
	{Key: "KeepAliveTimeSec", PropType: PropNumber, Min: 0, Max: 0, Step: 1},
	{Key: "KeepAliveIntervalSec", PropType: PropNumber, Min: 0, Max: 0, Step: 1},
	{Key: "KeepAliveProbes", PropType: PropNumber, Min: 0, Max: 0, Step: 1},
	{Key: "NoDelay", PropType: PropBoolean},
	{Key: "Priority", PropType: PropNumber, Min: 0, Max: 0, Step: 1},
	{Key: "DeferAcceptSec", PropType: PropNumber, Min: 0, Max: 0, Step: 1},
	{Key: "ReceiveBuffer", PropType: PropNumber, Min: 0, Max: 0, Step: 1},
	{Key: "SendBuffer", PropType: PropNumber, Min: 0, Max: 0, Step: 1},
	{Key: "IPTOS", PropType: PropNumber, Min: 0, Max: 0, Step: 1},
	{Key: "IPTTL", PropType: PropNumber, Min: 0, Max: 0, Step: 1},
	{Key: "PipeSize", PropType: PropNumber, Min: 0, Max: 0, Step: 1},
	{Key: "FreeBind", PropType: PropBoolean},
	{Key: "Transparent", PropType: PropBoolean},
	{Key: "Broadcast", PropType: PropBoolean},
	{Key: "PassCredentials", PropType: PropBoolean},
	{Key: "PassSecurity", PropType: PropBoolean},
	{Key: "PassPacketInfo", PropType: PropBoolean},
	{Key: "Timestamping", PropType: PropString},
	{Key: "RemoveOnStop", PropType: PropBoolean},
	{Key: "RemoveOnUnlink", PropType: PropBoolean},
	{Key: "Symlinks", PropType: PropString},
	{Key: "Mark", PropType: PropNumber, Min: 0, Max: 0, Step: 1},
	{Key: "TriggerLimitIntervalSec", PropType: PropNumber, Min: 0, Max: 0, Step: 1},
	{Key: "TriggerLimitBurst", PropType: PropNumber, Min: 0, Max: 0, Step: 1},
}

var InstallProperties = []PropertyMeta{
	{Key: "WantedBy", PropType: PropTarget},
	{Key: "RequiredBy", PropType: PropTarget},
	{Key: "UpheldBy", PropType: PropTarget},
	{Key: "Alias", PropType: PropString},
	{Key: "Also", PropType: PropString},
	{Key: "DefaultInstance", PropType: PropString},
}

var propertyMap = func() map[string]map[string]PropertyMeta {
	m := make(map[string]map[string]PropertyMeta)
	for _, list := range []struct {
		section string
		props   []PropertyMeta
	}{
		{"Unit", UnitProperties},
		{"Service", ServiceProperties},
		{"Socket", SocketProperties},
		{"Install", InstallProperties},
	} {
		inner := make(map[string]PropertyMeta)
		for _, p := range list.props {
			inner[p.Key] = p
		}
		m[list.section] = inner
	}
	return m
}()

func GetPropertyMeta(section, key string) (PropertyMeta, bool) {
	if sm, ok := propertyMap[section]; ok {
		if p, ok := sm[key]; ok {
			return p, true
		}
	}
	return PropertyMeta{}, false
}
