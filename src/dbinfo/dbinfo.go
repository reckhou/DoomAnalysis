package dbinfo

import (
  "database/sql"
  _ "github.com/go-sql-driver/mysql"
  "log"
  "os/exec"
  "strconv"
  "strings"
  "time"
)

type DumpMysql struct {
  db *sql.DB
}

var sql_instance *DumpMysql
var check_sql bool

/* 初始化数据库引擎 */
func Init() (*DumpMysql, error) {
  if sql_instance == nil {
    sql_instance = new(DumpMysql)
    db, err := sql.Open("mysql", "crash:crash2014@tcp(rdsiznueuzzvezb.mysql.rds.aliyuncs.com:3306)/crash?charset=utf8")

    if err != nil {
      log.Println("database initialize error : ", err.Error())
      return nil, err
    }
    sql_instance.db = db
    sql_instance.db.SetMaxIdleConns(10)
    check_sql = true
  }

  return sql_instance, nil
}

func Check_Sql_Connect() {
  for check_sql {
    test, _ := Init()
    if test.db != nil {
      test.db.Ping()
    }
    time.Sleep(1000 * 60 * time.Millisecond)
  }
}

/* 关闭数据库 */
func (test *DumpMysql) Close() {
  if test.db == nil {
    return
  }
  check_sql = false
  test.db.Close()
}

func (test *DumpMysql) CreateDB(pro string, ver string, address string, info string, uuid string, lianyun string) {

  if test.db == nil {
    return
  }

  if len(address) <= 0 {
    address = "No Address in libgame.so"
  }

  select_sql := "select count,ndk,filelist from " + pro + " where address ='" + address + "' and version =" + ver
  rows, err := test.db.Query(select_sql)
  if err != nil {
    log.Println(err.Error())
    return
  }
  defer rows.Close()

  var count_val int
  var ndk_val string
  var filelist_val string

  for rows.Next() {
    err := rows.Scan(&count_val, &ndk_val, &filelist_val)
    if err == nil {
      stmt, stmt_err := test.db.Prepare("update " + pro + " set count=?,filelist=? where address=? and version =?")
      defer stmt.Close()
      if stmt_err != nil {
        return
      }
      count_val++
      filelist_val = filelist_val + " " + uuid
      if result, err := stmt.Exec(count_val, filelist_val, address, ver); err == nil {
        if _, err := result.RowsAffected(); err == nil {
          return
        }
      }
    }
  }

  stmt, err := test.db.Prepare("insert into " + pro + "(address,version,count,ndk,filelist,lianyun)values(?,?,?,?,?,?)")
  if err != nil {
    log.Println("Prepare :", err.Error())
    return
  }
  defer stmt.Close()
  _, err = stmt.Exec([]byte(address), []byte(ver), 1, []byte(info), []byte(uuid), []byte(lianyun))
  if err != nil {
    log.Println(err.Error())
    return
  }

}

func GetListInfoDB(pro string, ver string) string {
  test, _ := Init()
  if test.db == nil {
    return ""
  }

  // 计算 count 总数
  dump_count := 0
  count_rows, err := test.db.Query("select sum(count),count(address) from " + pro + " where version =" + ver)
  defer count_rows.Close()
  if err != nil {
    log.Println(err.Error())
    return ""
  }
  dump_type_count := 0
  for count_rows.Next() {
    err := count_rows.Scan(&dump_count, &dump_type_count)
    if err != nil {
      log.Println(err.Error())
      return ""
    }
  }

  if dump_count <= 0 {
    return ""
  }

  // 排序输出
  select_sql := "select id,address,count,ndk,filelist from " + pro + " where version =" + ver + " order by count desc"
  select_rows, err := test.db.Query(select_sql)
  defer select_rows.Close()
  if err != nil {
    log.Println(err.Error())
    return ""
  }

  return_val := "<html>\n<body>\n<table border=\"1\">\n"
  return_val = return_val + CheckFreedisk()
  return_val = return_val + "<tr>\n"
  return_val = return_val + "<th align=\"left\">ID (" + strconv.Itoa(dump_type_count) + ") </th>\n"
  return_val = return_val + "<th align=\"right\">ADDRESS</th>\n"
  return_val = return_val + "<th align=\"right\">COUNT (" + strconv.Itoa(dump_count) + ") </th>\n"
  return_val = return_val + "<th align=\"center\">INFO</th>\n"
  return_val = return_val + "<th align=\"center\">BREAKPAD</th>\n"
  return_val = return_val + "</tr>\n"

  var id int
  var address string
  var count_val int
  var ndk_val string
  var filelist_val string

  index_val := 1
  for select_rows.Next() {
    if err := select_rows.Scan(&id, &address, &count_val, &ndk_val, &filelist_val); err == nil {
      id_val := strconv.Itoa(id)
      count_id_val := strconv.Itoa(count_val)
      percent := float64(count_val) * 100.0 / float64(dump_count)

      uuid := strings.Split(filelist_val, " ")[0]
      info_val := "<a href=\"/file/" + pro + "/" + ver + "/" + uuid + ".tar.bz2 \">" + uuid + ".tar.bz2" + "</a><br>"
      info_val = info_val + "<a href=\"?pat=detail&ver=" + ver + "&pro=" + pro + "&id=" + id_val + " \">" + "more..." + "</a><br>"

      color := ""
      if percent > 0.5 {
        color = " style=\"color:#F00\" "
      }

      return_val = return_val + "<tr>\n"
      return_val = return_val + "<th align=\"left\"><div " + color + ">" + id_val + "</div></th>\n"
      return_val = return_val + "<th align=\"right\"><div " + color + ">" + address + "</div></th>\n"
      return_val = return_val + "<th align=\"right\"><div " + color + ">" + count_id_val + " (" + strconv.FormatFloat(percent, 'f', 2, 64) + "%%)" + " </div></th>\n"
      return_val = return_val + "<th align=\"left\"><div " + color + ">" + ndk_val + " </div></th>\n"
      return_val = return_val + "<th align=\"left\">" + info_val + " </th>\n"
      return_val = return_val + "</tr>\n"
      index_val++
    }
  }

  return_val = return_val + "</table>\n</body>\n</html>"
  return return_val
}

