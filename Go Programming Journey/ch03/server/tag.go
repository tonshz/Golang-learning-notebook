package server

import (
	"ch03/pkg/bapi"
	"encoding/json"

	//"ch03/pkg/bapi"
	"ch03/pkg/errcode"
	pb "ch03/proto"
	"context"

	"google.golang.org/grpc/metadata"
	"log"
)

type TagServer struct {
	auth *Auth
}

type Auth struct{}

func (a *Auth) GetAppKey() string {
	return "admin"
}

func (a *Auth) GetAppSecret() string {
	return "go-learning"
}

func (a *Auth) Check(ctx context.Context) error {
	// 调用 metadata.FromIncomingContext() 从上下文中获取 metadata
	md, _ := metadata.FromIncomingContext(ctx)

	var appKey, appSecret string
	if value, ok := md["app_key"]; ok {
		appKey = value[0]
	}
	if value, ok := md["app_secret"]; ok {
		appSecret = value[0]
	}

	log.Printf("test: %s, %s", appKey, appSecret)
	if appKey != a.GetAppKey() || appSecret != a.GetAppSecret() {
		return errcode.TogRPCError(errcode.Unauthorized)
	}
	return nil
}

func NewTagServer() *TagServer {
	return &TagServer{}
}

func (t *TagServer) GetTagList(ctx context.Context, r *pb.GetTagListRequest) (*pb.GetTagListResponse, error) {
	if err := t.auth.Check(ctx); err != nil {
		return nil, err
	}

	// localhost：不通过网卡传输，不受网络防火墙和网卡相关的限制。
	// 127.0.0.1：通过网卡传输，依赖网卡，并受到网卡和防火墙相关的限制。
	api := bapi.NewAPI("http://127.0.0.1:8000")
	body, err := api.GetTagList(ctx, r.GetName())
	if err != nil {
		// 填入业务错误码
		return nil, errcode.TogRPCError(errcode.ErrorGetTagListFail)
	}

	tagList := pb.GetTagListResponse{}
	err = json.Unmarshal(body, &tagList)
	if err != nil {
		// 填入错误码
		return nil, errcode.TogRPCError(errcode.Fail)
	}

	return &tagList, nil
}
