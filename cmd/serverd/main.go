package main

import (
	"fmt"
	"os"

	"github.com/virsavik/sample-azure-func-app/internal/adapter/queue/message"
	"github.com/virsavik/sample-azure-func-app/internal/adapter/repository"
	rest "github.com/virsavik/sample-azure-func-app/internal/adapter/rest/v1"
	"github.com/virsavik/sample-azure-func-app/internal/adapter/storage"
	"github.com/virsavik/sample-azure-func-app/internal/config"
	"github.com/virsavik/sample-azure-func-app/internal/core/services"
	"github.com/virsavik/sample-azure-func-app/internal/httpio"
	"github.com/virsavik/sample-azure-func-app/internal/system"
)

func main() {
	if err := run(); err != nil {
		fmt.Printf("service exitted abnormally: %s\n", err.Error())
		os.Exit(1)
	}
}

func run() error {
	// Read config from env
	cfg, err := config.ReadConfigFromEnv()
	if err != nil {
		return err
	}

	// Setup system
	s, err := system.New(cfg)
	if err != nil {
		return err
	}

	// Setup module
	repo := repository.New(s.DB().Database(s.Config().MongoDB.DBName))
	publisher := message.New(s.StorageQueue())
	blobStorage := storage.New(s.BlobStorage(), s.Config().AzBlob.ContainerName)

	svc := services.New(repo.Files(), publisher, blobStorage)

	hdl := rest.New(svc)

	setupRouter(s, hdl)

	fmt.Println("started service")
	defer fmt.Println("stopped service")

	s.Waiter().Add(
		s.WaitForWeb,
		//s.WaitForRPC,
		//s.WaitForStream,
	)

	//go func() {
	//	for {
	//		var mem runtime.MemStats
	//		runtime.ReadMemStats(&mem)
	//		s.Logger().Infof("Alloc = %v, TotalAlloc = %v, Sys = %v, NumGC = %v", mem.Alloc/1024, mem.TotalAlloc/1024, mem.Sys/1024, mem.NumGC)
	//		time.Sleep(10 * time.Second)
	//	}
	//}()

	return s.Waiter().Wait()
}

func setupRouter(svc system.Service, hdl rest.Handler) {
	svc.Mux().Use(httpio.Middleware(svc.Logger()))

	svc.Mux().Get("/test/ping", hdl.Ping())
	svc.Mux().Get("/test/panic", hdl.Panic())

	svc.Mux().Post("/api/UploadFile", hdl.Upload())
}
