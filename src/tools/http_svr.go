package main

import (
  "fmt"
  "log"
  "net/http"
  "strings"
)

func parseURL(url *url.URL) (params map[string]string) {

  params = make(map[string]string)

  paramStr := url.RawQuery
  paramSlice := strings.Split(paramStr, "&")

  for _, v := range paramSlice {
    index := strings.LastIndex(v, "=")
    var key, val string
    if index >= 0 && index+1 < len(v) {
      key = v[:index]
      val = v[index+1:]
    } else {
      key = v
      val = ""
    }

    params[key] = val
  }

  return
}

func sayhelloName(w http.ResponseWriter, r *http.Request) {

  //r.ParseForm()
  log.Println("Form:", r.Form)
  log.Println("Url", r.URL)
  log.Println("RequestURL", r.RequestURI)

  params := parseURL(r.URL)

  url_list := strings.Split(r.RequestURI, "/")

  if len(url_list) >= 2 {
    if url_list[1] == "file" && len(url_list) == 5 {
      pro := url_list[2]
      ver := url_list[3]
      file_name := url_list[4]

      path := "./" + pro + "/dump/" + ver + "/" + file_name
      http.ServeFile(w, r, path)
      return
    }
  }

  log.Println("params:", params)
  // 新版参数解析
  if len(params) >= 2 {
    pat := params["pat"]
    if pat == "" {
      return
    }

    pro := params["pro"]
    if pro == "" {
      return
    }

    if pat == "post" {

      lianyun := params["lianyun"]
      if lianyun == "" {
        lianyun = pro
      }

      reqContent, err := ioutil.ReadAll(r.Body)
      //log.Println("reqContent ", string(reqContent))
      if err != nil {
        log.Println(err)
        return
      } else if reqContent == nil || len(reqContent) < 1 {
        log.Println("empty body!")
        return
      }
      result := CheckLegal(string(reqContent))
      if result {
        proname := "test_project"
        log.Println("ProcessDumpFil :", proname, reqContent, lianyun)
      } else {
        log.Println("error check md5 :", r.Form)
      }
    }

  }

  // 这个写入到w的信息是输出到客户端的
  fmt.Fprintf(w, "Hello gerryyang!\n")
}

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
  if s[0:4] == "java" {
    md5_str = s[5 : index-1]
  }

  if s[0:2] == "js" {
    md5_str = s[3 : index-1]
  }

  check_str := s[index:check_str_len]
  h := md5.New()
  h.Write([]byte(check_str))
  result_str := hex.EncodeToString(h.Sum(nil))
  if md5_str == result_str {

    if s[0:3] == "LOG" {
      return true
    }

    if s[0:4] == "java" {
      return true
    }

    if s[0:2] == "js" {
      return true
    }

    if string(s[check_str_len]) != "M" || string(s[check_str_len+1]) != "D" || string(s[check_str_len+2]) != "M" || string(s[check_str_len+3]) != "P" {
      return false
    }

    return true
  }

  return false
}

func main() {

  // 设置访问的路由
  http.HandleFunc("/", sayhelloName)

  // 设置监听的端口
  err := http.ListenAndServe(":9090", nil)
  if err != nil {
    log.Fatal("ListenAndServe: ", err)
  }
}
