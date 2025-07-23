package main

import (
	"context"
	"flag"
	"fmt"
	"goexamples/poem-stream/proto"
	"goexamples/poem-stream/testdata"
	"io"
	"log"
	"net"
	"os"
	"path/filepath"
	"sync"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
)

type Server struct {
	db testdata.DB
	mu sync.Mutex
	proto.UnimplementedPoemServiceServer
}

func (s *Server) SetDB(db testdata.DB) {
	s.db = db
}

func (s *Server) Start(port int) error {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	server := grpc.NewServer()
	proto.RegisterPoemServiceServer(server, s)
	log.Printf("server listening at %v", lis.Addr())
	return server.Serve(lis)
}

func (s *Server) GetPoem(_ context.Context, in *proto.GetPoemRequest) (*proto.Poem, error) {
	return s.db.GetPoem(in.GetTitle())
}

func (s *Server) GetPoemStream(in *proto.GetPoemRequest, sout grpc.ServerStreamingServer[proto.StreamPoem]) error {
	var poem *proto.Poem
	if p, err := s.db.GetPoem(in.GetTitle()); err != nil {
		return err
	} else {
		poem = p
	}

	if err := sout.Send(&proto.StreamPoem{OneOf: &proto.StreamPoem_Title{Title: poem.GetTitle()}}); err != nil {
		return err
	}
	if err := sout.Send(&proto.StreamPoem{OneOf: &proto.StreamPoem_Author{Author: poem.GetAuthor()}}); err != nil {
		return err
	}
	for _, content := range poem.GetContents() {
		if err := sout.Send(&proto.StreamPoem{OneOf: &proto.StreamPoem_Content{Content: content}}); err != nil {
			return err
		}
	}
	return nil
}

func (s *Server) GetPoemAll(_ context.Context, _ *emptypb.Empty) (*proto.PoemCollection, error) {
	return &proto.PoemCollection{Value: s.db.GetPoemCollection()}, nil
}

func (s *Server) GetPoemAllStream(_ *emptypb.Empty, sout grpc.ServerStreamingServer[proto.Poem]) error {
	for _, p := range s.db.GetPoemCollection() {
		if err := sout.Send(p); err != nil {
			return err
		}
	}
	return nil
}

func (s *Server) UploadPoem(_ context.Context, in *proto.Poem) (*proto.UploadPoemResponse, error) {
	s.db.SetPoem(in.GetTitle(), in)
	log.Printf("uploaded poem: %s\n", in.GetTitle())
	return &proto.UploadPoemResponse{EndTime: time.Now().Format(time.DateTime), Success: true, Data: []*proto.Poem{in}}, nil
}

func (s *Server) UploadPoemStream(sin grpc.ClientStreamingServer[proto.StreamPoem, proto.UploadPoemResponse]) error {
	poem := new(proto.Poem)
	for {
		in, err := sin.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		switch in.GetOneOf().(type) {
		case *proto.StreamPoem_Title:
			poem.Title = in.GetTitle()
		case *proto.StreamPoem_Author:
			poem.Author = in.GetAuthor()
		case *proto.StreamPoem_Content:
			poem.Contents = append(poem.Contents, in.GetContent())
		}
	}
	s.db.SetPoem(poem.GetTitle(), poem)
	log.Printf("uploaded poem: %s\n", poem.GetTitle())
	return sin.SendAndClose(&proto.UploadPoemResponse{EndTime: time.Now().Format(time.DateTime), Success: true, Data: []*proto.Poem{poem}})
}

func (s *Server) BatchUploadPoem(_ context.Context, in *proto.PoemCollection) (*proto.UploadPoemResponse, error) {
	for _, p := range in.GetValue() {
		s.db.SetPoem(p.GetTitle(), p)
		log.Printf("uploaded poem: %s\n", p.GetTitle())
	}
	return &proto.UploadPoemResponse{EndTime: time.Now().Format(time.DateTime), Success: true, Data: in.GetValue()}, nil
}

func (s *Server) BatchUploadPoemStream(stream grpc.BidiStreamingServer[proto.Poem, proto.UploadPoemResponse]) error {
	for {
		in, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}
		s.mu.Lock()
		s.db.SetPoem(in.GetTitle(), in)
		log.Printf("uploaded poem: %s\n", in.GetTitle())
		s.mu.Unlock()
		if err := stream.Send(&proto.UploadPoemResponse{EndTime: time.Now().Format(time.DateTime), Success: true, Data: []*proto.Poem{in}}); err != nil {
			return err
		}
	}
}

func NewServer(port int) *Server {
	return new(Server)
}

var (
	port     = flag.Int("port", 50051, "port to listen on")
	jsonFile = flag.String("json_file", "", "server poem json file")
)

func main() {
	flag.Parse()
	if *jsonFile == "" {
		if file, err := os.Getwd(); err != nil {
			log.Fatalf("failed to get work dir: %v", err)
		} else {
			*jsonFile = filepath.Join(file, "testdata", "server_poem.json")
		}
	}

	s := NewServer(*port)
	s.SetDB(testdata.NewDB(*jsonFile))
	if err := s.Start(*port); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
