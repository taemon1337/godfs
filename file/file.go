package file

import (
  "path/filepath"

  "github.com/taemon1337/godfs/block"
)

type File struct {
  Name      string
  Path      string          // path does not including File:Name
  Blocks    []block.Block
}

func (f *File) GetPath() string {
  return filepath.Join(f.Path, f.Name)
}

