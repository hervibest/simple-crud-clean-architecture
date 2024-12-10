package helper

import (
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/gofiber/fiber/v2"
)

func StartServer(a *fiber.App, port int) {

	err := a.Listen(fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func StartServerWithGracefulShutdown(a *fiber.App, port int) {
	idleConnectionsClosed := make(chan struct{})

	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt)
		<-sigint

		if err := a.Shutdown(); err != nil {
			log.Printf("Cannot shutdown server %v", err)
		}

		close(idleConnectionsClosed)
	}()

	if err := a.Listen(fmt.Sprintf(":%d", port)); err != nil {
		log.Printf("Cannot run server %v", err)
	}

	<-idleConnectionsClosed
}
