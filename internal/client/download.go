package client

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"sync"
	"time"

	"github.com/zy3101176/gpt-4o-image-go/pkg/models"
)

// DownloadImages 下载任务结果中的图片
func (c *client) DownloadImages(ctx context.Context, result *models.TaskResult, options *models.DownloadOptions) error {
	if result == nil || !result.Success {
		return fmt.Errorf("无法下载图片：任务未成功完成")
	}

	// 如果未提供下载选项，则使用默认值
	if options == nil {
		options = &models.DownloadOptions{
			OutputDir: "output",
			Timeout:   60,
		}
	}

	// 确保输出目录存在
	outputDir := options.OutputDir
	if outputDir == "" {
		outputDir = "output"
	}
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return fmt.Errorf("创建输出目录失败: %w", err)
	}

	// 如果内容中已有提取的图片URL，直接使用
	imagesToDownload := result.ImageURLs

	// 如果没有已提取的URL，则尝试从内容中提取
	if len(imagesToDownload) == 0 && result.Content != "" {
		// 使用正则表达式提取markdown中的图片地址
		re := regexp.MustCompile(`!\[.*?\]\((https?://[^\s]+)\)`)
		matches := re.FindAllStringSubmatch(result.Content, -1)
		for _, match := range matches {
			if len(match) > 1 {
				imagesToDownload = append(imagesToDownload, match[1])
			}
		}
	}

	if len(imagesToDownload) == 0 {
		return fmt.Errorf("未找到可下载的图片URL")
	}

	// 设置HTTP客户端超时
	httpClient := &http.Client{
		Timeout: time.Duration(options.Timeout) * time.Second,
	}

	// 用于存储下载的图片路径
	var imagePaths []string
	var downloadErrors []error
	var mutex sync.Mutex
	var wg sync.WaitGroup

	// 并发下载图片
	for i, imageURL := range imagesToDownload {
		wg.Add(1)
		go func(idx int, url string) {
			defer wg.Done()

			// 创建下载文件名
			fileName := fmt.Sprintf("%s-%d.png", result.TaskID, idx)
			filePath := filepath.Join(outputDir, fileName)

			// 下载图片
			req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
			if err != nil {
				mutex.Lock()
				downloadErrors = append(downloadErrors, fmt.Errorf("创建请求失败: %w", err))
				mutex.Unlock()
				return
			}

			resp, err := httpClient.Do(req)
			if err != nil {
				mutex.Lock()
				downloadErrors = append(downloadErrors, fmt.Errorf("下载图片失败: %w", err))
				mutex.Unlock()
				return
			}
			defer resp.Body.Close()

			if resp.StatusCode != http.StatusOK {
				mutex.Lock()
				downloadErrors = append(downloadErrors, fmt.Errorf("下载图片失败，状态码: %d", resp.StatusCode))
				mutex.Unlock()
				return
			}

			// 创建文件
			file, err := os.Create(filePath)
			if err != nil {
				mutex.Lock()
				downloadErrors = append(downloadErrors, fmt.Errorf("创建文件失败: %w", err))
				mutex.Unlock()
				return
			}
			defer file.Close()

			// 写入图片数据
			_, err = io.Copy(file, resp.Body)
			if err != nil {
				mutex.Lock()
				downloadErrors = append(downloadErrors, fmt.Errorf("写入文件失败: %w", err))
				mutex.Unlock()
				return
			}

			// 添加到结果路径列表
			mutex.Lock()
			imagePaths = append(imagePaths, filePath)
			mutex.Unlock()
		}(i, imageURL)
	}

	// 等待所有下载完成
	wg.Wait()

	// 保存下载的图片路径
	result.ImagePaths = imagePaths

	// 如果有错误，返回第一个错误
	if len(downloadErrors) > 0 {
		return fmt.Errorf("下载过程中发生 %d 个错误，第一个错误: %w", len(downloadErrors), downloadErrors[0])
	}

	return nil
}
