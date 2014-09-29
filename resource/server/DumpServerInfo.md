

#DUMP 服务器使用,搭建
- - - 
##使用
###URL查看接口
访问相应地url
查看DUMP 参数: pat=模式 pro=项目名
( http://ip:port/?pat=get&pro=项目名 ) 
例如:
http://127.0.0.1:8888/?pat=get&pro=proname       // c++版查看
http://127.0.0.1:8888/?pat=get&pro=proname_java  // java版查看
http://127.0.0.1:8888/?pat=get&pro=proname_js    // javascript版查看

有些情况下需要重建某一个版本的dump信息的
重建特定版本DUMP信息 参数: pat=模式 pro=项目名 ver=版本号 lianyun=联运平台
( http://ip:port/?pat=recreate&pro=项目名&ver=版本号&lianyun=联运平台 ) 
例如:
http://127.0.0.1:8888/?pat=recreate&pro=proname&ver=1234&lianyun=lianyunname

针对tencent dump的特殊处理
首先将dump的zip包上传倒 服务器根目录(root)/${proname}/tencentdump/ 下
然后通过访问
http://127.0.0.1:8888/tencent 正确填写项目名和 版本号
后台会自动读取  tencentdump 目录下所有的zip包 自动生成

- - -
##服务器搭建
###1.配置目录
首先建立服务器根目录(root)
将 (doomAnalysis gen_dump_info.sh gen_ndk_info.sh gensym.sh createproject.sh) 5个文件复制到根目录下
根目录下 建立 tools 文件夹, 将 (dump_syms minidump_stackwalk ndk-stack) 3个文件复制到根目录下

根目录下建立 config 文件夹, 在此文件夹下 建立config.json 配置文件
config.json
```
{
  "basic": {
    "Port": "8888",      // 本地服务器端口
    "Host": "127.0.0.1"  // 本地服务器IP
  },
  "mysql":{
    "Port": "3306",     // 数据库端口
    "Host": "127.0.0.1", // 数据库IP
    "User": "123456",    // 数据库用户名
    "PassWord": "123456",// 用户密码
    "DataBase": "crash"  // 要连接的数据库
  },
  "project":{           // 项目DUMP归纳设置
      /** 默认项目会统计到和本身的项目名字一致的目录和数据库下
       对于需要特殊合并统计的在这里设置:
       */
      "proname":{          // 默认属于 proname 这个项目所有的联运平台版本都会统计到 proname下面
      "proname_tr":"proname_tr", // 例外1 对于 属于 proname 这个项目proname_tr联运平台版本 会统计倒 proname_tr 下面
      "proname_kr":"proname_kr" // 例外2 对于 属于 proname 这个项目proname_kr联运平台版本 会统计倒 proname_kr 下面
    }
  }
}
```

运行 createproject.sh [项目名] [option:a] 建立项目文件夹
例如
```
./createproject.sh proname    // 建立 proname c++ 一个项目文件夹
./createproject.sh proname a    // 建立 proname c++,java,javascript三个项目文件夹
```
#### 注意:
#### 1.项目名,查看接口的项目名,config.json配置的统计输出 三个命名要一致 
#### 3.同一项目下项目的 java,javascript Dump查看是通过 ( 项目名_java,项目名_js) 当做不同的项目来区分对待



####目录结构
```
--root -- 服务器根目录
       |-doomAnalysis      -- 服务器执行程序
       |-gen_dump_info.sh  -- 生成breakpad dump脚本
       |-gen_ndk_info.sh   -- 生成ndk dump脚本
       |-gensym.sh         -- 生成中间调试信息脚本
       |-tools --          -- 第三方工具
       |       |-dump_syms          -- 生成中间调试信息程序
       |       |-minidump_stackwalk -- 生成breakpad dump程序
       |       |-ndk-stack          -- 生成ndk dump程序
       |
       |-sxda ---         -- 项目文件夹
               |-lib --   -- 存放对应版本的so
               |-dump--   -- 存放玩家上传的dump
```
- - -
###配置数据库
在配置文件config.json配置的数据库下面建立相应项目的数据库表 

// 统计dump信息
CREATE TABLE proname
(
id int NOT NULL auto_increment,  
address varchar(40) NOT NULL,   // 取崩溃dump前3个地址作为同一DUMP对比
version varchar(10) NOT NULL,   // 版本
count int NOT NULL,             // 数量
ndk text NOT NULL,              // 简易堆栈信息
filelist text NOT NULL,         // 对应的dump文本 uuid
lianyun varchar(20) NOT NULL default 'proname', // 联运平台标记
PRIMARY KEY (id)
);

// 统计设备信息(项目名+_device)
CREATE TABLE proname_device
(
id int NOT NULL auto_increment,  
address varchar(40) NOT NULL,   // 取崩溃dump前3个地址作为同一DUMP对比
version varchar(10) NOT NULL,   // 版本
device varchar(40) NOT NULL,    // 设备
lianyun varchar(20) NOT NULL default 'proname', // 联运平台标记
PRIMARY KEY (id)
);

#### 注意:
#### 1.数据库表名要和 项目名一致
#### 2.同一项目下的 java,javascript dump 需要分别建立 ( 项目名_java,项目名_js) 新表进行统计
- - -
### c++ dump分析数据的上传
对于分析c++的dump需要对应版本的动态链接库进行分析, 对于需要分析的版本,请将版本的动态链接库 上传到
root/proname/lib/ 下面 并且重命名为 版本号_libgame.so
- - -
##搭建完成
- - -
##工作流程

1. 客户端通过url post log和minidump文件到服务器 (文件名为UUID),服务器首先检测合法性,不合法的文件则丢弃
2. 对于c++ dump首先检测是否已经生成对应的调试信息文件:
    + 查找是否在root/项目/lib/ 下已经生成对应版本号 xxxx.txt 存在 txt 说明已经生成好
    不存在则调用dump_syms生成在root/项目/dump/版本号/symbols/下 ,并创建xxxx.txt标识已经生成好


3. 依此在 dumpserver/项目/dump/版本号/ 下生成
UUID.txt            -- 原始dump文件
UUID.txt.info  -- 解析出的breakpad信息文件
UUID.txt.ndk -- ndk dump文件
UUID.txt.ndk.info -- 解析出的ndk信息文件

4. 通过解析堆栈信息,插入到数据库中
5. 打包dump文件

- - -
###API列表
####doomAnalysis
描述: http web服务启动设置,url解析分发
```
Start: 入口函数,定义服务器配置
```
getUrlParameter: 获取URL参数值
```
in:
    [key string] 需要获取的URL参数key
    [form url.Values] URL参数对象
out:
    [string] 返回key所对应的值,出错则返回空字符串
    [bool] 返回处理结果 正确返回true 错误返回false
```
ServeHTTP: 处理URL
```
get:
    pro - 项目名称 
    ver - 版本号
    lianyun - 联运平台标识
    http://ip:port/file/${pro}/${ver}/${filename} //请求文件 
    http://ip:port/?pat=get&pro=${pro}  // 获取对应项目的版本列表
    http://ip:port/?pat=get&pro=${pro}&ver=${ver}  // 获取对应项目的版本列表
    http://ip:port/?pat=recreate&pro=${pro}&lianyun=${lianyun}&ver=${ver}  // 重建dump信息
post:
    http://ip:port/?pat=post&pro=${pro}&lianyun=${lianyun}  // 上传的DUMP
    post格式
    {
        MG5:XXXXXXXX
        UUID:XXXXXXXX
        device:xxxx
        version:xxxxxx
        product_name:xxxxxx
        symbol_file:(DUMP二进制文件)
    }
```
CheckLegal: 检验上传的DUMP文件是否合法
```
in:
    [context string] 要检验的内容
out:
    [bool] true- 合法 flase- 非法
```
GetProName: 通过配置获取实际DUMP对应的目录
```
in:
    [pro string] 项目名
    [lianyun strubg] 联运平台名
out:
    [string] 对应的目录名字
```
- - -
####dumpfile
DumpFileInfo : dump文件结构 
```
var key_arr [6]string = [...]string{"MD5", "UUID", "device", "version", "product_name", "file"}
type DumpFileInfo struct {
  info_          map[string]string - 存储文件基本信息(对应上面的key_arr)
  file_name_     string - 文件名 UUID
  stack_lib_name []string - 堆栈对应的动态链接库名称
  stack_address  []int64 - 堆栈对应的动态链接库地址
  block_in       bool   - 解析标识
  so_address     int64  - 
  project        string - 项目名
  ndk_stack_info string - 存储生成的 ndk 堆栈信息(文本信息)
  lianyun        string - 联运平台名
}
```

DumpFileInfo.GenInfo
描述:根据上传的dump文件填充 DumpFileInfo.info_ 数据
```
in : 上传的DUMP文件
```
DumpFileInfo.GenLogInfo
描述:将上传的log文件保存下来
```
in : 上传的log文件
```
DumpFileInfo.GenSym
描述:调用SH脚本gensym.sh生成so的sym文件

DumpFileInfo.GenBreakpadDumpInfo
描述:调用SH脚本gen_dump_info.sh生成BreakpadDump解析后的文件

DumpFileInfo.GenNdkDumpInfo
描述:通过解析BreakpadDump文件 生成 NDK dump信息

DumpFileInfo.GenNdkStack
描述:生成 NDK 堆栈信息

DumpFileInfo.GenNdkSoAddress
描述:生成NDK堆栈地址

DumpFileInfo.GenNdkfile
描述:生成 NDK dump文件并解析

DumpFileInfo.GenDbInfo
描述:将数据记录到数据库

DumpFileInfo.GenTar
描述:打包/解包DUMP文件
```
in:
    "c" 打包数据
    "x" 解包数据
```

ProcessDumpFile
描述: 开始dump处理流程入口函数
```
in:
    [project string] 项目名称
    [co []byte] 上传的DUMP信息
    [lianyun string] 联运平台
```

ListFileName
描述: 重建某一版本的所有DUMP文件
```
in:
    [path  string] 项目名称
    [project string] 项目名称
    [ver string] 版本号
    [lianyun string] 联运平台
```

---
####dbinfo
描述:数据库交互类
DB链接包装
```
type DumpMysql struct {
  db *sql.DB
}
```
Init
描述:初始化数据库链接
```
out:返回DumpMysql单件实例
```
Check_Sql_Connect
描述:每1分钟循环ping数据库,保持数据库链接

Close
描述:关闭数据库链接

CreateDB
描述:根据dump信息更新数据库
```
in:
    [pro string] 项目名称
    [ver string] 版本号
    [address string] 堆栈地址
    [info string] 堆栈信息
    [uuid string] UUID
    [lianyun string] 联运平台
```

GetListInfoDB
描述:显示指定版本的DUMP信息
```
in:
    [pro string] 项目名称
    [ver string] 版本号
```

DeleteInfoDB
描述:删除特定版本DUMP信息
```
in:
    [pro string] 项目名称
    [ver string] 版本号
```

VerInfoDB
描述:获取特定项目DUMP版本列表
```
in:
    [pro string] 项目名称
```

CheckFreedisk
描述:检测剩余磁盘空间