package config

import (
  "flag"
)

func load() *Config {
  cfg := &Config{}

  flag.Parse()

  return cfg
}


