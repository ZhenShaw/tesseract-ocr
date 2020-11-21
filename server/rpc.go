package server

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"net/rpc"

	"github.com/ZhenShaw/tesseract-ocr/orc"
	"github.com/astaxie/beego/logs"
)

func NetRPCServer(ln net.Listener, token string) {
	err := rpc.RegisterName("RPCOcrService", &RPCService{Token: token})
	if err != nil {
		log.Fatalf("failed to register: %v", err)
	}
	//将Rpc绑定到HTTP协议上
	//rpc.HandleHTTP()

	if err := http.Serve(ln, nil); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

type RPCService struct {
	Token string `json:"token"` //访问密码
	Data  []byte `json:"data"`  //识别数据
}

// 不使用gRPC的接口，使用net/rpc或net/rpc/jsonrpc 包的方式
func (s *RPCService) Recognize(req *RPCService, reply *string) (err error) {
	if req == nil || req.Token != s.Token {
		err = fmt.Errorf("rpc unauthorized")
		logs.Error(err)
		return
	}
	*reply, err = orc.Recognize(req.Data)
	if err != nil {
		logs.Error(err)
		return
	}
	logs.Debug("识别结果：", *reply)
	return
}
