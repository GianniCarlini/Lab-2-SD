package main

import (
	"io"
	"log"
	"net"
	"fmt"
	"io/ioutil"
	"os"
	//"strconv"
	"context"
	"time"

	pb "github.com/GianniCarlini/Lab-2-SD/proto"
	"google.golang.org/grpc"
)
const (
	ip1 = "IPDATA1"
	//ip2 = "IPDATA1"
	//ip3 = "IPDATA1"
	port = ":50051" //puerto de data1server
	address = "localhost:50052" //namenode
	address2 = "localhost:50054" //data2
	//address3 = "localhost:50053" //namenode


)

type server struct {
}
func CrearDistribucion (prop []string) (p1 []string, p2 []string, p3 []string){
	p1 = append(p1,prop[0])
	p2 = append(p2,prop[1])
	p3 = prop[2:]
	return p1,p2,p3
}

func (s *server) EnviarLibro(stream pb.Packet_EnviarLibroServer) error {
	log.Println("Started stream Centralizado")
	var contador = 1
	var b [][]byte
	var xd []string
	for {
		in, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}

		partBuffer := in.Id
		fileName := in.Name+"_"+ip1
		resp := pb.EnviarLibroReply{Id: in.Name}
		if err := stream.Send(&resp); err != nil { 
			log.Printf("send error %v", err)
		}
		/*file, err := os.Create(fileName)
		defer file.Close()
		if err != nil {
				fmt.Println(err)
				os.Exit(1)
		}
		// write/save buffer to disk
		ioutil.WriteFile(fileName, partBuffer, os.ModeAppend)*/
		contador++
		fmt.Println("Split to : ", fileName)
		b = append(b,partBuffer) //arreglo con los bytes de los archivos
		xd = append(xd,fileName) //arreglo con los nombres de los archivos
		if uint64(contador) == in.Numero+1{
			b1 := b[:in.Numero/3]
			b2 := b[in.Numero/3:2*in.Numero/3]
			//b3 := b[2*in.Numero/3:]
			conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
			if err != nil {
				log.Fatalf("did not connect: %v", err)
			}
			defer conn.Close()
			c := pb.NewPropuestaCentralizadoClient(conn)
			ctx, cancel := context.WithTimeout(context.Background(), time.Second)
			defer cancel()
			fmt.Println("toy creando una distri")
			p1,p2,p3 := CrearDistribucion(xd)

			r, err := c.EnviarPropuestaCentralizado(ctx, &pb.PropuestaRequestC{Propuesta: xd, Propuesta1: p1, Propuesta2: p2, Propuesta3: p3})
			if err != nil {
				log.Fatalf("could not greet: %v", err)
			}
			fmt.Println(r.GetDistribucion1())
			dist2 :=r.GetDistribucion2()
			//-------------------creo mis archivos---------------------------------
			for i := range r.GetDistribucion1(){
				file, err := os.Create(r.GetDistribucion1()[i])
				defer file.Close()
				if err != nil {
						fmt.Println(err)
						os.Exit(1)
				}
				// write/save buffer to disk
				ioutil.WriteFile(r.GetDistribucion1()[i], b1[i], os.ModeAppend)
			}
			conn2, err2 := grpc.Dial(address2, grpc.WithInsecure(), grpc.WithBlock())
			if err2 != nil {
				log.Fatalf("did not connect: %v", err)
			}
			defer conn2.Close()
			c2 := pb.NewLibroDatasClient(conn2)
			ctx2, cancel := context.WithTimeout(context.Background(), time.Second)
			defer cancel()
			r2, err3 := c2.EnviarLibroData(ctx2, &pb.DataRequestC{Distribucion: dist2, Bites: b2, Numero: int64(in.Numero/3)})
			if err3 != nil {
				log.Fatalf("could not greet: %v", err)
			}
			fmt.Println(r2.GetEstado())
		}
	}
}
func (s *server) EnviarLibro2(stream pb.Distribuido_EnviarLibro2Server) error {
	log.Println("Started stream Distribuido")
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
		fileName := in.Name+"_"+ip1
		resp := pb.EnviarLibroReply2{Id: in.Name}
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
	var comportamiento int
	fmt.Println("Ingrese 1 para modo centralizado")
	fmt.Println("Ingrese 2 para modo distribuida")
	fmt.Scanln(&comportamiento)
	switch comportamiento{
		case 1:
			lis, err := net.Listen("tcp", port)

			if err != nil {
				log.Fatalf("failed to listen: %v", err)
			}
		
			s := grpc.NewServer()
		
			pb.RegisterPacketServer(s, &server{})
		
			if err := s.Serve(lis); err != nil {
				log.Fatalf("failed to serve: %v", err)
			}
		case 2:
			lis, err := net.Listen("tcp", port)

			if err != nil {
				log.Fatalf("failed to listen: %v", err)
			}
		
			s := grpc.NewServer()
		
			pb.RegisterDistribuidoServer(s, &server{})
		
			if err := s.Serve(lis); err != nil {
				log.Fatalf("failed to serve: %v", err)
			}
	}
}