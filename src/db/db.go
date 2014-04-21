package db

import "C"

import (
  "github.com/syndtr/goleveldb/leveldb"
  "log"
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
  log.Println("ver: ", ver)
  version, _ := db_version.Get([]byte(ver), nil)
  if version == nil {
    db_version.Put([]byte(ver), []byte("0"), nil)
  }

  db_id, id_err := leveldb.OpenFile("./"+pro+"/dump/"+ver+"/db_id.db", nil)
  if id_err != nil {
    log.Println("db_id OpenFile err: ", id_err)
    <-db_channel
    return
  }
  db_info, info_err := leveldb.OpenFile("./"+pro+"/dump/"+ver+"/db_info.db", nil)
  if info_err != nil {
    log.Println("db_info OpenFile err: ", info_err)
    <-db_channel
    return
  }
  db_breakpad_info, breakpad_info_err := leveldb.OpenFile("./"+pro+"/dump/"+ver+"/db_breakpad_info.db", nil)
  if breakpad_info_err != nil {
    log.Println("db_breakpad_info OpenFile err: ", breakpad_info_err)
    <-db_channel
    return
  }
  db_count, db_count_err := leveldb.OpenFile("./"+pro+"/dump/"+ver+"/db_count.db", nil)
  if db_count_err != nil {
    log.Println("db_count OpenFile err: ", db_count_err)
    <-db_channel
    return
  }
  defer db_id.Close()
  defer db_info.Close()
  defer db_count.Close()
  defer db_breakpad_info.Close()

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
    log.Println("db_breakpad_info Put:", address)
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
    log.Println("db_info Put:", address)
    db_info.Put([]byte(address), []byte(info), nil)
  }

  s, _ = db_count.Get([]byte(address), nil)
  if s == nil {
    log.Println("db_info Put:", address)
    db_count.Put([]byte(address), []byte("1"), nil)
  } else {
    num_val, err := strconv.Atoi(string(s))
    if err != nil {
      log.Println("db_count error:", ver)
      return
    }
    num_val++
    address_num_val := strconv.Itoa(num_val)
    db_count.Put([]byte(address), []byte(address_num_val), nil)
  }

  <-db_channel
}

func GetListInfoDB(pro string, ver string) string {

  db_id, _ := leveldb.OpenFile("./"+pro+"/dump/"+ver+"/db_id.db", nil)
  db_info, _ := leveldb.OpenFile("./"+pro+"/dump/"+ver+"/db_info.db", nil)
  db_count, _ := leveldb.OpenFile("./"+pro+"/dump/"+ver+"/db_count.db", nil)
  db_breakpad_info, _ := leveldb.OpenFile("./"+pro+"/dump/"+ver+"/db_breakpad_info.db", nil)
  defer db_id.Close()
  defer db_info.Close()
  defer db_count.Close()
  defer db_breakpad_info.Close()

  return_val := "<html>\n<body>\n<table border=\"1\">\n"
  return_val = return_val + "<tr>\n"
  return_val = return_val + "<th align=\"left\">ID</th>\n"
  return_val = return_val + "<th align=\"right\">ADDRESS</th>\n"
  return_val = return_val + "<th align=\"right\">COUNT</th>\n"
  return_val = return_val + "<th align=\"center\">INFO</th>\n"
  return_val = return_val + "<th align=\"center\">BREAKPAD</th>\n"
  return_val = return_val + "</tr>\n"
  max_num_val, _ := db_id.Get([]byte("MAX_NUM"), nil)
  if max_num_val != nil {
    num_val, err := strconv.Atoi(string(max_num_val))
    if err != nil {
      log.Println("MAX_NUM error:", ver)
      return ""
    }
    max_num := num_val
    log.Println("MAX_NUM:", max_num)
    for i := 1; i <= max_num; i++ {

      id_val := strconv.Itoa(i)
      address, _ := db_id.Get([]byte(id_val), nil)

      if address != nil {

        info, _ := db_info.Get([]byte(address), nil)
        count, _ := db_count.Get([]byte(address), nil)
        info_pad, _ := db_breakpad_info.Get([]byte(address), nil)
        return_val = return_val + "<tr>\n"
        return_val = return_val + "<th align=\"left\">" + id_val + "</th>\n"
        return_val = return_val + "<th align=\"right\">" + string(address) + "</th>\n"
        return_val = return_val + "<th align=\"right\">" + string(count) + " </th>\n"
        return_val = return_val + "<th align=\"left\">" + string(info) + " </th>\n"
        return_val = return_val + "<th align=\"left\">" + string(info_pad) + " </th>\n"
        return_val = return_val + "</tr>\n"
      }
    }
  }
  return_val = return_val + "</table>\n</body>\n</html>"
  return return_val
}

func VerInfoDB(pro string) string {

  db_version, _ := leveldb.OpenFile("./"+pro+"/version_id.db", nil)
  defer db_version.Close()
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
