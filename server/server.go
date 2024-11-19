package main

import (
	"context"
	proto "distributed-auction-system/auction"
	"log"
	"net"
	"sync"
	"time"

	"google.golang.org/grpc"
)

type auctionServer struct{
	proto.UnimplementedAuctionServer
	lamportTime int32
	mu sync.Mutex
}

var highestBid float32 = 0;
var highestBidUsername string = "";
var auctionTime float32 = 0;
var openForBids bool = true;

func main(){

	listener, err := net.Listen("tcp", ":1234")
	if err != nil {
		log.Fatalf("Could not start listener: %s", err)
		return
	}

	server := grpc.NewServer()

	service := &auctionServer{
		lamportTime: 0,
	}

	proto.RegisterAuctionServer(server, service)

	go auctionTimer()
	err = server.Serve(listener)
	if err != nil {
		log.Fatalf("Could not serve: %s", err)
		return
	}
}

func auctionTimer(){
	openForBids = true
	time.Sleep(time.Duration(auctionTime))
	openForBids = false
}

func (s *auctionServer) bid(ctx context.Context, amount *proto.Amount) (ack *proto.Ack, err error){
	ackResponse := proto.Ack{}

	if amount.Amount <= highestBid{
		ackResponse.Acknowledge = false

	}else {
		highestBid = amount.Amount
		highestBidUsername = amount.Username
		ackResponse.Acknowledge = true
	}
	
	return &ackResponse, nil
}

func (s *auctionServer) result(ctx context.Context, ah *proto.AuctionHouse) (outcome *proto.Outcome, err error){
	outcome = &proto.Outcome{}
	outcome.Price = highestBid
	outcome.Username = highestBidUsername
	outcome.IsFinished = !openForBids
	return outcome, nil
}