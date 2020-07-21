package main

import (
	"context"
	"fmt"
	"github.com/codezork/pghandler/src"
	"log"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"
	"gopkg.in/natefinch/lumberjack.v2"
)

func handler(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	command := query.Get("cmd")
	if command == "" {
		command = "backup"
	}
	dumpfile := ""
	if command == "restore" {
		log.Printf("Received restore\n")
		dumpfile = query.Get("file")
		log.Printf("file:%s\n", dumpfile)
	}
	scalefactor := "1"
	if command == "fill" {
		log.Printf("Received fill\n")
		scalefactor = query.Get("factor")
		if scalefactor == "" {
			log.Printf("fill factor fallback\n")
			scalefactor = "1"
		}
	}
	log.Printf("Received command %s dumpfile %s scalefactor %s\n", command, dumpfile, scalefactor)
	cmdExec(w, command, dumpfile, scalefactor)
}

func cmdExec(w http.ResponseWriter, command string, dumpfile string, scalefactor string) {
	var path string = "/usr/bin/" + command + ".sh"
	log.Printf("path %s\n", path)
	_, err := exec.LookPath(path)
	if err != nil {
		log.Printf("Unknown command %s\n", command)
		w.Write([]byte(fmt.Sprintf("Unknown command %s\n", command)))
		return
	}
	arg := ""
	if len(dumpfile) > 0 {
		arg = dumpfile
	} else if scalefactor != "1" {
		arg = scalefactor
	}
	log.Printf("arg: %s\n", arg)
	cmd := &exec.Cmd{
		Path:   path,
		Args:   []string{arg},
		Stdout: nil,
		Stderr: nil,
	}
	w.Write([]byte(fmt.Sprintf("Executing %s\n", command)))
	cmd.Start()
	cmd.Wait()
	w.Write([]byte(fmt.Sprintf("%s done\n", command)))
	//output, err := cmd.Output()
	//if err != nil {
	//	log.Printf("err:%s", err.Error())
	//}
	//log.Printf("output:%s", output)
	//w.Write([]byte(fmt.Sprintf("%s\n", output)))
	log.Printf("%s done\n", command)
}

func main() {
	// Create Server and Route Handlers
	r := mux.NewRouter()

	r.HandleFunc("/", handler)

	cfg, err := src.ReadConfiguration()
	if err != nil {
		panic(err)
	}

	readtimeout, err := time.ParseDuration(cfg.Server.ReadTimeout)
	if err != nil {
		readtimeout = 10 * time.Second
	}
	writetimeout, err := time.ParseDuration(cfg.Server.WriteTimeout)
	if err != nil {
		writetimeout = 10 * time.Second
	}
	srv := &http.Server{
		Handler:      r,
		Addr:         cfg.Server.Host + ":" + cfg.Server.Port,
		ReadTimeout:  readtimeout,
		WriteTimeout: writetimeout,
	}

	// Configure Logging
	LOG_FILE_LOCATION := os.Getenv("LOG_FILE_LOCATION")
	if LOG_FILE_LOCATION != "" {
		log.SetOutput(&lumberjack.Logger{
			Filename:   LOG_FILE_LOCATION,
			MaxSize:    500, // megabytes
			MaxBackups: 3,
			MaxAge:     28,   //days
			Compress:   true, // disabled by default
		})
	}

	// Start Server
	go func() {
		log.Println("Starting Server")
		if err := srv.ListenAndServe(); err != nil {
			log.Fatal(err)
		}
	}()

	// Graceful Shutdown
	waitForShutdown(srv)
}

func waitForShutdown(srv *http.Server) {
	interruptChan := make(chan os.Signal, 1)
	signal.Notify(interruptChan, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	// Block until we receive our signal.
	<-interruptChan

	// Create a deadline to wait for.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	srv.Shutdown(ctx)

	log.Println("Shutting down")
	os.Exit(0)
}
