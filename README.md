# tesseract-ocr

使用grpc/http方式对tesseract图片识别接口进行暴露，镜像内只下载了eng的训练数据，支持一般英文文本或数字验证码的识别。


tesseract预置的训练模型对复杂验证码的识别准确率不高，可通过挂载目录`/usr/share/tesseract-ocr/4.00/tessdata`方式进行识别数据模型更换自己训练的模型，也可以直接修改Dockerfile文件。

> 感谢 `github.com/otiai10/gosseract` 对tesseract api的go 语言封装

## 使用方式

### server

可用环境变量：PORT端口、SERVER服务类型、TOKEN连接秘钥

- 方式一：直接启动
```bash
docker run -it -p 8080:8080 zhenshaw/tesseract:ocr
```

- 方式二：编译生成镜像并启动
```bash
# 克隆源码，进入目录
docker-compose up --build
```

### client

```go
package main

import (
	"bytes"
	"context"
	"fmt"
	pb "github.com/zhenshaw/tesseract-ocr/proto"
	"google.golang.org/grpc"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"time"
)

func main() {
	path := "./0060.png" //验证码图片
	data, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}
	gRpcClient("localhost:8080", data)
	httpClient("http://localhost:8080/ocr?token=", data)
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

func httpClient(addr string, reqData []byte) {

	body := new(bytes.Buffer)

	writer := multipart.NewWriter(body)
	formFile, err := writer.CreateFormFile("file", "pic.png")
	if err != nil {
		log.Fatal(err)
		return
	}
	_, _ = formFile.Write(reqData)
	err = writer.Close()
	if err != nil {
		log.Fatal(err)
		return
	}

	req, err := http.NewRequest("POST", addr, body)
	if err != nil {
		log.Fatal(err)
		return
	}
	//req.Header.Set("Content-Type","multipart/form-data")
	req.Header.Add("Content-Type", writer.FormDataContentType())

	HttpClient := &http.Client{
		Timeout: 3 * time.Second,
	}
	resp, err := HttpClient.Do(req)
	if err != nil {
		log.Fatal(err)
		return
	}
	defer resp.Body.Close()

	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
		return
	}
	fmt.Println("http reply:", string(content))
}

```