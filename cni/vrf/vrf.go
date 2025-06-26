package vrf

import (
	"fmt"
	"github.com/datum-cloud/galactic/cni/sysctl"
	"github.com/vishvananda/netlink"
)

const VrfNameTemplate = "galactic%d-vrf"

func Add(id int) error {
	name := fmt.Sprintf(VrfNameTemplate, id)

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
	link, err := netlink.LinkByName(fmt.Sprintf(VrfNameTemplate, id))
	if err != nil {
		return err
	}

	return netlink.LinkDel(link)
}
