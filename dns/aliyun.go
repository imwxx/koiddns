package dns

import (
	"fmt"
	"log"
	"strconv"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/alidns"
	"github.com/imwxx/koiddns/config"
)

func UpdateAliyunDNS(config config.ProviderConfig, subDomain, primaryDomain, recordType, recordId, line, priority, ip string) {
	client, err := alidns.NewClientWithAccessKey("cn-hangzhou", config.AccessKeyId, config.AccessKeySecret)
	if err != nil {
		log.Fatal(err)
	}

	cacheKey := generateCacheKey("aliyun", subDomain, primaryDomain, recordType)
	cachedRecord, _ := GetRecord(cacheKey)

	var existingValue string

	if cachedRecord != nil {
		recordId = cachedRecord.RecordId
		existingValue = cachedRecord.RecordValue
	}

	if recordId == "" {
		// 尝试获取现有的记录ID
		recordId = getAliyunRecordId(client, subDomain, primaryDomain, recordType, line)
		if recordId != "" {
			// 保存到缓存
			SetRecord(cacheKey, recordId, ip)
		}
	}

	if recordId == "" {
		// 新增域名解析
		resp, err := addAliyunRecord(client, subDomain, primaryDomain, recordType, line, priority, ip)
		if err != nil {
			log.Fatal(err)
		}

		SetRecord(cacheKey, resp.RecordId, ip)
	} else {
		// 更新域名解析
		if existingValue != ip && existingValue != "" {
			updateAliyunRecord(client, recordId, subDomain, primaryDomain, recordType, line, priority, ip)
		} else {
			log.Printf("Aliyun DNS record %s (%s) with value %s, don't need to update", fmt.Sprintf(`%s.%s`, subDomain, primaryDomain), recordType, ip)
		}
		SetRecord(cacheKey, recordId, ip)
	}
}

func getAliyunRecordId(client *alidns.Client, subDomain, primaryDomain, recordType, line string) string {
	request := alidns.CreateDescribeDomainRecordsRequest()
	request.DomainName = primaryDomain
	request.RRKeyWord = subDomain
	request.Type = recordType
	request.Line = line

	response, err := client.DescribeDomainRecords(request)
	if err != nil {
		log.Fatal(err)
	}

	for _, record := range response.DomainRecords.Record {
		if record.RR == subDomain {
			return record.RecordId
		}
	}

	return ""
}

func updateAliyunRecord(client *alidns.Client, recordId, subDomain, primaryDomain, recordType, line, priority, ip string) {
	request := alidns.CreateUpdateDomainRecordRequest()
	request.Scheme = "https"
	request.RecordId = recordId
	request.RR = subDomain
	request.Type = recordType
	request.Value = ip
	request.Line = line
	fmt.Println(request)
	if priority != "" {
		prio, err := strconv.Atoi(priority)
		if err != nil {
			log.Fatalf("无法转换优先级为整数: %v", err)
		}
		request.Priority = requests.NewInteger(prio)
	}

	_, err := client.UpdateDomainRecord(request)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Updated Aliyun DNS record %s (%s) to %s with line %s and priority %s", fmt.Sprintf(`%s.%s`, subDomain, primaryDomain), recordType, ip, line, priority)
}

func addAliyunRecord(client *alidns.Client, subDomain, primaryDomain, recordType, line, priority, ip string) (*alidns.AddDomainRecordResponse, error) {
	request := alidns.CreateAddDomainRecordRequest()
	request.DomainName = primaryDomain
	request.RR = subDomain
	request.Type = recordType
	request.Value = ip
	request.Line = line
	if priority != "" {
		prio, err := strconv.Atoi(priority)
		if err != nil {
			log.Fatalf("无法转换优先级为整数: %v", err)
		}
		request.Priority = requests.NewInteger(prio)
	}

	response, err := client.AddDomainRecord(request)
	if err != nil {
		return response, err
	}

	log.Printf("Added Aliyun DNS record %s (%s) with value %s, line %s and priority %s", fmt.Sprintf(`%s.%s`, subDomain, primaryDomain), recordType, ip, line, priority)

	return response, nil
}
