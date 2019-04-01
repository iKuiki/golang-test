package benchmark_test

import (
	"crypto/tls"
	"crypto/x509"
	pb "github.com/ikuiki/golang-test/grpc-conn-create-benchmark/helloworld"
	"github.com/ikuiki/golang-test/grpc-conn-create-benchmark/server"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"io/ioutil"
	"log"
	"net"
	"testing"
)

var (
	certPool *x509.CertPool
)

func init() {
	// 将根证书加入证书词
	// 测试证书的根如果不加入可信池，那么测试证书将视为不可惜，无法通过验证。
	certPool = x509.NewCertPool()
	rootBuf, err := ioutil.ReadFile("certs/myssl_root.cer")
	if err != nil {
		panic(err)
	}

	if !certPool.AppendCertsFromPEM(rootBuf) {
		panic("fail to append test ca")
	}
	go createServer()
}

func createServer() {
	// 创建一个Server并运行
	go func() {
		lis, err := net.Listen("tcp", ":8080")
		if err != nil {
			panic(err)
		}

		// 加载证书和密钥 （同时能验证证书与私钥是否匹配）
		cert, err := tls.LoadX509KeyPair("certs/server_chain.pem", "certs/server.key")
		if err != nil {
			panic(err)
		}

		tlsConf := &tls.Config{
			ClientAuth:   tls.RequireAndVerifyClientCert,
			Certificates: []tls.Certificate{cert},
			ClientCAs:    certPool,
		}

		serverOpt := grpc.Creds(credentials.NewTLS(tlsConf))
		grpcServer := grpc.NewServer(serverOpt)

		pb.RegisterHelloWorldServiceServer(grpcServer, &server.SayHelloServer{})

		log.Println("Server Start...")
		grpcServer.Serve(lis)
	}()
}

func createClientTransportCreds() credentials.TransportCredentials {

	cert, err := tls.LoadX509KeyPair("certs/client_chain.pem", "certs/client.key")
	if err != nil {
		panic(err)
	}

	// 将根证书加入证书池
	certPool := x509.NewCertPool()
	bs, err := ioutil.ReadFile("certs/myssl_root.cer")
	if err != nil {
		panic(err)
	}

	if !certPool.AppendCertsFromPEM(bs) {
		panic("cc")
	}

	// 新建凭证
	transportCreds := credentials.NewTLS(&tls.Config{
		ServerName:   "metal.kuiki.cn",
		Certificates: []tls.Certificate{cert},
		RootCAs:      certPool,
	})
	return transportCreds
}

func BenchmarkNewConn(b *testing.B) {
	transportCreds := createClientTransportCreds()
	dialOpt := grpc.WithTransportCredentials(transportCreds)
	for i := 0; i < b.N; i++ {
		conn, err := grpc.Dial("localhost:8080", dialOpt)
		if err != nil {
			b.Fatalf("Dial failed:%v", err)
		}
		defer conn.Close()

		client := pb.NewHelloWorldServiceClient(conn)
		resp1, err := client.SayHelloWorld(context.Background(), &pb.HelloWorldRequest{
			Greeting: "Hello Server 1 !!",
			Infos:    map[string]string{"hello": "world"},
		})
		if err != nil {
			b.Fatal(err)
		}
		b.Logf("Resp1:%+v", resp1)
	}
}

func BenchmarkNewClient(b *testing.B) {
	transportCreds := createClientTransportCreds()
	dialOpt := grpc.WithTransportCredentials(transportCreds)
	conn, err := grpc.Dial("localhost:8080", dialOpt)
	if err != nil {
		b.Fatalf("Dial failed:%v", err)
	}
	defer conn.Close()
	for i := 0; i < b.N; i++ {
		client := pb.NewHelloWorldServiceClient(conn)
		resp1, err := client.SayHelloWorld(context.Background(), &pb.HelloWorldRequest{
			Greeting: "Hello Server 1 !!",
			Infos:    map[string]string{"hello": "world"},
		})
		if err != nil {
			b.Fatal(err)
		}
		b.Logf("Resp1:%+v", resp1)
	}
}

func BenchmarkReuseClient(b *testing.B) {
	transportCreds := createClientTransportCreds()
	dialOpt := grpc.WithTransportCredentials(transportCreds)
	conn, err := grpc.Dial("localhost:8080", dialOpt)
	if err != nil {
		b.Fatalf("Dial failed:%v", err)
	}
	defer conn.Close()
	client := pb.NewHelloWorldServiceClient(conn)
	for i := 0; i < b.N; i++ {
		resp1, err := client.SayHelloWorld(context.Background(), &pb.HelloWorldRequest{
			Greeting: "Hello Server 1 !!",
			Infos:    map[string]string{"hello": "world"},
		})
		if err != nil {
			b.Fatal(err)
		}
		b.Logf("Resp1:%+v", resp1)
	}
}
