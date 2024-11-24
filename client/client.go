package main

import (
	"bufio"
	"context"
	"distributed-auction-system/auction"
	proto "distributed-auction-system/auction"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var username string

func main() {
	envFile, _ := godotenv.Read(".env")
	ports := strings.Split(envFile["PORTS"], ",")

	var clients []proto.AuctionClient

	for _, port := range ports {
		grpcOptions := grpc.WithTransportCredentials(insecure.NewCredentials())
		conn, err := grpc.NewClient(":"+port, grpcOptions)
		if err != nil {
			log.Fatalf("Cannot create client: %s", err)
		}

		clients = append(clients, auction.NewAuctionClient(conn))
	}

	log.Println("Bidding started! Come closer everybody!")
	input := bufio.NewScanner(os.Stdin)
	log.Println("Write your username (no spaces):")
	input.Scan()
	username = input.Text()

	for {
		log.Println("Options: bid [amount] ; results")
		input.Scan()
		inputList := strings.Split(input.Text(), " ")
		if inputList[0] == "bid" {
			bid := proto.Amount{}
			bidAmount, err := strconv.ParseInt(inputList[1], 10, 32)
			if err != nil {
				log.Println("The bid is not a number: 'bid [decimal number]'")
				continue
			}

			bid.Amount = int32(bidAmount)
			bid.Username = username

			sendBid(&bid, clients)
		} else if inputList[0] == "results" {
			outcome, err := getResults(clients)
			if err != nil {
				log.Fatalln("Some error happened when getting results")
				continue
			}

			if outcome.IsFinished {
				log.Println("The auction is over!")
				log.Println("The winner is: " + outcome.Username)
				log.Printf("Payed the amount: %d \n", outcome.Price)
			} else {
				log.Println("The auction is not over yet!")

				log.Printf("The highest bid: %d \n", outcome.Price)
			}
		} else {
			log.Println("Unknown command!")
		}
	}
}

func sendBid(bid *proto.Amount, clients []proto.AuctionClient) {
	for _, client := range clients {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		ack, err := client.Bid(ctx, bid)

		if err != nil {
			log.Println("Couldn't send bid to one of the nodes!")
		} else if !ack.Acknowledge {
			log.Println("Wrongful bid; Either too low or not a number!")
		} else {
			log.Println("Your bid went through!")
		}
	}
}

func getResults(clients []proto.AuctionClient) (*proto.Outcome, error) {
	var results []*proto.Outcome

	for _, client := range clients {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		outcome, err := client.Result(ctx, &proto.AuctionHouse{})
		if err != nil {
			log.Println("Couldnt get result from one node!")
		}

		results = append(results, outcome)
	}

	if len(results) == 0 {
		return nil, fmt.Errorf("no results received")
	}

	return results[0], nil
}
