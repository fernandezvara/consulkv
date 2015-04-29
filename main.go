package main

import (
	"fmt"
	"os"

	"github.com/codegangsta/cli"
	"github.com/hashicorp/consul/api"
)

func main() {
	os.Exit(realMain())
}

func realMain() int {

	app := cli.NewApp()
	app.Name = "consulkv"
	app.Author = "Antonio Fdez"
	app.Email = "antoniofernandezvara@gmail.com"
	app.Version = "0.2"
	app.Usage = "Helper utility to set and get to/from consul"

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "server, s",
			Value:  "localhost:8500",
			Usage:  "Server/Cluster to connect to. Defaults to `localhost:8500`",
			EnvVar: "CONF_CONSUL",
		},
	}

	app.Commands = []cli.Command{
		{
			Name:   "get",
			Usage:  "Get a key",
			Action: get,
		},
		{
			Name:   "set",
			Usage:  "Sets the desired value on the key",
			Action: set,
		},
	}

	app.Run(os.Args)
	return 0
}

func get(c *cli.Context) {
	_, kv := connectToConsul(c.GlobalString("server"))
	pair, _, err := kv.Get(c.Args().First(), nil)
	if err != nil {
		fmt.Println(err)
	}
	if pair != nil {
		fmt.Printf("%s", string(pair.Value))
	}
}

func set(c *cli.Context) {
	_, kv := connectToConsul(c.GlobalString("server"))
	key := c.Args().Get(0)
	value := c.Args().Get(1)
	if key == "" || value == "" {
		fmt.Println("Missing key or value")
	} else {
		p := &api.KVPair{Key: key, Value: []byte(value)}
		_, err := kv.Put(p, nil)
		if err != nil {
			fmt.Println(err)
		}
	}
}

func connectToConsul(serverAddress string) (client *api.Client, kv *api.KV) {
	var (
		err    error
		config *api.Config
	)
	config = &api.Config{
		Address: serverAddress,
	}

	client, err = api.NewClient(config)
	if err != nil {
		fmt.Println(err)
		os.Exit(2)
	}
	kv = client.KV()
	return
}
