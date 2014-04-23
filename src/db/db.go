package db

import (
	"github.com/syndtr/goleveldb/leveldb"
	"log"
	"regexp"
	"strconv"
)

var db_channel chan int

func InitChannelFlag() {
	db_channel = make(chan int, 1)

}

func CreateDB(pro string, ver string, address string, info string, uuid string) {

	db_channel <- 1
	db_version, err := leveldb.OpenFile("./"+pro+"/version_id.db", nil)

	if err != nil {
		log.Println("db_version OpenFile err: ", err)
		<-db_channel
		return
	}
	defer db_version.Close()

	version, version_err := db_version.Get([]byte(ver), nil)
	if version_err != nil {
		log.Println("db_version ger err: ", version_err)
	}
	if version == nil {
		db_version.Put([]byte(ver), []byte("0"), nil)
	}

	db_id, id_err := leveldb.OpenFile("./"+pro+"/dump/"+ver+"/db_id.db", nil)

	if id_err != nil {
		log.Println("db_id OpenFile err: ", id_err)
		<-db_channel
		return
	}
	defer db_id.Close()

	db_info, info_err := leveldb.OpenFile("./"+pro+"/dump/"+ver+"/db_info.db", nil)

	if info_err != nil {
		log.Println("db_info OpenFile err: ", info_err)
		<-db_channel
		return
	}
	defer db_info.Close()

	db_breakpad_info, breakpad_info_err := leveldb.OpenFile("./"+pro+"/dump/"+ver+"/db_breakpad_info.db", nil)

	if breakpad_info_err != nil {
		log.Println("db_breakpad_info OpenFile err: ", breakpad_info_err)
		<-db_channel
		return
	}
	defer db_breakpad_info.Close()

	db_count, db_count_err := leveldb.OpenFile("./"+pro+"/dump/"+ver+"/db_count.db", nil)

	if db_count_err != nil {
		log.Println("db_count OpenFile err: ", db_count_err)
		<-db_channel
		return
	}
	defer db_count.Close()

	// if ro and wo are not used again, be sure to Close them.

	max_num := 0
	s, _ := db_id.Get([]byte("MAX_NUM"), nil)
	if s != nil {
		//log.Println("MAX_NUM: ", string(s))
		num_val, err := strconv.Atoi(string(s))
		if err != nil {
			log.Println("MAX_NUM error:", ver)
			<-db_channel
			return
		}
		max_num = num_val
	}
	max_num++

	info_val := "<a href=\"?par=file&ver=" + ver + "&pro=sxd&filename=" + uuid + ".log \">" + uuid + ".log" + "</a><br>"
	info_val = info_val + "<a href=\"?par=file&ver=" + ver + "&pro=" + pro + "&filename=" + uuid + ".txt.info \">" + uuid + ".info" + "</a><br>"
	info_val = info_val + "<a href=\"?par=file&ver=" + ver + "&pro=" + pro + "&filename=" + uuid + ".txt.ndk.info \">" + uuid + ".ndk" + "</a><br>"
	s, _ = db_breakpad_info.Get([]byte(address), nil)
	if s == nil {
		db_breakpad_info.Put([]byte(address), []byte(info_val), nil)
	} else {
		valueall := string(s) + info_val
		db_breakpad_info.Put([]byte(address), []byte(valueall), nil)
	}

	s, _ = db_id.Get([]byte(address), nil)
	if s == nil {
		max_num_val := strconv.Itoa(max_num)
		db_id.Put([]byte("MAX_NUM"), []byte(max_num_val), nil)
		db_id.Put([]byte(address), []byte(max_num_val), nil)
		db_id.Put([]byte(max_num_val), []byte(address), nil)

	}

	s, _ = db_info.Get([]byte(address), nil)
	if s == nil {
		db_info.Put([]byte(address), []byte(info), nil)
	}

	s, _ = db_count.Get([]byte(address), nil)
	if s == nil {
		db_count.Put([]byte(address), []byte("1"), nil)
	} else {
		num_val, err := strconv.Atoi(string(s))
		if err != nil {
			return
		}
		num_val++
		address_num_val := strconv.Itoa(num_val)
		db_count.Put([]byte(address), []byte(address_num_val), nil)
	}

	<-db_channel
}

