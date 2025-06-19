package main

import (
	"context"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/cloudflare/cloudflare-go"
)

type DNSRecord struct {
	ID       string
	ZoneID   string
	ZoneName string
	Name     string
	Type     string
	Content  string
	TTL      int
	Proxied  bool
	Created  time.Time
	Modified time.Time
}

type CloudflareManager struct {
	Client *cloudflare.API
}

func NewCloudflareManager() (*CloudflareManager, error) {
	if err := initConfig(); err != nil {
		return nil, err
	}

	client, err := cloudflare.NewWithAPIToken(config.CloudflareToken)
	if err != nil {
		return nil, fmt.Errorf("创建Cloudflare客户端失败: %w", err)
	}

	return &CloudflareManager{Client: client}, nil
}

// GetAllDNSRecords 获取所有域名的DNS记录
func (cf *CloudflareManager) GetAllDNSRecords() ([]DNSRecord, error) {
	ctx := context.Background()
	
	// 获取所有域名
	zones, err := cf.Client.ListZones(ctx)
	if err != nil {
		return nil, fmt.Errorf("获取域名列表失败: %w", err)
	}

	var allRecords []DNSRecord

	for _, zone := range zones {
		// 获取该域名的所有DNS记录
		records, _, err := cf.Client.ListDNSRecords(ctx, cloudflare.ZoneIdentifier(zone.ID), cloudflare.ListDNSRecordsParams{})
		if err != nil {
			fmt.Printf("警告: 获取域名 %s 的DNS记录失败: %v\n", zone.Name, err)
			continue
		}

		for _, record := range records {
			dnsRecord := DNSRecord{
				ID:       record.ID,
				ZoneID:   zone.ID,
				ZoneName: zone.Name,
				Name:     record.Name,
				Type:     record.Type,
				Content:  record.Content,
				TTL:      record.TTL,
				Proxied:  record.Proxied != nil && *record.Proxied,
				Created:  record.CreatedOn,
				Modified: record.ModifiedOn,
			}
			allRecords = append(allRecords, dnsRecord)
		}
	}

	return allRecords, nil
}

// AddDNSRecord 添加DNS记录
func (cf *CloudflareManager) AddDNSRecord(zoneName, name, recordType, content string, ttl int, proxied bool) error {
	ctx := context.Background()

	// 查找域名
	zones, err := cf.Client.ListZones(ctx)
	if err != nil {
		return fmt.Errorf("获取域名列表失败: %w", err)
	}

	var zoneID string
	for _, zone := range zones {
		if zone.Name == zoneName {
			zoneID = zone.ID
			break
		}
	}

	if zoneID == "" {
		return fmt.Errorf("域名 %s 不存在", zoneName)
	}

	// 创建DNS记录
	params := cloudflare.CreateDNSRecordParams{
		Type:    recordType,
		Name:    name,
		Content: content,
		TTL:     ttl,
		Proxied: &proxied,
	}

	_, err = cf.Client.CreateDNSRecord(ctx, cloudflare.ZoneIdentifier(zoneID), params)
	if err != nil {
		return fmt.Errorf("创建DNS记录失败: %w", err)
	}

	return nil
}

// UpdateDNSRecord 更新DNS记录
func (cf *CloudflareManager) UpdateDNSRecord(recordID, zoneID, name, recordType, content string, ttl int, proxied bool) error {
	ctx := context.Background()

	_, err := cf.Client.UpdateDNSRecord(ctx, cloudflare.ZoneIdentifier(zoneID), cloudflare.UpdateDNSRecordParams{
		ID:      recordID,
		Type:    recordType,
		Name:    name,
		Content: content,
		TTL:     ttl,
		Proxied: &proxied,
	})

	if err != nil {
		return fmt.Errorf("更新DNS记录失败: %w", err)
	}

	return nil
}

// DeleteDNSRecord 删除DNS记录
func (cf *CloudflareManager) DeleteDNSRecord(recordID, zoneID string) error {
	ctx := context.Background()

	err := cf.Client.DeleteDNSRecord(ctx, cloudflare.ZoneIdentifier(zoneID), recordID)
	if err != nil {
		return fmt.Errorf("删除DNS记录失败: %w", err)
	}

	return nil
}

// SortRecords 排序DNS记录
func SortRecords(records []DNSRecord, sortBy string, ascending bool) {
	sort.Slice(records, func(i, j int) bool {
		var result bool
		switch strings.ToLower(sortBy) {
		case "name":
			result = records[i].Name < records[j].Name
		case "type":
			result = records[i].Type < records[j].Type
		case "zone":
			result = records[i].ZoneName < records[j].ZoneName
		case "content":
			result = records[i].Content < records[j].Content
		case "ttl":
			result = records[i].TTL < records[j].TTL
		case "created":
			result = records[i].Created.Before(records[j].Created)
		case "modified":
			result = records[i].Modified.Before(records[j].Modified)
		default:
			result = records[i].ZoneName < records[j].ZoneName
		}
		
		if !ascending {
			result = !result
		}
		return result
	})
}

// FilterRecords 筛选DNS记录
func FilterRecords(records []DNSRecord, filters map[string]string) []DNSRecord {
	var filtered []DNSRecord

	for _, record := range records {
		match := true
		for key, value := range filters {
			switch strings.ToLower(key) {
			case "name":
				if !strings.Contains(strings.ToLower(record.Name), strings.ToLower(value)) {
					match = false
				}
			case "type":
				if !strings.EqualFold(record.Type, value) {
					match = false
				}
			case "zone":
				if !strings.Contains(strings.ToLower(record.ZoneName), strings.ToLower(value)) {
					match = false
				}
			case "content":
				if !strings.Contains(strings.ToLower(record.Content), strings.ToLower(value)) {
					match = false
				}
			}
		}
		if match {
			filtered = append(filtered, record)
		}
	}

	return filtered
} 