package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"net"
	"os"
	"strings"
	"time"

	_ "github.com/lib/pq"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/encoding"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"

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

type orderService struct {
	db *sql.DB
}

func main() {
	encoding.RegisterCodec(jsoncodec.Codec{})

	addr := envDefault("GRPC_ADDR", ":9090")
	db, err := openDB(context.Background(), envDefault("DATABASE_URL", "postgres://protocol:protocol@localhost:5432/protocols?sslmode=disable"))
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	lis, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatal(err)
	}

	s := grpc.NewServer(grpc.UnaryInterceptor(unaryLog))
	registerOrderService(s, &orderService{db: db})
	reflection.Register(s)

	log.Printf("grpc listening on %s", addr)
	log.Fatal(s.Serve(lis))
}

func unaryLog(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
	start := time.Now()
	resp, err := handler(ctx, req)
	log.Printf("%s %s", info.FullMethod, time.Since(start))
	return resp, err
}

func registerOrderService(s *grpc.Server, impl *orderService) {
	s.RegisterService(&grpc.ServiceDesc{
		ServiceName: "orders.v1.OrderService",
		HandlerType: (*orderService)(nil),
		Methods: []grpc.MethodDesc{
			{
				MethodName: "CreateOrder",
				Handler:    impl.createOrderHandler,
			},
			{
				MethodName: "GetOrder",
				Handler:    impl.getOrderHandler,
			},
			{
				MethodName: "ListOrders",
				Handler:    impl.listOrdersHandler,
			},
		},
		Streams:  []grpc.StreamDesc{},
		Metadata: "no-proto",
	}, impl)
}

func (s *orderService) createOrderHandler(_ any, ctx context.Context, dec func(any) error, interceptor grpc.UnaryServerInterceptor) (any, error) {
	var in CreateOrderRequest
	if err := dec(&in); err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid_json")
	}
	handler := func(ctx context.Context, req any) (any, error) {
		r := req.(*CreateOrderRequest)
		r.Customer = strings.TrimSpace(r.Customer)
		if r.Customer == "" || r.AmountCents <= 0 {
			return nil, status.Error(codes.InvalidArgument, "invalid_payload")
		}

		o, err := createOrder(ctx, s.db, r.Customer, int(r.AmountCents))
		if err != nil {
			return nil, status.Error(codes.Internal, "db_error")
		}
		return &CreateOrderResponse{Order: o}, nil
	}
	if interceptor == nil {
		return handler(ctx, &in)
	}
	info := &grpc.UnaryServerInfo{Server: s, FullMethod: "/orders.v1.OrderService/CreateOrder"}
	return interceptor(ctx, &in, info, handler)
}

func (s *orderService) getOrderHandler(_ any, ctx context.Context, dec func(any) error, interceptor grpc.UnaryServerInterceptor) (any, error) {
	var in GetOrderRequest
	if err := dec(&in); err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid_json")
	}
	handler := func(ctx context.Context, req any) (any, error) {
		r := req.(*GetOrderRequest)
		r.ID = strings.TrimSpace(r.ID)
		if r.ID == "" {
			return nil, status.Error(codes.InvalidArgument, "invalid_id")
		}
		o, err := getOrder(ctx, s.db, r.ID)
		if errors.Is(err, sql.ErrNoRows) {
			return nil, status.Error(codes.NotFound, "not_found")
		}
		if err != nil {
			return nil, status.Error(codes.Internal, "db_error")
		}
		return &GetOrderResponse{Order: o}, nil
	}
	if interceptor == nil {
		return handler(ctx, &in)
	}
	info := &grpc.UnaryServerInfo{Server: s, FullMethod: "/orders.v1.OrderService/GetOrder"}
	return interceptor(ctx, &in, info, handler)
}

func (s *orderService) listOrdersHandler(_ any, ctx context.Context, dec func(any) error, interceptor grpc.UnaryServerInterceptor) (any, error) {
	var in ListOrdersRequest
	if err := dec(&in); err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid_json")
	}
	handler := func(ctx context.Context, req any) (any, error) {
		r := req.(*ListOrdersRequest)
		limit := int(r.Limit)
		if limit == 0 {
			limit = 25
		}
		if limit < 1 || limit > 200 {
			return nil, status.Error(codes.InvalidArgument, "invalid_limit")
		}
		items, err := listOrders(ctx, s.db, limit)
		if err != nil {
			return nil, status.Error(codes.Internal, "db_error")
		}
		return &ListOrdersResponse{Items: items}, nil
	}
	if interceptor == nil {
		return handler(ctx, &in)
	}
	info := &grpc.UnaryServerInfo{Server: s, FullMethod: "/orders.v1.OrderService/ListOrders"}
	return interceptor(ctx, &in, info, handler)
}

func openDB(ctx context.Context, url string) (*sql.DB, error) {
	db, err := sql.Open("postgres", url)
	if err != nil {
		return nil, err
	}
	db.SetMaxOpenConns(20)
	db.SetMaxIdleConns(20)
	db.SetConnMaxLifetime(30 * time.Minute)

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	if err := db.PingContext(ctx); err != nil {
		return nil, err
	}
	return db, nil
}

func createOrder(ctx context.Context, db *sql.DB, customer string, amountCents int) (Order, error) {
	var o Order
	var createdAt time.Time
	row := db.QueryRowContext(ctx, `
		INSERT INTO orders(customer, amount_cents, status)
		VALUES ($1, $2, 'created')
		RETURNING id::text, customer, amount_cents, status, created_at
	`, customer, amountCents)
	if err := row.Scan(&o.ID, &o.Customer, &o.AmountCents, &o.Status, &createdAt); err != nil {
		return Order{}, err
	}
	o.CreatedAt = createdAt.UTC().Format(time.RFC3339Nano)

	payload, _ := json.Marshal(map[string]any{
		"order_id":      o.ID,
		"customer":      o.Customer,
		"amount_cents":  o.AmountCents,
		"status":        o.Status,
		"occurred_at":   time.Now().UTC().Format(time.RFC3339Nano),
		"protocol_demo": "grpc",
	})
	_, _ = db.ExecContext(ctx, `
		INSERT INTO order_events(order_id, event_type, payload)
		VALUES ($1, $2, $3::jsonb)
	`, o.ID, "order.created", string(payload))

	return o, nil
}

func getOrder(ctx context.Context, db *sql.DB, id string) (Order, error) {
	var o Order
	var createdAt time.Time
	row := db.QueryRowContext(ctx, `
		SELECT id::text, customer, amount_cents, status, created_at
		FROM orders
		WHERE id = $1::uuid
	`, id)
	if err := row.Scan(&o.ID, &o.Customer, &o.AmountCents, &o.Status, &createdAt); err != nil {
		return Order{}, err
	}
	o.CreatedAt = createdAt.UTC().Format(time.RFC3339Nano)
	return o, nil
}

func listOrders(ctx context.Context, db *sql.DB, limit int) ([]Order, error) {
	rows, err := db.QueryContext(ctx, `
		SELECT id::text, customer, amount_cents, status, created_at
		FROM orders
		ORDER BY created_at DESC
		LIMIT $1
	`, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []Order
	for rows.Next() {
		var o Order
		var createdAt time.Time
		if err := rows.Scan(&o.ID, &o.Customer, &o.AmountCents, &o.Status, &createdAt); err != nil {
			return nil, err
		}
		o.CreatedAt = createdAt.UTC().Format(time.RFC3339Nano)
		out = append(out, o)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return out, nil
}

func envDefault(key, def string) string {
	if v := strings.TrimSpace(os.Getenv(key)); v != "" {
		return v
	}
	return def
}
