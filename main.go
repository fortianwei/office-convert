package main

import (
	"flag"
	"fmt"

	"office-convert/web"
	"os"
	"path/filepath"

	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	log "github.com/sirupsen/logrus"

	"github.com/chai2010/winsvc"
)

var (
	appPath string

	flagServiceName = flag.String("service-name", "OfficeConvertService", "Set service name")
	flagServiceDesc = flag.String("service-desc", "office convert service", "Set service description")

	flagServiceInstall   = flag.Bool("service-install", false, "Install service")
	flagServiceUninstall = flag.Bool("service-remove", false, "Remove service")
	flagServiceStart     = flag.Bool("service-start", false, "Start service")
	flagServiceStop      = flag.Bool("service-stop", false, "Stop service")

	flagHelp = flag.Bool("help", false, "Show usage and exit.")
)

func init() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, `Usage:
  hello [options]...
Options:
`)
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, "%s\n", `
Example:
  # run hello server
  $ go build -o hello.exe hello.go
  $ hello.exe
  # install hello as windows service
  $ hello.exe -service-install
  # start/stop hello service
  $ hello.exe -service-start
  $ hello.exe -service-stop
  # remove hello service
  $ hello.exe -service-remove
  # help
  $ hello.exe -h
Report bugs to <chaishushan{AT}gmail.com>.`)
	}

	// change to current dir
	var err error
	if appPath, err = winsvc.GetAppPath(); err != nil {
		log.Fatal(err)
	}
	if err := os.Chdir(filepath.Dir(appPath)); err != nil {
		log.Fatal(err)
	}
	file := "./server"
	writer, _ := rotatelogs.New(
		file+".%Y%m%d.log",
		rotatelogs.WithRotationSize(20*1024*1024),
		rotatelogs.WithRotationCount(10),
	)
	log.SetOutput(writer)
	// logFile, err := os.OpenFile(file, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0766)
	// if err != nil {
	// 	panic(err)
	// }
	// log.SetOutput(logFile)
	// log.SetFlags(log.Flags() | log.Lshortfile)
}

/// 下面这个main不支持出错自动恢复，暂时不用。选择windows定时任务启动代替即可
func mainBak() {
	flag.Parse()

	// install service
	if *flagServiceInstall {
		if err := winsvc.InstallService(appPath, *flagServiceName, *flagServiceDesc); err != nil {
			log.Fatalf("installService(%s, %s): %v\n", *flagServiceName, *flagServiceDesc, err)
		}
		fmt.Printf("Done\n")
		return
	}

	// remove service
	if *flagServiceUninstall {
		if err := winsvc.RemoveService(*flagServiceName); err != nil {
			log.Fatalln("removeService:", err)
		}
		fmt.Printf("remove service Done\n")
		log.Printf("remove service Done\n")
		return
	}

	// start service
	if *flagServiceStart {
		if err := winsvc.StartService(*flagServiceName); err != nil {
			log.Fatalln("startService:", err)
		}
		log.Printf("start service Done\n")
		fmt.Printf("start service Done\n")
		return
	}

	// stop service
	if *flagServiceStop {
		if err := winsvc.StopService(*flagServiceName); err != nil {
			log.Fatalln("stopService:", err)
		}
		fmt.Printf("stop service Done\n")
		log.Printf("stop service Done\n")
		return
	}

	// run as service
	if !winsvc.IsAnInteractiveSession() {
		log.Println("main:", "runService")
		if err := winsvc.RunAsService(*flagServiceName, StartServer, StopServer, false); err != nil {
			log.Fatalf("svc.Run: %v\n", err)
		}
		return
	}

	// run as normal
	StartServer()
}
func main() {
	StartServer()
}
func StartServer() {
	web.StartServer()
}

func StopServer() {
	log.Println("StopServer")
}
