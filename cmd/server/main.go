package main

import (
	"context"
	"errors"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/gorilla/mux"
	"github.com/soheilhy/cmux"
	"google.golang.org/grpc"

	"github.com/stevenferrer/cmux-http-grpc/pb"
	"github.com/stevenferrer/cmux-http-grpc/pingserver"
)

const (
	port = 8080
)

func main() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGKILL, syscall.SIGTERM)

	l, err := net.Listen("tcp", ":"+strconv.Itoa(port))
	if err != nil {
		log.Fatalf("net listen: %v", err)
	}

	m := cmux.New(l)

	grpcl := m.MatchWithWriters(cmux.HTTP2MatchHeaderFieldSendSettings("content-type", "application/grpc"))
	httpl := m.Match(cmux.HTTP1Fast())

	grpcServer := newGRPCServer()
	go func() {
		if err := grpcServer.Serve(grpcl); err != nil &&
			!errors.Is(err, grpc.ErrServerStopped) {
			log.Fatalf("grpc serve: %v", err)
		}
	}()
	log.Print("grpc server started.")

	httpServer := newHTTPServer(port)
	go func() {
		if err := httpServer.Serve(httpl); err != nil &&
			!errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("http listen: %s\n", err)
		}
	}()
	log.Print("http server started.")

	go func() {
		if err := m.Serve(); !strings.Contains(err.Error(), "use of closed network connection") {
			log.Fatalf("cmux serve: %v", err)
		}
	}()
	log.Print("cmux started.")

	<-c
	log.Print("exiting...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer func() { cancel() }()

	if err := httpServer.Shutdown(ctx); err != nil {
		log.Fatalf("http server shutdown: %v", err)
	}

	grpcServer.Stop()
	m.Close()

	log.Print("exited.")
}

func newGRPCServer() *grpc.Server {
	grpcServer := grpc.NewServer()
	pb.RegisterPingServer(grpcServer, pingserver.New())
	return grpcServer
}

func newHTTPServer(port int) *http.Server {
	r := mux.NewRouter()
	r.HandleFunc("/ping", indexHandler).Methods("GET")
	return &http.Server{Handler: r}
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("pong"))
}
