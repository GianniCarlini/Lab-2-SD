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
	"math/rand"


	pb "github.com/GianniCarlini/Lab-2-SD/proto"
	"google.golang.org/grpc"
)
const (
	//ip1 = "IPDATA1"
	//ip2 = "IPDATA1"
	//ip3 = "IPDATA1"
	port2 = ":50055"
	port = ":50051" //puerto de data1server
	address = "localhost:50052" //namenode
	address2 = "localhost:50054" //data2
	address3 = "localhost:50053" //data 3 no cordinador


)

type server struct {
}

var state = false
var idnode = 5

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
		fileName := in.Name
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
func (s *server) EnviarLibro2(stream pb.Distribuido_EnviarLibro2Server) error {
	log.Println("Started stream Distribuido")
	var contador = 1
	var b [][]byte
	var xd []string

	state = true
	for !RyA(){} 

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
			//--------------------DATA2-----------------------------------------------
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
			//--------------------DATA3-----------------------------------------------
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
					//----------------Data2Acept-----------------------------------------------
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
					file, err := os.Create(xd[0])
					defer file.Close()
					if err != nil {
							fmt.Println(err)
							os.Exit(1)
					}
					// write/save buffer to disk
					ioutil.WriteFile(xd[0], b[0], os.ModeAppend)
					//---------------------log--------------------------------------
					var log1 []string
					var log2 []string
					log1 = append(log1,xd[0])
					log2 = append(log2,xd[1])
					conn6, err6 := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
					if err6 != nil {
						log.Fatalf("did not connect: %v", err6)
					}
					defer conn.Close()
					c6 := pb.NewDistribuidoClient(conn6)
					ctx6, cancel := context.WithTimeout(context.Background(), time.Second)
					defer cancel()
		
					r6, err7 := c6.Escribir(ctx6, &pb.EscribirRequest{Propuesta: xd, Propuesta1: log1, Propuesta2: log2, Propuesta3: xd[2:], Nombre: in.GetLibro()})
					if err7 != nil {
						log.Fatalf("could not greet: %v", err7)
					}
	
					fmt.Println(r6.GetEstate())
			}else{
				fmt.Println("Propuesta inicial rechazada")
				fmt.Println("Generando nueva propuesta")

				p1 := b[:in.Numero/3]
				p2 := b[in.Numero/3:2*in.Numero/3]
				p3 := b[2*in.Numero/3:]
				
				//-----------------------DATA2 NEW---------------------------------------
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
				names := xd[:len(xd)/3]
				bita := b[:len(xd)/3]
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
				conn5, err5 := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
				if err5 != nil {
					log.Fatalf("did not connect: %v", err5)
				}
				defer conn.Close()
				c5 := pb.NewDistribuidoClient(conn5)
				ctx5, cancel := context.WithTimeout(context.Background(), time.Second)
				defer cancel()
	
				r5, err6 := c5.Escribir(ctx5, &pb.EscribirRequest{Propuesta: xd, Propuesta1: xd[:len(xd)/3], Propuesta2: xd[len(xd)/3:2*len(xd)/3], Propuesta3: xd[2*len(xd)/3:], Nombre: in.GetLibro()})
				if err6 != nil {
					log.Fatalf("could not greet: %v", err6)
				}

				fmt.Println(r5.GetEstate())
			}
	
		}

	}
}
//-------------------------------------------------------
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

	return &pb.DataReplyC{Estado: "OK DATA1"}, nil
}
//-------------------------------------------------------
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
//--------------------------------------------------------------
func (s *server) EnviarPropuestaDistribuido(ctx context.Context, in *pb.DistribuidoRequest) (*pb.DistribuidoReply, error) {
	log.Printf("Received: %v", in.GetPropuestaini())
	var resp string

	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)
	aceptar := r1.Intn(101)

	if in.GetFlag() == int64(1){
		if in.GetAceptada() == true{
			fmt.Println("Propuesta incial aceptada uwu")
			resp = "OK1"

			file, err := os.Create(in.GetPropuestaini()[0])
				defer file.Close()
				if err != nil {
						fmt.Println(err)
						os.Exit(1)
				}
				// write/save buffer to disk
				ioutil.WriteFile(in.GetPropuestaini()[0], in.GetPropuestaini1()[0], os.ModeAppend)
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
		names := in.GetPropuestaini()[:len(in.GetPropuestaini())/3]
		for i := range names{
			file, err := os.Create(names[i])
			defer file.Close()
			if err != nil {
					fmt.Println(err)
					os.Exit(1)
			}
			// write/save buffer to disk
			ioutil.WriteFile(names[i], in.GetPropuestaini1()[i], os.ModeAppend)
		}
	}
	return &pb.DistribuidoReply{Respuesta: resp}, nil
}
func (s *server) EnviarChunk3(ctx context.Context, in *pb.DataChunkRequestDist) (*pb.DataChunkReplyDist, error) {
	log.Printf("Received: %v", in.GetFilechunk2())
	b, err := ioutil.ReadFile(in.GetFilechunk2()) // just pass the file name
       if err != nil {
           fmt.Print(err)
       }
	return &pb.DataChunkReplyDist{Bitaso2: b}, nil
}
func RyA()bool{
	//-------------DATA2--------------------------------
	conn, err := grpc.Dial(address2, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewDistribuidoClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := c.RA(ctx, &pb.RARequest{Id: int64(idnode)})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	respuestad2 := r.GetReply()
	respuestad2id := r.GetId()
	//------------------DATA3---------------------------------------------
	conn2, err2 := grpc.Dial(address3, grpc.WithInsecure(), grpc.WithBlock())
	if err2 != nil {
		log.Fatalf("did not connect: %v", err2)
	}
	defer conn2.Close()
	c2 := pb.NewDistribuidoClient(conn2)
	ctx2, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r2, err3 := c2.RA(ctx2, &pb.RARequest{Id: int64(idnode)})
	if err3 != nil {
		log.Fatalf("could not greet: %v", err3)
	}
	respuestad3 := r2.GetReply()
	respuestad3id := r2.GetId()
	//--------------------
	if (respuestad2 == "Held" || respuestad3 == "Held"){
		if (respuestad2id > int64(idnode) || respuestad3id > int64(idnode)){
			return false
		}
		return true
	}else{
		return true
	}

}
func (s *server) RA(ctx context.Context, msg *pb.RARequest) (*pb.RAReply, error) {
	if state == false {
		return &pb.RAReply{Reply: "Released" , Id: int64(idnode)}, nil
	} else{
		return &pb.RAReply{Reply: "Held" , Id: int64(idnode)}, nil
	}
}
//-----------------------AAAAAAAAAAAAAAAAAA--------------------
func (s *server) Escribir(ctx context.Context, in *pb.EscribirRequest) (*pb.EscribirReply, error) {
	return &pb.EscribirReply{Estate: "Escrito"}, nil
}
func (s *server) PedirLista2(ctx context.Context, in *pb.ListaRequest2) (*pb.ListaReply2, error) {
	fmt.Println("Enviando listado de libros")
	lista := []string{"A", "B", "C"}
	return &pb.ListaReply2{Lista2: lista}, nil
}
func (s *server) PedirChunks2(ctx context.Context, in *pb.ChunkRequest2) (*pb.ChunkReply2, error) {
	fmt.Println("Enviando locacion chunks")
	retorno := []string{"A", "B", "C"}
	return &pb.ChunkReply2{Location2: retorno}, nil
}
func main() {
	var comportamiento int
	fmt.Println("Ingrese 1 para modo centralizado")
	fmt.Println("Ingrese 2 para modo distribuida")
	fmt.Scanln(&comportamiento)
	switch comportamiento{
		case 1:
			//------------------------------------------
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