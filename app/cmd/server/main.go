package main

import (
	"fmt"
	"os"
	"log"
	"context"
	"net"
	"time"
	"errors"
	"io"
	"os/signal"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/reflection"
	// "google.golang.org/grpc/codes"
	// "google.golang.org/grpc/status"
	// "google.golang.org/genproto/googleapis/rpc/errdetails"
	hellopb "mygrpc/pkg/grpc"
)

type myServer struct {
	hellopb.UnimplementedGreetingServiceServer
}

func NewMyServer() *myServer {
	return &myServer{}
}

func myUnaryServerInterceptor1(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	log.Println("[pre] my unary server interceptor 1: ", info.FullMethod)
	res, err := handler(ctx, req)
	log.Println("[post] my unary server interceptor 1: ", res)
	return res, err
}

func myUnaryServerInterceptor2(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	log.Println("[pre] my unary server interceptor 2: ", info.FullMethod, req)
	res, err := handler(ctx, req)
	log.Println("[post] my unary server interceptor 2: ", res)
	return res, err
}

func myStreamServerInterceptor1(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
	log.Println("[pre stream] my stream server interceptor 1: ", info.FullMethod)

	err := handler(srv, &myServerStreamWrapper1{ss})

	log.Println("[post stream] my stream server interceptor 1: ")
	return err
}

func myStreamServerInterceptor2(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
	log.Println("[pre stream] my stream server interceptor 2: ", info.FullMethod)

	err := handler(srv, &myServerStreamWrapper1{ss})

	log.Println("[post stream] my stream server interceptor 2: ")
	return err
}

type myServerStreamWrapper1 struct {
	grpc.ServerStream
}

func (s *myServerStreamWrapper1) RecvMsg(m interface{}) error {
	err := s.ServerStream.RecvMsg(m)
	if !errors.Is(err, io.EOF) {
		log.Println("[pre message] my stream server interceptor 1: ", m)
	}
	return err
}

func (s *myServerStreamWrapper1) SendMsg(m interface{}) error {
	log.Println("[post message] my stream server interceptor 1: ", m)
	return s.ServerStream.SendMsg(m)
}

func (s *myServer) Hello(ctx context.Context, req *hellopb.HelloRequest) (*hellopb.HelloResponse, error) {
	// err := status.Error(codes.Unknown, "unknown error occured")

	if md, ok := metadata.FromIncomingContext(ctx); ok {
		log.Println(md)
	}

	headerMD := metadata.New(map[string]string{"type": "unary", "from": "server", "in": "header"})
	if err := grpc.SetHeader(ctx, headerMD); err != nil {
		return nil, err
	}

	trailerMD := metadata.New(map[string]string{"type": "unary", "from": "server", "in": "trailer"})
	if err := grpc.SetTrailer(ctx, trailerMD); err != nil {
		return nil, err
	}

	return &hellopb.HelloResponse{
		Message: fmt.Sprintf("Hello, %s!", req.GetName()),
	}, nil

	// stat := status.New(codes.Unknown, "unknown error occurred")
	// // stat, _ = stat.WithDetails()
	// stat, _ = stat.WithDetails(&errdetails.DebugInfo{
	// 	Detail: "detail reason of err",
	// })
	// err := stat.Err()
	// return nil, err
}

func (s *myServer) HelloServerStream(req *hellopb.HelloRequest, stream hellopb.GreetingService_HelloServerStreamServer) error {
	resCount := 5
	for i := 0; i < resCount; i++ {
		if err := stream.Send(&hellopb.HelloResponse{
			Message: fmt.Sprintf("[%d] Hello, %s!", i, req.GetName()),
		}); err != nil {
			return err
		}
		time.Sleep(time.Second * 1)
	}
	return nil
}

func (s *myServer) HelloClientStream(stream hellopb.GreetingService_HelloClientStreamServer) error {
	nameList := make([]string, 0)
	for {
		req, err := stream.Recv()
		if errors.Is(err, io.EOF) {
			message := fmt.Sprintf("Hello, %v!", nameList)
			return stream.SendAndClose(&hellopb.HelloResponse{
				Message: message,
			})
		}
		if err != nil {
			return err
		}
		nameList = append(nameList, req.GetName())
	}
}

func (s *myServer) HelloBiStreams(stream hellopb.GreetingService_HelloBiStreamsServer) error {
	if md, ok := metadata.FromIncomingContext(stream.Context()); ok {
		log.Println(md)
	}

	headerMD := metadata.New(map[string]string{"type": "stream", "from": "server", "in": "header"})
	if err := stream.SetHeader(headerMD); err != nil {
		return err
	}

	trailerMD := metadata.New(map[string]string{"type": "stream", "from": "server", "in": "trailer"})
	stream.SetTrailer(trailerMD)

	for {
		req, err := stream.Recv()
		if errors.Is(err, io.EOF) {
			return nil
		}
		if err != nil {
			return err
		}
		message := fmt.Sprintf("Hello, %v!", req.GetName())
		if err := stream.Send(&hellopb.HelloResponse{
			Message: message,
		}); err != nil {
			return err
		}
	}
}

func main() {
	port := 50051
	listner, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		panic(err)
	}

	// s := grpc.NewServer(
	// 	grpc.UnaryInterceptor(myUnaryServerInterceptor1),
	// )
	s := grpc.NewServer(
		grpc.ChainStreamInterceptor(
			myStreamServerInterceptor1,
			myStreamServerInterceptor2,
		),
	)

	hellopb.RegisterGreetingServiceServer(s, NewMyServer())

	reflection.Register(s)

	go func() {
		log.Printf("Start gRPC server port: %v", port)
		s.Serve(listner)
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	log.Println("Stopping gRPC server...")
	s.GracefulStop()
}
