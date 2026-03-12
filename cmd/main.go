package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/imwxx/koiddns/config"
	"github.com/imwxx/koiddns/dns"
	"github.com/imwxx/koiddns/util"
	daemon "github.com/sevlyar/go-daemon"
)

func main() {
	help := flag.Bool("help", false, "Show this help message and exit")
	demo := flag.String("generate-config", "", "Generate a sample configuration file")
	doValidate := flag.Bool("validate", false, "Validate config file and exit (use with --config)")
	configFile := flag.String("config", "/etc/config/koiddns", "Specify the configuration file")
	daemonize := flag.Bool("daemon", false, "Run as a daemon")
	pidFile := flag.String("pidfile", "/var/run/koiddns.pid", "PID file when running as daemon")
	logFile := flag.String("logfile", "/var/log/koiddns.log", "Log file when running as daemon")

	flag.Parse()

	if *help {
		util.ShowHelp()
		return
	}

	if *demo != "" {
		err := config.GenerateSampleConfig(*demo)
		if err != nil {
			log.Fatalf("Error generating sample configuration file: %v", err)
		}
		log.Printf("Sample configuration file generated at: %s\n", *demo)
		return
	}

	if *doValidate {
		_, err := config.LoadConfig(*configFile)
		if err != nil {
			log.Printf("Config validation failed: %v", err)
			os.Exit(1)
		}
		log.Print("Config is valid")
		return
	}

	if *daemonize {
		cntxt := &daemon.Context{
			PidFileName: *pidFile,
			PidFilePerm: 0644,
			LogFileName: *logFile,
			LogFilePerm: 0640,
			WorkDir:     "./",
			Umask:       027,
			Args:        []string{"[koiddns daemon]"},
		}

		d, err := cntxt.Reborn()
		if err != nil {
			log.Fatal("Unable to run: ", err)
		}
		if d != nil {
			return
		}
		defer cntxt.Release()

		log.Print("- - - - - - - - - - - - - - -")
		log.Print("Daemon started")
	}

	run(configFile)
}

func run(configFile *string) {
	cfg, err := config.LoadConfig(*configFile)
	if err != nil {
		log.Fatal(err)
	}

	ctx, stop := context.WithCancel(context.Background())
	defer stop()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigCh
		log.Print("Received signal, shutting down...")
		stop()
	}()

	interval := time.Duration(cfg.Main.ExecutionCycleMinutes) * time.Minute
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	updateDNS(cfg)

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			updateDNS(cfg)
		}
	}
}

func updateDNS(cfg *config.Config) {
	ip, err := util.GetPublicIP()
	if err != nil {
		log.Printf("GetPublicIP failed: %v", err)
		return
	}

	providerConfigs := make(map[string]config.ProviderConfig)
	for _, p := range cfg.Providers {
		providerConfigs[p.Name] = p
	}

	for _, domain := range cfg.Domains {
		providerConfig, ok := providerConfigs[domain.Provider]
		if !ok {
			log.Printf("Skip domain %s.%s: unknown provider %q", domain.SubDomain, domain.PrimaryDomain, domain.Provider)
			continue
		}

		recordValue := ip
		if domain.Value != "" {
			recordValue = domain.Value
		}

		var err error
		switch domain.Provider {
		case "aliyun":
			err = dns.UpdateAliyunDNS(providerConfig, domain.SubDomain, domain.PrimaryDomain, domain.RecordType, domain.RecordId, domain.Line, domain.Priority, recordValue)
		case "tencent":
			err = dns.UpdateTencentDNS(providerConfig, domain.SubDomain, domain.PrimaryDomain, domain.RecordType, domain.RecordId, domain.Line, domain.Priority, recordValue)
		default:
			log.Printf("Skip domain %s.%s: unsupported provider %q", domain.SubDomain, domain.PrimaryDomain, domain.Provider)
			continue
		}
		if err != nil {
			log.Printf("Update DNS %s.%s (%s): %v", domain.SubDomain, domain.PrimaryDomain, domain.Provider, err)
		}
	}
}
