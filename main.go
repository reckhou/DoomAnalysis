package main

import (
  "bitbucket.org/reckhou/DoomAnalysis/src"
  "os"
)

func main() {
  os.Chdir("/data/pressureserver/")
  DoomAnalysis.Start()
}
