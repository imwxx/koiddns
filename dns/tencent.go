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

func UpdateTencentDNS(config config.ProviderConfig, subDomain, primaryDomain, recordType, recordId, line, priority, ip string) {
	credential := common.NewCredential(config.SecretId, config.SecretKey)
	cpf := profile.NewClientProfile()
	cpf.HttpProfile.Endpoint = "dnspod.tencentcloudapi.com"
	client, _ := dnspod.NewClient(credential, "", cpf)

	cacheKey := generateCacheKey("tencent", subDomain, primaryDomain, recordType)
	cachedRecord, _ := GetRecord(cacheKey)

	var existingValue string

	if cachedRecord != nil {
		recordId = cachedRecord.RecordId
		existingValue = cachedRecord.RecordValue
	}

	if recordId == "" {
		// 尝试获取现有的记录ID
		recordId = getTencentRecordId(client, subDomain, primaryDomain, recordType, line)
		if recordId != "" {
			// 保存到缓存
			SetRecord(cacheKey, recordId, ip)
		}
	}

	if recordId == "" {
		// 新增域名解析
		resp, err := addTencentRecord(client, subDomain, primaryDomain, recordType, line, priority, ip)
		if err != nil {
			log.Fatal(err)
		}
		recordId = strconv.FormatUint(*resp.Response.RecordId, 10)

		SetRecord(cacheKey, recordId, ip)
	} else {
		// 更新域名解析
		if existingValue != ip {
			updateTencentRecord(client, recordId, subDomain, primaryDomain, recordType, line, priority, ip)
			SetRecord(cacheKey, recordId, ip)
		}
	}
}

func getTencentRecordId(client *dnspod.Client, subDomain, primaryDomain, recordType, line string) string {
	request := dnspod.NewDescribeRecordListRequest()

	request.Domain = common.StringPtr(primaryDomain)
	request.RecordType = common.StringPtr(recordType)
	request.RecordLine = common.StringPtr(line)

	response, err := client.DescribeRecordList(request)
	if err != nil {
		log.Fatal(err)
	}

	for _, record := range response.Response.RecordList {
		if *record.Name == subDomain {
			return strconv.FormatUint(*record.RecordId, 10)
		}
	}
	return ""
}

func updateTencentRecord(client *dnspod.Client, recordId, subDomain, primaryDomain, recordType, line, priority, ip string) {
	recordIdUint64, err := strconv.ParseUint(recordId, 10, 64)
	if err != nil {
		log.Fatalf("无法将 recordId 转换为 uint64: %v", err)
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
			log.Fatalf("无法将 priority 转换为 uint64: %v", err)
		}
		request.Weight = common.Uint64Ptr(uint64(prio))
	}

	_, err = client.ModifyRecord(request)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Updated Tencent DNS record %s (%s) to %s with line %s and priority %s", fmt.Sprintf(`%s.%s`, subDomain, primaryDomain), recordType, ip, line, priority)
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
			log.Fatalf("无法将 priority 转换为 uint64: %v", err)
		}
		request.Weight = common.Uint64Ptr(uint64(prio))
	}

	resp, err := client.CreateRecord(request)
	if err != nil {
		return resp, err
	}

	log.Printf("Added Tencent DNS record %s (%s) with value %s, line %s and priority %s", fmt.Sprintf(`%s.%s`, subDomain, primaryDomain), recordType, ip, line, priority)

	return resp, nil
}
