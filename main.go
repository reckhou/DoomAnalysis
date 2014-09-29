package main

import (
  "bitbucket.org/kardianos/osext"
  "bitbucket.org/reckhou/DoomAnalysis/src"
  "os"
)

func main() {
  path, _ := osext.ExecutableFolder()
  os.Chdir(path)
  DoomAnalysis.Start()
}
