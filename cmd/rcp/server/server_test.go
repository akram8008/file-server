package main

import (
	"bufio"
	bytes2 "bytes"
	"fileServer/pkg/rcp"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net"
	"testing"
	"time"
)

func Test_DownloadInServerOk(t *testing.T) {
	host := "localhost"
	port := rand.Intn(999) + 9000
	addr := fmt.Sprintf("%s:%d", host, port)
	go func() {
		err := start(addr)
		if err != nil {
			t.Fatalf("can't start server: %v", err)
		}
	}()
	time.Sleep(rcp.TimeSleep)
	conn, err := net.Dial(rcp.Tcp, addr)
	if err != nil {
		t.Fatalf("can't connect to server: %v", err)
	}
	writer := bufio.NewWriter(conn)
	options := "test.txt"
	line := rcp.DownLoad + ":" + options
	err = rcp.WriteLine(line, writer)
	if err != nil {
		t.Fatalf("can't send command %s to server: %v", line, err)
	}
	reader := bufio.NewReader(conn)
	line, err = rcp.ReadLine(reader)
	src, err := ioutil.ReadFile( rcp.PathTestData + "test.txt")
	if err != nil {
		log.Fatalf("Can't read file: %v",err)
	}
	dst, err := ioutil.ReadFile(rcp.PathFileServer + options)
	if err != nil {
			log.Fatalf("can't Read file: %v",err)
	}
	if !bytes2.Equal(src, dst) {
		t.Fatalf("files are not equal: %v", err)
	}
}

func Test_UploadToServerOk(t *testing.T) {
	host := "localhost"
	port := rand.Intn(999) + 9000
	addr := fmt.Sprintf("%s:%d", host, port)
	go func() {
		err := start(addr)
		if err != nil {
			t.Fatalf("can't start server: %v", err)
		}
	}()
	time.Sleep(rcp.TimeSleep)
	conn, err := net.Dial(rcp.Tcp, addr)
	if err != nil {
		t.Fatalf("can't connect to server: %v", err)
	}
	writer := bufio.NewWriter(conn)
	options := "test.txt"
	line := rcp.Upload + ":" + options
	err = rcp.WriteLine(line, writer)
	if err != nil {
		t.Fatalf("can't send command %s to server: %v", line, err)
	}
	src, err := ioutil.ReadFile(rcp.PathTestData + "test.txt")
	if err != nil {
		log.Fatalf("Can't read file: %v",err)
	}
	_, err = writer.Write(src)
	if err != nil {
		log.Fatalf("Can't write: %v", err)
	}
	err = writer.Flush()
	if err != nil {
		log.Fatalf("Can't flush: %v", err)
	}
	err = conn.Close()
	if err != nil {
		log.Fatalf("Can't close conn: %v", err)
	}
	dst, err := ioutil.ReadFile(rcp.PathFileServer + "test.txt")
	if err != nil {
		log.Fatalf("can't Read file: %v",err)
	}
	if !bytes2.Equal(src, dst) {
		t.Fatalf("files are not equal: %v", err)
	}
}

func Test_ListInServerOk(t *testing.T)  {
	host := "localhost"
	port := rand.Intn(999) + 9000
	addr := fmt.Sprintf("%s:%d", host, port)
	go func() {
		err := start(addr)
		if err != nil {
			t.Fatalf("can't start server: %v", err)
		}
	}()
	time.Sleep(rcp.TimeSleep)
	conn, err := net.Dial(rcp.Tcp, addr)
	if err != nil {
		t.Fatalf("can't connect to server: %v", err)
	}
	writer := bufio.NewWriter(conn)
	options := ""
	line := rcp.List + ":" + options
	err = rcp.WriteLine(line, writer)
	if err != nil {
		t.Fatalf("can't send command %s to server: %v", line, err)
	}
	reader := bufio.NewReader(conn)
	line, err = rcp.ReadLine(reader)
	if line != "index.html test.txt\n" {
		t.Fatalf("result not ok: %s %v", line, err)
	}
}