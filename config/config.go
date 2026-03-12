package config

import (
	"errors"
	"fmt"
	"os"

	"gopkg.in/yaml.v2"
)

type ProviderConfig struct {
	Name            string `yaml:"name"`
	AccessKeyId     string `yaml:"accessKeyId,omitempty"`
	AccessKeySecret string `yaml:"accessKeySecret,omitempty"`
	SecretId        string `yaml:"secretId,omitempty"`
	SecretKey       string `yaml:"secretKey,omitempty"`
}

type DomainConfig struct {
	SubDomain     string `yaml:"subDomain"`
	PrimaryDomain string `yaml:"primaryDomain"`
	Value         string `yaml:"value"`
	Provider      string `yaml:"provider"`
	RecordType    string `yaml:"recordType"`
	RecordId      string `yaml:"recordId"`
	Line          string `yaml:"line"`
	Priority      string `yaml:"priority"`
}

type MainConfig struct {
	ExecutionCycleMinutes int `yaml:"executionCycleMinutes"`
}

type Config struct {
	Main      MainConfig       `yaml:"main"`
	Providers []ProviderConfig `yaml:"providers"`
	Domains   []DomainConfig   `yaml:"domains"`
}

func LoadConfig(file string) (*Config, error) {
	data, err := os.ReadFile(file)
	if err != nil {
		return nil, err
	}

	var cfg Config
	err = yaml.Unmarshal(data, &cfg)
	if err != nil {
		return nil, err
	}
	if err := Validate(&cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}

// Validate 校验配置合法性。
func Validate(cfg *Config) error {
	if cfg.Main.ExecutionCycleMinutes <= 0 {
		return errors.New("main.executionCycleMinutes must be positive")
	}
	if len(cfg.Domains) == 0 {
		return errors.New("domains cannot be empty")
	}
	providerNames := make(map[string]bool)
	for _, p := range cfg.Providers {
		if p.Name == "" {
			return errors.New("provider name cannot be empty")
		}
		providerNames[p.Name] = true
	}
	for i, d := range cfg.Domains {
		if d.SubDomain == "" {
			return fmt.Errorf("domains[%d]: subDomain is required", i)
		}
		if d.PrimaryDomain == "" {
			return fmt.Errorf("domains[%d]: primaryDomain is required", i)
		}
		if d.RecordType == "" {
			return fmt.Errorf("domains[%d]: recordType is required", i)
		}
		if d.Provider == "" {
			return fmt.Errorf("domains[%d]: provider is required", i)
		}
		if !providerNames[d.Provider] {
			return fmt.Errorf("domains[%d]: unknown provider %q", i, d.Provider)
		}
	}
	return nil
}

type Line struct {
	Id          string `yaml:"id"`
	Value       string `yaml:"value"`
	Description string `yaml:"description"`
}

func LoadProviderLines() map[string][]Line {
	data := map[string][]Line{}

	var aliyun []Line
	ali := map[string]string{
		"default":  "默认",
		"telecom":  "中国电信",
		"unicom":   "中国联通",
		"mobile":   "中国移动",
		"oversea":  "境外",
		"edu":      "中国教育网",
		"drpeng":   "中国鹏博士",
		"btvn":     "中国广电网",
		"aliyun":   "阿里云",
		"search":   "搜索引擎",
		"internal": "中国地区",
	}
	for k, v := range ali {
		t := Line{
			Id:          k,
			Value:       k,
			Description: v,
		}
		aliyun = append(aliyun, t)
	}

	var tencent []Line
	tens := map[string]string{
		"default":  "默认",
		"telecom":  "电信",
		"unicom":   "联通",
		"mobile":   "移动",
		"oversea":  "境外",
		"edu":      "教育网",
		"drpeng":   "鹏博士",
		"btvn":     "广电",
		"aliyun":   "阿里云",
		"search":   "搜索引擎",
		"internal": "中国",
	}

	for k, v := range tens {
		t := Line{
			Id:          k,
			Value:       k,
			Description: v,
		}
		tencent = append(tencent, t)
	}
	data["aliyun"] = aliyun
	data["tencent"] = tencent

	return data
}
