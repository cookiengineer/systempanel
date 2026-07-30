package systemdunit

import (
	"bufio"
	"os"
	"sort"
	"strings"
)

type UnitFile struct {
	Path     string
	Sections map[string]Section
}

type Section struct {
	Name     string
	Keys     map[string]string
	KeyOrder []string
}

func Parse(path string) (*UnitFile, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	uf := &UnitFile{
		Path:     path,
		Sections: make(map[string]Section),
	}

	var current *Section
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") || strings.HasPrefix(line, ";") {
			continue
		}
		if strings.HasPrefix(line, "[") && strings.HasSuffix(line, "]") {
			name := line[1 : len(line)-1]
			if sec, ok := uf.Sections[name]; ok {
				current = &sec
			} else {
				current = &Section{
					Name: name,
					Keys: make(map[string]string),
				}
				uf.Sections[name] = *current
			}
			continue
		}
		if current == nil {
			continue
		}
		eq := strings.Index(line, "=")
		if eq < 0 {
			continue
		}
		key := strings.TrimSpace(line[:eq])
		val := strings.TrimSpace(line[eq+1:])
		sec := uf.Sections[current.Name]
		if _, exists := sec.Keys[key]; !exists {
			sec.KeyOrder = append(sec.KeyOrder, key)
		}
		sec.Keys[key] = val
		uf.Sections[current.Name] = sec
	}

	return uf, scanner.Err()
}

func (uf *UnitFile) Serialize() string {
	var b strings.Builder

	sectionOrder := []string{"Unit", "Service", "Socket", "Timer", "Mount", "Automount", "Swap", "Path", "Slice", "Scope", "Install"}

	for _, name := range sectionOrder {
		sec, ok := uf.Sections[name]
		if !ok || len(sec.Keys) == 0 {
			continue
		}
		b.WriteString("[")
		b.WriteString(name)
		b.WriteString("]\n")
		for _, key := range sec.KeyOrder {
			b.WriteString(key)
			b.WriteString("=")
			b.WriteString(sec.Keys[key])
			b.WriteString("\n")
		}
		b.WriteString("\n")
	}

	for _, sec := range uf.Sections {
		if contains(sectionOrder, sec.Name) {
			continue
		}
		if len(sec.Keys) == 0 {
			continue
		}
		b.WriteString("[")
		b.WriteString(sec.Name)
		b.WriteString("]\n")
		for _, key := range sec.KeyOrder {
			b.WriteString(key)
			b.WriteString("=")
			b.WriteString(sec.Keys[key])
			b.WriteString("\n")
		}
		b.WriteString("\n")
	}

	return b.String()
}

func (uf *UnitFile) Get(section, key string) string {
	if sec, ok := uf.Sections[section]; ok {
		return sec.Keys[key]
	}
	return ""
}

func (uf *UnitFile) Set(section, key, value string) {
	sec, ok := uf.Sections[section]
	if !ok {
		sec = Section{Name: section, Keys: make(map[string]string)}
	}
	if value == "" {
		delete(sec.Keys, key)
		sec.KeyOrder = removeKey(sec.KeyOrder, key)
	} else {
		if _, exists := sec.Keys[key]; !exists {
			sec.KeyOrder = append(sec.KeyOrder, key)
		}
		sec.Keys[key] = value
	}
	uf.Sections[section] = sec
}

func (uf *UnitFile) SectionKeys(section string) []string {
	if sec, ok := uf.Sections[section]; ok {
		return sec.KeyOrder
	}
	return nil
}

