package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"log/slog"
	"os"
	"path"
	"runtime"
	"time"

	"github.com/kardianos/service"
	"gopkg.in/natefinch/lumberjack.v2"
)

var (
	version bool
	help    bool
	svcFlag string
)

func init() {
	flag.BoolVar(&version, "version", false, "show version")
	flag.BoolVar(&help, "help", false, "show help")
	flag.StringVar(&svcFlag, "svc", "", "[install/uninstall/start/stop/restart]contorl the zagent")

	// Customizable output directory.
	logFilePath := "./logs/"
	if err := os.MkdirAll(logFilePath, 0o777); err != nil {
		log.Println(err.Error())
		return
	}

	// Set filename to date
	logFileName := time.Now().Format("2006-01-02") + ".log"
	fileName := path.Join(logFilePath, logFileName)
	if _, err := os.Stat(fileName); err != nil {
		if _, err := os.Create(fileName); err != nil {
			log.Println(err.Error())
			return
		}
	}

	lumberjackLogger := &lumberjack.Logger{
		Filename:   fileName,
		MaxSize:    20,   // A file can be up to 20M.
		MaxBackups: 5,    // Save up to 5 files at the same time
		MaxAge:     7,    // A file can be saved for up to 7 days.
		Compress:   true, // Compress with gzip.
	}

	logger := slog.New(slog.NewTextHandler(lumberjackLogger, &slog.HandlerOptions{
		AddSource: true,           // 输出日志语句的位置信息
		Level:     slog.LevelInfo, // 设置最低日志等级
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			if a.Key == slog.TimeKey { // 格式化 key 为 "time" 的属性值
				if t, ok := a.Value.Any().(time.Time); ok {
					a.Value = slog.StringValue(t.Format(time.DateTime))
				}
			}
			return a
		},
	}))

	slog.SetDefault(logger)
}

var (
	Version      string = ""
	gitBranch    string = ""
	gitTag       string = ""
	gitCommit    string = "$Format:%H$"          // sha1 from git, output of $(git rev-parse HEAD)
	gitTreeState string = "not a git tree"       // state of git tree, either "clean" or "dirty"
	buildDate    string = "1970-01-01T00:00:00Z" // build date in ISO8601 format, output of $(date -u +'%Y-%m-%dT%H:%M:%SZ')
)

func GetCommit() string {
	if gitCommit != "" {
		h := gitCommit
		if len(h) > 7 {
			h = h[:7]
		}
		return h
	}
	return gitCommit
}

// Info contains versioning information.
type Info struct {
	Version      string `json:"Version"`
	GitBranch    string `json:"gitBranch"`
	GitTag       string `json:"gitTag"`
	GitCommit    string `json:"gitCommit"`
	GitTreeState string `json:"gitTreeState"`
	BuildDate    string `json:"buildDate"`
	GoVersion    string `json:"goVersion"`
	Compiler     string `json:"compiler"`
	Platform     string `json:"platform"`
}

func (info Info) String() string {
	return info.GitTag
}

func Get() Info {
	return Info{
		Version:      Version,
		GitBranch:    gitBranch,
		GitTag:       gitTag,
		GitCommit:    GetCommit(),
		GitTreeState: gitTreeState,
		BuildDate:    buildDate,
		GoVersion:    runtime.Version(),
		Compiler:     runtime.Compiler,
		Platform:     fmt.Sprintf("%s/%s", runtime.GOOS, runtime.GOARCH),
	}
}

type program struct{}

func (p *program) Start(s service.Service) error {
	slog.Info("start service ...")
	if service.Interactive() {
		slog.Info("Running in terminal.")
	} else {
		slog.Info("Running under service manager.")
	}
	go p.run()
	return nil
}

func (p *program) run() {
	slog.Info("run ...")
	Work()
}

func (p *program) Stop(s service.Service) error {
	slog.Info("stop service ...")
	return nil
}

func Work() {
	slog.Info("Work ...")
	for {
		time.Sleep(1 * time.Second)
		slog.Info("ping")
	}
}

func main() {
	flag.Parse()

	svcConfig := &service.Config{
		Name:        "zagent",
		DisplayName: "zagent",
		Description: "This is a agent service.",
	}

	prg := &program{}
	s, err := service.New(prg, svcConfig)
	if err != nil {
		slog.Error("%s", err)
	}

	if version {
		v := Get()
		marshalled, err := json.MarshalIndent(&v, "", "  ")
		if err != nil {
			fmt.Printf("%v\n", err)
			os.Exit(1)
		}

		fmt.Println(string(marshalled))
		os.Exit(0)
	}

	if help {
		fmt.Println("Usage: your_program [OPTIONS]")
		fmt.Println("Options:")
		flag.PrintDefaults()
		os.Exit(0)
	}

	if len(svcFlag) != 0 {
		err := service.Control(s, svcFlag)
		if err != nil {
			log.Printf("Valid actions: %q\n", service.ControlAction)
			log.Fatal(err)
		}
		return
	}

	err = s.Run()
	if err != nil {
		slog.Error("%s", err)
	}
}
