package server

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"sync"

	"github.com/afteralec/grpc-user/db"
	pb "github.com/afteralec/grpc-user/proto"
	"github.com/afteralec/grpc-user/services/user"
	"github.com/spf13/viper"

	"google.golang.org/grpc"
)

func Run(ctx context.Context, config *viper.Viper) error {
	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt)
	defer cancel()

	lis, err := net.Listen("tcp", ":8009")
	if err != nil {
		return err
	}

	db, err := db.Open("/var/db/mirror.db")
	if err != nil {
		return err
	}

	us, err := user.New(db, user.WithConfig(config))
	if err != nil {
		return err
	}
	us.SyncRootPermissions(ctx)

	s := grpc.NewServer()
	pb.RegisterMirrorServer(s, &server{user: &us})

	go func() {
		log.Printf("gRPC server listening at %v", lis.Addr())
		if err := s.Serve(lis); err != nil {
			fmt.Fprintf(os.Stderr, "error listening and serving: %s\n", err)
		}
	}()

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		<-ctx.Done()
		defer cancel()
		s.GracefulStop()
	}()
	wg.Wait()

	return nil
}
