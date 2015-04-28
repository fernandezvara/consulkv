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
	app.Name = "consul-cli"
	app.Usage = "make an explosive entrance"
	app.Action = func(c *cli.Context) {
		println("Please tell me what to do.")
	}

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "server, s",
			Value:  "localhost:8500",
			Usage:  "server/cluster to connect to. defaults to `localhost:8500`",
			EnvVar: "CONSUL_SERVER",
		},
		cli.BoolFlag{
			Name:  "debug, d",
			Usage: "shows error messages",
		},
	}

	app.Commands = []cli.Command{
		{
			Name:  "get",
			Usage: "get a key",
			Action: func(c *cli.Context) {
				_, kv := connectToConsul(c.String("server"))
				pair, _, err := kv.Get(c.Args().First(), nil)
				if err != nil {
					fmt.Println(err)
				}
				fmt.Printf("%s", string(pair.Value))
			},
		},
		{
			Name:  "set",
			Usage: "sets the desired value on the key",
			Action: func(c *cli.Context) {
				_, kv := connectToConsul(c.String("server"))
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
			},
		},
	}

	app.Run(os.Args)
	return 0
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
