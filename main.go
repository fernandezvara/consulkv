package main

import (
	"fmt"
	"os"

	"github.com/alyu/encrypt"
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
		cli.StringFlag{
			Name:   "crypt, c",
			Value:  "",
			Usage:  "Encryption/Decryption string",
			EnvVar: "CONF_CRYPT",
		},
	}

	app.Commands = []cli.Command{
		{
			Name:   "get",
			Usage:  "Get a key",
			Action: getKey,
		},
		{
			Name:   "set",
			Usage:  "Sets the desired value on the key",
			Action: setKey,
		},
		{
			Name:   "delete",
			Usage:  "Deletes the desired key",
			Action: deleteKey,
		},
	}

	app.Run(os.Args)
	return 0
}

func getKey(c *cli.Context) {
	cryptKey := c.GlobalString("crypt")
	_, kv := connectToConsul(c.GlobalString("server"))
	pair, _, err := kv.Get(c.Args().First(), nil)
	if err != nil {
		fmt.Println(err)
	}
	if pair != nil {
		fmt.Printf("%s", getDecrypted(pair.Value, cryptKey))
	}
}

func setKey(c *cli.Context) {
	_, kv := connectToConsul(c.GlobalString("server"))
	key := c.Args().Get(0)
	value := c.Args().Get(1)
	if key == "" || value == "" {
		fmt.Println("Missing key or value")
	} else {
		p := kvPair(key, value, c.GlobalString("crypt"))
		_, err := kv.Put(p, nil)
		if err != nil {
			fmt.Println(err)
		}
	}
}

func deleteKey(c *cli.Context) {
	_, kv := connectToConsul(c.GlobalString("server"))
	key := c.Args().Get(0)
	if key == "" {
		fmt.Println("Missing key")
	} else {
		_, err := kv.Delete(key, nil)
		if err != nil {
			fmt.Println(err)
		}
	}
}

func getDecrypted(value []byte, cryptKey string) string {
	if cryptKey != "" {
		return encrypt.Decrypt([]byte(cryptKey), value)
	}
	return string(value)
}

func kvPair(key, value, cryptKey string) *api.KVPair {
	if cryptKey != "" {
		return &api.KVPair{Key: key, Value: encrypt.Encrypt([]byte(cryptKey), []byte(value))}
	}
	return &api.KVPair{Key: key, Value: []byte(value)}
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
