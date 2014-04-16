package DoomAnalysis

import (
  "./db"
  "./dumpfile"
  "./file"
  "crypto/md5"
  "encoding/hex"
  "fmt"
  "io/ioutil"
  "log"
  "net/http"
)

func Start() {
  /*
  	  s := &http.Server{
  			Addr:           ":10010",
  			Handler:        http.HandlerFunc(httpHandler),
  			ReadTimeout:    5 * time.Second,
  			WriteTimeout:   5 * time.Second,
  			MaxHeaderBytes: 1 << 20,
  		}
  */
  db.InitChannelFlag()
  http.HandleFunc("/sxd", httpHandlerSxd)
  err := http.ListenAndServe(":10010", nil)
  if err != nil {
    log.Fatal("ListenAndServe :", err)
  }

  //log.Println(s.ListenAndServe())
}

func httpHandlerSxd(w http.ResponseWriter, r *http.Request) {

  r.ParseForm()
  log.Println("Form:", r.Form)

  if len(r.Form) > 0 {
    val := r.Form["par"]
    if len(val) > 0 && val[0] == "get" {
      if len(r.Form) > 1 {
        ver := r.Form["ver"][0]
        fmt.Fprintf(w, db.GetListInfoDB("sxd", ver))
      } else {
        fmt.Fprintf(w, db.VerInfoDB("sxd"))
      }
    } else if len(val) > 0 && val[0] == "post" {
      reqContent, err := ioutil.ReadAll(r.Body)
      if err != nil {
        log.Println(err)
        return
      } else if reqContent == nil || len(reqContent) < 1 {
        log.Println("empty body!")
        return
      }

      result := CheckLegal(string(reqContent))
      if result {
        go dumpfile.ProcessDumpFile("sxd", reqContent)
      } else {
        log.Println("error check md5 :", r.Form)
      }

    } else if len(val) > 0 && val[0] == "file" {
      if len(r.Form) > 2 {
        pro := r.Form["pro"][0]
        filename := r.Form["filename"][0]
        ver := r.Form["ver"][0]
        path := "./" + pro + "/dump/" + ver + "/" + filename
        fmt.Fprintf(w, string(file.ReadFile(path)))
      }
    } else if len(val) > 0 && val[0] == "recreate" {
      if len(r.Form) > 2 {
        pro := r.Form["pro"][0]
        ver := r.Form["ver"][0]
        path := "./" + pro + "/dump/" + ver + "/"

        go dumpfile.ListFileName(path, ver, pro)

      }

    }
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
