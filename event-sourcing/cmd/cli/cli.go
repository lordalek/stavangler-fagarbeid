package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/lordalek/stavangler-fagarbeid/event-sourcing/pkg/order"

	"google.golang.org/grpc"
)

func main() {

	conn, err := grpc.Dial("0.0.0.0:8877", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("fail to dial: %v", err)
	}

	defer conn.Close()

	client := order.NewOrderClient(conn)

	// reader := bufio.NewReader(os.Stdin)

	fmt.Print("<create id restaurant> or <get id>")
	// var cliCommand string
	// fmt.Scanln(&cliCommand)

	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	cliCommand := scanner.Text()

	if strings.Contains(cliCommand, "create") {

		res, err := client.CreateOrder(context.TODO(), &order.CreateOrderRequest{
			Restaurant: strings.Split(cliCommand, " ")[2],
			Id:         strings.Split(cliCommand, " ")[1],
		})

		if err != nil {
			log.Fatalf("failed to create: %v", err)
		}
		fmt.Printf("Create Results: %v\n", res)

	} else {
		getRes, err := client.GetOrder(context.TODO(), &order.FindOrderRequest{
			Id: strings.Split(cliCommand, " ")[1],
		})

		if err != nil {
			log.Fatalf("failed to get: %v", err)
		}

		fmt.Printf("ID: %s\nRestaurant: %s\nOrderLines: %s\n", getRes.Id, getRes.Restaurant, getRes.Orderlines)

	}
}
