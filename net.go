package main

import (
	"errors"
	"fmt"
	"net"
	"time"

	"github.com/vishvananda/netlink"
)

func waitForNetIFUp(name string, timeout time.Duration) error {
	return waitForNetIF(name, timeout, func(netif *net.Interface) bool {
		if netif != nil {
			addr, _ := netif.Addrs()
			if len(addr) > 0 {
				return true
			}
		}
		return false
	})
}
func waitForNetIFDown(name string, timeout time.Duration) error {
	return waitForNetIF(name, timeout, func(netif *net.Interface) bool {
		if netif == nil {
			return true
		}
		addrs, _ := netif.Addrs()
		if len(addrs) == 0 {
			return true
		}
		return false
	})
}

func waitForNetIF(name string, timeout time.Duration, f func(netif *net.Interface) bool) error {
	t := time.Now()
	deadline := t.Add(timeout)
	for {
		netif, _ := net.InterfaceByName(name)
		if f(netif) {
			return nil
		}

		if time.Now().After(deadline) {
			break
		}
		time.Sleep(time.Millisecond)
	}
	return errors.New("timeouted")
}

func addRouteForVPNServer(server string) error {
	defaultGw, err := getDefaultRoute()
	if err != nil {
		return fmt.Errorf("failed to get default route: %w", err)
	}
	serverIP, err := net.ResolveIPAddr("", server)
	if err != nil {
		return fmt.Errorf("failed to resolve server address: %w", err)
	}

	ok, err := hasVPNSesrverRoute(serverIP.String(), defaultGw)
	if err != nil {
		log.Fatalf("failed to check vpn route: %s", err)
	}
	if !ok {
		if err := netlink.RouteAdd(&netlink.Route{
			Dst: &net.IPNet{
				IP:   net.ParseIP(serverIP.String()),
				Mask: net.CIDRMask(32, 32),
			},
			Gw: net.ParseIP(defaultGw),
		}); err != nil {
			log.Fatalf("failed to add route: %s", err)
		}
	}
	return nil
}

func hasVPNSesrverRoute(vpnserver, gw string) (bool, error) {
	links, err := netlink.LinkList()
	if err != nil {
		return false, fmt.Errorf("failed to get link: %w", err)
	}
	for _, link := range links {
		routes, err := netlink.RouteList(link, 4)
		if err != nil {
			return false, fmt.Errorf("failed to get route list: %w", err)
		}
		for _, r := range routes {
			if r.Gw == nil {
				continue
			}
			if r.Dst == nil {
				continue
			}
			if r.Src != nil {
				continue
			}
			if r.Gw.String() == gw && r.Dst.String() == vpnserver+"/32" {
				return true, nil
			}
		}
	}
	return false, nil
}

func getDefaultRoute() (string, error) {
	link, err := netlink.LinkByName(*ifName)
	if err != nil {
		return "", fmt.Errorf("failed to get link for %s: %w", *ifName, err)
	}
	routes, err := netlink.RouteList(link, 4)
	if err != nil {
		return "", fmt.Errorf("failed to get route list: %w", err)
	}
	for _, r := range routes {
		if r.Dst == nil && r.Gw != nil {
			return r.Gw.String(), nil
		}
	}
	return "", nil
}
