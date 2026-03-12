package dns

import (
	"fmt"
	"log"
	"strconv"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/alidns"
	"github.com/imwxx/koiddns/config"
)

func UpdateAliyunDNS(c config.ProviderConfig, subDomain, primaryDomain, recordType, recordId, line, priority, ip string) error {
	client, err := alidns.NewClientWithAccessKey("cn-hangzhou", c.AccessKeyId, c.AccessKeySecret)
	if err != nil {
		return fmt.Errorf("create aliyun client: %w", err)
	}

	cacheKey := generateCacheKey("aliyun", subDomain, primaryDomain, recordType)
	cachedRecord, _ := GetRecord(cacheKey)

	var existingValue string
	if cachedRecord != nil {
		recordId = cachedRecord.RecordId
		existingValue = cachedRecord.RecordValue
	}

	if recordId == "" {
		recordId, err = getAliyunRecordId(client, subDomain, primaryDomain, recordType, line)
		if err != nil {
			return err
		}
		if recordId != "" {
			SetRecord(cacheKey, recordId, ip)
		}
	}

	if recordId == "" {
		resp, err := addAliyunRecord(client, subDomain, primaryDomain, recordType, line, priority, ip)
		if err != nil {
			return err
		}
		SetRecord(cacheKey, resp.RecordId, ip)
		return nil
	}

	if existingValue != ip && existingValue != "" {
		if err := updateAliyunRecord(client, recordId, subDomain, primaryDomain, recordType, line, priority, ip); err != nil {
			return err
		}
		log.Printf("Updated Aliyun DNS record %s.%s (%s) to %s", subDomain, primaryDomain, recordType, ip)
	} else {
		log.Printf("Aliyun DNS record %s.%s (%s) value %s, skip update", subDomain, primaryDomain, recordType, ip)
	}
	SetRecord(cacheKey, recordId, ip)
	return nil
}

func getAliyunRecordId(client *alidns.Client, subDomain, primaryDomain, recordType, line string) (string, error) {
	request := alidns.CreateDescribeDomainRecordsRequest()
	request.DomainName = primaryDomain
	request.RRKeyWord = subDomain
	request.Type = recordType
	request.Line = line

	response, err := client.DescribeDomainRecords(request)
	if err != nil {
		return "", fmt.Errorf("describe domain records: %w", err)
	}

	for _, record := range response.DomainRecords.Record {
		if record.RR == subDomain {
			return record.RecordId, nil
		}
	}
	return "", nil
}

func updateAliyunRecord(client *alidns.Client, recordId, subDomain, primaryDomain, recordType, line, priority, ip string) error {
	request := alidns.CreateUpdateDomainRecordRequest()
	request.Scheme = "https"
	request.RecordId = recordId
	request.RR = subDomain
	request.Type = recordType
	request.Value = ip
	request.Line = line
	if priority != "" {
		prio, err := strconv.Atoi(priority)
		if err != nil {
			return fmt.Errorf("invalid priority %q: %w", priority, err)
		}
		request.Priority = requests.NewInteger(prio)
	}

	_, err := client.UpdateDomainRecord(request)
	if err != nil {
		return fmt.Errorf("update domain record: %w", err)
	}
	return nil
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
			return nil, fmt.Errorf("invalid priority %q: %w", priority, err)
		}
		request.Priority = requests.NewInteger(prio)
	}

	response, err := client.AddDomainRecord(request)
	if err != nil {
		return nil, fmt.Errorf("add domain record: %w", err)
	}
	log.Printf("Added Aliyun DNS record %s.%s (%s) value %s", subDomain, primaryDomain, recordType, ip)
	return response, nil
}
