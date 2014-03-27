package DoomAnalysis

import (
  "./dumpfile"
  "crypto/md5"
  "encoding/hex"
  "io/ioutil"
  "log"
  "net/http"
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
  //log.Println("Header:", r.Header)
  //log.Println("Body:", string(reqContent))
  //log.Println("Form:", r.Form)

  result := CheckLegal(string(reqContent))
  if result {
    dumpfile.ProcessDumpFile(reqContent)
  } else {
    log.Println("error check md5 :", r.Form)
  }

}

// 验证MD5是否匹配
func CheckLegal(s string) bool {
  index := 0
  check_str_len := 0
  line_num := 0
  for i := 0; i < len(s); i++ {
    if s[i] == '\n' {

      line_num++

      if index == 0 {
        index = i + 1
      }

      if line_num == 5 {
        check_str_len = i + 1
        break
      }
    }
  }

  md5_str := s[4 : index-1]
  check_str := s[index:check_str_len]
  h := md5.New()
  h.Write([]byte(check_str))
  result_str := hex.EncodeToString(h.Sum(nil))

  if md5_str == result_str {
    return true
  } else {
    return false
  }
  return false
}
