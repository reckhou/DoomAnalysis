package main

import (
	"bitbucket.org/reckhou/DoomAnalysis/src/file"
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"github.com/reckhou/goCfgMgr"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"time"
)

var c chan int

func TestJava(url string, ver string) {
	dump_context := file.ReadFile("./testresourse/test.java")
	out, err := exec.Command("uuidgen").Output()
	if err != nil {
		log.Println("uuidgen err :", err)
	}
	msg := "UUID:" + string(out) + "device:test_client_tool\nversion:" + ver + "\nproduct_name:test\n"

	dump := "file:" + string(dump_context)
	h := md5.New()
	h.Write([]byte(msg))
	result_str := hex.EncodeToString(h.Sum(nil))

	msg_sender := "java:" + result_str + "\n" + msg + dump
	body := bytes.NewBufferString(msg_sender)
	log.Println(url)
	resp, err_http := http.Post(url, "", body)
	if err != nil {
		log.Println("TestJava err :", err_http)
	}
	if resp != nil {
		resp.Body.Close()
	}

	c <- 1
}

func TestJs(url string, ver string) {
	dump_context := file.ReadFile("./testresourse/test.js")
	out, err := exec.Command("uuidgen").Output()
	if err != nil {
		log.Println("uuidgen err :", err)
	}
	msg := "UUID:" + string(out) + "device:test_client_tool\nversion:" + ver + "\nproduct_name:test\n"
	dump := "file:" + string(dump_context)
	h := md5.New()
	h.Write([]byte(msg))
	result_str := hex.EncodeToString(h.Sum(nil))

	msg_sender := "js:" + result_str + "\n" + msg + dump
	body := bytes.NewBufferString(msg_sender)
	resp, err_http := http.Post(url, "", body)
	if err_http != nil {
		log.Println(err_http)
	}
	defer resp.Body.Close()
	c <- 1
}

func TestC(url string, ver string) {
	dump_context := file.ReadFile("./testresourse/c.js")
	out, err := exec.Command("uuidgen").Output()
	if err != nil {
		log.Println("uuidgen err :", err)
	}
	msg := "UUID:" + string(out) + "device:test_client_tool\nversion:" + ver + "\nproduct_name:test\n"
	dump := "file:" + string(dump_context)
	h := md5.New()
	h.Write([]byte(msg))
	result_str := hex.EncodeToString(h.Sum(nil))

	msg_sender := "MD5:" + result_str + "\n" + msg + dump
	body := bytes.NewBufferString(msg_sender)
	resp, err_http := http.Post(url, "", body)
	if err_http != nil {
		log.Println(err_http)
	}
	defer resp.Body.Close()
	c <- 1
}

func main() {
	arg_num := len(os.Args)
	if arg_num < 4 {
		log.Printf("[a|java|c|js] [并发数量] [测试次数]\n")
		return
	}

	mod := os.Args[1]

	test_thread_num, err := strconv.Atoi(os.Args[2])
	if err != nil {
		log.Println("error happened ,exit")
		return
	}

	test_num, err := strconv.Atoi(os.Args[3])
	if err != nil {
		log.Println("error happened ,exit")
		return
	}
	ver := goCfgMgr.Get("project", "ver").(string)
	server_address := "http://" + goCfgMgr.Get("basic", "Host").(string) + ":" +
		goCfgMgr.Get("basic", "Port").(string) + "/?pat=post&pro=" +
		goCfgMgr.Get("project", "name").(string) + "&ver=" + ver
	c = make(chan int)
	for j := 0; j < test_num; j++ {
		for i := 0; i < test_thread_num; i++ {
			if mod == "java" {
				go TestJava(server_address, ver)
			} else if mod == "js" {
				go TestJs(server_address, ver)
			} else if mod == "c" {
				go TestC(server_address, ver)
			} else if mod == "a" {
				go TestJava(server_address, ver)
				go TestJs(server_address, ver)
				go TestC(server_address, ver)
			}
		}

		for i := 0; i < test_thread_num; i++ {
			<-c
		}
		time.Sleep(100 * time.Millisecond)
	}
}
