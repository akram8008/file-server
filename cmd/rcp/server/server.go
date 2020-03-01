package main

import (
	"bufio"
	"fileServer/pkg/rcp"
	"io"
	"io/ioutil"
	"log"
	"net"
	"os"
	"strings"
)

func main()  {
	file, err := os.Create("server-log.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		err := file.Close()
		if err != nil {
			log.Printf("Can't close file: %v", err)
		}
	}()
	log.SetOutput(file)
	log.Print("server starting")
	host := "0.0.0.0"
	port, ok := os.LookupEnv("PORT")
	if !ok {
		port = "9999"
	}
	err = start(host + ":" + port)
	if err != nil {
		log.Fatal(err)
	}
}


func start(addr string) (err error) {

	listener, err := net.Listen(rcp.Tcp, addr)
	if err != nil {
		log.Fatalf("can't listen %s: %v", addr, err)
		return err
	}
	defer func() {
		err := listener.Close()
		if err != nil {
			log.Fatalf("Can't close conn: %v", err)
		}
	}()
	for {
		conn, err := listener.Accept()
		log.Print("accept connection")
		if err != nil {
			log.Fatalf("can't accept: %v", err)
		}
		log.Print("handle connection")
		go handleConn(conn)
	}
}

func handleConn(conn net.Conn) {
	defer func() {
		err := conn.Close()
		if err != nil {
			log.Fatalf("Can't close conn: %v", err)
		}
	}()
	reader := bufio.NewReader(conn)
	line, err := rcp.ReadLine(reader)
	if err != nil {
		log.Fatalf("error while reading: %v", err)
		return
	}
	index := strings.IndexByte(line, ':')
	writer := bufio.NewWriter(conn)
	if index == -1 {
		log.Printf("invalid line received %s", line)
		err := rcp.WriteLine("error: invalid line", writer)
		if err != nil {
			log.Printf("error while writing: %v", err)
			return
		}
		return
	}
	cmd, options := line[:index], line[index+1:]
	log.Printf("command received: %s", cmd)
	log.Printf("options received: %s", options)
	switch cmd {
	case rcp.Upload:
		uploadFromBuffer (options, reader, writer)
	case rcp.DownLoad:
		downloadToBuffer (options, reader, writer)
	case rcp.List:
		listFileToBuffer (options, reader, writer)
	default:
		err := rcp.WriteLine(rcp.ForError, writer)
		if err != nil {
			log.Printf("error while writing: %v", err)
			return
		}
	}
}

func uploadFromBuffer (options string, reader *bufio.Reader, writer *bufio.Writer){
	options = strings.TrimSuffix(options, rcp.Endl)
	line, err := rcp.ReadLine(reader)
	if err != nil {
		log.Printf("can't read: %v", err)
		return
	}
	if line == rcp.ForError + rcp.Endl {
		log.Printf("file not such: %v", err)
		return
	}
	bytes, err := ioutil.ReadAll(reader)
	if err != nil {
		if err != io.EOF {
			log.Printf("can't read data: %v", err)
			return
		}
	}
	err = ioutil.WriteFile(rcp.PathFileServer+options, bytes, 0666)
	if err != nil {
		log.Printf("can't write file: %v", err)
		return
	}
	err = rcp.WriteLine(rcp.CheckOk, writer)
	if err != nil {
		log.Printf("error while writing: %v", err)
		return
	}
}

func downloadToBuffer (options string, reader *bufio.Reader, writer *bufio.Writer) {
	options = strings.TrimSuffix(options, rcp.Endl)
	file, err := os.Open(rcp.PathFileServer + options)

	if err != nil {
		log.Printf("file does not exist %v ", rcp.PathFileServer + options)
		err = rcp.WriteLine(rcp.ForError, writer)
		return
	}
	err = rcp.WriteLine(rcp.CheckOk, writer)
	if err != nil {
		log.Printf("error while writing: %v", err)
		return
	}
	_, err = io.Copy(writer, file)
	err = writer.Flush()
	if err != nil {
		log.Printf("Can't flush: %v", err)
		return
	}
}

func listFileToBuffer (options string, reader *bufio.Reader, writer *bufio.Writer) {
	options = strings.TrimSuffix(options, rcp.Endl)
	fileName := rcp.ReadDir(rcp.PathFileServer)
	err := rcp.WriteLine(fileName, writer)
	if err != nil {
		log.Printf("error while writing: %v", err)
	}
}