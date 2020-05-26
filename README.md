# tesseract-rpc

使用rpc方式对tesseract图片识别接口进行简单暴露，供其它服务调用，镜像内只下载了eng的训练数据，支持一般英文文本或数字验证码的识别。


tesseract预置的训练模型对验证码的识别准确率不高，可通过挂载目录`/usr/share/tesseract-ocr/4.00/tessdata`方式进行识别数据模型更换自己训练的模型，也可以直接修改Dockerfile文件。

> 感谢 `github.com/otiai10/gosseract` 对tesseract api的go 语言封装

## 使用方式

### server

可用环境变量：PORT端口、SERVER服务类型、TOKEN连接秘钥

- 方式一：直接启动(可选服务类型)
```bash
# 可选类型
docker run -it -p 8080:8080 -e SERVER=grpc zhenshaw/tesseract:rpc
docker run -it -p 8080:8080 -e SERVER=rpc zhenshaw/tesseract:rpc
docker run -it -p 8080:8080 -e SERVER=jsonrpc zhenshaw/tesseract:rpc
docker run -it -p 8080-8082:8080-8082 -e SERVER=all zhenshaw/tesseract:rpc
```

- 方式二：编译生成镜像并启动
```bash
# 克隆源码，进入目录
docker-compose up --build
```

### client

不同类型rpc的验证码识别调用示例

```go
package main

import (
	"context"
	"fmt"
	pb "github.com/ZhenShaw/tesseract-rpc/proto"
	"google.golang.org/grpc"
	"io/ioutil"
	"log"
	"net"
	"net/rpc"
	"net/rpc/jsonrpc"
)

func main() {
	path := "./0060.png" //验证码图片
	data, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}
	gRpcClient("localhost:8080", data)
	netRpcClient("localhost:8081", data)
	jsonRpcClient("localhost:8082", data)
}

func gRpcClient(addr string, reqData []byte) {

	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("dial fail: %v", err)
	}
	defer conn.Close()

	client := pb.NewCaptureOCRClient(conn)

	req := &pb.OCRRequest{
		Data:  reqData,
		Token: "",
	}
	r, err := client.Recognize(context.Background(), req)
	if err != nil {
		log.Fatalf("call client err:%s\n", err)
	}
	fmt.Printf("grpc reply: %s\n", r.Code)
}

type Req struct {
	Token string `json:"token"` //访问密码
	Data  []byte `json:"data"`
}

func netRpcClient(addr string, reqData []byte) {
	client, err := rpc.DialHTTP("tcp", addr)
	if err != nil {
		log.Fatalf("dial fail: %v", err)
	}
	defer client.Close()

	req := &Req{
		Token: "",
		Data:  reqData,
	}

	var reply string
	err = client.Call("RPCOcrService.Recognize", req, &reply)
	if err != nil {
		log.Fatalf("call client err:%s\n", err)
	}
	fmt.Printf("net/rpc reply: %s\n", reply)
}

func jsonRpcClient(addr string, reqData []byte) {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		log.Fatalf("dial fail: %v", err)
	}
	defer conn.Close()

	client := jsonrpc.NewClient(conn)

	req := &Req{
		Token: "",
		Data:  reqData,
	}
	var reply string
	err = client.Call("JSONRPCOcrService.Recognize", req, &reply)
	if err != nil {
		log.Fatalf("call client err:%s\n", err)
	}
	fmt.Printf("jsonrpc reply: %s\n", reply)
}


```