package dumpfile

import (
  "bitbucket.org/reckhou/DoomAnalysis/src/cplus"
  "bitbucket.org/reckhou/DoomAnalysis/src/dbinfo"
  "bitbucket.org/reckhou/DoomAnalysis/src/javainfo"
  "bitbucket.org/reckhou/DoomAnalysis/src/js"
  "log"
  "os"
  "path/filepath"
  "strings"
)

func ProcessDumpFile(project string, co []byte, lianyun string) {

  s := string(co)

  var info cplus.DumpFileInfo
  info.InitData(project, lianyun)

  if s[0:3] == "LOG" {
    info.GenLogInfo(s)
  } else if s[0:4] == "java" {

    var javainfo_obj javainfo.JavaFileInfo
    pro := project + "_java"
    javainfo_obj.SetProjectInfo(pro, lianyun)
    javainfo_obj.GenJavaInfo(s)
    javainfo_obj.GenJavaDBInfo()
    javainfo_obj.GenTar("c")

  } else if s[0:2] == "js" {

    var js_obj js.JsFileInfo
    pro := project + "_js"
    js_obj.SetProjectInfo(pro, lianyun)
    js_obj.GenJsInfo(s)
    js_obj.GenJsDBInfo()
    js_obj.GenTar("c")

  } else {
    info.GenInfo(s)
    result := info.GenSym()
    if result {
      info.GenBreakpadDumpInfo()
      info.GenNdkDumpInfo()
      info.GenDbInfo()
      // tar
      info.GenTar("c")
    } else {
      log.Println("c++ dump error: ", info.GetVersion())
    }
  }

}

func ListFileName(path string, ver string, pro string, lianyun string) {
  dbinfo.DeleteInfo(pro, ver)
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

    if len(file_list) <= 2 && file_list[1] == "txt" {
      cplus.RecreateDumpInfo(pro, lianyun, filename, ver, name)
    }

    return nil
  })
}

func ListTencentFileName(path string, ver string, pro string, lianyun string) {
  fullPath, _ := filepath.Abs(path)
  log.Println("ListTencentFileName Path: ", fullPath)
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
    if len(file_list) <= 2 && file_list[1] == "zip" {
      cplus.CreateTencentDumpInfo(pro, lianyun, filename, ver, name)
    }

    return nil
  })
}
