package main

import (
	"context"
	"flag"
	"fmt"
	"goexamples/poem-stream/pb"
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
	pb.UnimplementedPoemServiceServer
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
	pb.RegisterPoemServiceServer(server, s)
	log.Printf("server listening at %v", lis.Addr())
	return server.Serve(lis)
}

func (s *Server) GetPoem(_ context.Context, in *pb.GetPoemRequest) (*pb.Poem, error) {
	return s.db.GetPoem(in.GetTitle())
}

func (s *Server) GetPoemStream(in *pb.GetPoemRequest, sout grpc.ServerStreamingServer[pb.StreamPoem]) error {
	var poem *pb.Poem
	if p, err := s.db.GetPoem(in.GetTitle()); err != nil {
		return err
	} else {
		poem = p
	}

	if err := sout.Send(&pb.StreamPoem{OneOf: &pb.StreamPoem_Title{Title: poem.GetTitle()}}); err != nil {
		return err
	}
	if err := sout.Send(&pb.StreamPoem{OneOf: &pb.StreamPoem_Author{Author: poem.GetAuthor()}}); err != nil {
		return err
	}
	for _, content := range poem.GetContents() {
		if err := sout.Send(&pb.StreamPoem{OneOf: &pb.StreamPoem_Content{Content: content}}); err != nil {
			return err
		}
	}
	return nil
}

func (s *Server) GetPoemAll(_ context.Context, _ *emptypb.Empty) (*pb.PoemCollection, error) {
	return &pb.PoemCollection{Value: s.db.GetPoemCollection()}, nil
}

func (s *Server) GetPoemAllStream(_ *emptypb.Empty, sout grpc.ServerStreamingServer[pb.Poem]) error {
	for _, p := range s.db.GetPoemCollection() {
		if err := sout.Send(p); err != nil {
			return err
		}
	}
	return nil
}

func (s *Server) UploadPoem(_ context.Context, in *pb.Poem) (*pb.UploadPoemResponse, error) {
	s.db.SetPoem(in.GetTitle(), in)
	log.Printf("uploaded poem: %s\n", in.GetTitle())
	return &pb.UploadPoemResponse{EndTime: time.Now().Format(time.DateTime), Success: true, Data: []*pb.Poem{in}}, nil
}

func (s *Server) UploadPoemStream(sin grpc.ClientStreamingServer[pb.StreamPoem, pb.UploadPoemResponse]) error {
	poem := new(pb.Poem)
	for {
		in, err := sin.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		switch in.GetOneOf().(type) {
		case *pb.StreamPoem_Title:
			poem.Title = in.GetTitle()
		case *pb.StreamPoem_Author:
			poem.Author = in.GetAuthor()
		case *pb.StreamPoem_Content:
			poem.Contents = append(poem.Contents, in.GetContent())
		}
	}
	s.db.SetPoem(poem.GetTitle(), poem)
	log.Printf("uploaded poem: %s\n", poem.GetTitle())
	return sin.SendAndClose(&pb.UploadPoemResponse{EndTime: time.Now().Format(time.DateTime), Success: true, Data: []*pb.Poem{poem}})
}

func (s *Server) BatchUploadPoem(_ context.Context, in *pb.PoemCollection) (*pb.UploadPoemResponse, error) {
	for _, p := range in.GetValue() {
		s.db.SetPoem(p.GetTitle(), p)
		log.Printf("uploaded poem: %s\n", p.GetTitle())
	}
	return &pb.UploadPoemResponse{EndTime: time.Now().Format(time.DateTime), Success: true, Data: in.GetValue()}, nil
}

func (s *Server) BatchUploadPoemStream(stream grpc.BidiStreamingServer[pb.Poem, pb.UploadPoemResponse]) error {
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
		if err := stream.Send(&pb.UploadPoemResponse{EndTime: time.Now().Format(time.DateTime), Success: true, Data: []*pb.Poem{in}}); err != nil {
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
