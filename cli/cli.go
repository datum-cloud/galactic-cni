package cli

import (
	"log"
	"strconv"

	"github.com/spf13/cobra"

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
		Use:   "add [ip] [id]",
		Short: "Add an ingress route",
		Args:  cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			ip, err := netlink.ParseIPNet(args[0])
			if err != nil {
				log.Fatalf("Invalid ip: %v", err)
			}
			id, err := strconv.Atoi(args[1])
			if err != nil {
				log.Fatalf("Invalid id: %v", err)
			}

			if err := routeingress.Add(ip, id); err != nil {
				log.Fatalf("routeingress add failed: %v", err)
			}
		},
	}
	routeIngressDelCmd := &cobra.Command{
		Use:   "del [ip] [id]",
		Short: "Delete an ingress route",
		Args:  cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			ip, err := netlink.ParseIPNet(args[0])
			if err != nil {
				log.Fatalf("Invalid ip: %v", err)
			}
			id, err := strconv.Atoi(args[1])
			if err != nil {
				log.Fatalf("Invalid id: %v", err)
			}

			if err := routeingress.Delete(ip, id); err != nil {
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
		Use:   "add <id> <prefix> <segments> [proxy]",
		Short: "Add an egress route",
		Args:  cobra.RangeArgs(3, 4),
		Run: func(cmd *cobra.Command, args []string) {
			id, err := strconv.Atoi(args[0])
			if err != nil {
				log.Fatalf("Invalid id: %v", err)
			}
			prefix, err := netlink.ParseIPNet(args[1])
			if err != nil {
				log.Fatalf("Invalid prefix: %v", err)
			}
			segments, err := util.ParseSegments(args[2])
			if err != nil {
				log.Fatalf("Invalid segments: %v", err)
			}

			if len(args) == 4 && args[3] == "proxy" {
				if err := neighborproxy.Add(prefix, id); err != nil {
					log.Fatalf("neighborproxy add failed: %v", err)
				}
			}
			if err := routeegress.Add(id, prefix, segments); err != nil {
				log.Fatalf("routeegress add failed: %v", err)
			}
		},
	}
	routeEgressDelCmd := &cobra.Command{
		Use:   "del <id> <prefix> <segments> [proxy]",
		Short: "Delete an egress route",
		Args:  cobra.RangeArgs(3, 4),
		Run: func(cmd *cobra.Command, args []string) {
			id, err := strconv.Atoi(args[0])
			if err != nil {
				log.Fatalf("Invalid id: %v", err)
			}
			prefix, err := netlink.ParseIPNet(args[1])
			if err != nil {
				log.Fatalf("Invalid prefix: %v", err)
			}
			segments, err := util.ParseSegments(args[2])
			if err != nil {
				log.Fatalf("Invalid segments: %v", err)
			}

			if len(args) == 4 && args[3] == "proxy" {
				if err := neighborproxy.Delete(prefix, id); err != nil {
					log.Fatalf("neighborproxy delete failed: %v", err)
				}
			}
			if err := routeegress.Delete(id, prefix, segments); err != nil {
				log.Fatalf("routeegress delete failed: %v", err)
			}
		},
	}
	routeEgressCmd.AddCommand(routeEgressAddCmd, routeEgressDelCmd)
	rootCmd.AddCommand(routeEgressCmd)

	return rootCmd
}
