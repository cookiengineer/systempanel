package systemdunit

type PropertyGroup struct {
	Name string
	Keys []string
}

var UnitGroups = []PropertyGroup{
	{
		Name: "Description", Keys: []string{"Description", "Documentation"},
	},
	{
		Name: "Dependencies",
		Keys: []string{"Wants", "Requires", "Requisite", "BindsTo", "PartOf", "Upholds", "Conflicts", "Before", "After", "DefaultDependencies"},
	},
	{
		Name: "Error Handling",
		Keys: []string{"OnFailure", "OnSuccess", "OnFailureJobMode", "FailureAction", "SuccessAction", "FailureActionExitStatus", "SuccessActionExitStatus"},
	},
	{
		Name: "Lifecycle",
		Keys: []string{"IgnoreOnIsolate", "StopWhenUnneeded", "RefuseManualStart", "RefuseManualStop", "AllowIsolate"},
	},
	{
		Name: "Job Timeouts",
		Keys: []string{"CollectMode", "JobTimeoutSec", "JobRunningTimeoutSec", "JobTimeoutAction", "JobTimeoutRebootArgument"},
	},
	{
		Name: "Start Limits",
		Keys: []string{"StartLimitIntervalSec", "StartLimitBurst", "StartLimitAction", "RebootArgument"},
	},
	{
		Name: "Source", Keys: []string{"SourcePath"},
	},
	{
		Name: "Conditions",
		Keys: []string{"ConditionArchitecture", "ConditionVirtualization", "ConditionHost", "ConditionKernelCommandLine", "ConditionKernelVersion", "ConditionSecurity", "ConditionCapability", "ConditionACPower", "ConditionNeedsUpdate", "ConditionFirstBoot", "ConditionFileNotEmpty", "ConditionFileIsExecutable", "ConditionUser", "ConditionGroup", "ConditionControlGroupController"},
	},
	{
		Name: "Assertions",
		Keys: []string{"AssertArchitecture", "AssertVirtualization", "AssertHost", "AssertKernelCommandLine", "AssertKernelVersion", "AssertSecurity", "AssertCapability", "AssertACPower", "AssertNeedsUpdate", "AssertFirstBoot", "AssertFileNotEmpty", "AssertFileIsExecutable", "AssertUser", "AssertGroup", "AssertControlGroupController"},
	},
}

