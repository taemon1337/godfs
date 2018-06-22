package cluster

import (
  "log"
  "time"
  "strconv"

  "github.com/hashicorp/serf/client"
  "github.com/hashicorp/mapstructure"
)

func handler(name string, c chan map[string]interface{}) {
  for {
    data := <-c
    log.Printf("%s: %v\n", name, data)
  }
}

type Node struct {
  Name        string
  Client      *client.RPCClient
}

func (n *Node) Handle(c chan map[string]interface{}) {
  log.Printf("Node Handling")
  ch := make(chan map[string]interface{})
  n.Client.Stream("user", ch)
  go handler(n.Name, ch)

  i := 0
  for {
    err := n.Client.UserEvent("test", []byte(strconv.Itoa(i)), false)
    if err != nil {
      panic(err)
    }

    time.Sleep(5 * time.Second)
    i++
  }
}

func NewNode(name string, opts map[string]interface{}) *Node {
  cfg := client.Config{
    Addr: "",
    AuthKey: "",
    Timeout: 5 * time.Second,
  }

  err := mapstructure.Decode(opts, &cfg)
  if err != nil {
    log.Fatal("Cannot parse Etcd ClientV3 Config - %v: %v", opts, err)
  }

  client, err := client.ClientFromConfig(&cfg)
  if err != nil {
    log.Fatal("Cannot connect node to Cluster: %v", err)
  }

  return &Node{
    Name: name,
    Client: client,
  }
}
