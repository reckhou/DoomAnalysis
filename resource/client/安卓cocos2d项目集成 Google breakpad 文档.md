
#安卓cocos2d项目集成 Google breakpad 文档

- - -
####简明步骤
* 集成 breakpad 的 libbreakpad_client.a 到项目里
* 修改代码初始化 breakpad_client
* 添加java文件 让native崩溃的process调用java进行发送服务端处理
* 服务端收集信息转成 minidump 文件
* 使用 breakpad 的 **dump_syms** 解析响应崩溃版本的 so 二进制文件 生成调试信息
* 使用 breakpad 的 **minidump_stackwalk** 通过 上一步生成的调试信息和 minidump 文件 定位出错代码位置

- - -
####具体集成方法

 1. 简单原理描述
    linux 平台程序崩溃时候会触发一个signal,googole的机制是在触发signal后另外启动一个process,将崩溃的process复制出来,输出dump信息.
	具体实现技术文档可以阅读google官方文档

 2. 事先准备
   * svn拉取 breakpad 最新版本
   * 集成到安卓项目客户端生成MINIDUMP文件需要  libbreakpad_client.a 文件
      + 如何编译安卓的lib文件
	          我们需要在安卓下的lib文件,就需要用NDK的库来编译它,这里使用
	          https://vilimpoc.org/blog/2013/10/08/ndk-configure-sh/
	          这篇文章的脚本 来配置 NDK 进行编译
	          - 将脚本复制到 breakpad 下的根目录(和 configure 同级)
	          - 编辑脚本,配置NDK路径
	          	 eg. 
	          	 将 NDK="/NDK绝对路径/BuildEnv/android-ndk-r7c" 配置在第8行
	          	 需要将 26-32行删除(没用)
	          	 34行 配置 COMPILER_VERSION 为 4.4.3 (NDK对应的gcc,大家可以在NDK/toolchains找到对应的版本号)
	          - 由于我们这里使用的是老版本NDK r7c 编译一些路径和新版本不匹配,所以下面是r7c版本特别设置
	          	 - 50行 去除 COMPILER_VERSION 
	             - 运行 ndk-configure.sh 假如当中还有路径错误,请大家根据NDK实际情况配置路径
	             - 运行 make 命令进行编译. (老版本NDK在编译到后面可能遇到错误,不过我们只需要 libbreakpad_client.a,只要在{$breakpadrootpath}/src/client/linux/ 下有编译 libbreakpad_client.a 即可 )  
