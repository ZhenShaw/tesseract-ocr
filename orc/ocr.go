package ocr

import (
	"context"
	"fmt"
	pb "github.com/ZhenShaw/tesseract-rpc/proto"

	"github.com/astaxie/beego/logs"
	"github.com/otiai10/gosseract"
)

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

	req.Data, err = RemoveBackground(req.Data)
	if err != nil {
		logs.Error(err)
		return
	}

	client := gosseract.NewClient()
	defer client.Close()

	err = client.SetImageFromBytes(req.Data)
	if err != nil {
		logs.Error(err)
		return
	}
	res = new(pb.OCRReply)
	res.Code, err = client.Text()
	if err != nil {
		logs.Error(err)
		return
	}
	logs.Debug("识别结果：", res.Code)

	return
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

	client := gosseract.NewClient()
	defer client.Close()

	err = client.SetImageFromBytes(req.Data)
	if err != nil {
		logs.Error(err)
		return
	}
	*reply, err = client.Text()
	if err != nil {
		logs.Error(err)
		return
	}
	logs.Debug("识别结果：", *reply)

	return
}
