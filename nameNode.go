package main

import (
	"context"
	"log"
	"net"
	"fmt"
	"time"
	"math/rand"

	pb "github.com/GianniCarlini/Lab-2-SD/proto"
	"google.golang.org/grpc"
)

const (
	port = ":50052"
)

type server struct {
}

func (s *server) EnviarPropuestaCentralizado(ctx context.Context, in *pb.PropuestaRequestC) (*pb.PropuestaReplyC, error) {
	p1 := in.GetPropuesta1()
	p2 := in.GetPropuesta2()
	p3 := in.GetPropuesta3()

	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)
	aceptar := r1.Intn(101)
	if aceptar <= 50{
		fmt.Println("Propuesta inicial aceptada")
		return &pb.PropuestaReplyC{Distribucion1: p1,Distribucion2: p2,Distribucion3: p3}, nil
	}else{
		for{
			s2 := rand.NewSource(time.Now().UnixNano())
			r2 := rand.New(s2)
			aceptar := r2.Intn(101)
			fmt.Println(aceptar)
			var d1 []string
			var d2 []string
			var d3 []string
			if aceptar <= 90{
				fmt.Println("Propuesta Aceptada")
				d1 = in.GetPropuesta()[:len(in.GetPropuesta())/3]
				d2 = in.GetPropuesta()[len(in.GetPropuesta())/3:2*len(in.GetPropuesta())/3]
				d3 = in.GetPropuesta()[2*len(in.GetPropuesta())/3:]
				return &pb.PropuestaReplyC{Distribucion1: d1,Distribucion2: d2,Distribucion3: d3}, nil
			}else{
				fmt.Println("Propuesta rechazada")
				continue
			}
		}
	}
	
}

func main() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterPropuestaCentralizadoServer(s, &server{})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
