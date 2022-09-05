##  1.protobuf优点
1.  序列化后体积相比xml和json更小，适合网络传输
2.  支持跨平台多语言
3.  消息格式升级和兼容性还不错
4.  序列化和饭序列化速度更快，快于json的处理速度

##  2.protobuf下载
1.  wget https://github.com/protocolbuffers/protobuf/releases/download/v3.9.1/protoc-3.9.1-linux-x86_64.zip
2.  mkdir protoc
3.  cd protoc
4.  unzip protoc-3.9.1-linux-x86_64.zip
5.  cd bin
6.  cp protoc /usr/bin
7.  protoc --version

##  3.生成go代码
-   protoc --go_out=plugins=grpc:OUT_FILEPATH PROTO_FILE

##  4.空值和基本类型的传参
-   RPC接口的参数和返回值必须是message类型。即使接口不需要参数或者返回值，也得指定一个message。这也导致不能不传参数，也不能不返回结果，也不能使用基本类型操作参数和返回结果(因为基本类型不是message类型)，而项目中经常需要这样做。如何解决？方法如下：
    1.  自定义message接口：
        ```
        message Int32 {
        int32 value = 1;
        }

        message Bool {
            bool value = 1;
        }
        ```
    2.  使用proto封装好的message，导入google封装好的wrappers.proto包：
        ```
        import "google/protobuf/wrappers.proto";

        service UserSevice {
            rpc getById(google.protobuf.Int32Value) returns (User);
        }
        ```
    3.  传入和返回空值，导入google封装好的empty.proto包：
        ```
        import "google/protobuf/empty.proto";
        service UserSevice {
            rpc list(google.protobuf.Empty) returns (google.protobuf.Empty);
        }
        ```
    