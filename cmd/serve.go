package cmd

import (
	"context"
	"errors"
	"github.com/ShaDream/library-api-practice/handler"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)
import "github.com/spf13/cobra"

var serveCommand = &cobra.Command{
	Use:   "serve",
	Short: "",
	Long:  "",
	RunE:  executeServe,
}

func init() {

}

func executeServe(cmd *cobra.Command, args []string) error {
	s := handler.NewServer()

	stop := make(chan os.Signal)
	signal.Notify(stop, os.Interrupt)

	go func() {
		if err := s.Start(); err != nil {
			if err != http.ErrServerClosed {
				log.Fatal(err)
			}
		}
	}()
	<-stop

	log.Println("Shutting down...")

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	if err := s.Shutdown(ctx); err != nil {
		if !errors.Is(err, context.DeadlineExceeded) {
			log.Fatalf("Context deadline exceeded. Shut down")
		}
	}

	log.Println("down")
	return nil
}
