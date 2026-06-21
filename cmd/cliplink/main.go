package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/mereska0/cliplink/api/gen/linkpb"
	"github.com/mereska0/cliplink/internal/config"
	"github.com/mereska0/cliplink/internal/encoder"
	"github.com/mereska0/cliplink/internal/grpcclient"
	"github.com/mereska0/cliplink/internal/grpcserver"
	"github.com/mereska0/cliplink/internal/httpserver"
	"github.com/mereska0/cliplink/internal/repository/postgres"
	"github.com/mereska0/cliplink/internal/service"
	"github.com/mereska0/cliplink/internal/tui"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	if len(os.Args) < 2 {
		runAuto()
		return
	}

	switch os.Args[1] {
	case "server":
		runServer()
	case "tui":
		runTUI()
	default:
		fmt.Println("unknown command:", os.Args[1])
		fmt.Println("usage: cliplink [server|tui]")
		os.Exit(1)
	}
}

func runServer() {
	ctx := context.Background()
	cfg := config.Load()

	pool, err := pgxpool.New(ctx, cfg.DatabaseURL)
	if err != nil {
		log.Fatal(err)
	}
	defer pool.Close()

	if err := pool.Ping(ctx); err != nil {
		log.Fatal(err)
	}

	repo := postgres.NewPostgresLinkRepository(pool)
	codeEncoder := encoder.NewBase62Encoder()
	linkService := service.NewLinkService(repo, codeEncoder)

	go func() {
		redirectHandler := httpserver.NewRedirectHandler(linkService)

		fmt.Println("HTTP redirect server started on", cfg.HTTPAddr)

		if err := http.ListenAndServe(cfg.HTTPAddr, redirectHandler); err != nil {
			log.Fatal(err)
		}
	}()

	listener, err := net.Listen("tcp", cfg.GRPCAddr)
	if err != nil {
		log.Fatal(err)
	}

	grpcServer := grpc.NewServer()

	linkpb.RegisterLinkServiceServer(
		grpcServer,
		grpcserver.NewLinkServer(linkService),
	)

	reflection.Register(grpcServer)

	fmt.Println("gRPC server started on", cfg.GRPCAddr)

	if err := grpcServer.Serve(listener); err != nil {
		log.Fatal(err)
	}
}

func runTUI() {
	cfg := config.Load()

	client, err := grpcclient.NewLinkClient(cfg.GRPCClientAddr)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	model := tui.NewModel(client)

	if _, err := tea.NewProgram(model).Run(); err != nil {
		log.Fatal(err)
	}
}
func runAuto() {
	cfg := config.Load()

	if isGRPCAvailable(cfg.GRPCClientAddr) {
		runTUI()
		return
	}

	fmt.Println("ClipLink server is not running. Starting local server...")

	go runServer()

	if err := waitForGRPC(cfg.GRPCClientAddr, 10*time.Second); err != nil {
		fmt.Println("failed to start ClipLink server:", err)
		return
	}

	runTUI()
}

func isGRPCAvailable(addr string) bool {
	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()

	client, err := grpcclient.NewLinkClient(addr)
	if err != nil {
		return false
	}
	defer client.Close()

	_, err = client.ListLinks(ctx)
	return err == nil
}

func waitForGRPC(addr string, timeout time.Duration) error {
	deadline := time.Now().Add(timeout)

	for time.Now().Before(deadline) {
		if isGRPCAvailable(addr) {
			return nil
		}

		time.Sleep(300 * time.Millisecond)
	}

	return context.DeadlineExceeded
}
