package main

import (
  "log"
  "flag"
  "os"
  "time"

  "crypto/rand"

//  "github.com/taemon1337/godfs/file"
  "github.com/taemon1337/godfs/kvstore"
  "github.com/taemon1337/godfs/cluster"
//  "github.com/taemon1337/godfs/policy"
)

func getNodeName() string {
  hostname, err := os.Hostname()
  if err != nil {
    c := 10
    b := make([]byte, c)
    _, err := rand.Read(b)
    if err != nil {
      log.Fatal("Could not generate random node name: %v", err)
    }

    return string(b)
  }

  return hostname
}

func main() {
  var name, authkey string
  kvopts := make(map[string]interface{})
  nodeopts := make(map[string]interface{})
//  policyopts := make(map[string]interface{})

  flag.StringVar(&name, "name", getNodeName(), "The name of this node")
  flag.StringVar(&authkey, "authkey", "", "The cluster node authentication key")

  log.Printf("Starting godfs...")

  kvopts["Endpoints"] = []string{"localhost:2379", "localhost:22379", "localhost:32379"}
  kvopts["DialTimeout"] = 5 * time.Second
  nodeopts["Addr"] = "localhost:7373"
  nodeopts["AuthKey"] = authkey
//  policyopts["Verbose"] = true

  kvclient := kvstore.NewKVClient(kvopts)
  nodeclient := cluster.NewNode(name, nodeopts)
//  policyclient := policy.NewPolicyEngineClient(policyopts)

  log.Printf("KV CLIENT: %v", kvclient)
  log.Printf("CLUSTER NODE: %v", nodeclient)
//  log.Printf("POLICY CLIENT: %v", policyclient)
}
