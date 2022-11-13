package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"

	pb "github.com/iron-ore-26k/ore-pb-gen"
	"google.golang.org/grpc"
)

var (
	port = flag.Int("port", 50051, "Server port")
)

type server struct {
	pb.UnimplementedOreServiceServer
}

func (s *server) ReceiveSong(ctx context.Context, in *pb.PlaySongRequest) {
	log.Printf("Recieved: %v", in.GetSongToPlay())
}

func main() {
	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterOreServiceServer(s, &server{})
	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

}
