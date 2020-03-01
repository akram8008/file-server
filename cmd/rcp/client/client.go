package main

import (
	"bufio"
	"fileServer/pkg/rcp"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"os"
	"strings"
)

var download = flag.String("download", "default", "Download")
var upload = flag.String("upload", "default", "Upload")
var list = flag.Bool("list", false, "List")

func main() {
	file, err := os.Create("client-log.txt")
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
	flag.Parse()
	var cmd, fileName string
	if *download != "default" {
		fileName = *download
		cmd = rcp.DownLoad
	} else if *upload != "default" {
		cmd = rcp.Upload
		fileName = *upload
	} else if *list != false {
		cmd = rcp.List
		fileName = ""
	} else{
		return}
	StartingOperationsLoop(cmd, fileName)
}

func StartingOperationsLoop(cmd string, fileName string)  {
	log.Print("client connecting")
	conn, err := net.Dial(rcp.Tcp, rcp.AddressClient)
	if err != nil {
		log.Fatalf("can't connect to %s: %v", rcp.AddressClient, err)
	}
	defer func() {
		err := conn.Close()
		if err != nil {
			log.Printf("Can't close conn: %v", err)
		}
	}()
	log.Print("client connected")
	writer := bufio.NewWriter(conn)
	line := cmd + ":" + fileName
	log.Print("command sending")
	err = rcp.WriteLine(line, writer)
	if err != nil {
		log.Fatalf("can't send command %s to server: %v", line, err)
	}
	log.Print("command sent")
	switch cmd {
	case rcp.DownLoad:
		downloadFromServer(conn, fileName)
	case rcp.Upload:
		uploadInServer(conn, fileName)
	case rcp.List:
		listFile(conn)
	default:
		fmt.Printf("Вы выбрали неверную команду: %s\n", cmd)
	}
	return
}

func downloadFromServer(conn net.Conn, fileName string) {
	reader := bufio.NewReader(conn)
	line, err := rcp.ReadLine(reader)
	if err != nil {
		log.Printf("can't read: %v", err)
		return
	}
	if line == rcp.ForError + rcp.Endl {
		log.Printf("file not such: %v", err)
		fmt.Printf("Файл с название %s на сервере не существует\n", fileName)
		return
	}
	log.Print(line)
	bytes, err := ioutil.ReadAll(reader) // while not EOF
	if err != nil {
		if err != io.EOF {
			log.Printf("can't read data: %v", err)
		}
	}
	log.Print(len(bytes))
	err = ioutil.WriteFile(rcp.PathFileClient + fileName, bytes, 0666)
	if err != nil {
		log.Printf("can't write file: %v", err)
	}
	fmt.Printf("Файл с название %s успешно скаченно\n", fileName)
}

func uploadInServer(conn net.Conn, fileName string) {
	options := strings.TrimSuffix(fileName, rcp.Endl)
	file, err := os.Open(rcp.PathFileClient + options)
	writer := bufio.NewWriter(conn)
	if err != nil {
		log.Print("file does not exist")
		err = rcp.WriteLine(rcp.ForError, writer)
		fmt.Printf("Файл с название %s не существует\n", fileName)
		return
	}
	err = rcp.WriteLine(rcp.CheckOk, writer)
	if err != nil {
		log.Printf("error while writing: %v", err)
		return
	}
	log.Print(fileName)
	fileByte, err := io.Copy(writer, file)
	log.Print(fileByte)
	err = writer.Flush()
	if err != nil {
		log.Printf("Can't flush: %v", err)
	}
	fmt.Printf("Файл с название %s успешно загруженно на сервер\n", fileName)
}


func listFile(conn net.Conn) {
	reader := bufio.NewReader(conn)
	line, err := rcp.ReadLine(reader)
	if err != nil {
		log.Printf("can't read: %v", err)
		return
	}
	fmt.Printf("Список доступных файлов в сервере:\n")

	for _, val := range line {
		if string(val) == " " {
			fmt.Println()
		} else {
			fmt.Print(string(val))
		}
	}
}