package main

import (
	"bitbucket.org/reckhou/DoomAnalysis/src"
	"os"
)

// import (
// 	"bitbucket.org/reckhou/DoomAnalysis/src/dbinfo"
// 	"github.com/syndtr/goleveldb/leveldb"
// 	"log"
// 	"os"
// 	"regexp"
// 	"strconv"
// )

func main() {
	os.Chdir("/root/dumpserver/")
	DoomAnalysis.Start()

	// mysqldb, _ := dbinfo.Init()
	// defer mysqldb.Close()

	// db_id, db_id_err := leveldb.OpenFile("./sxd/dump/5572/db_id.db", nil)
	// db_info, db_info_err := leveldb.OpenFile("./sxd/dump/5572/db_info.db", nil)
	// db_count, db_count_err := leveldb.OpenFile("./sxd/dump/5572/db_count.db", nil)
	// db_breakpad_info, db_breakpad_info_err := leveldb.OpenFile("./sxd/dump/5572/db_breakpad_info.db", nil)
	// defer db_id.Close()
	// defer db_info.Close()
	// defer db_count.Close()
	// defer db_breakpad_info.Close()

	// if db_id_err != nil {
	// 	log.Print("GetListInfoDB db_id_err: ", db_id_err)
	// 	return
	// }

	// if db_info_err != nil {
	// 	log.Print("GetListInfoDB db_info_err: ", db_info_err)
	// 	return
	// }

	// if db_count_err != nil {
	// 	log.Print("GetListInfoDB db_count_err: ", db_count_err)
	// 	return
	// }

	// if db_breakpad_info_err != nil {
	// 	log.Print("GetListInfoDB db_breakpad_info_err: ", db_breakpad_info_err)
	// 	return
	// }

	// max_num_val, _ := db_id.Get([]byte("MAX_NUM"), nil)
	// if max_num_val != nil {
	// 	num_val, err := strconv.Atoi(string(max_num_val))
	// 	if err != nil {
	// 		log.Println("MAX_NUM error:", err)
	// 		return
	// 	}
	// 	max_num := num_val
	// 	for i := 1; i <= max_num; i++ {

	// 		id_val := strconv.Itoa(i)
	// 		address, _ := db_id.Get([]byte(id_val), nil)

	// 		if address != nil {

	// 			info, _ := db_info.Get([]byte(address), nil)
	// 			info_pad, _ := db_breakpad_info.Get([]byte(address), nil)

	// 			re := regexp.MustCompile(">[a-z|0-9|.-]{0,100}</")
	// 			matched := re.FindAllString(string(info_pad), -1)
	// 			uuid := ""
	// 			for _, value := range matched {
	// 				str_last := value[len(value)-6 : len(value)-2]
	// 				if string(str_last) == ".log" {
	// 					uuid = string(value[1 : len(value)-6])
	// 					mysqldb.CreateDB("sxd", "5572", string(address), string(info), uuid)
	// 				}
	// 			}

	// 		}
	// 	}
	// }

}
