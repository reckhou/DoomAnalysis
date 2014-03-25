package DoomAnalysis

import (
  "bitbucket.org/reckhou/DoomAnalysis/src/file"
  "io/ioutil"
  "log"
  "net/http"
  "os"
  "time"
)

func Start() {
  s := &http.Server{
    Addr:           ":10010",
    Handler:        http.HandlerFunc(httpHandler),
    ReadTimeout:    5 * time.Second,
    WriteTimeout:   5 * time.Second,
    MaxHeaderBytes: 1 << 20,
  }

  log.Println(s.ListenAndServe())
}

func httpHandler(w http.ResponseWriter, r *http.Request) {
  reqContent, err := ioutil.ReadAll(r.Body)
  if err != nil {
    log.Println(err)
    return
  } else if reqContent == nil || len(reqContent) < 1 {
    log.Println("empty body!")
    return
  }
  r.ParseForm()
  log.Println("Header:", r.Header)
  log.Println("Body:", string(reqContent))
  log.Println("Form:", r.Form)

  fileName := "crash_" + time.Now().Format(time.RFC3339) + ".txt"
  file.WriteFile("./"+fileName, reqContent, os.O_TRUNC)
}
