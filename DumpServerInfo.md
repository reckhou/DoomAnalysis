

#DUMP 服务器
- - - 
####URL
查看DUMP 参数: pat=模式 pro=项目名
http://115.28.202.194:10010/?pat=get&pro=sxda 
重建特定版本DUMP信息 参数: pat=模式 pro=项目名 ver=版本号 lianyun=联运平台
http://115.28.202.194:10010/?pat=recreate&pro=sxda&ver=5572&lianyun=sxda
- - -
####目录结构
```
root -- 根目录
     |--dumpserver -- 服务器根目录
                   |-doomAnalysis      -- 服务器执行程序
                   |-gen_dump_info.sh  -- 生成breakpad dump脚本
                   |-gen_ndk_info.sh   -- 生成ndk dump脚本
                   |-gensym.sh         -- 生成中间调试信息脚本
                   |-tools --          -- 第三方工具
                   |       |-dump_syms          -- 生成中间调试信息程序
                   |       |-minidump_stackwalk -- 生成breakpad dump程序
                   |       |-ndk-stack          -- 生成ndk dump程序
                   |
                   |-sxda ---           -- 项目文件夹
                           |-lib --       -- 存放对应版本的so
                           |     |-5146_libgame.so
                           |      ....
                           |     |-5526_libgame.so
                           |-dump--   -- 存放玩家上传的dump
                                 |-5416--     -- 根据版本号分别存放dump
                                 ....
                                 |-5526--
```
- - -
####数据库
服务器数据库 用户名:root 密码: pindump123 数据库sxddump

CREATE TABLE sxd
(
id int NOT NULL auto_increment,  
address varchar(40) NOT NULL,   // 取崩溃dump前3个地址作为同一DUMP对比
version varchar(10) NOT NULL,   // 版本
count int NOT NULL,             // 数量
ndk text NOT NULL,              // 简易堆栈信息
filelist text NOT NULL,         // 对应的dump文本 uuid
lianyun varchar(20) NOT NULL default 'sxda', // 联运平台标记
PRIMARY KEY (id)
);
- - -
####工作流程

1. 客户端通过url post log和minidump文件到服务器 (文件名为UUID),检测合法性
2. 是否已经生成对应的调试信息文件
    + 查找是否在dumpserver/项目/lib/ 下已经生成对应版本号 xxxx.txt 存在 txt 说明已经生成好
    不存在则调用dump_syms生成在dumpserver/项目/dump/版本号/symbols/下 ,并创建xxxx.txt标识已经生成好


3. 依此在 dumpserver/项目/dump/版本号/ 下生成
UUID.txt            -- 原始dump文件
UUID.txt.info  -- 解析出的breakpad信息文件
UUID.txt.ndk -- ndk dump文件
UUID.txt.ndk.info -- 解析出的ndk信息文件

4. 通过解析堆栈信息,插入到数据库中
5. 打包dump文件







