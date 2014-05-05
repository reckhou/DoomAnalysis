package debug

import (
  "log"
  "net/http"
  "runtime"
  "time"
)

const (
  Enabled = true
)

func HTTPRequest(r *http.Request) {
  if r == nil || !Enabled {
    return
  }
  log.Println("Request:")
  log.Println("Method:", r.Method)
  log.Println("Path:", r.URL.Path)
  log.Println("RawQuery:", r.URL.RawQuery)
  log.Println("Protocol:", r.Proto)
  log.Println("Host:", r.Host)
  log.Println("RemoteAddr:", r.RemoteAddr)
  log.Println("RequestURI:", r.RequestURI)
  log.Println("Header:\n", r.Header)
  log.Println("Body:\n", r.Body)
}

func MemStats() {
  var memStat runtime.MemStats
  runtime.ReadMemStats(&memStat)
  if !Enabled {
    return
  }
  log.Println("MemStats:")
  log.Println("Alloc", memStat.Alloc)
  log.Println("TotalAlloc", memStat.TotalAlloc)
  log.Println("Sys", memStat.Sys)
  log.Println("HeapAlloc", memStat.HeapAlloc)
  log.Println("HeapSys", memStat.HeapSys)
  log.Println("HeapReleased", memStat.HeapReleased)
  log.Println("TotalAlloc", memStat.TotalAlloc)
  log.Println("NumGC", memStat.NumGC)
}

func CheckMemStats() {
  for Enabled {
    MemStats()
    time.Sleep(1000 * 60 * 30 * time.Millisecond)
  }
}