func KnownSectionKeys(section string) []string {
	switch section {
	case "Unit":
		return []string{"Description", "Documentation", "Wants", "Requires", "Requisite", "BindsTo", "PartOf", "Upholds", "Conflicts", "Before", "After", "OnFailure", "OnSuccess", "OnFailureJobMode", "IgnoreOnIsolate", "StopWhenUnneeded", "RefuseManualStart", "RefuseManualStop", "AllowIsolate", "DefaultDependencies", "CollectMode", "FailureAction", "SuccessAction", "FailureActionExitStatus", "SuccessActionExitStatus", "JobTimeoutSec", "JobRunningTimeoutSec", "JobTimeoutAction", "JobTimeoutRebootArgument", "StartLimitIntervalSec", "StartLimitBurst", "StartLimitAction", "RebootArgument", "SourcePath", "ConditionArchitecture", "ConditionVirtualization", "ConditionHost", "ConditionKernelCommandLine", "ConditionKernelVersion", "ConditionSecurity", "ConditionCapability", "ConditionACPower", "ConditionNeedsUpdate", "ConditionFirstBoot", "ConditionFileNotEmpty", "ConditionFileIsExecutable", "ConditionUser", "ConditionGroup", "ConditionControlGroupController", "AssertArchitecture", "AssertVirtualization", "AssertHost", "AssertKernelCommandLine", "AssertKernelVersion", "AssertSecurity", "AssertCapability", "AssertACPower", "AssertNeedsUpdate", "AssertFirstBoot", "AssertFileNotEmpty", "AssertFileIsExecutable", "AssertUser", "AssertGroup", "AssertControlGroupController"}
	case "Service":
		return []string{"Type", "ExitType", "RemainAfterExit", "GuessMainPID", "PIDFile", "BusName", "ExecStart", "ExecStartPre", "ExecStartPost", "ExecCondition", "ExecReload", "ExecStop", "ExecStopPost", "RestartSec", "TimeoutStartSec", "TimeoutStopSec", "TimeoutAbortSec", "TimeoutSec", "RuntimeMaxSec", "RuntimeRandomizedExtraSec", "WatchdogSec", "Restart", "SuccessExitStatus", "RestartPreventExitStatus", "RestartForceExitStatus", "RootDirectoryStartOnly", "NonBlocking", "NotifyAccess", "Sockets", "FileDescriptorStoreMax", "FileDescriptorStorePreserve", "USBFunctionDescriptors", "USBFunctionStrings", "OOMPolicy", "OpenFile", "ReloadSignal", "Environment", "EnvironmentFile", "PassEnvironment", "UnsetEnvironment", "WorkingDirectory", "RootDirectory", "RootImage", "RootImageOptions", "RootHash", "RootHashSignature", "RootVerity", "MountAPIVFS", "ProtectProc", "ProcSubset", "BindPaths", "BindReadOnlyPaths", "User", "Group", "DynamicUser", "SupplementaryGroups", "PAMName", "CapabilityBoundingSet", "AmbientCapabilities", "NoNewPrivileges", "SecureBits", "SELinuxContext", "AppArmorProfile", "SmackProcessLabel", "LimitCPU", "LimitFSIZE", "LimitDATA", "LimitSTACK", "LimitCORE", "LimitRSS", "LimitNOFILE", "LimitAS", "LimitNPROC", "LimitMEMLOCK", "LimitLOCKS", "LimitSIGPENDING", "LimitMSGQUEUE", "LimitNICE", "LimitRTPRIO", "LimitRTTIME", "ReadWritePaths", "ReadOnlyPaths", "InaccessiblePaths", "ExecPaths", "NoExecPaths", "ExecSearchPath", "PrivateTmp", "PrivateDevices", "PrivateNetwork", "PrivateUsers", "PrivateMounts", "PrivateIPC", "ProtectHome", "ProtectSystem", "ProtectHostname", "ProtectKernelTunables", "ProtectKernelModules", "ProtectKernelLogs", "ProtectClock", "ProtectControlGroups", "RestrictAddressFamilies", "RestrictFileSystems", "LockPersonality", "MemoryDenyWriteExecute", "RestrictRealtime", "RestrictSUIDSGID", "RemoveIPC", "SystemCallFilter", "SystemCallArchitectures", "SystemCallErrorNumber", "SystemCallLog", "MemoryMax", "MemoryHigh", "MemorySwapMax", "TasksMax", "IOWeight", "IOReadBandwidthMax", "IOWriteBandwidthMax", "IOReadIOPSMax", "IOWriteIOPSMax", "DeviceAllow", "DevicePolicy", "Slice", "Delegate", "CPUAccounting", "CPUWeight", "StartupCPUWeight", "CPUShares", "StartupCPUShares", "CPUQuotaPeriodSec", "CPUQuota", "AllowedCPUs", "AllowedMemoryNodes", "MemoryAccounting", "MemoryMin", "MemoryLow", "MemoryCurrent", "MemoryZSwapMax", "TasksAccounting", "IOAccounting", "IPAccounting", "ManagedOOMMemoryPressure", "ManagedOOMMemoryPressureLimit", "ManagedOOMSwap", "ManagedOOMPreference", "BPFProgram", "SocketBindAllow", "SocketBindDeny", "RestrictNamespaces", "MountImage", "ExtensionDirectories", "ExtensionImages", "LogLevelMax", "LogRateLimitIntervalSec", "LogRateLimitBurst", "LogExtraFields", "LogFilterPatterns", "LogNamespace", "StandardInput", "StandardOutput", "StandardError", "TTYPath", "TTYReset", "TTYVHangup", "TTYVTDisallocate", "SyslogIdentifier", "SyslogFacility", "SyslogLevel", "SyslogLevelPrefix", "UtmpIdentifier", "UtmpMode"}
	case "Install":
		return []string{"WantedBy", "RequiredBy", "UpheldBy", "Alias", "Also", "DefaultInstance"}
	case "Socket":
		return []string{"ListenStream", "ListenDatagram", "ListenSequentialPacket", "ListenFIFO", "ListenSpecial", "ListenNetlink", "ListenMessageQueue", "ListenUSBFunction", "SocketProtocol", "BindIPv6Only", "Backlog", "BindToDevice", "SocketUser", "SocketGroup", "SocketMode", "DirectoryMode", "Accept", "Writable", "FlushPending", "MaxConnections", "MaxConnectionsPerSource", "KeepAlive", "KeepAliveTimeSec", "KeepAliveIntervalSec", "KeepAliveProbes", "NoDelay", "Priority", "DeferAcceptSec", "ReceiveBuffer", "SendBuffer", "IPTOS", "IPTTL", "PipeSize", "FreeBind", "Transparent", "Broadcast", "PassCredentials", "PassSecurity", "PassPacketInfo", "Timestamping", "RemoveOnStop", "RemoveOnUnlink", "Symlinks", "Mark", "MaxConnectionsPerSource", "TriggerLimitIntervalSec", "TriggerLimitBurst"}
	}
	return nil
}

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

func removeKey(keys []string, key string) []string {
	var result []string
	for _, k := range keys {
		if k != key {
			result = append(result, k)
		}
	}
	return result
}

func SortKeysByKnown(section string, keys []string) []string {
	known := KnownSectionKeys(section)
	if known == nil {
		sort.Strings(keys)
		return keys
	}
	knownSet := make(map[string]bool)
	for _, k := range known {
		knownSet[k] = true
	}
	var knownKeys, unknown []string
	for _, k := range keys {
		if knownSet[k] {
			knownKeys = append(knownKeys, k)
		} else {
			unknown = append(unknown, k)
		}
	}
	sort.Strings(unknown)
	result := make([]string, 0, len(knownKeys)+len(unknown))
	for _, k := range known {
		if inSlice(knownKeys, k) {
			result = append(result, k)
		}
	}
	result = append(result, unknown...)
	return result
}

func inSlice(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}
