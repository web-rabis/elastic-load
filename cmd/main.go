package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/hashicorp/logutils"
	"golang.org/x/sync/errgroup"

	"github.com/web-rabis/elastic-load/internal/adapter/database/drivers"
	"github.com/web-rabis/elastic-load/internal/adapter/database/drivers/pgsql"
	"github.com/web-rabis/elastic-load/internal/config"
	"github.com/web-rabis/elastic-load/internal/server/http"
)

var (
	version = "unknown"
)

func main() {
	fmt.Printf("ELK export %s\n", version)

	opts := config.Parse(new(config.APIServer)).(*config.APIServer)

	setupLogUtils(opts.Dbg)

	appCtx, appCtxCancel := context.WithCancel(context.Background())
	defer appCtxCancel()

	go catchForTermination(appCtxCancel, os.Interrupt, syscall.SIGTERM)

	ds, err := SetupDatabase(appCtx, opts.DSURL, opts.DSName, opts.DSDB)
	if err != nil {
		log.Println(err)
		return
	}
	defer ds.Close(appCtx)

	servers, serversCtx := errgroup.WithContext(appCtx)
	servers.Go(func() error {
		return http.Run(serversCtx, opts, ds, version)
	})

	if err := servers.Wait(); err != nil {
		log.Printf("[INFO] process terminated, %s", err)
		return
	}

}

func setupLogUtils(inDebugMode bool) {
	filter := &logutils.LevelFilter{
		Levels:   []logutils.LogLevel{"DEBUG", "INFO", "WARN", "ERROR"},
		MinLevel: logutils.LogLevel("INFO"),
		Writer:   os.Stdout,
	}

	log.SetFlags(log.Ldate | log.Ltime)

	if inDebugMode {
		log.SetFlags(log.Ldate | log.Ltime | log.Lmicroseconds | log.Lshortfile)
		filter.MinLevel = "DEBUG"
	}

	log.SetOutput(filter)
}

func catchForTermination(cancel context.CancelFunc, signals ...os.Signal) {
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, signals...)
	<-stop
	log.Print("[WARN] interrupt signal")
	cancel()
}

func SetupDatabase(ctx context.Context, url, dsName, dbName string) (drivers.DataStore, error) {
	ds, err := pgsql.New(drivers.DataStoreConfig{
		URL:           url,
		DataStoreName: dsName,
		DataBaseName:  dbName,
	})
	if err != nil {
		return nil, err
	}

	if err := ds.Connect(ctx); err != nil {
		errText := fmt.Sprintf("[ERROR] cannot connect to datastore %s: %v", dsName, err)
		return nil, errors.New(errText)
	}

	fmt.Println("Connected to", ds.Name())

	return ds, nil
}
