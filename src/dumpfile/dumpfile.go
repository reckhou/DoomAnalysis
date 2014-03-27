package dumpfile

import (
	//"./src"
	"../file"
	"log"
	"os"
	"os/exec"
	"time"
)

// 上传的dump文件行
var key_arr [6]string = [...]string{"MD5", "UUID", "device", "version", "product_name", "file"}

type DumpFileInfo struct {
	info_      map[string]string
	file_name_ string
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

	path := "./dump/" + info.info_["version"]
	file.CreateDir(path)
	info.file_name_ = "crash_" + info.info_["UUID"] + "_" + time.Now().Format(time.RFC3339) + ".txt"
	file.WriteFile(path+"/"+info.file_name_, []byte(info.info_["file"]), os.O_TRUNC)

}

func (info *DumpFileInfo) GenSym() bool {
	// 查找是否有对应的 sym文件

	result := file.IsFileExists("./lib/" + info.info_["version"] + ".txt")
	if result {
		return true
	}

	lib_name := "lib/" + info.info_["version"] + "_libgame.so"
	result = file.IsFileExists(lib_name)
	if result {
		cmd := exec.Command("/bin/sh", "gensym.sh", info.info_["version"])
		b, err := cmd.Output()
		if err != nil {
			log.Println("GenSym err:" + err.Error())
		}
		log.Println("GenSym info:" + string(b))
	}

	return true
}

func (info *DumpFileInfo) GenBreakpadDumpInfo() {
	cmd := exec.Command("/bin/sh", "./gen_dump_info.sh", info.info_["version"], info.file_name_)
	b, err := cmd.Output()
	if err != nil {
		log.Println("GenSym err:" + err.Error())
	}
	log.Println("GenSym info:" + string(b))

}

func (info *DumpFileInfo) GenNdkDumpInfo() {

}

func ProcessDumpFile(co []byte) {
	log.Println("ProcessDumpFile start")
	context := file.ReadFile("./a.txt")

	s := string(context)
	var info DumpFileInfo
	info.GenInfo(s)
	result := info.GenSym()
	if result {
		log.Println("GenBreakpadDumpInfo start")
		info.GenBreakpadDumpInfo()
		info.GenNdkDumpInfo()
	}

	log.Println("ProcessDumpFile end")
}
