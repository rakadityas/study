package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/encoding"

	"protocols/grpc/internal/jsoncodec"
)

type Order struct {
	ID          string `json:"id"`
	Customer    string `json:"customer"`
	AmountCents int32  `json:"amount_cents"`
	Status      string `json:"status"`
	CreatedAt   string `json:"created_at"`
}

type CreateOrderRequest struct {
	Customer    string `json:"customer"`
	AmountCents int32  `json:"amount_cents"`
}

type CreateOrderResponse struct {
	Order Order `json:"order"`
}

type GetOrderRequest struct {
	ID string `json:"id"`
}

type GetOrderResponse struct {
	Order Order `json:"order"`
}

type ListOrdersRequest struct {
	Limit int32 `json:"limit"`
}

type ListOrdersResponse struct {
	Items []Order `json:"items"`
}

func main() {
	encoding.RegisterCodec(jsoncodec.Codec{})

	addr := envDefault("GRPC_TARGET", "localhost:9090")
	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithDefaultCallOptions(grpc.CallContentSubtype("json")))
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	createReq := &CreateOrderRequest{Customer: "alice", AmountCents: 2599}
	var createResp CreateOrderResponse
	if err := conn.Invoke(ctx, "/orders.v1.OrderService/CreateOrder", createReq, &createResp); err != nil {
		log.Fatal(err)
	}
	printJSON("CreateOrder", createResp)

	getReq := &GetOrderRequest{ID: createResp.Order.ID}
	var getResp GetOrderResponse
	if err := conn.Invoke(ctx, "/orders.v1.OrderService/GetOrder", getReq, &getResp); err != nil {
		log.Fatal(err)
	}
	printJSON("GetOrder", getResp)

	listReq := &ListOrdersRequest{Limit: 5}
	var listResp ListOrdersResponse
	if err := conn.Invoke(ctx, "/orders.v1.OrderService/ListOrders", listReq, &listResp); err != nil {
		log.Fatal(err)
	}
	printJSON("ListOrders", listResp)
}

func printJSON(label string, v any) {
	b, _ := json.MarshalIndent(v, "", "  ")
	fmt.Printf("%s:\n%s\n\n", label, string(b))
}

func envDefault(key, def string) string {
	if v := strings.TrimSpace(os.Getenv(key)); v != "" {
		return v
	}
	return def
}
