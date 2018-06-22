package kvstore

import (
  "github.com/coreos/etcd/clientv3"
)

type Config struct {
  clientv3.Config
}
