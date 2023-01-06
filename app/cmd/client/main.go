package main

import (
	"bufio"
	"fmt"
	"os"
	"log"
	"context"
	"errors"
	"io"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/credentials/insecure"
	_ "google.golang.org/genproto/googleapis/rpc/errdetails"
	hellopb "mygrpc/pkg/grpc"
)

var (
	scanner *bufio.Scanner
	client hellopb.GreetingServiceClient
)

func Hello() {
	fmt.Println("Please enter your name.")
	scanner.Scan()
	name := scanner.Text()

	req := &hellopb.HelloRequest{
		Name: name,
	}
	ctx := context.Background()
	md := metadata.New(map[string]string{"type": "unary", "from": "client"})
	ctx = metadata.NewOutgoingContext(ctx, md)
	// res, err := client.Hello(context.Background(), req)
	var header, trailer metadata.MD
	// res, err := client.Hello(ctx, req)
	res, err := client.Hello(ctx, req, grpc.Header(&header), grpc.Trailer(&trailer))
	if err != nil {
		// fmt.Println(err)
		if stat, ok := status.FromError(err); ok {
			fmt.Printf("code: %s\n", stat.Code())
			fmt.Printf("message: %s\n", stat.Message())
			fmt.Printf("details: %s\n", stat.Details())
		} else {
			fmt.Println(err)
		}
	} else {
		fmt.Println(header)
		fmt.Println(trailer)
		fmt.Println(res.GetMessage())
	}
}

func HelloServerStream() {
	fmt.Println("Please enter your name.")
	scanner.Scan()
	name := scanner.Text()

	req := &hellopb.HelloRequest{
		Name: name,
	}
	stream, err := client.HelloServerStream(context.Background(), req)
	if err != nil {
		fmt.Println(err)
		return
	}

	for {
		res, err := stream.Recv()
		if errors.Is(err, io.EOF) {
			fmt.Println("all th responses have already received.")
			break
		}

		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(res)
	}
}

func HelloClientStream() {
	stream, err := client.HelloClientStream(context.Background())
	if err != nil {
		fmt.Println(err)
		return
	}

	sendCount := 5
	fmt.Printf("Please enter %d names.\n", sendCount)
	for i := 0; i < sendCount; i++ {
		scanner.Scan()
		name := scanner.Text()

		if err := stream.Send(&hellopb.HelloRequest{
			Name: name,
		}); err != nil {
			fmt.Println(err)
			return
		}
	}

	res, err := stream.CloseAndRecv()
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(res.GetMessage())
	}
}

func HelloBiStreams() {
	ctx := context.Background()
	md := metadata.New(map[string]string{"type": "stream", "from": "client"})
	ctx = metadata.NewOutgoingContext(ctx, md)
	// stream, err := client.HelloBiStreams(context.Background())
	stream, err := client.HelloBiStreams(ctx)
	if err != nil {
		fmt.Println(err)
		return
	}

	sendNum := 5
	fmt.Printf("Please enter %d names.\n", sendNum)

	var sendEnd, recvEnd bool
	sendCount := 0
	for !(sendEnd && recvEnd) {
		if !sendEnd {
			scanner.Scan()
			name := scanner.Text()

			sendCount++
			if err := stream.Send(&hellopb.HelloRequest{
				Name: name,
			}); err != nil {
				fmt.Println(err)
				sendEnd = true
			}

			if sendCount == sendNum {
				sendEnd = true
				if err := stream.CloseSend(); err != nil {
					fmt.Println(err)
				}
			}
		}

		var headerMD metadata.MD
		if !recvEnd {
			if headerMD == nil {
				headerMD, err = stream.Header()
				if err != nil {
					fmt.Println(err)
				} else {
					fmt.Println(headerMD)
				}
			}

			if res, err := stream.Recv(); err != nil {
				if !errors.Is(err, io.EOF) {
					fmt.Println(err)
				}
				recvEnd = true
			} else {
				fmt.Println(res.GetMessage())
			}
		}
	}

	trailerMD := stream.Trailer()
	fmt.Println(trailerMD)
}

func myUnaryClientInteceptor1(ctx context.Context, method string, req, res interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
	fmt.Println("[pre] my unary client interceptor 1", method, req)
	err := invoker(ctx, method, req, res, cc, opts...)
	fmt.Println("[post] my unary client interceptor 1", res)
	return err
}

func myStreamClientInteceptor1(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn, method string, streamer grpc.Streamer, opts ...grpc.CallOption) (grpc.ClientStream, error) {
	log.Println("[pre] my stream client interceptor 1", method)

	stream, err := streamer(ctx, desc, cc, method, opts...)
	return &myClientStreamWrapper1{stream}, err
}

type myClientStreamWrapper1 struct {
	grpc.ClientStream
}

func (s *myClientStreamWrapper1) SendMsg(m interface{}) error {
	log.Println("[pre message] my stream client interceptor 1: ", m)
	return s.ClientStream.SendMsg(m)
}

func (s *myClientStreamWrapper1) RecvMsg(m interface{}) error {
	err := s.ClientStream.RecvMsg(m)

	if !errors.Is(err, io.EOF) {
		log.Println("[post message] my stream client interceptor 1: ", m)
	}
	return err
}

func (s *myClientStreamWrapper1) CloseSend() error {
	err := s.ClientStream.CloseSend()

	log.Println("[post] my stream client interceptor 1")
	return err
}

func main() {
	fmt.Println("Start gRPC client.")

	scanner = bufio.NewScanner(os.Stdin)

	address := "localhost:50051"
	conn, err := grpc.Dial(
		address,
		grpc.WithStreamInterceptor(myStreamClientInteceptor1),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
	)

	if err != nil {
		log.Fatal("Connection failed!")
		return
	}
	defer conn.Close()

	client = hellopb.NewGreetingServiceClient(conn)

	for {
		fmt.Println("1: send Request")
		fmt.Println("2: HelloServerStream")
		fmt.Println("3: HelloClientStream")
		fmt.Println("4: HelloBiStream")
		fmt.Println("5: exit")
		fmt.Print("please enter >")

		scanner.Scan()
		in := scanner.Text()

		switch in {
		case "1":
			Hello()
		case "2":
			HelloServerStream()
		case "3":
			HelloClientStream()
		case "4":
			HelloBiStreams()
		case "5":
			fmt.Println("bye.")
			goto M
		}
	}
M:

}
