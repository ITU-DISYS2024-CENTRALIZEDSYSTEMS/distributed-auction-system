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
	_ sync.Mutex
}

var highestBid float32 = 0;
var highestBidUsername string = "";
var auctionTime float32 = 0;
var openForBids bool = true;

func main(){

	listener, err := net.Listen("tcp", ":50051")
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

func (s *auctionServer) bid(_ context.Context, amount *proto.Amount) (ack *proto.Ack, err error){
	ackResponse := proto.Ack{}

	if amount.Amount <= highestBid{
		ackResponse.Acknowledge = false
		log.Printf("\nBid refused from: %s of the amount: %f \n", amount.Username, amount.Amount)

	}else {
		highestBid = amount.Amount
		highestBidUsername = amount.Username
		ackResponse.Acknowledge = true
		log.Printf("\nBid acknowledged from: %s of the amount: %f \n", amount.Username, amount.Amount)
	}
	
	return &ackResponse, nil
}

func (s *auctionServer) result(_ context.Context, _ *proto.AuctionHouse) (outcome *proto.Outcome, err error){
	outcome = &proto.Outcome{}
	outcome.Price = highestBid
	outcome.Username = highestBidUsername
	outcome.IsFinished = !openForBids
	return outcome, nil
}