package main

import (
	"context"
	"flag"
	"fmt"
	"goexamples/poem-stream/proto"
	"goexamples/poem-stream/testdata"
	"io"
	"log"
	"os"
	"path/filepath"
	"sync"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/emptypb"
)

type Client struct {
	conn   *grpc.ClientConn
	client proto.PoemServiceClient
}

func (c *Client) GetPoem(ctx context.Context, in *proto.GetPoemRequest, opts ...grpc.CallOption) (*proto.Poem, error) {
	return c.client.GetPoem(ctx, in, opts...)
}

func (c *Client) GetPoemStream(ctx context.Context, in *proto.GetPoemRequest, opts ...grpc.CallOption) (*proto.Poem, error) {
	sout, err := c.client.GetPoemStream(ctx, in, opts...)
	if err != nil {
		return nil, err
	}

	p := new(proto.Poem)
	for {
		r, err := sout.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		switch r.OneOf.(type) {
		case *proto.StreamPoem_Title:
			p.Title = r.OneOf.(*proto.StreamPoem_Title).Title
		case *proto.StreamPoem_Author:
			p.Author = r.OneOf.(*proto.StreamPoem_Author).Author
		case *proto.StreamPoem_Content:
			p.Contents = append(p.GetContents(), r.OneOf.(*proto.StreamPoem_Content).Content)
		}
	}
	return p, nil
}

func (c *Client) GetPoemAll(ctx context.Context, opts ...grpc.CallOption) ([]*proto.Poem, error) {
	if r, err := c.client.GetPoemAll(ctx, new(emptypb.Empty), opts...); err != nil {
		return nil, err
	} else {
		return r.GetValue(), nil
	}
}

func (c *Client) GetPoemAllStream(ctx context.Context, opts ...grpc.CallOption) ([]*proto.Poem, error) {
	sin, err := c.client.GetPoemAllStream(ctx, new(emptypb.Empty), opts...)
	if err != nil {
		return nil, err
	}
	poems := []*proto.Poem{}
	for {
		r, err := sin.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		poems = append(poems, r)
	}
	return poems, nil
}

func (c *Client) UploadPoem(ctx context.Context, in *proto.Poem, opts ...grpc.CallOption) (*proto.UploadPoemResponse, error) {
	return c.client.UploadPoem(ctx, in, opts...)
}

func (c *Client) UploadPoemStream(ctx context.Context, in *proto.Poem, opts ...grpc.CallOption) (*proto.UploadPoemResponse, error) {
	sin, err := c.client.UploadPoemStream(ctx, opts...)
	if err != nil {
		return nil, err
	}
	if err := sin.Send(&proto.StreamPoem{OneOf: &proto.StreamPoem_Title{Title: in.GetTitle()}}); err != nil {
		return nil, err
	}
	if err := sin.Send(&proto.StreamPoem{OneOf: &proto.StreamPoem_Author{Author: in.GetAuthor()}}); err != nil {
		return nil, err
	}
	for _, content := range in.GetContents() {
		if err := sin.Send(&proto.StreamPoem{OneOf: &proto.StreamPoem_Content{Content: content}}); err != nil {
			return nil, err
		}
	}
	return sin.CloseAndRecv()
}

func (c *Client) BatchUploadPoem(ctx context.Context, in []*proto.Poem, opts ...grpc.CallOption) (*proto.UploadPoemResponse, error) {
	return c.client.BatchUploadPoem(ctx, &proto.PoemCollection{Value: in}, opts...)
}

func (c *Client) BatchUploadPoemStream(ctx context.Context, in []*proto.Poem, afterUpload func(*proto.UploadPoemResponse), opts ...grpc.CallOption) error {
	stream, err := c.client.BatchUploadPoemStream(ctx, opts...)
	if err != nil {
		return err
	}

	var mu sync.Mutex
	safeAfterUpload := func(r *proto.UploadPoemResponse) error {
		mu.Lock()
		defer mu.Unlock()

		var err error
		defer func() {
			if r := recover(); r != nil {
				err = fmt.Errorf("after upload err: %v", r)
			}
		}()
		afterUpload(r)
		return err
	}

	ch := make(chan error, 2)
	defer close(ch)

	done := sync.WaitGroup{}
	done.Add(1)
	go func() {
		defer done.Done()
		for {
			r, err := stream.Recv()
			if err == io.EOF {
				return
			}
			if err != nil {
				ch <- err
				return
			}
			if err := safeAfterUpload(r); err != nil {
				ch <- err
				return
			}
		}
	}()

	done.Add(1)
	go func() {
		defer done.Done()
		for _, p := range in {
			if err := stream.Send(p); err != nil {
				ch <- err
				return
			}
		}
		if err := stream.CloseSend(); err != nil {
			ch <- err
		}
	}()

	done.Wait()
	select {
	case err := <-ch:
		return err
	default:
		return nil
	}
}

