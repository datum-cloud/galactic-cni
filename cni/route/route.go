package route

import (
	"net"

	"golang.org/x/sys/unix"

	"github.com/vishvananda/netlink"

	gutil "github.com/datum-cloud/galactic/util"
)

func assembleRoute(id int, prefix, nextHop, dev string) (*netlink.Route, error) {
	_, routeDst, err := net.ParseCIDR(prefix)
	if err != nil {
		return nil, err
	}

	if nextHop != "" {
		routeGw, err := gutil.ParseIP(nextHop)
		if err != nil {
			return nil, err
		}
		return &netlink.Route{
			Dst:   routeDst,
			Gw:    routeGw,
			Table: id,
		}, nil
	}

	link, err := netlink.LinkByName(dev)
	if err != nil {
		return nil, err
	}
	return &netlink.Route{
		Dst:       routeDst,
		Table:     id,
		LinkIndex: link.Attrs().Index,
		Scope:     unix.RT_SCOPE_LINK,
	}, nil
}

func Add(id int, prefix, nextHop, dev string) error {
	route, err := assembleRoute(id, prefix, nextHop, dev)
	if err != nil {
		return err
	}
	return netlink.RouteAdd(route)
}

func Delete(id int, prefix, nextHop, dev string) error {
	route, err := assembleRoute(id, prefix, nextHop, dev)
	if err != nil {
		return err
	}
	return netlink.RouteDel(route)
}
