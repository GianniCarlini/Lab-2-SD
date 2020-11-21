package main

import (
	"log"
	"time"

	"golang.org/x/net/context"
	"google.golang.org/grpc"

	pb "github.com/GianniCarlini/Lab-2-SD/proto"
)

func main() {
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("failed to connect: %s", err)
	}
	defer conn.Close()

	client := pb.NewPacketClient(conn)
	stream, err := client.EnviarLibro(context.Background())
	waitc := make(chan struct{})

	msg := &pb.EnviarLibroRequest{Id: "soy una id"}
	go func() {
		for i := 1; i <= 3; i++{
			log.Println("Sending msg...")
			stream.Send(msg)
		}
	}()
	<-waitc
	stream.CloseSend()
}