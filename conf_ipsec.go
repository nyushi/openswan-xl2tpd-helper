package main

import (
	"fmt"
	"os"
	"text/template"
)

func ipsecConfFile(local, remote string) error {
	path := "/etc/ipsec.conf"
	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("failed to create %s: %w", path, err)
	}
	t := template.Must(template.New("ipsec_conf").Parse(ipsecConfTemplate))
	if err := t.Execute(f, struct {
		Local  string
		Remote string
	}{
		Local:  local,
		Remote: remote,
	}); err != nil {
		return fmt.Errorf("failed to fill template %s: %w", path, err)
	}
	return nil
}

var ipsecConfTemplate = `
# /etc/ipsec.conf - Openswan IPsec configuration file

# This file:  /usr/share/doc/openswan/ipsec.conf-sample
#
# Manual:     ipsec.conf.5


version	2.0	# conforms to second version of ipsec.conf specification

# basic configuration
config setup
	# Do not set debug options to debug configuration issues!
	# plutodebug / klipsdebug = "all", "none" or a combination from below:
	# "raw crypt parsing emitting control klips pfkey natt x509 dpd private"
	# eg:
	# plutodebug="control parsing"
	# Again: only enable plutodebug or klipsdebug when asked by a developer
	#
	# enable to get logs per-peer
	# plutoopts="--perpeerlog"
	#
	# Enable core dumps (might require system changes, like ulimit -C)
	# This is required for abrtd to work properly
	# Note: incorrect SElinux policies might prevent pluto writing the core
	dumpdir=/var/run/pluto/
	#
	# NAT-TRAVERSAL support, see README.NAT-Traversal
	nat_traversal=yes
	# exclude networks used on server side by adding %v4:!a.b.c.0/24
	# It seems that T-Mobile in the US and Rogers/Fido in Canada are
	# using 25/8 as "private" address space on their 3G network.
	# This range has not been announced via BGP (at least upto 2010-12-21)
	virtual_private=%v4:10.0.0.0/8,%v4:192.168.0.0/16,%v4:172.16.0.0/12,%v4:25.0.0.0/8,%v6:fd00::/8,%v6:fe80::/10
	# OE is now off by default. Uncomment and change to on, to enable.
	oe=off
	# which IPsec stack to use. auto will try netkey, then klips then mast
	#protostack=auto
	protostack=netkey
	# Use this to log to a file, or disable logging on embedded systems (like openwrt)
	#plutostderrlog=/dev/null
        plutoopts="--interface=ens33" # added


# Add connections here

# sample VPN connection
# for more examples, see /etc/ipsec.d/examples/
#conn sample
#		# Left security gateway, subnet behind it, nexthop toward right.
#		left=10.0.0.1
#		leftsubnet=172.16.0.0/24
#		leftnexthop=10.22.33.44
#		# Right security gateway, subnet behind it, nexthop toward left.
#		right=10.12.12.1
#		rightsubnet=192.168.0.0/24
#		rightnexthop=10.101.102.103
#		# To authorize this connection, but not actually start it, 
#		# at startup, uncomment this.
#		#auto=add


# added
conn L2TP-PSK
     authby=secret
     pfs=no
     auto=add
     keyingtries=3
     dpddelay=30
     dpdtimeout=120
     dpdaction=clear
     rekey=yes
     ikelifetime=8h
     keylife=1h
     type=transport
# Replace %any below with your local IP address (private, behind NAT IP is okay as well)
     left={{ .Local }}
     leftprotoport=17/1701
# Replace IP address with your VPN server's IP
     right={{ .Remote}}
     rightprotoport=17/1701
`
