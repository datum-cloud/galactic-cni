package routeegress

import (
	"net"

	"github.com/vishvananda/netlink"
	"github.com/vishvananda/netlink/nl"
)

const LoopbackDevice = "lo-srv6"

func Add(id int, prefix *net.IPNet, segments []net.IP) error {
	link, err := netlink.LinkByName(LoopbackDevice)
	if err != nil {
		return err
	}

	encap := &netlink.SEG6Encap{
		Mode:     nl.SEG6_IPTUN_MODE_ENCAP,
		Segments: segments,
	}
	route := &netlink.Route{
		Dst:       prefix,
		Table:     id,
		LinkIndex: link.Attrs().Index,
		Encap:     encap,
	}
	return netlink.RouteReplace(route)
}

func Delete(id int, prefix *net.IPNet, segments []net.IP) error {
	link, err := netlink.LinkByName(LoopbackDevice)
	if err != nil {
		return err
	}

	route := &netlink.Route{
		Dst:       prefix,
		Table:     id,
		LinkIndex: link.Attrs().Index,
	}
	return netlink.RouteDel(route)
}
