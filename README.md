# consulkv

consulkv is a simple getter and setter for consul.

The goal of this helper command is to get or set information in a simple way from scripts or fast lookups.


## Usage

```bash

consulkv <options> get _key_              # get a key by name/path

consulkv <options> set _key_ _value_      # sets the value on the desired key

Options:

    -debug, -d    Show errors to console if any.
    -server, -s   Consul server and port to connect to. Defaults to 'localhost:8500'.

                  This server can also be configured using the Environment Variable 'CONSUL_SERVER'

```

##Examples:

Getting a key

```bash

CONSUL_SERVER=1.2.3.4:8500

$MASTER_IP=$(consulkv get "/my/master/ip")

```

Setting the host IP from inside a docker container

```bash

HOST_IP=$(/sbin/ip route|awk '/default/ { print $3}')

consulkv set "/service/test/ip" $HOST_IP


```
