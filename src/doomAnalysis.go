package DoomAnalysis

import (
  gozd "bitbucket.org/PinIdea/zero-downtime-daemon"
  "bitbucket.org/reckhou/DoomAnalysis/src/dbinfo"
  "bitbucket.org/reckhou/DoomAnalysis/src/debug"
  "bitbucket.org/reckhou/DoomAnalysis/src/dumpfile"
  "bitbucket.org/reckhou/DoomAnalysis/src/file"
  "crypto/md5"
  "encoding/hex"
  "fmt"
  "github.com/reckhou/goCfgMgr"
  "io/ioutil"
  "log"
  "net"
  "net/http"
  "net/url"
  "os"
  "strings"
  "syscall"
)

import _ "net/http/pprof"

type HTTPServer struct{}

func handleListners(cl chan net.Listener) {

  for v := range cl {
    go func(l net.Listener) {
      handler := new(HTTPServer)
      http.Serve(l, handler)
    }(v)
  }
}

func Start() {
  server_address := goCfgMgr.Get("basic", "Host").(string) + ":" +
    goCfgMgr.Get("basic", "Port").(string)

  ctx := gozd.Context{
    Hash:    "pin_dump_test",
    Logfile: os.TempDir() + "/pin_dump_test.log",
    Directives: map[string]gozd.Server{
      "sock": gozd.Server{
        Network: "unix",
        Address: os.TempDir() + "/pin_dump_test.sock",
      },
      "port1": gozd.Server{
        Network: "tcp",
        Address: server_address,
      },
    },
  }

  cl := make(chan net.Listener, 1)
  go handleListners(cl)
  sig, err := gozd.Daemonize(ctx, cl) // returns channel that connects with daemon
  if err != nil {
    log.Println("error: ", err)
    return
  }

  // other initializations or config setting

  for s := range sig {
    switch s {
    case syscall.SIGHUP, syscall.SIGUSR2:
      // do some custom jobs while reload/hotupdate

    case syscall.SIGTERM:
      // do some clean up and exit
      return
    }
  }

  dbinfo.Init()

  go dbinfo.Check_Sql_Connect()
  go debug.CheckMemStats()
}

func getUrlParameter(key string, form url.Values) (string, bool) {

  value_array := form[key]
  if len(value_array) < 0 {
    log.Println("getUrlParameter error : ", key)
    return "", false
  }

  pat := form[key][0]

  return pat, true
}

// target=API name or FileName want to download
// project=folderSlice[1]
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

  /*if debug.Enabled {
    log.Println("Folders:\n", folderSlice)
    log.Println("Target:\n", target)
    log.Println("Params:\n", params)
  }*/

  return
}

func (s HTTPServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {

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
    } else if url_list[1] == "tencent" {
      context := file.ReadFile("./tencent_create.html")
      fmt.Fprintf(w, string(context))
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

    if pat == "get" {
      ver_array := params["ver"]
      if ver_array == "" {
        fmt.Fprintf(w, dbinfo.GerVersionList(pro))
      } else {
        ver := params["ver"]
        fmt.Fprintf(w, dbinfo.GetDumpList(pro, ver))
      }
    } else if pat == "post" {

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
        proname := GetProName(pro, lianyun)
        go dumpfile.ProcessDumpFile(proname, reqContent, lianyun)
      } else {
        log.Println("error check md5 :", r.Form)
      }
    } else if pat == "recreate" {

      ver := params["ver"]
      if ver == "" {
        return
      }

      lianyun := params["lianyun"]
      if lianyun == "" {
        return
      }

      proname := GetProName(pro, lianyun)
      path := "./" + proname + "/dump/" + ver + "/"

      go dumpfile.ListFileName(path, ver, pro, lianyun)

    } else if pat == "detail" {

      id := params["id"]
      if id == "" {
        return
      }

      ver := params["ver"]
      if ver == "" {
        return
      }

      fmt.Fprintf(w, dbinfo.GetDumpFileList(pro, ver, id))

    } else if pat == "create_tencent" {

      ver := params["ver"]
      if ver == "" {
        return
      }

      lianyun := params["lianyun"]
      if lianyun == "" {
        return
      }

      proname := GetProName(pro, lianyun)
      path := "./" + proname + "/tencentdump/"

      go dumpfile.ListTencentFileName(path, ver, pro, lianyun)

    } else if pat == "allversion" {

      fmt.Fprintf(w, dbinfo.GetAllVersionList(pro))

    }

  }

  /*if len(r.Form) == 1 {
  	val := r.Form["par"]

  	if len(val) > 0 && val[0] == "post" {
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
  			go dumpfile.ProcessDumpFile("sxda", reqContent, "sxda")
  		} else {
  			log.Println("error check md5 :", r.Form)
  		}
  	}
  }*/

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

func GetProName(pro string, lianyun string) string {
  log.Println("GetProName:", pro)
  pro_list := goCfgMgr.Get("project", pro).(map[string]interface{})

  for i, u := range pro_list {
    switch value := u.(type) {
    case string:
      if i == lianyun {
        return value
      }
    default:
      fmt.Println(value, "is of a type I don't know how to handle ")
    }

  }

  return pro
}
