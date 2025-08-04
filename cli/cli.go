package cli

import (
	"log"
	"net"
	"strings"

	"github.com/spf13/cobra"

	"github.com/kenshaw/baseconv"
	"github.com/vishvananda/netlink"

	"github.com/datum-cloud/galactic/cli/neighborproxy"
	"github.com/datum-cloud/galactic/cli/routeegress"
	"github.com/datum-cloud/galactic/cli/routeingress"
	"github.com/datum-cloud/galactic/util"
)

func NewCommand() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "galactic-cli",
		Short: "Galactic CLI tool",
	}

	routeIngressCmd := &cobra.Command{
		Use:   "route-ingress",
		Short: "Manage ingress routes",
	}
	routeIngressAddCmd := &cobra.Command{
		Use:   "add <ip> <vpc> <vpcattachment>",
		Short: "Add an ingress route",
		Args:  cobra.ExactArgs(3),
		Run: func(cmd *cobra.Command, args []string) {
			ip, err := netlink.ParseIPNet(args[0])
			if err != nil {
				log.Fatalf("Invalid ip: %v", err)
			}
			if !IsHost(ip) {
				log.Fatalf("ip is not a host route")
			}
			vpc, err := ToBase62(args[1])
			if err != nil {
				log.Fatalf("Invalid vpc: %v", err)
			}
			vpcAttachment, err := ToBase62(args[2])
			if err != nil {
				log.Fatalf("Invalid vpcattachment: %v", err)
			}

			if err := routeingress.Add(ip, vpc, vpcAttachment); err != nil {
				log.Fatalf("routeingress add failed: %v", err)
			}
		},
	}
	routeIngressDelCmd := &cobra.Command{
		Use:   "del <ip> <vpc> <vpcattachment>",
		Short: "Delete an ingress route",
		Args:  cobra.ExactArgs(3),
		Run: func(cmd *cobra.Command, args []string) {
			ip, err := netlink.ParseIPNet(args[0])
			if err != nil {
				log.Fatalf("Invalid ip: %v", err)
			}
			if !IsHost(ip) {
				log.Fatalf("ip is not a host route")
			}
			vpc, err := ToBase62(args[1])
			if err != nil {
				log.Fatalf("Invalid vpc: %v", err)
			}
			vpcAttachment, err := ToBase62(args[2])
			if err != nil {
				log.Fatalf("Invalid vpcattachment: %v", err)
			}

			if err := routeingress.Delete(ip, vpc, vpcAttachment); err != nil {
				log.Fatalf("routeingress delete failed: %v", err)
			}
		},
	}
	routeIngressCmd.AddCommand(routeIngressAddCmd, routeIngressDelCmd)
	rootCmd.AddCommand(routeIngressCmd)

	routeEgressCmd := &cobra.Command{
		Use:   "route-egress",
		Short: "Manage egress routes",
	}
	routeEgressAddCmd := &cobra.Command{
		Use:   "add <vpc> <vpcattachment> <prefix> <segments>",
		Short: "Add an egress route",
		Args:  cobra.ExactArgs(4),
		Run: func(cmd *cobra.Command, args []string) {
			vpc, err := ToBase62(args[0])
			if err != nil {
				log.Fatalf("Invalid vpc: %v", err)
			}
			vpcAttachment, err := ToBase62(args[1])
			if err != nil {
				log.Fatalf("Invalid vpcattachment: %v", err)
			}
			prefix, err := netlink.ParseIPNet(args[2])
			if err != nil {
				log.Fatalf("Invalid prefix: %v", err)
			}
			segments, err := util.ParseSegments(args[3])
			if err != nil {
				log.Fatalf("Invalid segments: %v", err)
			}

			if IsHost(prefix) {
				if err := neighborproxy.Add(prefix, vpc, vpcAttachment); err != nil {
					log.Fatalf("neighborproxy add failed: %v", err)
				}
			}
			if err := routeegress.Add(vpc, vpcAttachment, prefix, segments); err != nil {
				log.Fatalf("routeegress add failed: %v", err)
			}
		},
	}
	routeEgressDelCmd := &cobra.Command{
		Use:   "del <vpc> <vpcattachment> <prefix> <segments>",
		Short: "Delete an egress route",
		Args:  cobra.ExactArgs(4),
		Run: func(cmd *cobra.Command, args []string) {
			vpc, err := ToBase62(args[0])
			if err != nil {
				log.Fatalf("Invalid vpc: %v", err)
			}
			vpcAttachment, err := ToBase62(args[1])
			if err != nil {
				log.Fatalf("Invalid vpcattachment: %v", err)
			}
			prefix, err := netlink.ParseIPNet(args[2])
			if err != nil {
				log.Fatalf("Invalid prefix: %v", err)
			}
			segments, err := util.ParseSegments(args[3])
			if err != nil {
				log.Fatalf("Invalid segments: %v", err)
			}

			if IsHost(prefix) {
				if err := neighborproxy.Delete(prefix, vpc, vpcAttachment); err != nil {
					log.Fatalf("neighborproxy delete failed: %v", err)
				}
			}
			if err := routeegress.Delete(vpc, vpcAttachment, prefix, segments); err != nil {
				log.Fatalf("routeegress delete failed: %v", err)
			}
		},
	}
	routeEgressCmd.AddCommand(routeEgressAddCmd, routeEgressDelCmd)
	rootCmd.AddCommand(routeEgressCmd)

	return rootCmd
}

func ToBase62(value string) (string, error) {
	return baseconv.Convert(strings.ToLower(value), baseconv.DigitsHex, baseconv.Digits62)
}

func IsHost(ipNet *net.IPNet) bool {
	ones, bits := ipNet.Mask.Size()
	// host if mask is full length: /32 for IPv4, /128 for IPv6
	return ones == bits
}
