package main

import (
	"context"
	"log"
	"net"
	"fmt"
	"time"
	"math/rand"
	"io/ioutil"
	"strconv"

	pb "github.com/GianniCarlini/Lab-2-SD/proto"
	"google.golang.org/grpc"
)

const (
	ip2 = "localhost:50054" //data2
	ip3 = "localhost:50053" //data 3
	ip1 = "localhost:50051" //data1
	port = ":50052"
)

var lista []string
var datachunks [][]string
var state = false
var idnode = 1000

type server struct {
}
func existeEnArreglo(arreglo []string, busqueda string) bool {
	for _, numero := range arreglo {
		if numero == busqueda {
			return true
		}
	}
	return false
}
/*func buscarArreglo(nombre string, arreglo [][]string) []string {
	var retorno []string
	for _, numero := range arreglo {
		if numero[0] == nombre {
			retorno = numero
		}
	}
	return retorno
}
*/
func (s *server) EnviarPropuestaCentralizado(ctx context.Context, in *pb.PropuestaRequestC) (*pb.PropuestaReplyC, error) {
	var datachunk []string
	nombre := in.GetNombre()
	datachunk = append(datachunk, nombre)
	if !(existeEnArreglo(lista, nombre)){
		lista = append(lista,nombre)
	}
	fmt.Println(lista)
	p1 := in.GetPropuesta1()
	p2 := in.GetPropuesta2()
	p3 := in.GetPropuesta3()
	largo := len(in.GetPropuesta())

	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)
	aceptar := r1.Intn(101)

	content, err := ioutil.ReadFile("log.txt") // just pass the file name
	if err != nil {
		fmt.Print(err)
	}
	
	content = append(content, []byte(in.GetNombre()+" "+strconv.Itoa(largo)+"\n")...)
    err = ioutil.WriteFile("log.txt", content, 0644)
    if err != nil {
        log.Fatal(err)
	}	

	if aceptar <= 50{
		fmt.Println("Propuesta inicial aceptada")
		for i := range p1{
			content, err := ioutil.ReadFile("log.txt") // just pass the file name
			if err != nil {
				fmt.Print(err)
			}
			
			content = append(content, []byte(p1[i]+" "+ip1+"\n")...)
			datachunk = append(datachunk, p1[i]+","+ip1)

			err = ioutil.WriteFile("log.txt", content, 0644)
			if err != nil {
				log.Fatal(err)
			}
		}
		for i := range p2{
			content, err := ioutil.ReadFile("log.txt") // just pass the file name
			if err != nil {
				fmt.Print(err)
			}
			
			content = append(content, []byte(p2[i]+" "+ip2+"\n")...)
			datachunk = append(datachunk, p2[i]+","+ip2)
			
			err = ioutil.WriteFile("log.txt", content, 0644)
			if err != nil {
				log.Fatal(err)
			}
		}
		for i := range p3{
			content, err := ioutil.ReadFile("log.txt") // just pass the file name
			if err != nil {
				fmt.Print(err)
			}
			
			content = append(content, []byte(p3[i]+" "+ip3+"\n")...)
			datachunk = append(datachunk, p3[i]+","+ip3)

			err = ioutil.WriteFile("log.txt", content, 0644)
			if err != nil {
				log.Fatal(err)
			}
			
		}
		datachunks = append(datachunks, datachunk)
		return &pb.PropuestaReplyC{Distribucion1: p1,Distribucion2: p2,Distribucion3: p3, Tipo: 1}, nil
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
				d2 = in.GetPropuesta()[(len(in.GetPropuesta())/3):(2*len(in.GetPropuesta())/3)]
				d3 = in.GetPropuesta()[(2*len(in.GetPropuesta())/3):]
				fmt.Println(d1)
				fmt.Println(d2)
				fmt.Println(d3)
				for i := range d1{
					content, err := ioutil.ReadFile("log.txt") // just pass the file name
					if err != nil {
						fmt.Print(err)
					}
					
					content = append(content, []byte(d1[i]+" "+ip1+"\n")...)
					datachunk = append(datachunk, d1[i]+","+ip1)
					err = ioutil.WriteFile("log.txt", content, 0644)
					if err != nil {
						log.Fatal(err)
					}
				}
				for i := range d2{
					content, err := ioutil.ReadFile("log.txt") // just pass the file name
					if err != nil {
						fmt.Print(err)
					}
					
					content = append(content, []byte(d2[i]+" "+ip2+"\n")...)
					datachunk = append(datachunk, d2[i]+","+ip2)
					
					err = ioutil.WriteFile("log.txt", content, 0644)
					if err != nil {
						log.Fatal(err)
					}
				}
				for i := range d3{
					content, err := ioutil.ReadFile("log.txt") // just pass the file name
					if err != nil {
						fmt.Print(err)
					}
					
					content = append(content, []byte(d3[i]+" "+ip3+"\n")...)
					datachunk = append(datachunk, d3[i]+","+ip3)
		
					err = ioutil.WriteFile("log.txt", content, 0644)
					if err != nil {
						log.Fatal(err)
					}
				}
				datachunks = append(datachunks, datachunk)
				return &pb.PropuestaReplyC{Distribucion1: d1,Distribucion2: d2,Distribucion3: d3, Tipo: 2}, nil
			}else{
				fmt.Println("Propuesta rechazada")
				continue
			}
		}
	}
	
}
func (s *server) PedirLista(ctx context.Context, in *pb.ListaRequest) (*pb.ListaReply, error) {
	fmt.Println("Enviando listado de libros")
	return &pb.ListaReply{Lista: lista}, nil
}
func (s *server) PedirLista2(ctx context.Context, in *pb.ListaRequest2) (*pb.ListaReply2, error) {
	fmt.Println("Enviando listado de libros")
	return &pb.ListaReply2{Lista2: lista}, nil
}
func (s *server) PedirChunks(ctx context.Context, in *pb.ChunkRequest) (*pb.ChunkReply, error) {
	fmt.Println("Enviando locacion chunks")
	var retorno []string
	for _,numero := range datachunks {
		if numero[0] == in.GetBook() {
			retorno = numero
		}
	}
	return &pb.ChunkReply{Location: retorno}, nil
}

