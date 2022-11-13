package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"

	"github.com/iron-ore-26k/ore-pb-gen"
	pb "github.com/iron-ore-26k/ore-pb-gen"
	"google.golang.org/grpc"
)

var (
	port = flag.Int("port", 50051, "Server port")
)

type server struct {
	pb.UnimplementedOreServiceServer
}

func (s *server) PlaySong(ctx context.Context, in *pb.PlaySongRequest) (*ore.PlaySongResponse, error) {
	log.Printf("Recieved: %v", in.GetSongToPlay())
	return &pb.PlaySongResponse{}, nil
}

func (s *server) Pause(ctx context.Context, in *pb.PauseRequest) (*ore.PauseResponse, error) {
	log.Printf("Recieved: pause request")
	return &pb.PauseResponse{}, nil
}

func main() {
	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterOreServiceServer(s, &server{})
	log.Printf("server listening at %v, %v", lis.Addr(), lis.Addr().Network())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

}
