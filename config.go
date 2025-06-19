package main

import (
	"fmt"
	"os"

	"github.com/spf13/viper"
)

type Config struct {
	CloudflareToken string `mapstructure:"cloudflare_token"`
	AccountID       string `mapstructure:"account_id"`
}

var config Config

func initConfig() error {
	// 设置配置文件路径
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("$HOME/.cf-dns-manager")

	// 从环境变量读取
	viper.AutomaticEnv()
	viper.BindEnv("cloudflare_token", "CF_API_TOKEN")
	viper.BindEnv("account_id", "CF_ACCOUNT_ID")

	// 如果配置文件不存在，创建默认配置
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// 创建默认配置文件
			viper.Set("cloudflare_token", "")
			viper.Set("account_id", "")
			
			if err := viper.SafeWriteConfig(); err != nil {
				return fmt.Errorf("无法创建配置文件: %w", err)
			}
			
			fmt.Println("已创建默认配置文件 config.yaml")
			fmt.Println("请设置你的 Cloudflare API Token 和 Account ID")
			fmt.Println("或者设置环境变量 CF_API_TOKEN 和 CF_ACCOUNT_ID")
			os.Exit(1)
		}
		return fmt.Errorf("读取配置文件失败: %w", err)
	}

	if err := viper.Unmarshal(&config); err != nil {
		return fmt.Errorf("解析配置文件失败: %w", err)
	}

	// 检查必要的配置
	if config.CloudflareToken == "" {
		return fmt.Errorf("Cloudflare API Token 未设置")
	}

	return nil
} 