func (s *server) PedirChunks2(ctx context.Context, in *pb.ChunkRequest2) (*pb.ChunkReply2, error) {
	fmt.Println("Enviando locacion chunks")
	var retorno []string
	for _,numero := range datachunks {
		if numero[0] == in.GetBook2() {
			retorno = numero
		}
	}
	return &pb.ChunkReply2{Location2: retorno}, nil
}
//------------------AAAAAAAA----------------------------
func (s *server) EnviarChunk3(ctx context.Context, in *pb.DataChunkRequestDist) (*pb.DataChunkReplyDist, error) {
	log.Printf("Received: %v", in.GetFilechunk2())
	b, err := ioutil.ReadFile(in.GetFilechunk2()) // just pass the file name
       if err != nil {
           fmt.Print(err)
       }
	return &pb.DataChunkReplyDist{Bitaso2: b}, nil
}
func (s *server) Escribir(ctx context.Context, in *pb.EscribirRequest) (*pb.EscribirReply, error) {
	fmt.Println("Ingresando al log")
	var datachunk []string
	nombre := in.GetNombre()
	datachunk = append(datachunk, nombre)
	largo := len(in.GetPropuesta())
	if !(existeEnArreglo(lista, nombre)){
		lista = append(lista,nombre)
	}
	//----------nombre largo---------------------
	content, err := ioutil.ReadFile("log.txt") // just pass the file name
	if err != nil {
		fmt.Print(err)
	}
	content = append(content, []byte(in.GetNombre()+" "+strconv.Itoa(largo)+"\n")...)
    err = ioutil.WriteFile("log.txt", content, 0644)
    if err != nil {
        log.Fatal(err)
	}
	//------------distribucion-------------------------
	d1 := in.GetPropuesta1()
	d2 := in.GetPropuesta2()
	d3 := in.GetPropuesta3()

	for i := range d1{
		content, err := ioutil.ReadFile("log.txt") // just pass the file name
		if err != nil {
			fmt.Print(err)
		}
		
		content = append(content, []byte(d1[i]+" "+ip1+"\n")...)
		datachunk = append(datachunk, d1[i]+","+ip1)
		err = ioutil.WriteFile("log.txt", content, 0644)
		if err != nil {
			log.Fatal(err)
		}
	}
	for i := range d2{
		content, err := ioutil.ReadFile("log.txt") // just pass the file name
		if err != nil {
			fmt.Print(err)
		}
		
		content = append(content, []byte(d2[i]+" "+ip2+"\n")...)
		datachunk = append(datachunk, d2[i]+","+ip2)
		
		err = ioutil.WriteFile("log.txt", content, 0644)
		if err != nil {
			log.Fatal(err)
		}
	}
	for i := range d3{
		content, err := ioutil.ReadFile("log.txt") // just pass the file name
		if err != nil {
			fmt.Print(err)
		}
		
		content = append(content, []byte(d3[i]+" "+ip3+"\n")...)
		datachunk = append(datachunk, d3[i]+","+ip3)

		err = ioutil.WriteFile("log.txt", content, 0644)
		if err != nil {
			log.Fatal(err)
		}
	}
	datachunks = append(datachunks, datachunk)
	fmt.Println(datachunks)
	return &pb.EscribirReply{Estate: "Escrito"}, nil
}
//-----------------------AAAAAAAAAAAAAAAAAA--------------------
func RyA()bool{
	//-------------DATA2--------------------------------
	respuestad2 := "A"
	respuestad2id := int64(1)
	//------------------DATA3---------------------------------------------
	respuestad3 := "A"
	respuestad3id := int64(1)
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
func (s *server) EnviarLibro2(stream pb.Distribuido_EnviarLibro2Server) error {
	log.Println("Started stream Distribuido")
	for {
		in, err := stream.Recv()
		if err != nil {
			return err
		}
		resp := pb.EnviarLibroReply2{Id: in.Name}
		if err := stream.Send(&resp); err != nil { 
			log.Printf("send error %v", err)
		}

	}
}
func (s *server) EnviarPropuestaDistribuido(ctx context.Context, in *pb.DistribuidoRequest) (*pb.DistribuidoReply, error) {
	log.Printf("Received: %v", in.GetPropuestaini())
	return &pb.DistribuidoReply{Respuesta: "ok"}, nil
}

//-----------------------AAAAAAAAAAAAAAAAAA--------------------
func main() {
	var tipo int 
	fmt.Println("Ingrese 1 para modo centralizado")
	fmt.Println("Ingrese 2 para modo distribuida")
	fmt.Scanln(&tipo)
	switch tipo{
		case 1:
			lis, err := net.Listen("tcp", port)
			if err != nil {
				log.Fatalf("failed to listen: %v", err)
			}
			s := grpc.NewServer()
			pb.RegisterPropuestaCentralizadoServer(s, &server{})
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
