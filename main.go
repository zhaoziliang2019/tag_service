package main

import (
	"context"
	"flag"
	"github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/soheilhy/cmux"
	"github.com/zhaoziliang2019/tag_service/middleware"
	pb "github.com/zhaoziliang2019/tag_service/proto"
	"github.com/zhaoziliang2019/tag_service/server"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
	"net/http"
)

var grpcPort string

var httpPort string
var port string

func init() {
	flag.StringVar(&grpcPort, "grpcPort", "8001", "gRPC启动端口号")
	flag.StringVar(&httpPort, "httpPort", "9001", "http启动端口号")
	flag.StringVar(&port, "port", "8003", "启动端口号")
	flag.Parse()
}

//func main() {
//	errs := make(chan error)
//	go func() {
//		err := RunHttpServer(httpPort)
//		if err != nil {
//			errs <- err
//		}
//	}()
//	go func() {
//		err := RunGrpcServer(grpcPort)
//		if err != nil {
//			errs <- err
//		}
//	}()
//	select {
//	case err := <-errs:
//		log.Fatalf("Run Server err:%v", err)
//	}
//}
func RunTCPServer(port string) (net.Listener, error) {
	return net.Listen("tcp", ":"+port)
}

//func RunHttpServer(port string) error {
//	serveMux := http.NewServeMux()
//	serveMux.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
//		_, _ = w.Write([]byte(`pong`))
//	})
//	return http.ListenAndServe(":"+port, serveMux)
//}
//func RunGrpcServer(port string) error {
//	s := grpc.NewServer()
//	pb.RegisterTagServiceServer(s, server.NewTagServer())
//	reflection.Register(s)
//	lis, err := net.Listen("tcp", ":"+port)
//	if err != nil {
//		return err
//	}
//	return s.Serve(lis)
//}
func RunGrpcServer() *grpc.Server {
	opts := []grpc.ServerOption{
		//grpc.UnaryInterceptor(HelloInterceptor),
		//grpc.UnaryInterceptor(WorldInterceptor),
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(HelloInterceptor, WorldInterceptor,
			middleware.AccessLog, middleware.ErrorLog, middleware.Recovery),
			grpc_middleware.ChainUnaryClient(grpc_retry.un)),
	}
	s := grpc.NewServer(opts...)
	pb.RegisterTagServiceServer(s, server.NewTagServer())
	reflection.Register(s)
	return s
}
func RunHttpServer(port string) *http.Server {
	serveMux := http.NewServeMux()
	serveMux.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`pong`))
	})
	//prefix := "/swagger-ui/"
	//fileServer := http.FileServer(&assetfs.AssetFS{
	//	Asset: swagger.,
	//	AssetDir:  nil,
	//	AssetInfo: nil,
	//	Prefix:    "",
	//	Fallback:  "",
	//})
	return &http.Server{
		Addr:    ":" + port,
		Handler: serveMux,
	}
}
func main() {
	l, err := RunTCPServer(port)
	if err != nil {
		log.Fatalf("Run TCP Server err:%v", err)
	}
	m := cmux.New(l)
	grpcL := m.MatchWithWriters(cmux.HTTP2MatchHeaderFieldPrefixSendSettings("content-type", "application/grpc"))
	httpL := m.Match(cmux.HTTP1Fast())
	grpcS := RunGrpcServer()
	httpS := RunHttpServer(port)
	go grpcS.Serve(grpcL)
	go httpS.Serve(httpL)
	err = m.Serve()
	if err != nil {
		log.Fatalf("Run Server err:%v", err)
	}
}
func HelloInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	log.Println("你好")
	resp, err := handler(ctx, req)
	log.Println("再见")
	return resp, err
}
func WorldInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	log.Println("你好")
	resp, err := handler(ctx, req)
	log.Println("再见")
	return resp, err
}