- - -

 3. 集成到客户端
	参照 https://github.com/pixonic/google-breakpad 俄罗斯人的工程进行修改
	* 在需要编译的 so makefile文件下 加入 
			> 引用文件:
			```
            LOCAL_C_INCLUDES += $(LOCAL_PATH)/../third-party/breakpad/google-breakpad \
						$(LOCAL_PATH)/../third-party/breakpad/google-breakpad/src \
						$(LOCAL_PATH)/../third-party/breakpad/google-breakpad/src/common/android/include
            ```
			库文件
			   ```LOCAL_STATIC_LIBRARIES += breakpad_client```
            >**注意: 由于NDK对C++标准库只有有限的支持,可能在引用后遇到链接错误,这里需要把 breakpad_client 的链接放到所有库的最前面编译**

	* 代码调用
        + c++ : cocos2d项目的 main 文件
		    - 添加包含文件
			    ```#include "breakpad路径/src/client/linux/handler/exception_handler.h"```
				```#include "breakpad路径/src/client/linux/handler/minidump_descriptor.h"```
				
			- 声明全局变量
				```JavaVM *gJavaVMInstance = NULL;```
						   
			- 添加函数
				```
                void crashHandlerSetJavaVM(JavaVM *javaVM) {
			        gJavaVMInstance = javaVM;
				}
                ```
                ```
				JNIEnv* getJNIEnv() {
				    JNIEnv* env = NULL;
					if(gJavaVMInstance) {
					    if(gJavaVMInstance->GetEnv((void**)&env, JNI_VERSION_1_4) == JNI_EDETACHED) {
						     gJavaVMInstance->AttachCurrentThread(&env, NULL);
					    }
					}
				    else {
					    LOGD("getJNIEnv: invalid java vm");
				    }
				    return env;
				}
				```
				```
                bool crashHandlerDumpCallback(const google_breakpad::MinidumpDescriptor& descriptor, void* context, bool succeeded) {
                    LOGD("crashHandlerDumpCallback: application crashed");
    		        std::string dumpPath(descriptor.path());
                    std::string dumpFile = "crash.dmp";
				    size_t found = dumpPath.find_last_of("/\\");
				    if(found != std::string::npos) {
					    dumpFile = dumpPath.substr(found + 1);
					}
    				JNIEnv *env = getJNIEnv();
				    if(env) {
				        jclass classID = env->FindClass("com/pixonic/breakpadintergation/CrashHandler");
					    if(classID) {
						    jmethodID methodID = env->GetStaticMethodID(classID, "nativeCrashed", "(Ljava/lang/String;)V");
						    if(methodID) {
						        // create the first parameter for java method
						        jstring firstParameter = env->NewStringUTF(dumpFile.c_str());

						        env->CallStaticVoidMethod(classID, methodID, firstParameter);

						        // delete local references
						        env->DeleteLocalRef(firstParameter);
						    }
						    else {
						        LOGD("crashHandlerDumpCallback: java method not found");
						    }

						    env->DeleteLocalRef(classID);
					    }
					    else {
					        LOGD("crashHandlerDumpCallback: java class not found");
				        }
				    }
					else {
				        LOGD("crashHandlerDumpCallback: invalid Java environment");
					}
						
					// remove local file
					::remove(dumpPath.c_str());
						
					return succeeded;
				}
				```	
				```		    
				void setupCrashHandler(const std::string &path) {
				    LOGD("setupCrashHandler");
						
    			    std::string writeablePath = path;
					size_t found = writeablePath.find_last_of("/\\");
					if(found == writeablePath.length() - 1) {
					    writeablePath = writeablePath.substr(0, writeablePath.length() - 1);
					}
					    google_breakpad::MinidumpDescriptor dumpDescriptor(writeablePath);
					    static google_breakpad::ExceptionHandler exceptionHandler(dumpDescriptor, NULL, crashHandlerDumpCallback, NULL, true, -1);
				}
				```	
				```
				// 这个函数放在 extern "C" 作用域里
				void Java_com_pixonic_breakpadintergation_CrashHandler_nativeInit(JNIEnv *env, jobject self, jstring path) {
			        jboolean isCopy;
			        const char* chars = env->GetStringUTFChars(path, &isCopy);
				    string pathStr(chars);
				    if(isCopy) {
				        env->ReleaseStringUTFChars(path, chars);
				    }
				    setupCrashHandler(pathStr);
				}
				```
										
			-  修改函数
				```jint JNI_OnLoad(JavaVM *vm, void *reserved)```
				添加
				```crashHandlerSetJavaVM(vm);```
						   
	            **注意: 客户端编译时需添加 -g 参数,生成调试信息**

        + JAVA : 将 CrashHandler.java MultipartHttpEntity.java PIDefaultExceptionHandler.java 到cocos2d安卓项目中去 在项目入口的 Activity 的 onCreate 函数第一行调用 CrashHandler Init 方法
        CrashHandler 的构造函数里要去设置 mSubmitUrl 变量设置要上传的url
        url如下:
			```
			  要上传的url = http://ip:port/?pat=post&pro=项目名&lianyun=
			```
			
			```
			  CrashHandler.init(); // 注册 native 和 js crash
			  PIDefaultExceptionHandler defaultExceptionHandler = new PIDefaultExceptionHandler();
    defaultExceptionHandler.init(this); // 注册JAVA crash
            ```
 4. 发送到服务端数据格式
	原理:使用 java 的 HttpEntity post 到服务端
    在 CrashHandler.java 里
