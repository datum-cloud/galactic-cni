package vrf

import (
	"fmt"
	"math"
	"slices"

	"github.com/datum-cloud/galactic/cni/sysctl"
	"github.com/datum-cloud/galactic/util"
	"github.com/vishvananda/netlink"
)

const MinVRFId = uint32(1)
const MaxVRFId = uint32(math.MaxUint32 - 1)

func Add(id int) error {
	name := util.GenerateInterfaceNameVRF(id)

	vrfId, err := FindNextAvailableVRFId()
	if err != nil {
		return err
	}

	vrf := &netlink.Vrf{
		LinkAttrs: netlink.LinkAttrs{
			Name: name,
		},
		Table: uint32(vrfId),
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

func FindNextAvailableVRFId() (uint32, error) {
	links, err := netlink.LinkList()
	if err != nil {
		return 0, err
	}

	used := make([]uint32, 0, len(links))
	for _, link := range links {
		if vrf, ok := link.(*netlink.Vrf); ok {
			used = append(used, vrf.Table)
		}
	}

	for vrfId := MinVRFId; vrfId <= MaxVRFId; vrfId++ {
		if !slices.Contains(used, vrfId) {
			return vrfId, nil
		}
	}

	return 0, fmt.Errorf("could not find any available VRF id")
}
