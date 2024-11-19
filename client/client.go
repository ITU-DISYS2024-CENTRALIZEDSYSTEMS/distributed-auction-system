package main

import (
	"bufio"
	"context"
	proto "distributed-auction-system/auction"
	"log"
	"os"
	"strconv"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var username string;

func main() {
	log.Println("Bidding started! Come closer everybody!")
	input := bufio.NewScanner(os.Stdin)
	log.Println("Write your username (no spaces):")
	input.Scan()
	username = input.Text()

	
	for {
		
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
			sendBid(&bid, ":1234")
		}else if inputList[0] == "results"{
			outcome, err := getResults(":1234")
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
		}
	}
}

func sendBid(bid *proto.Amount, port string){
	
	grpcOptions := grpc.WithTransportCredentials(insecure.NewCredentials())
	conn, err := grpc.NewClient(port, grpcOptions)
	if err != nil {
		log.Fatalf("Couldn't open client")
		return
	}

	client := proto.NewAuctionClient(conn)
	ctx := context.Background()
	ack, err := client.Bid(ctx, bid)
	log.Println(ack)
	if !ack.Acknowledge{
		log.Println("Wrongful bid; Either too low or not a number!")
	}else if err != nil{
		log.Fatalf("The bid did not go through!")
	}else{
		log.Println("Your bid went through!")
	}
}

func getResults(port string) (proto.Outcome, error){
	grpcOptions := grpc.WithTransportCredentials(insecure.NewCredentials())
	conn, err := grpc.NewClient(port, grpcOptions)
	if err != nil{
		return proto.Outcome{}, err
	}

	client := proto.NewAuctionClient(conn)
	ctx := context.Background()
	outcome, err := client.Result(ctx, &proto.AuctionHouse{})
	if err != nil{
		log.Fatalln("Couldnt get result!")
		return proto.Outcome{}, err
	}

	return *outcome, nil
}