package message

import (
	"goexamples/utils"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func GetStreamHandler(method string) grpc.StreamHandler {
	switch method {
	case "ClientStream":
		return _MessageService_ClientStream_Handler
	case "ServerStream":
		return _MessageService_ServerStream_Handler
	case "BidirectionalStream":
		return _MessageService_BidirectionalStream_Handler
	}
	return nil
}

func GenerateServerMetadata(from string) metadata.MD {
	return metadata.Pairs("timestamp", time.Now().Format(time.DateTime), "from", from, "random", utils.RandString(8))
}
