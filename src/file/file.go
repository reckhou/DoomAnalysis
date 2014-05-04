package file

import (
  "log"
  "os"
  "os/exec"
)

func ReadFile(fullPath string) []byte {
  file, err := os.Open(fullPath)
  if err != nil {
    log.Println(err)
    return nil
  }

  fileLen, _ := file.Seek(0, 2)
  data := make([]byte, fileLen)
  file.Seek(0, 0)
  file.Read(data)
  //log.Printf("read %d bytes from %s", readLen, fullPath)

  file.Close()
  return data
}

func WriteFile(fullPath string, content []byte, flag int) bool {
  file, errFile := os.OpenFile(fullPath, os.O_WRONLY|os.O_CREATE|flag, 0777)
  if errFile != nil {
    log.Println(errFile)
    return false
  }

  file.Write(content)
  file.Close()

  return true
}

func IsFileExists(fullPath string) bool {
  if _, err := os.Stat(fullPath); os.IsNotExist(err) {
    log.Println("no such file or directory: ", fullPath)
    return false
  }

  return true
}

func CreateDir(path string) error {
  cmd := exec.Command("mkdir", "-p", path)
  out, err := cmd.Output()
  if err != nil {
    log.Println(string(out))
  }
  return err
}

func DeleteFile(path string) error {
  cmd := exec.Command("rm", "-f", path)
  out, err := cmd.Output()
  if err != nil {
    log.Println(string(out))
  }
  return err
}
