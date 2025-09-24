package main

import (
	"flag"
	"log"
	"time"

	"github.com/imwxx/koiddns/config"
	"github.com/imwxx/koiddns/dns"
	"github.com/imwxx/koiddns/util"
	daemon "github.com/sevlyar/go-daemon"
)

func main() {
	help := flag.Bool("help", false, "Show this help message and exit")
	demo := flag.String("generate-config", "", "Generate a sample configuration file")
	configFile := flag.String("config", "/etc/config/koiddns", "Specify the configuration file")
	daemonize := flag.Bool("daemon", false, "Run as a daemon")

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

	if *daemonize {
		cntxt := &daemon.Context{
			PidFileName: "/var/run/koiddns.pid",
			PidFilePerm: 0644,
			LogFileName: "/var/log/koiddns.log",
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

	// 正常运行代码
	run(configFile)

}

func run(configFile *string) {
	cfg, err := config.LoadConfig(*configFile)
	if err != nil {
		log.Fatal(err)
	}

	ticker := time.NewTicker(time.Duration(cfg.Main.ExecutionCycleMinutes) * time.Minute)
	defer ticker.Stop()

	updateDNS(cfg) // Initial update

	for range ticker.C {
		updateDNS(cfg)
	}

}

func updateDNS(cfg *config.Config) {

	ip := util.GetPublicIP()

	providerConfigs := map[string]config.ProviderConfig{}

	for _, provider := range cfg.Providers {
		providerConfigs[provider.Name] = provider
	}

	for _, domain := range cfg.Domains {
		providerConfig, ok := providerConfigs[domain.Provider]
		if !ok {
			log.Fatalf("Unknown DNS provider: %s", domain.Provider)
		}

		if domain.Value != "" {
			ip = domain.Value
		}

		switch domain.Provider {
		case "aliyun":
			dns.UpdateAliyunDNS(providerConfig, domain.SubDomain, domain.PrimaryDomain, domain.RecordType, domain.RecordId, domain.Line, domain.Priority, ip)
		case "tencent":
			dns.UpdateTencentDNS(providerConfig, domain.SubDomain, domain.PrimaryDomain, domain.RecordType, domain.RecordId, domain.Line, domain.Priority, ip)
		default:
			log.Fatalf("Unsupported DNS provider: %s", domain.Provider)
		}
	}
}
