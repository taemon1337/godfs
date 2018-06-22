package kvstore

import (
  "log"
  "time"

  "github.com/coreos/etcd/clientv3"
  "github.com/hashicorp/mapstructure"
)

type KVClient struct {
  Client    *clientv3.Client
}

func NewKVClient(opts map[string]interface{}) *KVClient {
  cfg := clientv3.Config{
    DialTimeout: 5 * time.Second,
  }

  err := mapstructure.Decode(opts, &cfg)
  if err != nil {
    log.Fatal("Cannot parse Etcd ClientV3 Config - %v: %v", opts, err)
  }

  client, err := clientv3.New(cfg)
  if err != nil {
    log.Fatal("Cannot create new etcd client: %v", err)
  }

  defer client.Close()

  return &KVClient{
    Client: client,
  }
}
