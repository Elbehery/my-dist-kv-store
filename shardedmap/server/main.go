package main

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
	"log"
	"net"
	"playground/my-dist-kv-store/shardedmap"
	pb "playground/my-dist-kv-store/shardedmap/grpc-kv-store"
)

type kvServer struct {
	pb.UnimplementedKeyValueStoreServer
	store shardedmap.ShardedMap
}

func (s *kvServer) Get(ctx context.Context, req *pb.GetRequest) (*pb.GetResponse, error) {
	key := req.Key
	log.Printf("Received GET key=%v\n", key)
	val := s.store.Get(key)
	return &pb.GetResponse{Value: val}, nil
}

func (s *kvServer) Put(ctx context.Context, req *pb.PutRequest) (*emptypb.Empty, error) {
	key := req.Key
	val := req.Value
	log.Printf("Received PUT key=%v ,value=%v\n", key, val)
	s.store.Set(key, val)
	return &emptypb.Empty{}, nil
}

func main() {

	srvr := grpc.NewServer()
	pb.RegisterKeyValueStoreServer(srvr, &kvServer{
		store: shardedmap.NewShardedMap(3),
	})

	lis, err := net.Listen("tcp", ":5050")
	if err != nil {
		log.Fatal(err)
	}

	if err := srvr.Serve(lis); err != nil {
		log.Fatal(err)
	}
}
