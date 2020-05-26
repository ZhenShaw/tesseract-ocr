package main

import (
	"flag"
	"fmt"
	ocr "github.com/ZhenShaw/tesseract-rpc/orc"
	pb "github.com/ZhenShaw/tesseract-rpc/proto"
	"log"
	"net"
	"net/http"
	"net/rpc"
	"net/rpc/jsonrpc"
	"os"
	"strconv"
	"strings"

	"github.com/astaxie/beego/logs"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const defaultPort = "8080"
const defaultServer = "grpc"
const defaultToken = ""

// 命令行参数
var argPort = flag.String("p", "", "port 指定端口")
var argServer = flag.String("s", "", "server 指定rpc类型，可选值[rpc|grpc/jsonrpc|all]")
var argToken = flag.String("t", "", "token 指定rpc请求token，防止非法连接，默认空")

func main() {

	flag.Parse()

	port := defaultPort
	if *argPort == "" {
		if os.Getenv("PORT") != "" {
			port = os.Getenv("PORT")
		}
	} else {
		port = *argPort
	}
	server := defaultServer
	if *argServer == "" {
		if os.Getenv("SERVER") != "" {
			server = os.Getenv("SERVER")
		}
	} else {
		server = *argServer
	}
	token := defaultToken
	if *argToken == "" {
		if os.Getenv("TOKEN") != "" {
			token = os.Getenv("TOKEN")
		}
	} else {
		token = *argToken
	}

	logs.Info("server: %s, port: %s, token: %s", server, port, token)

	switch strings.ToLower(server) {
	case "rpc":
		netRPCServer(port, token)

	case "grpc":
		gRPCServer(port, token)

	case "jsonrpc":
		jsonRPCServer(port, token)

	case "all":
		p, err := strconv.Atoi(port)
		if err != nil {
			log.Fatal(err)
		}
		port1 := port
		port2 := fmt.Sprint(p + 1)
		port3 := fmt.Sprint(p + 2)

		go netRPCServer(port2, token)
		go jsonRPCServer(port3, token)
		gRPCServer(port1, token)

	default:
		gRPCServer(port, token)
	}

}

// gRPCServer 使用grpc 对外提供服务
func gRPCServer(port string, token string) {
	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterCaptureOCRServer(s, &ocr.GRPCService{Token: token})

	// Register reflection service on gRPC server.
	reflection.Register(s)

	logs.Info("grpc server listening on port: %s", port)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

// netRPCServer 使用net/rpc 对外提供服务
func netRPCServer(port string, token string) {
	err := rpc.RegisterName("RPCOcrService", &ocr.RPCService{Token: token})
	if err != nil {
		log.Fatalf("failed to register: %v", err)
	}
	//将Rpc绑定到HTTP协议上
	rpc.HandleHTTP()

	listener, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	logs.Info("net/rpc server listening on port: %s", port)

	err = http.Serve(listener, nil)
	if err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

// jsonRPCServer 使用net/rpc/jsonrpc 对外提供服务
func jsonRPCServer(port string, token string) {

	err := rpc.RegisterName("JSONRPCOcrService", &ocr.RPCService{Token: token})
	if err != nil {
		log.Fatalf("failed to register: %v", err)
	}

	listener, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	logs.Info("jsonrpc server listening on port: %s", port)

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println(err)
			continue
		}
		go jsonrpc.ServeConn(conn)
	}
}
