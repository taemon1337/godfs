package main

import (
  "log"
  "time"

  "github.com/coreos/etcd/clientv3"
  "github.com/taemon1337/godfs/block"
)

func main() {
  log.Printf("Starting godfs...")

  cli, err := clientv3.New(clientv3.Config{
    Endpoints:   []string{"localhost:2379", "localhost:22379", "localhost:32379"},
    DialTimeout: 5 * time.Second,
  })

  if err != nil {
    log.Printf("ERROR: %v", err)
  }

  defer cli.Close()

  log.Println(block.NewBlock("test block"))
}