func GetFileListInfoDB(pro string, ver string, id string) string {

  test, _ := Init()
  if test.db == nil {
    return ""
  }

  // 输出
  select_sql := "select filelist from " + pro + " where id =" + id
  select_rows, err := test.db.Query(select_sql)
  defer select_rows.Close()
  if err != nil {
    log.Println(err.Error())
    return ""
  }

  var filelist_val string
  info_val := "<html>\n<body>\n"
  for select_rows.Next() {
    if err := select_rows.Scan(&filelist_val); err == nil {

      uuid := strings.Split(filelist_val, " ")

      for _, v := range uuid {
        info_val = info_val + "<a href=\"/file/" + pro + "/" + ver + "/" + v + ".tar.bz2 \">" + v + ".tar.bz2" + "</a><br>"
      }

    }
  }
  info_val = info_val + "</body>\n</html>"
  return info_val
}

func DeleteInfoDB(pro string, ver string) {
  test, _ := Init()
  if test.db == nil {
    return
  }

  stmt, err := test.db.Prepare("delete from " + pro + " where version=?")
  if err != nil {
    log.Println(err.Error())
    return
  }
  defer stmt.Close()
  if result, err := stmt.Exec(ver); err == nil {
    if c, err := result.RowsAffected(); err == nil {
      log.Println("remove count : ", c)
    }
  }
}

func VerInfoDB(pro string) string {
  test, _ := Init()
  if test.db == nil {
    return ""
  }

  rows, err := test.db.Query("SELECT DISTINCT version FROM " + pro)
  defer rows.Close()
  if err != nil {
    log.Println(err.Error())
    return ""
  }

  var id int

  return_val := "<html>\n<body>\n"
  return_val += CheckFreedisk()
  for rows.Next() {
    if err := rows.Scan(&id); err == nil {
      str1 := strconv.Itoa(id)
      return_val = return_val + "<a href=\"?pat=get&pro=" + pro + "&ver=" + str1 + "\">" + str1 + "</a><br>"
    }
  }

  return_val = return_val + "</body>\n</html>"
  return return_val
}

func CheckFreedisk() string {

  cmd := exec.Command("df", "-hk")
  b, err := cmd.Output()
  if err != nil {
    log.Println("err:" + err.Error())
    return ""
  }
  cmd_result := string(b)
  arr := strings.Split(cmd_result, "\n")
  title := strings.Split(arr[0], " ")

  free_index := 0
  for i := 0; i < len(title); i++ {
    if title[i] != "" {
      if title[i] == "可用" || title[i] == "Available" {
        break
      }
      free_index++
    }
  }

  free_disk_count := 0
  free_check_index := 0
  for i := 1; i < len(arr); i++ {
    content := strings.Split(arr[i], " ")
    free_check_index = 0
    for j := 0; j < len(content); j++ {
      if content[j] == "tmpfs" {
        break
      }
      if content[j] != "" {
        if free_index == free_check_index {
          count, _ := strconv.Atoi(content[j])
          free_disk_count += count
        }
        free_check_index++
      }
    }
  }
  result := "<div "
  if free_disk_count < 1024*1024*2 {
    result = result + " style=\"color:#F00\" "
  }
  result += ">"
  result += "磁盘空间剩余: "
  result += strconv.Itoa(free_disk_count) + "k ("
  result += strconv.Itoa(free_disk_count/1024/1024) + "G)</div><br>"
  return result
}