```
private void saveFile(String dumpFile); // 将生成的dump和日志转存到一个固定文件下,便于上传,客户端产生dump时会自动调用
public void UploadDumpFile(); // 将dump文件上传, 需要手动调用
```
需要手工调用
CrashHandler.getInstance().UploadDumpFile();
来上传文件


数据格式

第一行根据上传文件的类型选择不同的参数名称(大小写敏感)

支持的文件类型：

    MD5: C++部分crash dump
    java: Java部分crash dump
    js: js部分异常
    LOG: log信息，不做处理
    
其值为第2-4行内容首位相接后取MD5值，包括换行符。

```
[MD5|LOG|java|js]:(计算第二行 UUID 到第五行 product_name: 内容的md5值 (包含回车))
UUID:....\n
device:....\n
version:....\n
product_name:....\n
file:(dump的具体信息)
```
```
例如上传java dump信息
java:<MD5 Vale>\n
UUID:....\n
device:....\n
version:....\n
product_name:....\n
file:(dump的具体信息)
```
- - -
###配置结束
- - -
###分析步骤原理(服务端已经实现,不用配置)
 1. 解析dump数据
  参考 http://blog.csdn.net/brook0344/article/details/20126351
    * 准备文件
		-  需要linux版下breakpad编译出的 dump_syms 和 minidump_stackwalk
					      具体编译方法: linux环境下 获取breakpad源代码后,在根目录下执行 ./configure之后 执行 make
		-  需要dump相对应编译版本的 so 或者 可执行 文件, cocos2d则为 so 文件
		- 使用 dump_syms 解析 so 文件 生成 库调试信息
			eg. (假设为 libgame.so 文件)
			>**./dump_syms libgame.so > libgame.so.sym**

			创建文件夹 symbols/libgame.so/2D1C163A1347A1190B26F10560E926CE0
			后面那个一堆乱数字是前一步生成的“libgame.so.sym”文件的第一行复制出来的
			
			然后将生成的 libgame.so.sym 放入新文件夹中，最终它的路径：
			symbols/libgame.so/2D1C163A1347A1190B26F10560E926CE0/libgame.so.sym
						
			**这里由于需要手工操作,所以修改了 breakpad 源代码,使得 dump_syms 能够自动生成移动文件到新文件夹.**
			修改源代码如下:
						
			```
			// dump_symbols.cc
			// WriteSymbolFile 方法 在 delete module; 前添加如下代码`
			//自动生成文件夹,将解析的sym移到新文件夹下
			char result[ 2048 ];
			string real_path = std::string( getcwd( result, sizeof( result ) ) );
			string real_copypath = "";
			string copypath = "";
			
			mkdir("symbols", 0755);
			if(  chdir("symbols") !=-1 ) {  
			    mkdir(module->GetName().c_str(), 0755);
			    if( chdir(module->GetName().c_str()) !=-1) {  
			        mkdir(module->GetId().c_str(),   0755);
			        if( chdir(module->GetId().c_str()) !=-1 ) {  
			            copypath = "symbols/"+module->GetName()+"/"+module->GetId();
			            real_copypath = std::string( getcwd( result, sizeof( result ) ) );
			        }
			    } 
			}
			  
			chdir(real_path.c_str());
			  
			string copyfilename = obj_file + ".sym";
			copypath = copypath +"/"+ copyfilename;
			rename( copyfilename.c_str(),copypath.c_str() );
			```		  
		- 导出文本信息
			返回到第三步建立的那个文件夹的根部，执行
			>**./minidump_stackwalk XXX.dmp symbols/ > XXX.txt**

            查看新文件 XXX.txt 里面就有崩溃文本信
