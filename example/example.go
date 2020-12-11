/**
 * @File: example
 * @Author: Shaw
 * @Date: 2020/5/22 1:37 AM
 * @Desc

 */

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
