package DoomAnalysis

import (
	gozd "bitbucket.org/PinIdea/zero-downtime-daemon"
	"bitbucket.org/reckhou/DoomAnalysis/src/dbinfo"
	"bitbucket.org/reckhou/DoomAnalysis/src/dumpfile"
	"bitbucket.org/reckhou/DoomAnalysis/src/file"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"syscall"
)

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

	ctx := gozd.Context{
		Hash:    "pin_dump",
		Logfile: os.TempDir() + "/pin_dump.log",
		Directives: map[string]gozd.Server{
			"sock": gozd.Server{
				Network: "unix",
				Address: os.TempDir() + "/pin_dump.sock",
			},
			"port1": gozd.Server{
				Network: "tcp",
				Address: "115.28.202.194:10010",
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

	//log.Println(s.ListenAndServe())
}

func (s HTTPServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	r.ParseForm()
	log.Println("Form:", r.Form)

	// 新版参数解析
	if len(r.Form) >= 2 {
		pat_array := r.Form["pat"]
		if len(pat_array) < 0 {
			log.Println("url error")
			return
		}

		pat := r.Form["pat"][0]

		pro_array := r.Form["pro"][0]
		if len(pro_array) < 0 {
			log.Println("url error")
			return
		}

		pro := r.Form["pro"][0]

		if pat == "get" {
			ver_array := r.Form["ver"]
			if len(ver_array) <= 0 {
				fmt.Fprintf(w, dbinfo.VerInfoDB(pro))
			} else {
				ver := r.Form["ver"][0]
				fmt.Fprintf(w, dbinfo.GetListInfoDB(pro, ver))
			}
		} else if pat == "post" {
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
				go dumpfile.ProcessDumpFile(pro, reqContent)
			} else {
				log.Println("error check md5 :", r.Form)
			}
		} else if pat == "file" {
			filename_array := r.Form["filename"]

			if len(filename_array) < 0 {
				log.Println("url error")
				return
			}
			filename := r.Form["filename"][0]

			ver_array := r.Form["ver"]

			if len(ver_array) < 0 {
				log.Println("url error")
				return
			}
			ver := r.Form["ver"][0]

			path := "./" + pro + "/dump/" + ver + "/" + filename
			fmt.Fprintf(w, string(file.ReadFile(path)))
		} else if pat == "recreate" {
			ver_array := r.Form["ver"]

			if len(ver_array) < 0 {
				log.Println("url error")
				return
			}
			ver := r.Form["ver"][0]
			path := "./" + pro + "/dump/" + ver + "/"

			go dumpfile.ListFileName(path, ver, pro)
		} else if pat == "detail" {
			id_array := r.Form["id"]

			if len(id_array) < 0 {
				log.Println("url error")
				return
			}
			id := r.Form["id"][0]

			ver_array := r.Form["ver"]

			if len(ver_array) < 0 {
				log.Println("url error")
				return
			}
			ver := r.Form["ver"][0]

			fmt.Fprintf(w, dbinfo.GetFileListInfoDB(pro, ver, id))
		}

	}

	if len(r.Form) == 1 {
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
				go dumpfile.ProcessDumpFile("sxd", reqContent)
			} else {
				log.Println("error check md5 :", r.Form)
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

		if s[0:3] == "LOG" {
			return true
		}

		if string(s[check_str_len]) != "M" || string(s[check_str_len+1]) != "D" || string(s[check_str_len+2]) != "M" || string(s[check_str_len+3]) != "P" {
			return false
		}

		return true
	} else {
		return false
	}

	return false
}
