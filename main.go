package main

import (
	"context"
	"log"
	"net/http"
	"os/signal"
	"poc-batch-database/app/demo"
	"poc-batch-database/database"
	"sync"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT)
	defer stop()
	r, stop := router(ctx)
	defer stop()

	srv := &http.Server{
		Addr:              ":8080",
		Handler:           r.Handler(),
		ReadHeaderTimeout: 30 * time.Second,
		WriteTimeout:      10 * time.Second,
		MaxHeaderBytes:    1 << 20,
	}

	go func(ctx context.Context, srv *http.Server) {
		<-ctx.Done()
		d := time.Duration(10 * time.Second)
		log.Printf("Shutdown Server in %s ...\n", d)
		ctx, cancel := context.WithTimeout(context.Background(), d)
		defer cancel()
		if err := srv.Shutdown(ctx); err != nil {
			log.Fatal("Server Shutdown:", err)
		}
	}(ctx, srv)

	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Panic(err)
	}

}

func router(ctx context.Context) (r *gin.Engine, stop func()) {
	e := gin.New()
	e.Use(gin.Recovery())
	clientDatabse, err := database.New()
	if err != nil {
		log.Fatal(err)
	}
	var BufferRequest []demo.DemoRequest
	syncLock := &sync.Mutex{}
	// time ticker if buffer is morethan zero

	{

		demoStorage := demo.NewStorage(clientDatabse)
		demoService := demo.NewService(demoStorage)

		go func() {
			ticker := time.NewTicker(800 * time.Millisecond)
			for {
				select {
				case <-ticker.C:
					go func() {
						err := demoService.BatchSaveDatabase(context.Background(), &BufferRequest, syncLock)
						if err != nil {
							log.Println(err)
						}
					}()
				case <-ctx.Done():
					// do something before end service
					// EX. Send to one by one kafka, send to topic consume, Log to sre
					return
				}
			}
		}()
		e.POST("/demo", demoService.SaveToDatabase(&BufferRequest, syncLock))
	}

	return e, func() {
		clientDatabse.Close()
	}
}
