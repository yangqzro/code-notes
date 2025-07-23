package message

import (
	context "context"
	"goexamples/utils"
	"io"
	"log"

	"google.golang.org/grpc"
)

// Client Handler 实现：
//
//  1. 依据情况，准备元数据（metadata）、请求数据 ？
//  2. 发送 rpc 请求
//  3. 依据情况，接收响应 Header ？
//  4. 处理数据，处理响应数据、向服务器发送数据
//  5. 依据情况，接收响应 Trailer ？
type MessageSrvClient struct {
	conn   *grpc.ClientConn
	client MessageServiceClient
}

func (c *MessageSrvClient) GetConn() *grpc.ClientConn {
	return c.conn
}

func (c *MessageSrvClient) Unary(ctx context.Context, in *Message, opts ...grpc.CallOption) (*Message, error) {
	return c.client.Unary(ctx, in, opts...)
}

func (c *MessageSrvClient) ClientStream(ctx context.Context, in []*Message, opts ...grpc.CallOption) ([]*Message, error) {
	sin, err := c.client.ClientStream(ctx, opts...)
	if err != nil {
		return nil, err
	}

	// ClientStream.Header() 是个阻塞方法，需要注意使用时机。
	// ClientStream.Recv() 也是一个阻塞方法，用于接收响应数据。在首次接收响应数据时会自动等待 Header 的到达，因为 Header 总是先于数据流传输。
	// 在流式 rpc 通信中，若在服务端未显式发送 Header（即未调用 ServerStream.SendHeader()）时，
	// 且在客户端下 Header 读取（ClientStream.Header()）先于响应数据读取（ClientStream.Recv()、ClientStreamingClient.CloseAndRecv()）执行，会导致双向阻塞死锁，即客户端和服务端相互等待对方响应。
	if header, err := sin.Header(); err == nil {
		log.Printf("client.ClientStream header: %s\n", utils.String(header))
	}
	for _, m := range in {
		if err := sin.Send(m); err != nil {
			return nil, err
		}
	}

	if mc, err := sin.CloseAndRecv(); err != nil {
		return nil, err
	} else {
		log.Printf("client.ClientStream trailer: %s\n", utils.String(sin.Trailer()))
		return mc.GetValue(), nil
	}
}

func (c *MessageSrvClient) ServerStream(ctx context.Context, in []*Message, opts ...grpc.CallOption) ([]*Message, error) {
	sout, err := c.client.ServerStream(ctx, &MessageCollection{Value: in}, opts...)
	if err != nil {
		return nil, err
	}

	if header, err := sout.Header(); err == nil {
		log.Printf("client.ServerStream header: %s\n", utils.String(header))
	}

	cout := []*Message{}
	for {
		out, err := sout.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		cout = append(cout, out)
	}

	log.Printf("client.ServerStream trailer: %s\n", utils.String(sout.Trailer()))
	return cout, nil
}

func (c *MessageSrvClient) BidirectionalStream(ctx context.Context, in []*Message, opts ...grpc.CallOption) ([]*Message, error) {
	stream, err := c.client.BidirectionalStream(ctx, opts...)
	if err != nil {
		return nil, err
	}

	if header, err := stream.Header(); err == nil {
		log.Printf("client.BidirectionalStream header: %s\n", utils.String(header))
	}

	for _, m := range in {
		if err := stream.Send(m); err != nil {
			return nil, err
		}
	}
	if err := stream.CloseSend(); err != nil {
		return nil, err
	}

	cout := []*Message{}
	for {
		if out, err := stream.Recv(); err == nil {
			cout = append(cout, out)
		} else if err == io.EOF {
			break
		} else {
			return nil, err
		}
	}

	log.Printf("client.BidirectionalStream trailer: %s\n", utils.String(stream.Trailer()))
	return cout, nil
}

func (c *MessageSrvClient) Close() error {
	return c.conn.Close()
}

func NewMessageSrvClient(addr string, opts ...grpc.DialOption) (*MessageSrvClient, error) {
	conn, err := grpc.NewClient(addr, opts...)
	if err != nil {
		return nil, err
	}
	return &MessageSrvClient{conn: conn, client: NewMessageServiceClient(conn)}, nil
}
