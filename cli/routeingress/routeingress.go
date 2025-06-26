package routeingress

import (
	"net"

	"github.com/vishvananda/netlink"
	"github.com/vishvananda/netlink/nl"

	"github.com/datum-cloud/galactic/util"
)

func Add(ip *net.IPNet, id int) error {
	dev := util.GenerateInterfaceNameHost(id)
	link, err := netlink.LinkByName(dev)
	if err != nil {
		return err
	}

	var flags [nl.SEG6_LOCAL_MAX]bool
	flags[nl.SEG6_LOCAL_ACTION] = true
	flags[nl.SEG6_LOCAL_VRFTABLE] = true
	encap := &netlink.SEG6LocalEncap{
		Action:   nl.SEG6_LOCAL_ACTION_END_DT46,
		Flags:    flags,
		VrfTable: id,
	}
	route := &netlink.Route{
		Dst:       ip,
		LinkIndex: link.Attrs().Index,
		Encap:     encap,
	}
	return netlink.RouteAdd(route)
}

func Delete(ip *net.IPNet, id int) error {
	dev := util.GenerateInterfaceNameHost(id)
	link, err := netlink.LinkByName(dev)
	if err != nil {
		return err
	}

	route := &netlink.Route{
		Dst:       ip,
		LinkIndex: link.Attrs().Index,
		Encap:     &netlink.SEG6LocalEncap{},
	}
	return netlink.RouteDel(route)
}
