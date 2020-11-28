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
	port = ":50052"
)

var lista []string

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

func (s *server) EnviarPropuestaCentralizado(ctx context.Context, in *pb.PropuestaRequestC) (*pb.PropuestaReplyC, error) {

	nombre := in.GetNombre()
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
			
			content = append(content, []byte(p1[i]+" ip1"+"\n")...)

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
			
			content = append(content, []byte(p2[i]+" ip2"+"\n")...)
			
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
			
			content = append(content, []byte(p3[i]+" ip3"+"\n")...)

			err = ioutil.WriteFile("log.txt", content, 0644)
			if err != nil {
				log.Fatal(err)
			}
		}
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
				fmt.Println(d2)
				fmt.Println(d3)
				for i := range d1{
					content, err := ioutil.ReadFile("log.txt") // just pass the file name
					if err != nil {
						fmt.Print(err)
					}
					
					content = append(content, []byte(d1[i]+" ip1"+"\n")...)
		
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
					
					content = append(content, []byte(d2[i]+" ip2"+"\n")...)
					
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
					
					content = append(content, []byte(d3[i]+" ip3"+"\n")...)
		
					err = ioutil.WriteFile("log.txt", content, 0644)
					if err != nil {
						log.Fatal(err)
					}
				}
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
