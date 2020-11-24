package main

import (
	"fmt"
	//"io/ioutil"
	"math"
	"os"
	"strconv"
	"bufio"

	"log"

	"golang.org/x/net/context"
	"google.golang.org/grpc"

	pb "github.com/GianniCarlini/Lab-2-SD/proto"
)

 func main() {
	fmt.Println("Bienvenido querido cliente")
		var tipo int
		fmt.Println("Ingrese 1 para modo centralizado")
		fmt.Println("Ingrese 2 para modo distribuida")
		fmt.Scanln(&tipo)
		switch tipo{
			case 1:
			for{
				var comportamiento int
				fmt.Println("Ingrese 1 si desea subir un archivo")
				fmt.Println("Ingrese 2 si desea descargar un archivo")
				fmt.Scanln(&comportamiento)
				switch comportamiento {
					case 1:
						conn, err := grpc.Dial("localhost:50053", grpc.WithInsecure())
						if err != nil {
							log.Fatalf("failed to connect: %s", err)
						}
						defer conn.Close()

						client := pb.NewPacketClient(conn)
						stream, err := client.EnviarLibro(context.Background())
						//-----------------------------division de archivos------------------------
						var nameLibro string
						fmt.Println("Ingrese el nombre del libro")
						fmt.Scanln(&nameLibro)
						fileToBeChunked := "./Libros/"+nameLibro+".pdf"
			
						file, err := os.Open(fileToBeChunked)
			
						if err != nil {
								fmt.Println(err)
								os.Exit(1)
						}
			
						defer file.Close()
			
						fileInfo, _ := file.Stat()
			
						var fileSize int64 = fileInfo.Size()
			
						const fileChunk = 256000 // 1 MB, change this to your requirement
			
						// calculate total number of parts the file will be chunked into
			
						totalPartsNum := uint64(math.Ceil(float64(fileSize) / float64(fileChunk)))
			
						fmt.Printf("Splitting to %d pieces.\n", totalPartsNum)
			
						for i := uint64(0); i < totalPartsNum; i++ {
								
								partSize := int(math.Min(fileChunk, float64(fileSize-int64(i*fileChunk))))
								partBuffer := make([]byte, partSize)
								file.Read(partBuffer)
								//------------------envio de chunks------------------------------------
								// write to disk
								fileName := nameLibro+"_" + strconv.FormatUint(i, 10)
								msg := &pb.EnviarLibroRequest{Id: partBuffer, Name: fileName, Numero: totalPartsNum, Libro: nameLibro}
								stream.Send(msg)
								resp, err := stream.Recv()
								if err != nil {
									log.Fatalf("can not receive %v", err)
								}
								fmt.Println("Ingresado: "+resp.Id)
						}
						stream.CloseSend()
						
					case 2:
						//-----------------------------reconstruccion de archivos------------------------
						fmt.Println("Ingrese el nombre del libro")
						fmt.Println("Aca mostrar lista")
						var nameLibro string
						fmt.Scanln(&nameLibro)

						fileToBeChunked := "./Libros/"+nameLibro+".pdf"
			
						file, err := os.Open(fileToBeChunked)
			
						if err != nil {
								fmt.Println(err)
								os.Exit(1)
						}
			
						defer file.Close()
			
						fileInfo, _ := file.Stat()
			
						var fileSize int64 = fileInfo.Size()
			
						const fileChunk = 256000 // 1 MB, change this to your requirement
			
						// calculate total number of parts the file will be chunked into
			
						totalPartsNum := uint64(math.Ceil(float64(fileSize) / float64(fileChunk)))
			
						// just for fun, let's recombine back the chunked files in a new file
			
						newFileName := "./Descargas/"+nameLibro+"Copy.pdf"
						_, err = os.Create(newFileName)
			
						if err != nil {
								fmt.Println(err)
								os.Exit(1)
						}
			
						//set the newFileName file to APPEND MODE!!
						// open files r and w
			
						file, err = os.OpenFile(newFileName, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
			
						if err != nil {
								fmt.Println(err)
								os.Exit(1)
						}
			
						// IMPORTANT! do not defer a file.Close when opening a file for APPEND mode!
						// defer file.Close()
			
						// just information on which part of the new file we are appending
						var writePosition int64 = 0
			
						for j := uint64(0); j < totalPartsNum; j++ {
			
								//read a chunk
								currentChunkFileName := nameLibro+ "_" + strconv.FormatUint(j, 10)
			
								newFileChunk, err := os.Open(currentChunkFileName)
			
								if err != nil {
										fmt.Println(err)
										os.Exit(1)
								}
			
								defer newFileChunk.Close()
			
								chunkInfo, err := newFileChunk.Stat()
			
								if err != nil {
										fmt.Println(err)
										os.Exit(1)
								}
			
								// calculate the bytes size of each chunk
								// we are not going to rely on previous data and constant
			
								var chunkSize int64 = chunkInfo.Size()
								chunkBufferBytes := make([]byte, chunkSize)
			
								fmt.Println("Appending at position : [", writePosition, "] bytes")
								writePosition = writePosition + chunkSize
			
								// read into chunkBufferBytes
								reader := bufio.NewReader(newFileChunk)
								_, err = reader.Read(chunkBufferBytes)
			
								if err != nil {
										fmt.Println(err)
										os.Exit(1)
								}
			
								// DON't USE ioutil.WriteFile -- it will overwrite the previous bytes!
								// write/save buffer to disk
								//ioutil.WriteFile(newFileName, chunkBufferBytes, os.ModeAppend)
			
								n, err := file.Write(chunkBufferBytes)
			
								if err != nil {
										fmt.Println(err)
										os.Exit(1)
								}
			
								file.Sync() //flush to disk
			
								// free up the buffer for next cycle
								// should not be a problem if the chunk size is small, but
								// can be resource hogging if the chunk size is huge.
								// also a good practice to clean up your own plate after eating
			
								chunkBufferBytes = nil // reset or empty our buffer
			
								fmt.Println("Written ", n, " bytes")
			
								fmt.Println("Recombining part [", j, "] into : ", newFileName)
						}
						fmt.Println("Libro guardado en la carpeta de Descargas")
						// now, we close the newFileName
						file.Close()
				}
			}
			case 2:
			for{
				var comportamiento int
				fmt.Println("Ingrese 1 si desea subir un archivo")
				fmt.Println("Ingrese 2 si desea descargar un archivo")
				fmt.Scanln(&comportamiento)
				switch comportamiento {
					case 1:
						conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
						if err != nil {
							log.Fatalf("failed to connect: %s", err)
						}
						defer conn.Close()

						client := pb.NewDistribuidoClient(conn)
						stream, err := client.EnviarLibro2(context.Background())
						//-----------------------------division de archivos------------------------
						var nameLibro string
						fmt.Println("Ingrese el nombre del libro")
						fmt.Scanln(&nameLibro)
						fileToBeChunked := "./Libros/"+nameLibro+".pdf"
			
						file, err := os.Open(fileToBeChunked)
			
						if err != nil {
								fmt.Println(err)
								os.Exit(1)
						}
			
						defer file.Close()
			
						fileInfo, _ := file.Stat()
			
						var fileSize int64 = fileInfo.Size()
			
						const fileChunk = 256000 // 1 MB, change this to your requirement
			
						// calculate total number of parts the file will be chunked into
			
						totalPartsNum := uint64(math.Ceil(float64(fileSize) / float64(fileChunk)))
			
						fmt.Printf("Splitting to %d pieces.\n", totalPartsNum)
			
						for i := uint64(0); i < totalPartsNum; i++ {
								
								partSize := int(math.Min(fileChunk, float64(fileSize-int64(i*fileChunk))))
								partBuffer := make([]byte, partSize)
								file.Read(partBuffer)
								//------------------envio de chunks------------------------------------
								// write to disk
								fileName := nameLibro+"_" + strconv.FormatUint(i, 10)
								msg := &pb.EnviarLibroRequest2{Id: partBuffer, Name: fileName}
								stream.Send(msg)
								resp, err := stream.Recv()
								if err != nil {
									log.Fatalf("can not receive %v", err)
								}
								fmt.Println("Ingresado: "+resp.Id)
						}
						stream.CloseSend()
						
					case 2:
						//-----------------------------reconstruccion de archivos------------------------
						fmt.Println("Ingrese el nombre del libro")
						fmt.Println("Aca mostrar lista")
						var nameLibro string
						fmt.Scanln(&nameLibro)

						fileToBeChunked := "./Libros/"+nameLibro+".pdf"
			
						file, err := os.Open(fileToBeChunked)
			
						if err != nil {
								fmt.Println(err)
								os.Exit(1)
						}
			
						defer file.Close()
			
						fileInfo, _ := file.Stat()
			
						var fileSize int64 = fileInfo.Size()
			
						const fileChunk = 256000 // 1 MB, change this to your requirement
			
						// calculate total number of parts the file will be chunked into
			
						totalPartsNum := uint64(math.Ceil(float64(fileSize) / float64(fileChunk)))
			
						// just for fun, let's recombine back the chunked files in a new file
			
						newFileName := "./Descargas/"+nameLibro+"Copy.pdf"
						_, err = os.Create(newFileName)
			
						if err != nil {
								fmt.Println(err)
								os.Exit(1)
						}
			
						//set the newFileName file to APPEND MODE!!
						// open files r and w
			
						file, err = os.OpenFile(newFileName, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
			
						if err != nil {
								fmt.Println(err)
								os.Exit(1)
						}
			
						// IMPORTANT! do not defer a file.Close when opening a file for APPEND mode!
						// defer file.Close()
			
						// just information on which part of the new file we are appending
						var writePosition int64 = 0
			
						for j := uint64(0); j < totalPartsNum; j++ {
			
								//read a chunk
								currentChunkFileName := nameLibro+ "_" + strconv.FormatUint(j, 10)
			
								newFileChunk, err := os.Open(currentChunkFileName)
			
								if err != nil {
										fmt.Println(err)
										os.Exit(1)
								}
			
								defer newFileChunk.Close()
			
								chunkInfo, err := newFileChunk.Stat()
			
								if err != nil {
										fmt.Println(err)
										os.Exit(1)
								}
			
								// calculate the bytes size of each chunk
								// we are not going to rely on previous data and constant
			
								var chunkSize int64 = chunkInfo.Size()
								chunkBufferBytes := make([]byte, chunkSize)
			
								fmt.Println("Appending at position : [", writePosition, "] bytes")
								writePosition = writePosition + chunkSize
			
								// read into chunkBufferBytes
								reader := bufio.NewReader(newFileChunk)
								_, err = reader.Read(chunkBufferBytes)
			
								if err != nil {
										fmt.Println(err)
										os.Exit(1)
								}
			
								// DON't USE ioutil.WriteFile -- it will overwrite the previous bytes!
								// write/save buffer to disk
								//ioutil.WriteFile(newFileName, chunkBufferBytes, os.ModeAppend)
			
								n, err := file.Write(chunkBufferBytes)
			
								if err != nil {
										fmt.Println(err)
										os.Exit(1)
								}
			
								file.Sync() //flush to disk
			
								// free up the buffer for next cycle
								// should not be a problem if the chunk size is small, but
								// can be resource hogging if the chunk size is huge.
								// also a good practice to clean up your own plate after eating
			
								chunkBufferBytes = nil // reset or empty our buffer
			
								fmt.Println("Written ", n, " bytes")
			
								fmt.Println("Recombining part [", j, "] into : ", newFileName)
						}
						fmt.Println("Libro guardado en la carpeta de Descargas")
						// now, we close the newFileName
						file.Close()
				}
			}
		}

	}
