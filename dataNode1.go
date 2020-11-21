package main

import (
	"io"
	"log"
	"net"

	pb "github.com/GianniCarlini/Lab-2-SD/proto"
	"google.golang.org/grpc"
)
const (
	port = ":50051"

)

type server struct {
}

func (s *server) EnviarLibro(stream pb.Packet_EnviarLibroServer) error {
	log.Println("Started stream")
	var contador = 0
	for {
		contador++
		_, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}
		log.Println("recibi"+string(contador))
		resp := pb.EnviarLibroReply{Id: "xd"}
		if err := stream.Send(&resp); err != nil { 
			log.Printf("send error %v", err)
		}
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