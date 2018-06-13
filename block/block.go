package block

import (
)

type Block struct {
  Id string
}

func (b *Block) ToString() string {
  return "BLOCK ID: " + b.Id
}

func NewBlock(id string) *Block {
  return &Block{
    Id: id,
  }
}