func GetListInfoDB(pro string, ver string) string {

	db_id, db_id_err := leveldb.OpenFile("./"+pro+"/dump/"+ver+"/db_id.db", nil)
	db_info, db_info_err := leveldb.OpenFile("./"+pro+"/dump/"+ver+"/db_info.db", nil)
	db_count, db_count_err := leveldb.OpenFile("./"+pro+"/dump/"+ver+"/db_count.db", nil)
	db_breakpad_info, db_breakpad_info_err := leveldb.OpenFile("./"+pro+"/dump/"+ver+"/db_breakpad_info.db", nil)
	defer db_id.Close()
	defer db_info.Close()
	defer db_count.Close()
	defer db_breakpad_info.Close()

	if db_id_err != nil {
		log.Print("GetListInfoDB db_id_err: ", db_id_err)
		return ""
	}

	if db_info_err != nil {
		log.Print("GetListInfoDB db_info_err: ", db_info_err)
		return ""
	}

	if db_count_err != nil {
		log.Print("GetListInfoDB db_count_err: ", db_count_err)
		return ""
	}

	if db_breakpad_info_err != nil {
		log.Print("GetListInfoDB db_breakpad_info_err: ", db_breakpad_info_err)
		return ""
	}
	return_val := ""
	max_num_val, _ := db_id.Get([]byte("MAX_NUM"), nil)
	if max_num_val != nil {
		num_val, err := strconv.Atoi(string(max_num_val))
		if err != nil {
			log.Println("MAX_NUM error:", ver)
			return ""
		}
		max_num := num_val
		precent := 0

		for i := 1; i <= max_num; i++ {

			id_val := strconv.Itoa(i)
			address, _ := db_id.Get([]byte(id_val), nil)

			if address != nil {

				count, _ := db_count.Get([]byte(address), nil)
				num, err := strconv.Atoi(string(count))
				if err == nil {
					precent += num
				} else {
					return ""
				}

			}
		}

		return_val = "<html>\n<body>\n<table border=\"1\">\n"
		return_val = return_val + "<tr>\n"
		return_val = return_val + "<th align=\"left\">ID</th>\n"
		return_val = return_val + "<th align=\"right\">ADDRESS</th>\n"
		return_val = return_val + "<th align=\"right\">COUNT (" + strconv.Itoa(precent) + ") </th>\n"
		return_val = return_val + "<th align=\"center\">INFO</th>\n"
		return_val = return_val + "<th align=\"center\">BREAKPAD</th>\n"
		return_val = return_val + "</tr>\n"

		for i := 1; i <= max_num; i++ {

			id_val := strconv.Itoa(i)
			address, _ := db_id.Get([]byte(id_val), nil)

			if address != nil {

				info, _ := db_info.Get([]byte(address), nil)
				count, _ := db_count.Get([]byte(address), nil)
				info_pad, _ := db_breakpad_info.Get([]byte(address), nil)

				num, _ := strconv.Atoi(string(count))
				pre_value := float64(num) * 100.0 / float64(precent)

				pre_value_str := strconv.FormatFloat(pre_value, 'f', 2, 64)

				color := ""
				if pre_value > 0.5 {
					color = " style=\"color:#F00\" "
				}
				return_val = return_val + "<tr>\n"
				return_val = return_val + "<th align=\"left\"><div " + color + ">" + id_val + "</div></th>\n"
				return_val = return_val + "<th align=\"right\"><div " + color + ">" + string(address) + "</div></th>\n"
				return_val = return_val + "<th align=\"right\"><div " + color + ">" + string(count) + " (" + pre_value_str + "%%)" + " </div></th>\n"
				return_val = return_val + "<th align=\"left\"><div " + color + ">" + string(info) + " </div></th>\n"

				re := regexp.MustCompile(">[a-z|0-9|.-]{0,100}</")
				matched := re.FindAllString(string(info_pad), -1)
				for _, value := range matched {
					str_last := value[len(value)-6 : len(value)-2]
					if string(str_last) == ".log" {
						uuid := string(value[1 : len(value)-6])

						info_val := "<a href=\"?par=file&ver=" + ver + "&pro=sxd&filename=" + uuid + ".log \">" + uuid + ".log" + "</a><br>"
						info_val = info_val + "<a href=\"?par=file&ver=" + ver + "&pro=" + pro + "&filename=" + uuid + ".txt.info \">" + uuid + ".info" + "</a><br>"
						info_val = info_val + "<a href=\"?par=file&ver=" + ver + "&pro=" + pro + "&filename=" + uuid + ".txt.ndk.info \">" + uuid + ".ndk" + "</a><br>"

						return_val = return_val + "<th align=\"left\">" + info_val + " </th>\n"
						break
					}
				}

				return_val = return_val + "</tr>\n"
			}
		}
	}
	return_val = return_val + "</table>\n</body>\n</html>"
	return return_val
}

func VerInfoDB(pro string) string {

	db_version, db_version_err := leveldb.OpenFile("./"+pro+"/version_id.db", nil)
	defer db_version.Close()
	if db_version_err != nil {
		log.Print("VerInfoDB db_version_err: ", db_version_err)
		return ""
	}
	it := db_version.NewIterator(nil, nil)
	defer it.Release()

	return_val := "<html>\n<body>\n"
	for it.Next() {
		s := string(it.Key())
		return_val = return_val + "<a href=\"?par=get&ver=" + s + "\">" + s + "</a><br>"
	}
	return_val = return_val + "</body>\n</html>"
	return return_val
}
