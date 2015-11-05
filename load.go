package main

import (
	"flag"
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"net"
	"os"

	"github.com/hashicorp/consul/api"
)

// details about a box ...
type Box struct {
	Addresses_ []string `yaml:"addresses"`
	Address_   *string  `yaml:"address"`
	Model      string   `yaml:"model"`
	Code       *string  `yaml:"code"`
}

func (b Box) Address() (*net.IP, error) {
	if b.Address_ != nil {
		a, _, err := net.ParseCIDR(*b.Address_)
		if err != nil {
			return nil, err
		}
		return &a, nil
	}
	if len(b.Addresses_) > 0 {
		a, _, err := net.ParseCIDR(b.Addresses_[0])
		if err != nil {
			return nil, err
		}
		return &a, nil
	}
	return nil, nil
}

func (b Box) Addresses() ([]net.IP, error) {
	var ip []net.IP
	if b.Address_ != nil {
		a, _, err := net.ParseCIDR(*b.Address_)
		if err != nil {
			return nil, err
		}
		ip = append(ip, a)
	}
	for _, i := range b.Addresses_ {
		a, _, err := net.ParseCIDR(i)
		if err != nil {
			return nil, err
		}
		ip = append(ip, a)
	}

	return ip, nil
}

// collection of equipment at a given place
type Place struct {
	Equipment map[string]Box `yaml:"equipment"`
	Runnet    *string        `yaml:"runnet"`
	Tag       *string        `yaml:"tag"`
}

func load(args []string) {

	f := flag.NewFlagSet("load", flag.ExitOnError)
	f.Usage = func() {
		fmt.Fprintf(os.Stderr, "\n")
		fmt.Fprintf(os.Stderr, "Load equipment yaml files into consul\n")
		fmt.Fprintf(os.Stderr, "\n")
		fmt.Fprintf(os.Stderr, "Usage:\n")
		fmt.Fprintf(os.Stderr, "\n")
		fmt.Fprintf(os.Stderr, "  %s [options] load [options]\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "\n")
		fmt.Fprintf(os.Stderr, "General Options:\n")
		fmt.Fprintf(os.Stderr, "\n")
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, "\n")
		fmt.Fprintf(os.Stderr, "Equipment Load Options:\n")
		fmt.Fprintf(os.Stderr, "\n")
		f.PrintDefaults()
		fmt.Fprintf(os.Stderr, "\n")
	}

	var config string
	f.StringVar(&config, "config", "equipment.yaml", "equipment file to load")

	var consul string
	f.StringVar(&consul, "consul", "127.0.0.1:8500", "default consul server to connect to")

	if err := f.Parse(args); err != nil {
		f.Usage()

		log.Fatalf("Invalid option(s) given")
	}

	def := api.DefaultConfig()
	if def.Address != consul {
		def.Address = consul
	}

	client, err := api.NewClient(def)
	if err != nil {
		log.Fatal(err)
	}
	catalog := client.Catalog()
	services, _, err := catalog.Services(&api.QueryOptions{})
	if err != nil {
		log.Fatal(err)
	}
	for s, m := range services {
		for _, x := range m {
			if s == "netrs" {
				log.Println(s, x)
				c, _, err := catalog.Service(s, x, &api.QueryOptions{})
				if err != nil {
					log.Fatal(err)
				}
				log.Printf("%q\n", c)
				for _, cc := range c {
					_, err := catalog.Deregister(&api.CatalogDeregistration{Node: cc.Node, Address: cc.Address, Datacenter: "avc", ServiceID: cc.ServiceID}, &api.WriteOptions{Datacenter: "avc"})
					if err != nil {
						log.Fatal(err)
					}
				}
			}

		}
	}

	c, err := ioutil.ReadFile(config)
	if err != nil {
		log.Fatal(err)
	}

	var places map[string]Place
	err = yaml.Unmarshal(c, &places)
	if err != nil {
		log.Fatal(err)
	}

	for _, place := range places {
		for b, box := range place.Equipment {
			if i, err := box.Address(); err == nil {
				if box.Code != nil && *box.Code != "" {
					var r *api.CatalogRegistration
					switch box.Model {
					case "Quanterra Q330":
						r = &api.CatalogRegistration{
							Node:       b,
							Address:    b + ".wan.geonet.org.nz",
							Datacenter: "avc",
							Service: &api.AgentService{
								Service: "qdp",
								Tags: []string{
									*box.Code,
								},
								Port:    5330,
								Address: i.String(),
							},
						}
					case "Quanterra Q330+":
						r = &api.CatalogRegistration{
							Node:       b,
							Address:    b + ".wan.geonet.org.nz",
							Datacenter: "avc",
							Service: &api.AgentService{
								Service: "qdp",
								Tags: []string{
									*box.Code,
								},
								Port:    6330,
								Address: i.String(),
							},
						}
						//log.Println("dc", "node", b, "address", i.String(), "service", "qdp", "port", "6330")
					case "Trimble NetRS":
						r = &api.CatalogRegistration{
							Node:       b,
							Address:    b + ".wan.geonet.org.nz",
							Datacenter: "avc",
							Service: &api.AgentService{
								Service: "netrs",
								Tags: []string{
									*box.Code,
								},
								Port:    80,
								Address: i.String(),
							},
						}
						// tags = box.Code, netrs, p
						//log.Println("dc", "node", b, "address", i.String(), "service", "http", "port", "80")
					case "Trimble NetR9":
						r = &api.CatalogRegistration{
							Node:       b,
							Address:    b + ".wan.geonet.org.nz",
							Datacenter: "avc",
							Service: &api.AgentService{
								Service: "netr9",
								Tags: []string{
									*box.Code,
								},
								Port:    80,
								Address: i.String(),
							},
						}
						// tags = box.Code, netr9, p
						//log.Println("dc", "node", b, "address", i.String(), "service", "http", "port", "80")
					}
					if r != nil {
						_, err := catalog.Register(r, &api.WriteOptions{Datacenter: "avc"})
						if err != nil {
							log.Fatal(err)
						}
					}
				}
			}
		}
	}

}