func (c *Client) Close() {
	c.conn.Close()
}

func NewClient(addr string) *Client {
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	return &Client{conn: conn, client: proto.NewPoemServiceClient(conn)}
}

var (
	addr     = flag.String("addr", "localhost:50051", "port to connect to")
	jsonFile = flag.String("json_file", "", "client upload poem json file")
)

func main() {
	flag.Parse()

	c := NewClient(*addr)
	defer c.Close()

	if *jsonFile == "" {
		if file, err := os.Getwd(); err != nil {
			log.Fatalf("failed to get work dir: %v", err)
		} else {
			*jsonFile = filepath.Join(file, "testdata", "client_poem.json")
		}
	}
	db := testdata.NewDB(*jsonFile)

	func() {
		t := "无题"
		log.Printf("get poem: %s\n", t)
		if r, err := c.GetPoem(context.Background(), &proto.GetPoemRequest{Title: t}); err != nil {
			log.Fatalf("did not get poem: %v", err)
		} else {
			fmt.Println(proto.Serialize(r))
		}
	}()

	func() {
		t := "洛神赋"
		log.Printf("get poem: %s by stream\n", t)
		if r, err := c.GetPoemStream(context.Background(), &proto.GetPoemRequest{Title: t}); err != nil {
			log.Fatalf("did not get poem: %v", err)
		} else {
			fmt.Println(proto.Serialize(r))
		}
	}()

	func() {
		t := "醉花阴"
		log.Printf("upload poem: %s\n", t)
		if poem, err := db.GetPoem(t); err != nil {
			log.Fatalf("did not get poem: %v", err)
		} else if r, err := c.UploadPoem(context.Background(), poem); err != nil {
			log.Fatalf("did not upload poem: %v", err)
		} else {
			for _, p := range r.GetData() {
				fmt.Println(proto.Serialize(p))
			}
		}
	}()

	func() {
		t := "滕王阁序"
		log.Printf("upload poem: %s by stream\n", t)
		if poem, err := db.GetPoem(t); err != nil {
			log.Fatalf("did not get poem: %v", err)
		} else if r, err := c.UploadPoemStream(context.Background(), poem); err != nil {
			log.Fatalf("did not upload poem: %v", err)
		} else {
			for _, p := range r.GetData() {
				fmt.Println(proto.Serialize(p))
			}
		}
	}()

	func() {
		log.Println("batch upload poems")
		if r, err := c.BatchUploadPoem(context.Background(), db.GetPoemCollection()); err != nil {
			log.Fatalf("did not batch upload poem: %v", err)
		} else {
			for _, p := range r.GetData() {
				fmt.Println(proto.Serialize(p))
			}
		}
	}()

	func() {
		log.Println("batch upload poems by stream")
		afterFunc := func(r *proto.UploadPoemResponse) {
			for _, p := range r.GetData() {
				fmt.Println(proto.Serialize(p))
			}
		}
		if err := c.BatchUploadPoemStream(context.Background(), db.GetPoemCollection(), afterFunc); err != nil {
			log.Fatalf("did not batch upload poem: %v", err)
		}
	}()

	func() {
		log.Println("get all poems")
		if poems, err := c.GetPoemAll(context.Background()); err != nil {
			log.Fatalf("did not get poem: %v", err)
		} else {
			for _, p := range poems {
				fmt.Println(proto.Serialize(p))
			}
		}
	}()

	func() {
		log.Println("get all poems by stream")
		if poems, err := c.GetPoemAllStream(context.Background()); err != nil {
			log.Fatalf("did not get poem: %v", err)
		} else {
			for _, p := range poems {
				fmt.Println(proto.Serialize(p))
			}
		}
	}()
}
