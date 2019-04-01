package server

import (
	"context"
	"io"
	// "log"
	"time"

	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/any"

	pb "github.com/ikuiki/golang-test/grpc-conn-create-benchmark/helloworld"
)

// SayHelloServer SayHello服务
type SayHelloServer struct{}

// ListHello 服务器端流式 RPC, 接收一次客户端请求，返回一个流
func (s *SayHelloServer) ListHello(in *pb.HelloWorldRequest, stream pb.HelloWorldService_ListHelloServer) error {
	// log.Printf("Client Say: %v", in.Greeting)

	stream.Send(&pb.HelloWorldResponse{Reply: "ListHello Reply " + in.Greeting + " 1"})
	time.Sleep(1 * time.Second)
	stream.Send(&pb.HelloWorldResponse{Reply: "ListHello Reply " + in.Greeting + " 2"})
	time.Sleep(1 * time.Second)
	stream.Send(&pb.HelloWorldResponse{Reply: "ListHello Reply " + in.Greeting + " 3"})
	time.Sleep(1 * time.Second)
	return nil
}

// SayMoreHello 客户端流式 RPC， 客户端流式请求，服务器可返回一次
func (s *SayHelloServer) SayMoreHello(stream pb.HelloWorldService_SayMoreHelloServer) error {
	// 接受客户端请求
	for {
		_, err := stream.Recv()
		if err == io.EOF {
			break
		}

		if err != nil {
			return err
		}

		// log.Printf("SayMoreHello Client Say: %v", req.Greeting)
	}

	// 流读取完成后，返回
	return stream.SendAndClose(&pb.HelloWorldResponse{Reply: "SayMoreHello Recv Muti Greeting"})
}

// SayHelloChat 双向流式RPC，客户端和服务器可以分别发送请求
func (s *SayHelloServer) SayHelloChat(stream pb.HelloWorldService_SayHelloChatServer) error {

	go func() {
		for {
			// req, err := stream.Recv()
			_, err := stream.Recv()
			if err == io.EOF {
				break
			}

			if err != nil {
				return
			}

			// log.Printf("SayHelloChat Client Say: %v", req.Greeting)
		}
	}()

	stream.Send(&pb.HelloWorldRequest{Greeting: "SayHelloChat Server Say Hello 1"})
	time.Sleep(1 * time.Second)
	stream.Send(&pb.HelloWorldRequest{Greeting: "SayHelloChat Server Say Hello 2"})
	time.Sleep(1 * time.Second)
	stream.Send(&pb.HelloWorldRequest{Greeting: "SayHelloChat Server Say Hello 3"})
	time.Sleep(1 * time.Second)
	return nil
}

// SayHelloWorld 单次请求
func (s *SayHelloServer) SayHelloWorld(ctx context.Context, in *pb.HelloWorldRequest) (res *pb.HelloWorldResponse, err error) {
	// log.Printf("Client Greeting:%s", in.Greeting)
	// log.Printf("Client Info:%v", in.Infos)

	var an *any.Any
	if in.Infos["hello"] == "world" {
		an, err = ptypes.MarshalAny(&pb.HelloWorld{Msg: "Good Request"})
	} else {
		an, err = ptypes.MarshalAny(&pb.Error{Msg: []string{"Bad Request", "Wrong Info Msg"}})
	}

	if err != nil {
		return
	}
	return &pb.HelloWorldResponse{
		Reply:   "Hello World !!",
		Details: []*any.Any{an},
	}, nil
}
