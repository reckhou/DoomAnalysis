package dumpfile

import (
  //"bitbucket.org/reckhou/DoomAnalysis/src/db"
  "bitbucket.org/reckhou/DoomAnalysis/src/dbinfo"
  "bitbucket.org/reckhou/DoomAnalysis/src/file"
  "encoding/binary"
  "encoding/hex"
  "log"
  "os"
  "os/exec"
  "path/filepath"
  "regexp"
  "strconv"
  "strings"
)

// 上传的dump文件行
var key_arr [6]string = [...]string{"MD5", "UUID", "device", "version", "product_name", "file"}

// 上传的log文件行
var key_arr_log [6]string = [...]string{"LOG", "UUID", "device", "version", "product_name", "file"}

type DumpFileInfo struct {
  info_          map[string]string
  file_name_     string
  stack_lib_name []string
  stack_address  []int64
  block_in       bool
  so_address     int64
  project        string
  ndk_stack_info string
  lianyun        string
}

func (info *DumpFileInfo) GenInfo(s string) {

  info.info_ = make(map[string]string)
  line_num := 0
  start_index := 0
  context_start_index := 0
  for i := 0; i < len(s); i++ {
    if s[i] == '\n' {
      if start_index > 0 {
        start_index++
      }
      context_start_index = start_index + len(key_arr[line_num]) + 1

      if line_num == 5 {
        info.info_[key_arr[line_num]] = s[start_index:]
      } else {
        info.info_[key_arr[line_num]] = s[context_start_index:i]
      }

      start_index = i
      line_num++
      if line_num >= len(key_arr) {
        break
      }
    }
  }

  path := "./" + info.project + "/dump/" + info.info_["version"]
  file.CreateDir(path)
  //info.file_name_ = "crash_" + info.info_["UUID"] + "_" + time.Now().Format(time.RFC3339) + ".txt"
  info.file_name_ = info.info_["UUID"] + ".txt"
  file.WriteFile(path+"/"+info.file_name_, []byte(info.info_["file"]), os.O_TRUNC)

}

func (info *DumpFileInfo) GenLogInfo(s string) {

  info.info_ = make(map[string]string)
  line_num := 0
  start_index := 0
  context_start_index := 0
  for i := 0; i < len(s); i++ {
    if s[i] == '\n' {
      if start_index > 0 {
        start_index++
      }
      context_start_index = start_index + len(key_arr_log[line_num]) + 1

      if line_num == 5 {
        info.info_[key_arr_log[line_num]] = s[start_index:]
      } else {
        info.info_[key_arr_log[line_num]] = s[context_start_index:i]
      }

      start_index = i
      line_num++
      if line_num >= len(key_arr_log) {
        break
      }
    }
  }

  path := "./" + info.project + "/dump/" + info.info_["version"]
  file.CreateDir(path)
  //info.file_name_ = "crash_" + info.info_["UUID"] + "_" + time.Now().Format(time.RFC3339) + ".txt"
  info.file_name_ = info.info_["UUID"] + ".log"
  file.WriteFile(path+"/"+info.file_name_, []byte(info.info_["file"]), os.O_TRUNC)

}

func (info *DumpFileInfo) GenSym() bool {
  // 查找是否有对应的 sym文件

  result := file.IsFileExists("./" + info.project + "/lib/" + info.info_["version"] + ".txt")
  if result {
    return true
  }

  lib_name := "./" + info.project + "/lib/" + info.info_["version"] + "_libgame.so"
  result = file.IsFileExists(lib_name)
  if result {
    cmd := exec.Command("/bin/sh", "gensym.sh", info.info_["version"], info.project)
    _, err := cmd.Output()
    if err != nil {
      log.Println("GenSym err:" + err.Error())
      return false
    }
    return true
  }

  return false
}

func (info *DumpFileInfo) GenBreakpadDumpInfo() {
  cmd := exec.Command("/bin/sh", "./gen_dump_info.sh", info.info_["version"], info.file_name_, info.project)
  _, err := cmd.Output()
  if err != nil {
    log.Println("GenSym err:" + err.Error())
  }
}

