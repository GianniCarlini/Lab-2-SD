package main

import (
	"io"
	"log"
	"net"
	"fmt"
	"io/ioutil"
	"os"
	//"strconv"

	pb "github.com/GianniCarlini/Lab-2-SD/proto"
	"google.golang.org/grpc"
)
const (
	ip = "IPDATA1"
	port = ":50051"

)

type server struct {
}

func (s *server) EnviarLibro(stream pb.Packet_EnviarLibroServer) error {
	log.Println("Started stream")
	var contador = 0
	for {
		in, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}

		partBuffer := in.Id
		fileName := in.Name+"_"+ip
		resp := pb.EnviarLibroReply{Id: in.Name}
		if err := stream.Send(&resp); err != nil { 
			log.Printf("send error %v", err)
		}
		file, err := os.Create(fileName)
		defer file.Close()
		if err != nil {
				fmt.Println(err)
				os.Exit(1)
		}
		// write/save buffer to disk
		ioutil.WriteFile(fileName, partBuffer, os.ModeAppend)

		fmt.Println("Split to : ", fileName)
		contador++
	}
}

func main() {
	lis, err := net.Listen("tcp", port)

	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()

	pb.RegisterPacketServer(s, &server{})

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}