// SPDX-License-Identifier: Apache-2.0
// Copyright Authors of Cilium

package iptables

import (
	"github.com/spf13/pflag"

	"github.com/cilium/cilium/pkg/cidr"
	"github.com/cilium/cilium/pkg/hive/cell"
	"github.com/cilium/cilium/pkg/option"
	"github.com/cilium/cilium/pkg/time"
)

var Cell = cell.Module(
	"iptables",
	"Handle iptables-related configuration for Cilium",

	cell.Config(defaultConfig),
	cell.ProvidePrivate(func(
		cfg *option.DaemonConfig,
	) SharedConfig {
		return SharedConfig{
			TunnelingEnabled:                cfg.TunnelingEnabled(),
			NodeIpsetNeeded:                 cfg.NodeIpsetNeeded(),
			Devices:                         cfg.GetDevices(),
			IptablesMasqueradingIPv4Enabled: cfg.IptablesMasqueradingIPv4Enabled(),
			IptablesMasqueradingIPv6Enabled: cfg.IptablesMasqueradingIPv6Enabled(),
			IPv4NativeRoutingCIDR:           cfg.GetIPv4NativeRoutingCIDR(),

			EnableIPv4:                  cfg.EnableIPv4,
			EnableIPv6:                  cfg.EnableIPv6,
			EnableXTSocketFallback:      cfg.EnableXTSocketFallback,
			EnableBPFTProxy:             cfg.EnableBPFTProxy,
			InstallNoConntrackIptRules:  cfg.InstallNoConntrackIptRules,
			EnableEndpointRoutes:        cfg.EnableEndpointRoutes,
			IPAM:                        cfg.IPAM,
			EnableIPSec:                 cfg.EnableIPSec,
			MasqueradeInterfaces:        cfg.MasqueradeInterfaces,
			EnableMasqueradeRouteSource: cfg.EnableMasqueradeRouteSource,
		}
	}),
	cell.Provide(newIptablesManager),
)

type Config struct {
	// IPTablesLockTimeout defines the "-w" iptables option when the
	// iptables CLI is directly invoked from the Cilium agent.
	IPTablesLockTimeout time.Duration

	// DisableIptablesFeederRules specifies which chains will be excluded
	// when installing the feeder rules
	DisableIptablesFeederRules []string

	// IPTablesRandomFully defines the "--random-fully" iptables option when the
	// iptables CLI is directly invoked from the Cilium agent.
	IPTablesRandomFully bool
}

var defaultConfig = Config{
	IPTablesLockTimeout: 5 * time.Second,
}

func (def Config) Flags(flags *pflag.FlagSet) {
	flags.Duration("iptables-lock-timeout", def.IPTablesLockTimeout, "Time to pass to each iptables invocation to wait for xtables lock acquisition")
	flags.StringSlice("disable-iptables-feeder-rules", def.DisableIptablesFeederRules, "Chains to ignore when installing feeder rules.")
	flags.Bool("iptables-random-fully", def.IPTablesRandomFully, "Set iptables flag random-fully on masquerading rules")
}

type SharedConfig struct {
	TunnelingEnabled                bool
	NodeIpsetNeeded                 bool
	Devices                         []string
	IptablesMasqueradingIPv4Enabled bool
	IptablesMasqueradingIPv6Enabled bool
	IPv4NativeRoutingCIDR           *cidr.CIDR

	EnableIPv4                  bool
	EnableIPv6                  bool
	EnableXTSocketFallback      bool
	EnableBPFTProxy             bool
	InstallNoConntrackIptRules  bool
	EnableEndpointRoutes        bool
	IPAM                        string
	EnableIPSec                 bool
	MasqueradeInterfaces        []string
	EnableMasqueradeRouteSource bool
}
