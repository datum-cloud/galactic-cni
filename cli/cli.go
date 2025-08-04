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
		Use:   "add <ip>",
		Short: "Add an ingress route",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			ip, err := netlink.ParseIPNet(args[0])
			if err != nil {
				log.Fatalf("Invalid ip: %v", err)
			}
			if !IsHost(ip) {
				log.Fatalf("ip is not a host route")
			}
			vpc, vpcAttachment, err := util.ExtractVPCFromSRv6Endpoint(ip.IP)
			if err != nil {
				log.Fatalf("could not extract SRv6 endpoint: %v", err)
			}
			vpc, err = ToBase62(vpc)
			if err != nil {
				log.Fatalf("Invalid vpc: %v", err)
			}
			vpcAttachment, err = ToBase62(vpcAttachment)
			if err != nil {
				log.Fatalf("Invalid vpcattachment: %v", err)
			}

			if err := routeingress.Add(ip, vpc, vpcAttachment); err != nil {
				log.Fatalf("routeingress add failed: %v", err)
			}
		},
	}
	routeIngressDelCmd := &cobra.Command{
		Use:   "del <ip>",
		Short: "Delete an ingress route",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			ip, err := netlink.ParseIPNet(args[0])
			if err != nil {
				log.Fatalf("Invalid ip: %v", err)
			}
			if !IsHost(ip) {
				log.Fatalf("ip is not a host route")
			}
			vpc, vpcAttachment, err := util.ExtractVPCFromSRv6Endpoint(ip.IP)
			if err != nil {
				log.Fatalf("could not extract SRv6 endpoint: %v", err)
			}
			vpc, err = ToBase62(vpc)
			if err != nil {
				log.Fatalf("Invalid vpc: %v", err)
			}
			vpcAttachment, err = ToBase62(vpcAttachment)
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
		Use:   "add <prefix> <src> <segments>",
		Short: "Add an egress route",
		Args:  cobra.ExactArgs(3),
		Run: func(cmd *cobra.Command, args []string) {
			prefix, err := netlink.ParseIPNet(args[0])
			if err != nil {
				log.Fatalf("Invalid prefix: %v", err)
			}
			src, err := netlink.ParseIPNet(args[1])
			if err != nil {
				log.Fatalf("Invalid src: %v", err)
			}
			if !IsHost(src) {
				log.Fatalf("src is not a host route")
			}
			segments, err := util.ParseSegments(args[2])
			if err != nil {
				log.Fatalf("Invalid segments: %v", err)
			}

			vpc, vpcAttachment, err := util.ExtractVPCFromSRv6Endpoint(src.IP)
			if err != nil {
				log.Fatalf("could not extract SRv6 endpoint: %v", err)
			}
			vpc, err = ToBase62(vpc)
			if err != nil {
				log.Fatalf("Invalid vpc: %v", err)
			}
			vpcAttachment, err = ToBase62(vpcAttachment)
			if err != nil {
				log.Fatalf("Invalid vpcattachment: %v", err)
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
		Use:   "del <prefix> <src> <segments>",
		Short: "Delete an egress route",
		Args:  cobra.ExactArgs(3),
		Run: func(cmd *cobra.Command, args []string) {
			prefix, err := netlink.ParseIPNet(args[0])
			if err != nil {
				log.Fatalf("Invalid prefix: %v", err)
			}
			src, err := netlink.ParseIPNet(args[1])
			if err != nil {
				log.Fatalf("Invalid src: %v", err)
			}
			if !IsHost(src) {
				log.Fatalf("src is not a host route")
			}
			segments, err := util.ParseSegments(args[2])
			if err != nil {
				log.Fatalf("Invalid segments: %v", err)
			}

			vpc, vpcAttachment, err := util.ExtractVPCFromSRv6Endpoint(src.IP)
			if err != nil {
				log.Fatalf("could not extract SRv6 endpoint: %v", err)
			}
			vpc, err = ToBase62(vpc)
			if err != nil {
				log.Fatalf("Invalid vpc: %v", err)
			}
			vpcAttachment, err = ToBase62(vpcAttachment)
			if err != nil {
				log.Fatalf("Invalid vpcattachment: %v", err)
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