func (info *DumpFileInfo) GenNdkDumpInfo() {
  context := file.ReadFile("./" + info.project + "/dump/" + info.info_["version"] + "/" + info.file_name_ + ".info")
  //context := file.ReadFile("./a.txt.info")
  start_pos := 0

  open_stack := false
  open_lib_info := false

  for i := 0; i < len(context); i++ {
    if context[i] == '\n' {
      if i-start_pos > 1 {
        temp_str := string(context[start_pos:i])

        matched, _ := regexp.MatchString("(?i:^Thread).*(crashed)", temp_str)
        if matched {
          //log.Println("regexp : ", temp_str)
          open_stack = true
        }

        if open_stack {
          info.GenNdkStack(temp_str)
          matched, _ = regexp.MatchString("(?i:^Thread)[\\s]\\d{1,2}$", temp_str)
          if matched {
            //log.Println("regexp : ", temp_str)
            open_stack = false
          }
        }

        matched, _ = regexp.MatchString("(?i:^Loaded)[\\s](modules:)$", temp_str)
        if matched {
          //log.Println("regexp : ", temp_str)
          open_lib_info = true
        }

        if open_lib_info {
          matched, _ = regexp.MatchString("0x[0-9|a-f]{8}\\s-\\s0x[0-9|a-f]{8}\\s{2}libgame.so\\s{2}\\?{3}$", temp_str)
          if matched {
            //log.Println("regexp : ", temp_str)
            info.GenNdkSoAddress(temp_str)
            open_lib_info = false
          }
        }

      }

      start_pos = i + 1

    }
  }

  info.GenNdkfile()

}

func (info *DumpFileInfo) GenNdkStack(s string) {
  // pc\s=\s0x[0-9|a-z]{8}
  // \s\d{1,4}\s{2}[^\s]+\.so
  if !info.block_in {
    re := regexp.MustCompile("\\s\\d{1,4}\\s{2}[^\\s]+\\.so")
    matched := re.FindString(s)
    info.ndk_stack_info = info.ndk_stack_info + s + "<br>"
    if matched != "" {

      re = regexp.MustCompile("[a-z&A-Z&\\.]+")
      matched = re.FindString(s)

      if matched != "" {
        info.stack_lib_name = append(info.stack_lib_name, matched)
        info.block_in = true
      }

    }

  } else {
    re := regexp.MustCompile("pc\\s=\\s0x[0-9|a-z]{8}")
    matched := re.FindString(s)
    if matched != "" {

      re = regexp.MustCompile("0x[a-z&0-9]{8}")
      address := re.FindString(matched)

      if address != "" {

        hex_str := address[2:]
        re, err := hex.DecodeString(hex_str)

        if err != nil {
          log.Println("hex err :", err)
        }

        info.stack_address = append(info.stack_address, int64(binary.BigEndian.Uint32(re)))
        info.block_in = false
      }

    }
  }
}

func (info *DumpFileInfo) GenNdkSoAddress(s string) {
  // pc\s=\s0x[0-9|a-z]{8}
  // \s\d{1,4}\s{2}[^\s]+\.so
  re := regexp.MustCompile("0x[a-z&0-9]{8}")
  address := re.FindString(s)

  if address != "" {

    hex_str := address[2:]
    re, err := hex.DecodeString(hex_str)

    if err != nil {
      log.Println("hex err :", err)
    }

    info.so_address = int64(binary.BigEndian.Uint32(re))
  }
}

func (info *DumpFileInfo) GenNdkfile() {
  /* eg.
  03-24 15:34:32.361: I/DEBUG(130): *** *** *** *** *** *** *** *** *** *** *** *** *** *** *** ***
  03-24 15:34:32.361: I/DEBUG(130): Build fingerprint: '111'
  03-24 15:34:32.361: I/DEBUG(130): pid: 11142, tid: 11365, name: dfdsf  >>> sfsaf <<<
  03-24 15:34:32.794: I/DEBUG(130): backtrace:
  03-24 15:34:32.794: I/DEBUG(130):     #00  pc 0069D08F  libgame.so ()
  03-24 15:34:32.794: I/DEBUG(130):     #01  pc 5e4ade49  libgame.so ()
  */
  max_info_count := 3
  info_key := ""

  file_context := "03-24 15:34:32.361: I/DEBUG(130): *** *** *** *** *** *** *** *** *** *** *** *** *** *** *** ***\n"
  file_context += "03-24 15:34:32.361: I/DEBUG(130): Build fingerprint: '111'\n"
  file_context += "03-24 15:34:32.361: I/DEBUG(130): pid: 11142, tid: 11365, name: dfdsf  >>> sfsaf <<<\n"
  file_context += "03-24 15:34:32.794: I/DEBUG(130): backtrace:\n"
  stack_head := "03-24 15:34:32.794: I/DEBUG(130):     #"
  for i := 0; i < len(info.stack_lib_name); i++ {
    result_str := ""
    address := info.stack_address[i] - info.so_address

    var buf = make([]byte, 4)
    binary.BigEndian.PutUint32(buf, uint32(address))

    if i < 10 {
      result_str = "0" + strconv.Itoa(i) + "  pc " + hex.EncodeToString(buf) + "  " + info.info_["version"] + "_" + info.stack_lib_name[i] + " ()\n"
      file_context += (stack_head + result_str)
    }

    if max_info_count > 0 {
      info_key = info_key + hex.EncodeToString(buf) + "_"
      max_info_count--
    }

  }
  file_context += "\n"

  path := "./" + info.project + "/dump/" + info.info_["version"]
  file.WriteFile(path+"/"+info.file_name_+".ndk", []byte(file_context), os.O_TRUNC)

  cmd := exec.Command("/bin/sh", "./gen_ndk_info.sh", info.info_["version"], info.file_name_+".ndk", info.project)
  _, err := cmd.Output()
  if err != nil {
    log.Println("GenNdkfile err:" + err.Error())
  }
}

