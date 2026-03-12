package dns

import (
	"fmt"
	"log"
	"strconv"

	"github.com/imwxx/koiddns/config"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	dnspod "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/dnspod/v20210323"
)

func UpdateTencentDNS(c config.ProviderConfig, subDomain, primaryDomain, recordType, recordId, line, priority, ip string) error {
	credential := common.NewCredential(c.SecretId, c.SecretKey)
	cpf := profile.NewClientProfile()
	cpf.HttpProfile.Endpoint = "dnspod.tencentcloudapi.com"
	client, err := dnspod.NewClient(credential, "", cpf)
	if err != nil {
		return fmt.Errorf("create tencent dnspod client: %w", err)
	}

	cacheKey := generateCacheKey("tencent", subDomain, primaryDomain, recordType)
	cachedRecord, _ := GetRecord(cacheKey)

	var existingValue string
	if cachedRecord != nil {
		recordId = cachedRecord.RecordId
		existingValue = cachedRecord.RecordValue
	}

	if recordId == "" {
		recordId, err = getTencentRecordId(client, subDomain, primaryDomain, recordType, line)
		if err != nil {
			return err
		}
		if recordId != "" {
			SetRecord(cacheKey, recordId, ip)
		}
	}

	if recordId == "" {
		resp, err := addTencentRecord(client, subDomain, primaryDomain, recordType, line, priority, ip)
		if err != nil {
			return err
		}
		recordId = strconv.FormatUint(*resp.Response.RecordId, 10)
		SetRecord(cacheKey, recordId, ip)
		return nil
	}

	if existingValue != ip {
		if err := updateTencentRecord(client, recordId, subDomain, primaryDomain, recordType, line, priority, ip); err != nil {
			return err
		}
		log.Printf("Updated Tencent DNS record %s.%s (%s) to %s", subDomain, primaryDomain, recordType, ip)
	}
	SetRecord(cacheKey, recordId, ip)
	return nil
}

func getTencentRecordId(client *dnspod.Client, subDomain, primaryDomain, recordType, line string) (string, error) {
	request := dnspod.NewDescribeRecordListRequest()
	request.Domain = common.StringPtr(primaryDomain)
	request.RecordType = common.StringPtr(recordType)
	request.RecordLine = common.StringPtr(line)

	response, err := client.DescribeRecordList(request)
	if err != nil {
		return "", fmt.Errorf("describe record list: %w", err)
	}

	for _, record := range response.Response.RecordList {
		if record.Name != nil && *record.Name == subDomain && record.RecordId != nil {
			return strconv.FormatUint(*record.RecordId, 10), nil
		}
	}
	return "", nil
}

func updateTencentRecord(client *dnspod.Client, recordId, subDomain, primaryDomain, recordType, line, priority, ip string) error {
	recordIdUint64, err := strconv.ParseUint(recordId, 10, 64)
	if err != nil {
		return fmt.Errorf("invalid recordId %q: %w", recordId, err)
	}

	request := &dnspod.ModifyRecordRequest{}
	request.RecordId = common.Uint64Ptr(recordIdUint64)
	request.Domain = common.StringPtr(primaryDomain)
	request.RecordType = common.StringPtr(recordType)
	request.RecordLine = common.StringPtr(line)
	request.Value = common.StringPtr(ip)
	if priority != "" {
		prio, err := strconv.ParseUint(priority, 10, 64)
		if err != nil {
			return fmt.Errorf("invalid priority %q: %w", priority, err)
		}
		request.Weight = common.Uint64Ptr(prio)
	}

	_, err = client.ModifyRecord(request)
	if err != nil {
		return fmt.Errorf("modify record: %w", err)
	}
	return nil
}

func addTencentRecord(client *dnspod.Client, subDomain, primaryDomain, recordType, line, priority, ip string) (*dnspod.CreateRecordResponse, error) {
	request := &dnspod.CreateRecordRequest{}
	request.SubDomain = common.StringPtr(subDomain)
	request.Domain = common.StringPtr(primaryDomain)
	request.RecordType = common.StringPtr(recordType)
	request.RecordLine = common.StringPtr(line)
	request.Value = common.StringPtr(ip)
	if priority != "" {
		prio, err := strconv.ParseUint(priority, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid priority %q: %w", priority, err)
		}
		request.Weight = common.Uint64Ptr(prio)
	}

	resp, err := client.CreateRecord(request)
	if err != nil {
		return nil, fmt.Errorf("create record: %w", err)
	}
	log.Printf("Added Tencent DNS record %s.%s (%s) value %s", subDomain, primaryDomain, recordType, ip)
	return resp, nil
}
