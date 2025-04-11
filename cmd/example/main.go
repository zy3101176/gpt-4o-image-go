package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/zy3101176/gpt-4o-image-go/internal/client"
	"github.com/zy3101176/gpt-4o-image-go/pkg/models"
)

func main() {
	// 解析命令行参数
	apiKey := flag.String("api-key", "", "API密钥")
	apiURL := flag.String("api-url", "https://api.tu-zi.com/v1/chat/completions", "API URL")
	model := flag.String("model", "gpt-4o-image-vip", "模型名称")
	prompt := flag.String("prompt", "描述这张图片", "提示词")
	outputDir := flag.String("output", "output", "输出目录")
	timeout := flag.Int("timeout", 60, "超时时间（秒）")
	rate := flag.Int("rate", 500, "请求限流（毫秒）")
	maxWorkers := flag.Int("workers", 5, "最大工作协程数")
	flag.Parse()

	// 检查API密钥
	if *apiKey == "" {
		apiKey = getEnv("API_KEY", "")
		if *apiKey == "" {
			log.Fatal("未提供API密钥. 请使用 -api-key 选项或设置 API_KEY 环境变量")
		}
	}

	// 获取图片路径参数
	images := flag.Args()
	if len(images) == 0 {
		log.Fatal("未提供图片路径. 请作为命令行参数提供一个或多个图片路径")
	}

	// 检查文件是否存在
	for _, img := range images {
		if _, err := os.Stat(img); os.IsNotExist(err) {
			log.Fatalf("图片文件不存在: %s", img)
		}
	}

	// 创建客户端
	client, err := client.NewClient(
		client.WithAPIKey(*apiKey),
		client.WithAPIURL(*apiURL),
		client.WithRateLimit(time.Duration(*rate)*time.Millisecond),
		client.WithMaxWorkers(*maxWorkers),
	)
	if err != nil {
		log.Fatalf("创建客户端失败: %v", err)
	}

	// 创建任务
	tasks := make([]*models.Task, len(images))
	for i, img := range images {
		tasks[i] = &models.Task{
			Name:   filepath.Base(img),
			Prompt: *prompt,
			Images: []string{img},
			Model:  *model,
		}
	}

	// 创建下载选项
	downloadOptions := &models.DownloadOptions{
		OutputDir: *outputDir,
		Timeout:   *timeout,
	}

	// 创建上下文
	ctx := context.Background()

	log.Printf("正在处理 %d 个任务，请等待...\n", len(tasks))
	startTime := time.Now()

	// 第一步：处理任务
	results := client.ProcessTasks(ctx, tasks)

	// 第二步：下载图片
	var successCount, failCount, downloadCount int
	for _, result := range results {
		if !result.Success {
			failCount++
			continue
		}

		successCount++
		err := client.DownloadImages(ctx, result, downloadOptions)
		if err != nil {
			log.Printf("下载图片失败: %v", err)
		} else {
			downloadCount += len(result.ImagePaths)
		}
	}

	// 打印结果
	fmt.Printf("\n处理完成，耗时 %.2f 秒\n", time.Since(startTime).Seconds())
	fmt.Printf("任务统计：总计 %d 个，成功 %d 个，失败 %d 个，下载 %d 张图片\n\n",
		len(tasks), successCount, failCount, downloadCount)

	for i, result := range results {
		fmt.Printf("任务 %d: %s\n", i+1, result.TaskName)
		if !result.Success {
			fmt.Printf("  状态: 失败 (%v)\n", result.Error)
			continue
		}

		fmt.Printf("  状态: 成功\n")
		if len(result.ImagePaths) > 0 {
			fmt.Printf("  下载的图片:\n")
			for j, path := range result.ImagePaths {
				fmt.Printf("    图片 %d: %s\n", j+1, path)
			}
		} else {
			fmt.Printf("  没有下载任何图片\n")
		}
		fmt.Println()
	}
}

// 获取环境变量，如果不存在则返回默认值
func getEnv(key, defaultValue string) *string {
	value := os.Getenv(key)
	if value == "" {
		return &defaultValue
	}
	return &value
}
