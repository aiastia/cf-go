package main

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "列出所有DNS记录",
	Long:  `列出所有域名的DNS记录，支持排序和筛选`,
	Run: func(cmd *cobra.Command, args []string) {
		cf, err := NewCloudflareManager()
		if err != nil {
			fmt.Printf("错误: %v\n", err)
			return
		}

		records, err := cf.GetAllDNSRecords()
		if err != nil {
			fmt.Printf("获取DNS记录失败: %v\n", err)
			return
		}

		// 应用筛选
		filterName, _ := cmd.Flags().GetString("filter-name")
		filterType, _ := cmd.Flags().GetString("filter-type")
		filterZone, _ := cmd.Flags().GetString("filter-zone")
		filterContent, _ := cmd.Flags().GetString("filter-content")

		filters := make(map[string]string)
		if filterName != "" {
			filters["name"] = filterName
		}
		if filterType != "" {
			filters["type"] = filterType
		}
		if filterZone != "" {
			filters["zone"] = filterZone
		}
		if filterContent != "" {
			filters["content"] = filterContent
		}

		if len(filters) > 0 {
			records = FilterRecords(records, filters)
		}

		// 应用排序
		sortBy, _ := cmd.Flags().GetString("sort-by")
		ascending, _ := cmd.Flags().GetBool("ascending")
		SortRecords(records, sortBy, ascending)

		// 显示记录
		displayRecords(records)
	},
}

var addCmd = &cobra.Command{
	Use:   "add [域名] [记录名] [类型] [内容]",
	Short: "添加DNS记录",
	Long:  `为指定域名添加新的DNS记录`,
	Args:  cobra.ExactArgs(4),
	Run: func(cmd *cobra.Command, args []string) {
		cf, err := NewCloudflareManager()
		if err != nil {
			fmt.Printf("错误: %v\n", err)
			return
		}

		zoneName := args[0]
		name := args[1]
		recordType := strings.ToUpper(args[2])
		content := args[3]

		ttl, _ := cmd.Flags().GetInt("ttl")
		proxied, _ := cmd.Flags().GetBool("proxied")

		err = cf.AddDNSRecord(zoneName, name, recordType, content, ttl, proxied)
		if err != nil {
			fmt.Printf("添加DNS记录失败: %v\n", err)
			return
		}

		fmt.Printf("成功添加DNS记录: %s.%s -> %s\n", name, zoneName, content)
	},
}

var updateCmd = &cobra.Command{
	Use:   "update [记录ID] [域名] [记录名] [类型] [内容]",
	Short: "更新DNS记录",
	Long:  `更新指定的DNS记录`,
	Args:  cobra.ExactArgs(5),
	Run: func(cmd *cobra.Command, args []string) {
		cf, err := NewCloudflareManager()
		if err != nil {
			fmt.Printf("错误: %v\n", err)
			return
		}

		recordID := args[0]
		zoneName := args[1]
		name := args[2]
		recordType := strings.ToUpper(args[3])
		content := args[4]

		// 获取ZoneID
		zones, err := cf.Client.ListZones(cmd.Context())
		if err != nil {
			fmt.Printf("获取域名列表失败: %v\n", err)
			return
		}

		var zoneID string
		for _, zone := range zones {
			if zone.Name == zoneName {
				zoneID = zone.ID
				break
			}
		}

		if zoneID == "" {
			fmt.Printf("域名 %s 不存在\n", zoneName)
			return
		}

		ttl, _ := cmd.Flags().GetInt("ttl")
		proxied, _ := cmd.Flags().GetBool("proxied")

		err = cf.UpdateDNSRecord(recordID, zoneID, name, recordType, content, ttl, proxied)
		if err != nil {
			fmt.Printf("更新DNS记录失败: %v\n", err)
			return
		}

		fmt.Printf("成功更新DNS记录: %s\n", recordID)
	},
}

var deleteCmd = &cobra.Command{
	Use:   "delete [记录ID] [域名]",
	Short: "删除DNS记录",
	Long:  `删除指定的DNS记录`,
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		cf, err := NewCloudflareManager()
		if err != nil {
			fmt.Printf("错误: %v\n", err)
			return
		}

		recordID := args[0]
		zoneName := args[1]

		// 获取ZoneID
		zones, err := cf.Client.ListZones(cmd.Context())
		if err != nil {
			fmt.Printf("获取域名列表失败: %v\n", err)
			return
		}

		var zoneID string
		for _, zone := range zones {
			if zone.Name == zoneName {
				zoneID = zone.ID
				break
			}
		}

		if zoneID == "" {
			fmt.Printf("域名 %s 不存在\n", zoneName)
			return
		}

		err = cf.DeleteDNSRecord(recordID, zoneID)
		if err != nil {
			fmt.Printf("删除DNS记录失败: %v\n", err)
			return
		}

		fmt.Printf("成功删除DNS记录: %s\n", recordID)
	},
}

var interactiveCmd = &cobra.Command{
	Use:   "interactive",
	Short: "交互式DNS管理",
	Long:  `启动交互式界面来管理DNS记录`,
	Run: func(cmd *cobra.Command, args []string) {
		startInteractiveMode()
	},
}

func init() {
	// list命令的标志
	listCmd.Flags().String("filter-name", "", "按记录名筛选")
	listCmd.Flags().String("filter-type", "", "按记录类型筛选")
	listCmd.Flags().String("filter-zone", "", "按域名筛选")
	listCmd.Flags().String("filter-content", "", "按内容筛选")
	listCmd.Flags().String("sort-by", "zone", "排序字段 (name, type, zone, content, ttl, created, modified)")
	listCmd.Flags().Bool("ascending", true, "升序排序")

	// add命令的标志
	addCmd.Flags().Int("ttl", 1, "TTL值 (1=自动)")
	addCmd.Flags().Bool("proxied", false, "是否启用Cloudflare代理")

	// update命令的标志
	updateCmd.Flags().Int("ttl", 1, "TTL值 (1=自动)")
	updateCmd.Flags().Bool("proxied", false, "是否启用Cloudflare代理")
}

func displayRecords(records []DNSRecord) {
	if len(records) == 0 {
		fmt.Println("没有找到DNS记录")
		return
	}

	fmt.Printf("找到 %d 条DNS记录:\n\n", len(records))
	fmt.Printf("%-36s %-20s %-15s %-8s %-50s %-6s %-8s\n", 
		"记录ID", "域名", "记录名", "类型", "内容", "TTL", "代理")
	fmt.Println(strings.Repeat("-", 150))

	for _, record := range records {
		proxiedStr := "否"
		if record.Proxied {
			proxiedStr = "是"
		}

		ttlStr := strconv.Itoa(record.TTL)
		if record.TTL == 1 {
			ttlStr = "自动"
		}

		// 截断过长的内容
		content := record.Content
		if len(content) > 47 {
			content = content[:44] + "..."
		}

		fmt.Printf("%-36s %-20s %-15s %-8s %-50s %-6s %-8s\n",
			record.ID, record.ZoneName, record.Name, record.Type, content, ttlStr, proxiedStr)
	}
} 