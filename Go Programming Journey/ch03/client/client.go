package main

import (
	pb "ch03/proto"
	"context"
	"google.golang.org/grpc"
	"log"
)

func main() {
	/*
		context.Background()
		返回一个非零的空上下文。它永远不会被取消，没有价值，也没有最后期限。
		它通常由主函数、初始化和测试使用，并作为传入请求的顶级上下文。
	*/
	ctx := context.Background()

	//// 设置客户端拦截器
	//opts := []grpc.DialOption{
	//	grpc.WithUnaryInterceptor(
	//		grpc_middleware.ChainUnaryClient(
	//			middleware.UnaryContextTimeout(),
	//			grpc_retry.UnaryClientInterceptor(
	//				grpc_retry.WithMax(2),
	//				grpc_retry.WithCodes(
	//					codes.Unknown,
	//					codes.Internal,
	//					codes.DeadlineExceeded,
	//				),
	//			),
	//			middleware.ClientTracing(),
	//		)),
	//	grpc.WithStreamInterceptor(
	//		grpc_middleware.ChainStreamClient(
	//			middleware.StreamContextTimeout())),
	//}

	clientConn, _ := GetClientConn(ctx, "localhost:8004", nil)
	defer clientConn.Close()

	// 初始化指定 RPC Proto Service 的客户端实例对象
	tagServiceClient := pb.NewTagServiceClient(clientConn)
	newCtx := tagServiceClient.WithOrgCode(ctx, "Go 语言学习")
	// 发起指定 RPC 方法的调用
	resp, _ := tagServiceClient.GetTagList(newCtx, &pb.GetTagListRequest{Name: "Go"})

	log.Printf("resp: %v", resp)
}

func GetClientConn(ctx context.Context, target string, opts []grpc.DialOption) (*grpc.ClientConn, error) {
	// grpc.WithInsecure() 已弃用，作用为跳过对服务器证书的验证，此时客户端和服务端会使用明文通信
	// 使用 WithTransportCredentials 和 insecure.NewCredentials() 代替
	//opts = append(opts, grp c.WithInsecure())
	// insecure.NewCredentials 返回一个禁用传输安全的凭据
	opts = append(opts, grpc.WithInsecure())

	/*
		grpc.DialContext() 创建到给定目标的客户端连接。
		默认情况下，它是一个非阻塞拨号（该功能不会等待建立连接，并且连接发生在后台）。
		要使其成为阻塞拨号，请使用 WithBlock() 拨号选项。
	*/
	return grpc.DialContext(ctx, target, opts...)
}
