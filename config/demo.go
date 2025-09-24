package config

import (
	"os"

	"gopkg.in/yaml.v2"
)

func GenerateSampleConfig(file string) error {
	sampleConfig := Config{
		Main: MainConfig{
			ExecutionCycleMinutes: 5,
		},
		Providers: []ProviderConfig{
			{
				Name:            "aliyun",
				AccessKeyId:     "your_aliyun_access_key_id",
				AccessKeySecret: "your_aliyun_access_key_secret",
			},
			{
				Name:      "tencent",
				SecretId:  "your_tencent_secret_id",
				SecretKey: "your_tencent_secret_key",
			},
		},
		Domains: []DomainConfig{
			{
				SubDomain:     "www",
				PrimaryDomain: "example.com",
				Provider:      "aliyun",
				RecordType:    "A",
				RecordId:      "recordId1",
				Line:          "default",
				Priority:      "10",
			},
			{
				SubDomain:     "mail",
				PrimaryDomain: "anotherdomain.com",
				Provider:      "tencent",
				RecordType:    "CNAME",
				RecordId:      "recordId2",
				Line:          "unicom",
				Priority:      "20",
			},
		},
	}

	data, err := yaml.Marshal(&sampleConfig)
	if err != nil {
		return err
	}

	return os.WriteFile(file, data, 0644)
}
