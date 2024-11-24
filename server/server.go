package main

import (
	"context"
	proto "distributed-auction-system/auction"
	"errors"
	"log"
	"net"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
	"google.golang.org/grpc"
)

type auctionServer struct {
	proto.UnimplementedAuctionServer
}

var highestBid int32 = 0
var highestBidUsername string = ""
var auctionTime int32 = 30
var openForBids bool = true

// Finds the first available port from an array of ports in .env file. Then returns an listener on that port.
func findAvailablePort(ports []string) (listener net.Listener, err error) {
	for _, port := range ports {
		listener, err := net.Listen("tcp", ":"+port)
		if err == nil {
			return listener, nil
		}
	}

	return nil, err
}

// Returns the index where the current port used by the listener is, in the given array.
func selectedPortIndex(ports []string, listener net.Listener) (index int, err error) {
	selectedPort := strconv.Itoa(listener.Addr().(*net.TCPAddr).Port)

	for i, port := range ports {
		if selectedPort == port {
			return i, nil
		}
	}

	return -1, errors.New("port not found")
}

func auctionTimer() {
	openForBids = true
	time.Sleep(time.Duration(auctionTime * int32(time.Second)))
	openForBids = false
}

func (s *auctionServer) Bid(_ context.Context, amount *proto.Amount) (*proto.Ack, error) {
	ackResponse := proto.Ack{}

	if amount.Amount <= highestBid {
		ackResponse.Acknowledge = false
		log.Printf("\nBid refused from: %s of the amount: %d \n", amount.Username, amount.Amount)
	} else {
		highestBid = amount.Amount
		highestBidUsername = amount.Username
		ackResponse.Acknowledge = true
		log.Printf("\nBid acknowledged from: %s of the amount: %d \n", amount.Username, amount.Amount)
	}

	return &ackResponse, nil
}

func (s *auctionServer) Result(_ context.Context, _ *proto.AuctionHouse) (*proto.Outcome, error) {
	outcomeResult := &proto.Outcome{
		Price:      highestBid,
		Username:   highestBidUsername,
		IsFinished: !openForBids,
	}

	return outcomeResult, nil
}

func main() {
	// Read ports from .env file
	envFile, _ := godotenv.Read(".env")
	ports := strings.Split(envFile["PORTS"], ",")

	listener, err := findAvailablePort(ports)
	if err != nil {
		log.Fatalf("Could not start listener: %s", err)
		return
	}

	server := grpc.NewServer()

	service := &auctionServer{}

	proto.RegisterAuctionServer(server, service)

	go auctionTimer()
	err = server.Serve(listener)
	if err != nil {
		log.Fatalf("Could not serve: %s", err)
		return
	}
}
