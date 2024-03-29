package main

import (
	"flag"
	"io"
	"log"
	"math/rand"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"
	"unsafe"

	_ "github.com/joho/godotenv/autoload"

	"github.com/ivoras/bulletcache2/types"
	"github.com/ivoras/bulletcache2/util"
)

const (
	eventQuit = iota
)

type sysEventMessage struct {
	event int
	idata int
}

var sysEventChannel = make(chan sysEventMessage, 5)
var logOutput io.Writer
var startTime time.Time

var logFileName = flag.String("log", "/tmp/bulletcache2.log", "Log file ('-' for only stderr)")

func main() {
	os.Setenv("TZ", "UTC")
	startTime = time.Now()
	rand.Seed(startTime.UnixNano())

	if runtime.GOOS == "windows" {
		*logFileName = "c:\\temp\\bulletcache2.log"
	}
	flag.Parse()

	if *logFileName != "-" {
		f, err := os.OpenFile(*logFileName, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0664)
		if err == nil {
			defer f.Close()
			logOutput = io.MultiWriter(os.Stderr, f)
		} else {
			log.Println("Cannot open log file " + *logFileName)
		}
	} else {
		logOutput = os.Stderr
	}
	log.SetOutput(logOutput)

	log.Println("Starting up...")
	log.Printf("sizeof(CacheRecord)=%d\n", unsafe.Sizeof(types.CacheRecord{}))

	sigChannel := make(chan os.Signal, 1)
	signal.Notify(sigChannel, syscall.SIGINT)

	//go webServer()
	//go infraWebServer()

	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	oldAlloc := int64(m.Alloc)
	printMemStats(&m)

	for {
		select {
		case msg := <-sysEventChannel:
			switch msg.event {
			case eventQuit:
				log.Println("Exiting")
				os.Exit(msg.idata)
			}
		case sig := <-sigChannel:
			switch sig {
			case syscall.SIGINT:
				sysEventChannel <- sysEventMessage{event: eventQuit, idata: 0}
				log.Println("^C detected")
			}
		case <-time.After(60 * time.Second):

			runtime.ReadMemStats(&m)
			if util.Abs(int64(m.Alloc)-oldAlloc) > 1024*1024 {
				printMemStats(&m)
				oldAlloc = int64(m.Alloc)
			}
		case <-time.After(15 * time.Minute):
			//cleanupDb()
		}
	}
}

func printMemStats(m *runtime.MemStats) {
	// For info on each, see: https://golang.org/pkg/runtime/#MemStats
	log.Printf("Alloc: %v MiB\tTotalAlloc: %v MiB\tSys: %v MiB\tNumGC: %v\tUptime: %0.1fh\n",
		util.BToMB(m.Alloc), util.BToMB(m.TotalAlloc), util.BToMB(m.Sys), m.NumGC, time.Since(startTime).Hours())
}
