package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func main() {
	var rootCmd = &cobra.Command{
		Use:   "cf-dns-manager",
		Short: "Cloudflare DNS记录管理工具",
		Long: `一个用于管理Cloudflare DNS记录的命令行工具。
支持查询、添加、修改、删除DNS记录，并提供排序和筛选功能。`,
	}

	// 添加子命令
	rootCmd.AddCommand(listCmd)
	rootCmd.AddCommand(addCmd)
	rootCmd.AddCommand(updateCmd)
	rootCmd.AddCommand(deleteCmd)
	rootCmd.AddCommand(interactiveCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
} 