package server

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/astaxie/beego/logs"
	"github.com/soheilhy/cmux"
	"github.com/zhenshaw/tesseract-ocr/orc"
	pb "github.com/zhenshaw/tesseract-ocr/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func GRPCServer(ln net.Listener, token string) {

	srv := grpc.NewServer()
	pb.RegisterCaptureOCRServer(srv, &GRPCService{Token: token})

	// Register reflection service on gRPC server.
	reflection.Register(srv)

	if err := srv.Serve(ln); err != cmux.ErrListenerClosed {
		log.Fatalf("failed to serve grpc: %v", err)
	}
}

type GRPCService struct {
	Token string `json:"token"` //访问密码
}

// 实现gRPC proto定义的接口
func (s *GRPCService) Recognize(ctx context.Context, req *pb.OCRRequest) (res *pb.OCRReply, err error) {

	if req == nil || req.Token != s.Token {
		err = fmt.Errorf("rpc unauthorized")
		logs.Error(err)
		return
	}
	res = new(pb.OCRReply)
	res.Code, err = orc.Recognize(req.Data)
	if err != nil {
		logs.Error(err)
		return
	}
	logs.Debug("识别结果：", res.Code)
	return
}
