package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"log"
	pb "playground/my-dist-kv-store/shardedmap/grpc-kv-store"
)

func main() {
	conn, err := grpc.Dial("localhost:5050", grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("can not dial server: %v", err)
	}

	client := pb.NewKeyValueStoreClient(conn)
	resp, err := client.Get(context.Background(), &pb.GetRequest{Key: string("A")})
	if err != nil {
		log.Fatalf("expected no result, got %v", resp.Value)
	}
	_, err = client.Put(context.Background(), &pb.PutRequest{
		Key:   "A",
		Value: "1",
	})
	if err != nil {
		log.Fatalf("can not put key value")
	}
	fmt.Println("PUT OK!")
	resp, err = client.Get(context.Background(), &pb.GetRequest{Key: "A"})
	if err != nil {
		log.Fatalf("second GET failed, %v", err)
	}
	fmt.Printf("Get Ok!, value is %v\n", resp.Value)
	fmt.Println("client done")
}
