package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"log"
	"net"
	//"strconv"
	"context"
	"io"
	"time"
	"math/rand"

	pb "github.com/GianniCarlini/Lab-2-SD/proto"
	"google.golang.org/grpc"
)
const (

	port = ":50053" //puerto de data1server
	address = "localhost:50052" //namenode
	address2 = "localhost:50051" //data1 no cordinador
	address3 = "localhost:50054" //data 2 no cordinador


)

type server struct {
}

func CrearDistribucion (prop []string) (p1 []string, p2 []string, p3 []string){
	p1 = append(p1,prop[0])
	p2 = append(p2,prop[1])
	p3 = prop[2:]
	return p1,p2,p3
}

func (s *server) EnviarLibroData(ctx context.Context, in *pb.DataRequestC) (*pb.DataReplyC, error) {
	fmt.Println(in.GetDistribucion())
	fmt.Println(len(in.GetBites()))
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
	return &pb.DataReplyC{Estado: "OK DATA3"}, nil
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
		fileName := in.Name
		resp := pb.EnviarLibroReply{Id: in.Name}
		if err := stream.Send(&resp); err != nil { 
			log.Printf("send error %v", err)
		}
		contador++
		fmt.Println("Split to : ", fileName)
		b = append(b,partBuffer) //arreglo con los bytes de los archivos
		xd = append(xd,fileName) //arreglo con los nombres de los archivos
		if uint64(contador) == in.Numero+1{
			b1 := b[:in.Numero/3]
			b2 := b[in.Numero/3:2*in.Numero/3]
			b3 := b[2*in.Numero/3:]
			b3i := b[2:]
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

			r, err := c.EnviarPropuestaCentralizado(ctx, &pb.PropuestaRequestC{Propuesta: xd, Propuesta1: p1, Propuesta2: p2, Propuesta3: p3, Nombre: in.GetLibro()})
			if err != nil {
				log.Fatalf("could not greet: %v", err)
			}
			fmt.Println(r.GetDistribucion1())
			dist2 :=r.GetDistribucion2()
			dist3 :=r.GetDistribucion3()
			var distri3 [][]byte
			if r.GetTipo() == int64(1){
				fmt.Println("soy inicial")
				distri3 = b3i
			}else{
				fmt.Println("soy secundario")
				distri3 = b3
			}
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
			//------------------------------------------
			conn3, err3 := grpc.Dial(address3, grpc.WithInsecure(), grpc.WithBlock())
			if err3 != nil {
				log.Fatalf("did not connect: %v", err3)
			}
			defer conn3.Close()
			c3 := pb.NewLibroDatasClient(conn3)
			ctx3, cancel := context.WithTimeout(context.Background(), time.Second)
			defer cancel()
			r3, err4 := c3.EnviarLibroData(ctx3, &pb.DataRequestC{Distribucion: dist3, Bites: distri3, Numero: int64(in.Numero/3)})
			if err4 != nil {
				log.Fatalf("could not greet: %v", err)
			}
			fmt.Println(r3.GetEstado())
		}
	}
}
//------------------------- No cordinador
func (s *server) EnviarLibro2(stream pb.Distribuido_EnviarLibro2Server) error {
	log.Println("Started stream Distribuido")
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
		fileName := in.Name
		resp := pb.EnviarLibroReply2{Id: in.Name}
		if err := stream.Send(&resp); err != nil { 
			log.Printf("send error %v", err)
		}
		b = append(b,partBuffer) //arreglo con los bytes de los archivos
		xd = append(xd,fileName) //arreglo con los nombres de los archivos
		contador++
		fmt.Println("Split to : ", fileName)
		if uint64(contador) == in.Numero+1{
			var p1i [][]byte
			var p2i [][]byte
			var p3i [][]byte
			p1i = append(p1i,b[0])
			p2i = append(p2i,b[1])
			p3i = b[2:]
			//--------------------DATA1-----------------------------------------------
			conn, err := grpc.Dial(address2, grpc.WithInsecure(), grpc.WithBlock())
			if err != nil {
				log.Fatalf("did not connect: %v", err)
			}
			defer conn.Close()
			c := pb.NewDistribuidoClient(conn)
			ctx, cancel := context.WithTimeout(context.Background(), time.Second)
			defer cancel()
			r, err := c.EnviarPropuestaDistribuido(ctx, &pb.DistribuidoRequest{Propuestaini: xd, Propuestaini1: p1i, Propuestaini2: p2i, Propuestaini3: p3i, Flag: int64(1), Aceptada: false})
			if err != nil {
				log.Fatalf("could not greet: %v", err)
			}
			respinidata2 := r.GetRespuesta()
			//--------------------DATA2-----------------------------------------------
			conn2, err2 := grpc.Dial(address3, grpc.WithInsecure(), grpc.WithBlock())
			if err2 != nil {
				log.Fatalf("did not connect: %v", err2)
			}
			defer conn.Close()
			c2 := pb.NewDistribuidoClient(conn2)
			ctx2, cancel := context.WithTimeout(context.Background(), time.Second)
			defer cancel()
			r2, err3 := c2.EnviarPropuestaDistribuido(ctx2, &pb.DistribuidoRequest{Propuestaini: xd, Propuestaini1: p1i, Propuestaini2: p2i, Propuestaini3: p3i, Flag: int64(1), Aceptada: false})
			if err3 != nil {
				log.Fatalf("could not greet: %v", err3)
			}
			respinidata3 := r2.GetRespuesta()
			//-----------------------------------------------------------------------
			if (respinidata2 == "ACEPTADA" && respinidata3 == "ACEPTADA"){
					fmt.Println("Me aceptaron los 2")
					fmt.Println("Propuesta inicial aceptada")
					//----------------Data1Acept-----------------------------------------------
					conn, err := grpc.Dial(address2, grpc.WithInsecure(), grpc.WithBlock())
					if err != nil {
						log.Fatalf("did not connect: %v", err)
					}
					defer conn.Close()
					c := pb.NewDistribuidoClient(conn)
					ctx, cancel := context.WithTimeout(context.Background(), time.Second)
					defer cancel()
					r, err := c.EnviarPropuestaDistribuido(ctx, &pb.DistribuidoRequest{Propuestaini: xd, Propuestaini1: p1i, Propuestaini2: p2i, Propuestaini3: p3i, Flag: int64(1), Aceptada: true})
					if err != nil {
						log.Fatalf("could not greet: %v", err)
					}
					fmt.Println(r.GetRespuesta())
					//----------------Data3Acept-----------------------------------------------
					conn2, err2 := grpc.Dial(address3, grpc.WithInsecure(), grpc.WithBlock())
					if err2 != nil {
						log.Fatalf("did not connect: %v", err2)
					}
					defer conn.Close()
					c2 := pb.NewDistribuidoClient(conn2)
					ctx2, cancel := context.WithTimeout(context.Background(), time.Second)
					defer cancel()
					r2, err3 := c2.EnviarPropuestaDistribuido(ctx2, &pb.DistribuidoRequest{Propuestaini: xd, Propuestaini1: p1i, Propuestaini2: p2i, Propuestaini3: p3i, Flag: int64(1), Aceptada: true})
					if err3 != nil {
						log.Fatalf("could not greet: %v", err3)
					}
					fmt.Println(r2.GetRespuesta())
					//------------------escribo-----------------------------------
					names := xd[2:]
					bita := b[2:]
					for i := range names{
						file, err := os.Create(names[i])
						defer file.Close()
						if err != nil {
								fmt.Println(err)
								os.Exit(1)
						}
						// write/save buffer to disk
						ioutil.WriteFile(names[i], bita[i], os.ModeAppend)
					}
			}else{
				fmt.Println("Propuesta inicial rechazada")
				fmt.Println("Generando nueva propuesta")

				p1 := b[:in.Numero/3]
				p2 := b[in.Numero/3:2*in.Numero/3]
				p3 := b[2*in.Numero/3:]
				
				//-----------------------DATA1 NEW---------------------------------------
				conn, err := grpc.Dial(address2, grpc.WithInsecure(), grpc.WithBlock())
				if err != nil {
					log.Fatalf("did not connect: %v", err)
				}
				defer conn.Close()
				c := pb.NewDistribuidoClient(conn)
				ctx, cancel := context.WithTimeout(context.Background(), time.Second)
				defer cancel()
				r, err := c.EnviarPropuestaDistribuido(ctx, &pb.DistribuidoRequest{Propuestaini: xd, Propuestaini1: p1, Propuestaini2: p2, Propuestaini3: p3, Flag: int64(2), Aceptada: true})
				if err != nil {
					log.Fatalf("could not greet: %v", err)
				}
				fmt.Println(r.GetRespuesta())
				//-----------------------DATA3 NEW---------------------------------------
				conn2, err2 := grpc.Dial(address3, grpc.WithInsecure(), grpc.WithBlock())
				if err2 != nil {
					log.Fatalf("did not connect: %v", err2)
				}
				defer conn.Close()
				c2 := pb.NewDistribuidoClient(conn2)
				ctx2, cancel := context.WithTimeout(context.Background(), time.Second)
				defer cancel()
				r2, err3 := c2.EnviarPropuestaDistribuido(ctx2, &pb.DistribuidoRequest{Propuestaini: xd, Propuestaini1: p1, Propuestaini2: p2, Propuestaini3: p3, Flag: int64(2), Aceptada: true})
				if err3 != nil {
					log.Fatalf("could not greet: %v", err3)
				}
				fmt.Println(r2.GetRespuesta())
				//----------------------escribo--------------------------------------------
				names := xd[2*len(xd)/3:]
				bita := b[2*len(b)/3:]
				for i := range names{
					file, err := os.Create(names[i])
					defer file.Close()
					if err != nil {
							fmt.Println(err)
							os.Exit(1)
					}
					// write/save buffer to disk
					ioutil.WriteFile(names[i], bita[i], os.ModeAppend)
				}

				
			}
	
		}

	}
}
func (s *server) EnviarChunk(ctx context.Context, in *pb.DataChunkRequest) (*pb.DataChunkReply, error) {
	log.Printf("Received: %v", in.GetFilechunk())
	b, err := ioutil.ReadFile(in.GetFilechunk()) // just pass the file name
       if err != nil {
           fmt.Print(err)
       }
	return &pb.DataChunkReply{Bitaso: b}, nil
}
func (s *server) EnviarChunk2(ctx context.Context, in *pb.DataChunkRequest2) (*pb.DataChunkReply2, error) {
	log.Printf("Received: %v", in.GetFilechunk())
	b, err := ioutil.ReadFile(in.GetFilechunk()) // just pass the file name
       if err != nil {
           fmt.Print(err)
       }
	return &pb.DataChunkReply2{Bitaso: b}, nil
}

