package vrf

import (
	"github.com/datum-cloud/galactic/cni/sysctl"
	"github.com/datum-cloud/galactic/util"
	"github.com/vishvananda/netlink"
)

func Add(id int) error {
	name := util.GenerateInterfaceNameVRF(id)

	vrf := &netlink.Vrf{
		LinkAttrs: netlink.LinkAttrs{
			Name: name,
		},
		Table: uint32(id),
	}

	if err := netlink.LinkAdd(vrf); err != nil {
		return err
	}

	if err := sysctl.ConfigureInterfaceSysctls(name); err != nil {
		return err
	}

	return netlink.LinkSetUp(vrf)
}

func Delete(id int) error {
	link, err := netlink.LinkByName(util.GenerateInterfaceNameVRF(id))
	if err != nil {
		return err
	}

	return netlink.LinkDel(link)
}
