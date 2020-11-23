package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"log"
	"net"
	//"strconv"
	"context"

	pb "github.com/GianniCarlini/Lab-2-SD/proto"
	"google.golang.org/grpc"
)
const (
	//ip1 = "IPDATA1"
	//ip2 = "IPDATA1"
	//ip3 = "IPDATA1"
	port = ":50054" //puerto de data1server
	//address = "localhost:50052" //namenode
	//address2 = "localhost:50054" //data2
	//address3 = "localhost:50053" //namenode


)

type server struct {
}
func (s *server) EnviarLibroData(ctx context.Context, in *pb.DataRequestC) (*pb.DataReplyC, error) {
		
	for i := range in.GetDistribucion(){
		file, err := os.Create(in.GetDistribucion()[i])
		defer file.Close()
		if err != nil {
				fmt.Println(err)
				os.Exit(1)
		}
		// write/save buffer to disk
		ioutil.WriteFile(in.GetDistribucion()[i], in.GetBites()[i], os.ModeAppend)
	}

	return &pb.DataReplyC{Estado: "OK DATA2"}, nil
}
func main(){
	lis, err := net.Listen("tcp", port)

	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()

	pb.RegisterLibroDatasServer(s, &server{})

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}