package cni

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/containernetworking/cni/pkg/skel"
	"github.com/containernetworking/cni/pkg/types"
	type100 "github.com/containernetworking/cni/pkg/types/100"
	"github.com/containernetworking/cni/pkg/version"
	bv "github.com/containernetworking/plugins/pkg/utils/buildversion"

	"github.com/datum-cloud/galactic/cni/route"
	"github.com/datum-cloud/galactic/cni/veth"
	"github.com/datum-cloud/galactic/cni/vrf"
)

type Termination struct {
	Network string `json:"network"`
	Via     string `json:"via,omitempty"`
}

type PluginConf struct {
	types.PluginConf
	Id           int           `json:"id"`
	MTU          int           `json:"mtu,omitempty"`
	Terminations []Termination `json:"terminations,omitempty"`
}

func NewCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "galactic-cni",
		Short: "Galactic CNI plugin",
		Run: func(cmd *cobra.Command, args []string) {
			skel.PluginMainFuncs(skel.CNIFuncs{
				Add: cmdAdd,
				Del: cmdDel,
			}, version.All, bv.BuildString("none"))
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
	if err := vrf.Add(pluginConf.Id); err != nil {
		return err
	}
	if err := veth.Add(pluginConf.Id, pluginConf.MTU); err != nil {
		return err
	}
	dev := fmt.Sprintf(veth.VethNameTemplateHost, pluginConf.Id)
	for _, termination := range pluginConf.Terminations {
		if err := route.Add(pluginConf.Id, termination.Network, termination.Via, dev); err != nil {
			return err
		}
	}
	result := &type100.Result{}
	return types.PrintResult(result, pluginConf.CNIVersion)
}

func cmdDel(args *skel.CmdArgs) error {
	pluginConf, _ := parseConf(args.StdinData)
	dev := fmt.Sprintf(veth.VethNameTemplateHost, pluginConf.Id)
	for _, termination := range pluginConf.Terminations {
		if err := route.Delete(pluginConf.Id, termination.Network, termination.Via, dev); err != nil {
			return err
		}
	}
	if err := veth.Delete(pluginConf.Id, pluginConf.MTU); err != nil {
		return err
	}
	if err := vrf.Delete(pluginConf.Id); err != nil {
		return err
	}
	result := &type100.Result{}
	return types.PrintResult(result, pluginConf.CNIVersion)
}
