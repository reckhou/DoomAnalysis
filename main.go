package main

import (
  "bitbucket.org/kardianos/osext"
  "bitbucket.org/reckhou/DoomAnalysis/src"
  "github.com/reckhou/goCfgMgr"
  "log"
  "net/http"
  "os"
)

import _ "net/http/pprof"

func main() {
  path, _ := osext.ExecutableFolder()
  os.Chdir(path)

  go func() {
    log.Println(http.ListenAndServe(goCfgMgr.Get("basic", "Host").(string)+":10012", nil))
  }()
  DoomAnalysis.Start()
}
