package cni

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/containernetworking/cni/pkg/invoke"
	"github.com/containernetworking/cni/pkg/skel"
	"github.com/containernetworking/cni/pkg/types"
	type100 "github.com/containernetworking/cni/pkg/types/100"
	"github.com/containernetworking/cni/pkg/version"

	"github.com/datum-cloud/galactic/cni/route"
	"github.com/datum-cloud/galactic/cni/veth"
	"github.com/datum-cloud/galactic/cni/vrf"
	"github.com/datum-cloud/galactic/debug"
	"github.com/datum-cloud/galactic/util"
)

type Termination struct {
	Network string `json:"network"`
	Via     string `json:"via,omitempty"`
}

type PluginConf struct {
	types.PluginConf
	VPC           string        `json:"vpc"`
	VPCAttachment string        `json:"vpcattachment"`
	MTU           int           `json:"mtu,omitempty"`
	Terminations  []Termination `json:"terminations,omitempty"`
}

func NewCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "galactic-cni",
		Short: "Galactic CNI plugin",
		Run: func(cmd *cobra.Command, args []string) {
			skel.PluginMainFuncs(
				skel.CNIFuncs{
					Add: cmdAdd,
					Del: cmdDel,
				},
				version.All,
				fmt.Sprintf("CNI galactic plugin %s", debug.Version()),
			)
		},
	}
}

func parseConf(data []byte) (*PluginConf, error) {
	conf := &PluginConf{}
	if err := json.Unmarshal(data, &conf); err != nil {
		return nil, err
	}
	return conf, nil
}

func cmdAdd(args *skel.CmdArgs) error {
	pluginConf, _ := parseConf(args.StdinData)
	if err := vrf.Add(pluginConf.VPC, pluginConf.VPCAttachment); err != nil {
		return err
	}
	if err := veth.Add(pluginConf.VPC, pluginConf.VPCAttachment, pluginConf.MTU); err != nil {
		return err
	}
	dev := util.GenerateInterfaceNameHost(pluginConf.VPC, pluginConf.VPCAttachment)
	for _, termination := range pluginConf.Terminations {
		if err := route.Add(pluginConf.VPC, pluginConf.VPCAttachment, termination.Network, termination.Via, dev); err != nil {
			return err
		}
	}
	if err := hostDevice("ADD", args, pluginConf); err != nil {
		return err
	}
	result := &type100.Result{}
	return types.PrintResult(result, pluginConf.CNIVersion)
}

func cmdDel(args *skel.CmdArgs) error {
	pluginConf, _ := parseConf(args.StdinData)
	dev := util.GenerateInterfaceNameHost(pluginConf.VPC, pluginConf.VPCAttachment)
	if err := hostDevice("DEL", args, pluginConf); err != nil {
		return err
	}
	for _, termination := range pluginConf.Terminations {
		if err := route.Delete(pluginConf.VPC, pluginConf.VPCAttachment, termination.Network, termination.Via, dev); err != nil {
			return err
		}
	}
	if err := veth.Delete(pluginConf.VPC, pluginConf.VPCAttachment, pluginConf.MTU); err != nil {
		return err
	}
	if err := vrf.Delete(pluginConf.VPC, pluginConf.VPCAttachment); err != nil {
		return err
	}
	result := &type100.Result{}
	return types.PrintResult(result, pluginConf.CNIVersion)
}

type HostDevicePluginConf struct {
	types.PluginConf
	Device string          `json:"device"`
	IPAM   json.RawMessage `json:"ipam,omitempty"`
}

func getIPAM(input []byte) json.RawMessage {
	var m map[string]json.RawMessage
	if err := json.Unmarshal(input, &m); err != nil {
		return nil
	}
	ipam, ok := m["ipam"]
	if !ok {
		return nil
	}
	return ipam
}

func hostDeviceExecutable() string {
	path, _ := os.Executable()
	dir := filepath.Dir(path)
	return filepath.Join(dir, "host-device")
}

func hostDevice(command string, skelArgs *skel.CmdArgs, pluginConf *PluginConf) error {
	conf, err := json.Marshal(HostDevicePluginConf{
		PluginConf: types.PluginConf{
			CNIVersion: pluginConf.CNIVersion,
			Name:       pluginConf.Name,
			Type:       "host-device",
		},
		Device: util.GenerateInterfaceNameGuest(pluginConf.VPC, pluginConf.VPCAttachment),
		IPAM:   getIPAM(skelArgs.StdinData),
	})
	if err != nil {
		return err
	}

	invokeExec := &invoke.DefaultExec{
		RawExec:       &invoke.RawExec{Stderr: os.Stderr},
		PluginDecoder: version.PluginDecoder{},
	}
	invokeArgs := invoke.Args{
		Command:       command,
		ContainerID:   skelArgs.ContainerID,
		NetNS:         skelArgs.Netns,
		PluginArgsStr: skelArgs.Args,
		IfName:        skelArgs.IfName,
		Path:          skelArgs.Path,
	}
	if _, err := invokeExec.ExecPlugin(context.Background(), hostDeviceExecutable(), conf, invokeArgs.AsEnv()); err != nil {
		return err
	}
	return nil
}