func (s *server) EnviarPropuestaDistribuido(ctx context.Context, in *pb.DistribuidoRequest) (*pb.DistribuidoReply, error) {
	log.Printf("Received: %v", in.GetPropuestaini())
	var resp string

	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)
	aceptar := r1.Intn(101)

	if in.GetFlag() == int64(1){
		if in.GetAceptada() == true{
			fmt.Println("Propuesta incial aceptada uwu")
			resp = "OK3"
			for i := 2; i < len(in.GetPropuestaini()); i++{
				file, err := os.Create(in.GetPropuestaini()[i])
				defer file.Close()
				if err != nil {
						fmt.Println(err)
						os.Exit(1)
				}
				// write/save buffer to disk
				ioutil.WriteFile(in.GetPropuestaini()[i], in.GetPropuestaini3()[i-2], os.ModeAppend)
			}
		}else if in.GetAceptada() == false{
			fmt.Println("Me llego una propuesta inicial")
			if aceptar <= 50{
				resp = "ACEPTADA"
			}else{
				resp = "RECHAZADA"
			}
		}
	}else if in.GetFlag() == int64(2){
		fmt.Println("Esta es una segunda propuesta")
		resp = "ACEPTADA"
		names := in.GetPropuestaini()[2*len(in.GetPropuestaini())/3:]
		for i := range names{
			file, err := os.Create(names[i])
			defer file.Close()
			if err != nil {
					fmt.Println(err)
					os.Exit(1)
			}
			// write/save buffer to disk
			ioutil.WriteFile(names[i], in.GetPropuestaini3()[i], os.ModeAppend)
		}
	}
	return &pb.DistribuidoReply{Respuesta: resp}, nil
}
func main() {
		var comportamiento int
		fmt.Println("Ingrese 1 para modo centralizado")
		fmt.Println("Ingrese 2 para modo distribuida")
		fmt.Scanln(&comportamiento)
		switch comportamiento{
			case 1:
				var cordinador int
				fmt.Println("ingrese 1 si este es el nodo cordinador, de no ser ingrese 2")
				fmt.Scanln(&cordinador)
				switch cordinador{
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
				
					pb.RegisterLibroDatasServer(s, &server{})
				
					if err := s.Serve(lis); err != nil {
						log.Fatalf("failed to serve: %v", err)
					}
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