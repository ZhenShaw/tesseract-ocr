/**
 * @File: example
 * @Author: Shaw
 * @Date: 2020/5/22 1:37 AM
 * @Desc

 */

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