func (info *DumpFileInfo) GenDbInfo() {
  context := file.ReadFile("./" + info.project + "/dump/" + info.info_["version"] + "/" + info.file_name_ + ".ndk.info")
  //context := file.ReadFile("./a.txt.info")
  start_pos := 0

  address_info := 3
  info_str := ""
  info_key := ""
  for i := 0; i < len(context); i++ {
    if context[i] == '\n' {
      if i-start_pos > 1 {
        temp_str := string(context[start_pos:i])

        re := regexp.MustCompile("pc\\s{1}[0-9|a-f]{8}")
        matched := re.FindString(temp_str)

        if matched != "" && address_info > 0 {

          prostack_flag, _ := regexp.MatchString("libgame.so", temp_str)
          if prostack_flag {
            re = regexp.MustCompile("[0-9|a-f]{8}")
            key := re.FindString(matched)

            info_key = info_key + "_" + key
            address_info--
          }

        }

        info_str = info_str + temp_str + "<br>"

      }

      start_pos = i + 1

    }
  }

  //db.CreateDB(info.project, info.info_["version"], info_key, info_str, info.info_["UUID"])
  mysql_c, db_err := dbinfo.Init()
  if db_err == nil {
    defer mysql_c.Close()
    mysql_c.CreateDB(info.project, info.info_["version"], info_key, info_str, info.info_["UUID"], info.lianyun)
  }
}

func (info *DumpFileInfo) GenTar(mode string) {
  // info.info_["UUID"]
  cmd := exec.Command("/bin/sh", "gen_tar.sh", info.info_["version"], info.project, info.info_["UUID"], mode)
  _, err := cmd.Output()
  if err != nil {
    log.Println("GenSym err:" + err.Error())
  }
}

func ProcessDumpFile(project string, co []byte, lianyun string) {

  //context := file.ReadFile("./a.txt")
  //s := string(context)
  s := string(co)

  var info DumpFileInfo
  info.block_in = false
  info.project = project
  info.lianyun = lianyun

  if s[0:3] == "LOG" {
    info.GenLogInfo(s)
  } else {
    info.GenInfo(s)
    result := info.GenSym()
    if result {
      info.GenBreakpadDumpInfo()
      info.GenNdkDumpInfo()
      info.GenDbInfo()
      // tar
      info.GenTar("c")
    }
  }

}

func ListFileName(path string, ver string, pro string, lianyun string) {
  dbinfo.DeleteInfoDB(pro, ver)
  fullPath, _ := filepath.Abs(path)
  log.Println("ListFileName: ", fullPath)
  filepath.Walk(fullPath, func(path string, fi os.FileInfo, err error) error {
    if nil == fi {
      return err
    }
    if fi.IsDir() {
      return nil
    }

    name := fi.Name()
    file_list := strings.Split(name, ".")
    filename := file_list[0]

    if len(file_list) >= 3 && file_list[2] == "info" {
      var info DumpFileInfo
      info.block_in = false
      info.project = pro
      info.lianyun = lianyun

      info.info_ = make(map[string]string)
      info.info_[key_arr[1]] = filename
      info.info_[key_arr[3]] = ver
      info.info_[key_arr[4]] = pro

      info.file_name_ = info.info_["UUID"] + ".txt"
      log.Println("recreate: ", name)
      info.GenTar("x")

      result := info.GenSym()
      if result {
        info.GenBreakpadDumpInfo()
        info.GenNdkDumpInfo()
        info.GenDbInfo()
        info.GenTar("c")
      }
    }

    return nil
  })
}
