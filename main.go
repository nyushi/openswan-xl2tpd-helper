package main

import (
	"encoding/json"
	"flag"
	"io/ioutil"
	"net"
	"strings"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
)

var ifName = flag.String("interface", "", "network interface name to use vpn connection")
var confPath = flag.String("conf", "", "path to config json")
var debug = flag.Bool("debug", false, "debug")

var log = logrus.New()

type Config struct {
	Server string
	PSK    string
	User   string
	Pass   string
}

func serviceOperationParallel(svcs []string, op string) {
	wg := sync.WaitGroup{}
	for _, svc := range svcs {
		wg.Add(1)
		go func(svc string) {
			log.Debugf("%s service: %s", op, svc)
			var err error
			switch op {
			case "start":
				err = startSystemdService(svc)
			case "stop":
				err = stopSystemdService(svc)
			}
			if err != nil {
				log.Warnf("error at %s unit %s: %s", op, svc, err)
			}
			log.Infof("%s %s\n", svc, op)
			wg.Done()
		}(svc)
	}
	wg.Wait()
}

func main() {
	flag.Parse()
	if *debug {
		log.SetLevel(logrus.DebugLevel)
	}
	b, err := ioutil.ReadFile(*confPath)
	if err != nil {
		log.Fatalf("failed to read %s", confPath)
	}
	conf := Config{}
	if err := json.Unmarshal(b, &conf); err != nil {
		log.Fatalf("failed to parse %s: %s", confPath, err)
	}
	switch flag.Args()[0] {
	case "stop":
		commands := []string{"ip", "route", "del", "default", "dev", "ppp0"}
		if err := execCommand(commands...); err != nil {
			log.Warnf("failed to delete default gw to ppp0: %s", err)
		}
		if err := l2tpstop(); err != nil {
			log.Warnf("failed to disconnect l2tp: %s", err)
		}

		commands = []string{"ipsec", "auto", "--down", "L2TP-PSK"}
		if err := execCommand(commands...); err != nil {
			log.Warnf("failed to disconnect ipsec: %s", err)
		}

		serviceOperationParallel([]string{"xl2tpd.service", "openswan.service"}, "stop")

		log.Debugf("wait for ppp0 down")
		if err := waitForNetIFDown("ppp0", 10*time.Second); err != nil {
			log.Warnf("ppp0 stil exists: %s", err)
		}

	case "start":
		log.Debug("add route for vpn server")
		if err := addRouteForVPNServer(conf.Server); err != nil {
			log.Fatalf("failed to add vpn route: %w", err)
		}
		log.Info("Added route for vpn server")

		serviceOperationParallel([]string{"xl2tpd.service", "openswan.service"}, "stop")
		makeConfig(conf)
		serviceOperationParallel([]string{"xl2tpd.service", "openswan.service"}, "start")

		for {
			err := ipsecCommand("auto", "--up", "L2TP-PSK")
			if err != nil {
				if ipsecErr, ok := err.(*ipsecTemporaryError); ok {
					if ipsecErr.Temporary() {
						continue
					}
				}
				log.Fatalf("failed to execute ipsec start")
			}
			break
		}
		log.Info("ipsec started")

		log.Debug("start l2tp")
		if err := l2tpstart(); err != nil {
			log.Fatalf("error at xl2tp start: %s", err)
		}
		log.Info("l2tp connection started")

		log.Debug("wait for ppp0 up")
		if err := waitForNetIFUp("ppp0", 10*time.Second); err != nil {
			log.Fatalf("ppp0 not found: %s", err)
		}
		log.Info("ppp0 up")

		commands := []string{"ip", "route", "add", "default", "dev", "ppp0"}
		if err := execCommand(commands...); err != nil {
			log.Fatalf("failed to add default gw to ppp0")
		}
	}
}

func makeConfig(c Config) {
	i, err := net.InterfaceByName(*ifName)
	if err != nil {
		log.Fatalf("failed to get interface: %s", err)
	}
	addrs, err := i.Addrs()
	if err != nil {
		log.Fatalf("failed to get addresses: %s", err)
	}
	localAddr := strings.Split(addrs[0].String(), "/")[0]
	serverIP, err := net.ResolveIPAddr("", c.Server)
	if err != nil {
		log.Fatalf("failed to resolve server address: %s", err)
	}
	serverAddr := serverIP.String()

	if err := ipsecSecretFile(c.Server, c.PSK); err != nil {
		log.Fatalf("error at PSK file creation: %s", err)
	}
	log.Debug("psk file created")
	if err := ipsecConfFile(localAddr, serverAddr); err != nil {
		log.Fatalf("error at ipsec file creation: %s", err)
	}
	log.Debug("ipsec config created")
	if err := xl2tpdConfFile(serverAddr); err != nil {
		log.Fatalf("error at xl2tpd file creation: %s", err)
	}
	log.Debug("xl2tpd config created")
	if err := optionsL2TPDClientFile(c.User, c.Pass); err != nil {
		log.Fatalf("error at xl2tpd option file creation: %s", err)
	}
	log.Debug("xl2tpd ppp config created")
}