var ServiceGroups = []PropertyGroup{
	{
		Name: "Service Type",
		Keys: []string{"Type", "ExitType", "RemainAfterExit", "GuessMainPID", "PIDFile", "BusName"},
	},
	{
		Name: "Execution",
		Keys: []string{"ExecStart", "ExecStartPre", "ExecStartPost", "ExecCondition", "ExecReload", "ExecStop", "ExecStopPost"},
	},
	{
		Name: "Restart & Timeouts",
		Keys: []string{"Restart", "RestartSec", "TimeoutStartSec", "TimeoutStopSec", "TimeoutAbortSec", "TimeoutSec", "RuntimeMaxSec", "RuntimeRandomizedExtraSec", "WatchdogSec"},
	},
	{
		Name: "User & Group",
		Keys: []string{"User", "Group", "SupplementaryGroups", "DynamicUser", "PAMName"},
	},
	{
		Name: "Namespacing",
		Keys: []string{"PrivateUsers", "RootDirectory", "RootImage", "WorkingDirectory"},
	},
	{
		Name: "Filesystem Sandboxing",
		Keys: []string{"ProtectSystem", "ProtectHome", "TemporaryFileSystem", "ReadOnlyPaths", "ReadWritePaths", "InaccessiblePaths", "BindPaths", "BindReadOnlyPaths"},
	},
	{
		Name: "Runtime Directories",
		Keys: []string{"RuntimeDirectory", "StateDirectory", "CacheDirectory", "LogsDirectory", "ConfigurationDirectory"},
	},
	{
		Name: "Process Isolation",
		Keys: []string{"PrivateTmp", "PrivateDevices", "PrivateMounts", "PrivateIPC", "PrivatePIDs", "PrivateNetwork"},
	},
	{
		Name: "Device Access",
		Keys: []string{"DevicePolicy", "DeviceAllow"},
	},
	{
		Name: "Kernel Protection",
		Keys: []string{"ProtectKernelModules", "ProtectKernelTunables", "ProtectKernelLogs", "ProtectControlGroups", "ProtectClock", "ProtectHostname"},
	},
	{
		Name: "/proc Restrictions",
		Keys: []string{"ProtectProc", "ProcSubset", "MountAPIVFS"},
	},
	{
		Name: "Capabilities",
		Keys: []string{"CapabilityBoundingSet", "AmbientCapabilities", "NoNewPrivileges"},
	},
	{
		Name: "System Call Filtering",
		Keys: []string{"SystemCallFilter", "SystemCallArchitectures", "SystemCallErrorNumber"},
	},
	{
		Name: "Further Restrictions",
		Keys: []string{"RestrictNamespaces", "RestrictSUIDSGID", "RestrictRealtime", "RestrictAddressFamilies", "RestrictNetworkInterfaces", "RestrictFileSystems", "LockPersonality", "MemoryDenyWriteExecute", "RemoveIPC"},
	},
	{
		Name: "Resource Limits",
		Keys: []string{"LimitCPU", "LimitFSIZE", "LimitDATA", "LimitSTACK", "LimitCORE", "LimitRSS", "LimitNOFILE", "LimitAS", "LimitNPROC", "LimitMEMLOCK", "LimitLOCKS", "LimitSIGPENDING", "LimitMSGQUEUE", "LimitNICE", "LimitRTPRIO", "LimitRTTIME"},
	},
	{
		Name: "Memory & Tasks",
		Keys: []string{"MemoryMax", "MemoryHigh", "MemorySwapMax", "TasksMax", "IOWeight", "IODeviceWeight", "IOReadBandwidthMax", "IOWriteBandwidthMax", "IOReadIOPSMax", "IOWriteIOPSMax", "CPUWeight", "CPUShares", "CPUQuota", "Nice", "OOMScoreAdjust"},
	},
	{
		Name: "Accounting",
		Keys: []string{"IPAccounting", "CPUAccounting", "MemoryAccounting", "TasksAccounting", "IOAccounting"},
	},
	{
		Name: "Environment",
		Keys: []string{"Environment", "EnvironmentFile", "PassEnvironment", "UnsetEnvironment"},
	},
	{
		Name: "Notify & Sockets",
		Keys: []string{"NotifyAccess", "Sockets", "FileDescriptorStoreMax", "FileDescriptorStorePreserve", "OpenFile", "ReloadSignal"},
	},
	{
		Name: "Exit Codes",
		Keys: []string{"SuccessExitStatus", "RestartPreventExitStatus", "RestartForceExitStatus"},
	},
	{
		Name: "TTY",
		Keys: []string{"TTYReset", "TTYVHangup", "TTYVTDisallocate"},
	},
	{
		Name: "Security Labels",
		Keys: []string{"SELinuxContext", "AppArmorProfile", "SmackProcessLabel", "SecureBits"},
	},
	{
		Name: "Root Image Verification",
		Keys: []string{"RootImageOptions", "RootHash", "RootHashSignature", "RootVerity"},
	},
	{
		Name: "OOMPolicy & Misc",
		Keys: []string{"OOMPolicy", "NonBlocking", "RootDirectoryStartOnly", "USBFunctionDescriptors", "USBFunctionStrings"},
	},
}

var SocketGroups = []PropertyGroup{
	{
		Name: "Listen Addresses",
		Keys: []string{"ListenStream", "ListenDatagram", "ListenSequentialPacket", "ListenFIFO", "ListenSpecial", "ListenNetlink", "ListenMessageQueue", "ListenUSBFunction"},
	},
	{
		Name: "Socket Options",
		Keys: []string{"SocketProtocol", "BindIPv6Only", "Backlog", "BindToDevice", "SocketUser", "SocketGroup", "SocketMode", "DirectoryMode"},
	},
	{
		Name: "Connection Handling",
		Keys: []string{"Accept", "Writable", "FlushPending", "MaxConnections", "MaxConnectionsPerSource"},
	},
	{
		Name: "Keep-Alive",
		Keys: []string{"KeepAlive", "KeepAliveTimeSec", "KeepAliveIntervalSec", "KeepAliveProbes"},
	},
	{
		Name: "Network Tuning",
		Keys: []string{"NoDelay", "Priority", "DeferAcceptSec", "ReceiveBuffer", "SendBuffer", "IPTOS", "IPTTL", "PipeSize"},
	},
	{
		Name: "Permissions & Options",
		Keys: []string{"FreeBind", "Transparent", "Broadcast", "PassCredentials", "PassSecurity", "PassPacketInfo"},
	},
	{
		Name: "Cleanup",
		Keys: []string{"Timestamping", "RemoveOnStop", "RemoveOnUnlink", "Symlinks", "Mark"},
	},
	{
		Name: "Trigger Limits",
		Keys: []string{"TriggerLimitIntervalSec", "TriggerLimitBurst"},
	},
}

var InstallGroups = []PropertyGroup{
	{
		Name: "Installation Targets",
		Keys: []string{"WantedBy", "RequiredBy", "UpheldBy"},
	},
	{
		Name: "Aliases",
		Keys: []string{"Alias", "Also", "DefaultInstance"},
	},
}

var sectionGroups = map[string][]PropertyGroup{
	"Unit":    UnitGroups,
	"Service": ServiceGroups,
	"Socket":  SocketGroups,
	"Install": InstallGroups,
}

func GetSectionGroups(section string) []PropertyGroup {
	return sectionGroups[section]
}
