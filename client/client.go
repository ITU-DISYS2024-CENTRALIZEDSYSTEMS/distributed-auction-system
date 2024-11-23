package main

import (
	"bufio"
	"context"
	proto "distributed-auction-system/auction"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var username string;
var clt proto.AuctionClient;

func main() {
	clt = createClient(":50051")
	log.Println("Bidding started! Come closer everybody!")
	input := bufio.NewScanner(os.Stdin)
	log.Println("Write your username (no spaces):")
	input.Scan()
	username = input.Text()

	
	for {
		log.Println("Options: bid [amount] ; results")
		input.Scan()
		inputList := strings.Split(input.Text(), " ")
		if inputList[0] == "bid"{
			bid := proto.Amount{}
			bidAmount, err := strconv.ParseFloat(inputList[1], 32)
			if err != nil{
				log.Println("The bid is not a number: 'bid [decimal number]'")
				continue
			}
			bid.Amount = float32(bidAmount)
			bid.Username = username
			sendBid(&bid, clt)
		}else if inputList[0] == "results"{
			outcome, err := getResults(clt)
			if err != nil{
				log.Fatalln("Some error happened when getting results")
				continue
			}

			if outcome.IsFinished{
				log.Println("The auction is over!")
				log.Println("The winner is: " + outcome.Username)
				log.Printf("Payed the amount: %d \n", &outcome.Price)
			}else{
				log.Println("The auction is not over yet!")
				log.Printf("The highest bid: %d \n", &outcome.Price)
			}
		}else {
			log.Println("Unknown command!")
		}
	}
}

func createClient(port string) (proto.AuctionClient){
	grpcOptions := grpc.WithTransportCredentials(insecure.NewCredentials())
	conn, err := grpc.NewClient(port, grpcOptions)
	if err != nil {
		panic("Couldnt open client")
	}
	defer conn.Close()
	return proto.NewAuctionClient(conn)
}

func sendBid(bid *proto.Amount, client proto.AuctionClient){
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	ack, err := client.Bid(ctx, bid)
	log.Println(ack)
	if err != nil{
		log.Fatalf("The bid did not go through!")
	}else if !ack.Acknowledge{
		log.Println("Wrongful bid; Either too low or not a number!")
	}else{
		log.Println("Your bid went through!")
	}
}

func getResults(client proto.AuctionClient) (proto.Outcome, error){
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	outcome, err := client.Result(ctx, &proto.AuctionHouse{})
	if err != nil{
		log.Fatalln("Couldnt get result!")
		return proto.Outcome{}, err
	}

	return *outcome, nil
}