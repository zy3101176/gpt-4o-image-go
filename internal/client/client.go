package client

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"sync"
	"time"

	"github.com/zy3101176/gpt-4o-image-go/internal/utils"
	"github.com/zy3101176/gpt-4o-image-go/pkg/models"
)

type client struct {
	options     *Options
	tokenBucket *utils.TokenBucket
	httpClient  *http.Client
}

// NewClient 创建新的客户端
func NewClient(opts ...models.Option[Options]) (models.IClient, error) {
	c := &client{
		options: defaultOptions(),
		httpClient: &http.Client{
			Timeout: 120 * time.Second,
		},
	}

	for _, opt := range opts {
		opt(c.options)
	}

	if c.options.APIKey == "" {
		return nil, fmt.Errorf("API key is required")
	}

	c.tokenBucket = utils.NewTokenBucket(c.options.RateLimit)
	return c, nil
}

func (c *client) ProcessTask(ctx context.Context, task *models.Task) (*models.TaskResult, error) {
	if err := c.tokenBucket.Wait(ctx); err != nil {
		return nil, err
	}

	result := &models.TaskResult{
		TaskID:   fmt.Sprintf("%d", time.Now().UnixNano()),
		TaskName: task.Name,
	}

	// 准备请求数据
	messageContent := []models.ContentItem{
		{
			Type: "text",
			Text: task.Prompt,
		},
	}

	// 处理图片
	for _, imagePath := range task.Images {
		imageData, err := c.prepareImageData(imagePath)
		if err != nil {
			result.Success = false
			result.Error = err
			return result, err
		}

		messageContent = append(messageContent, models.ContentItem{
			Type: "image_url",
			ImageURL: &models.ImageURL{
				URL: imageData,
			},
		})
	}

	// 构建请求
	reqBody := map[string]interface{}{
		"model": task.Model,
		"messages": []models.Message{
			{
				Role:    "user",
				Content: messageContent,
			},
		},
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		result.Success = false
		result.Error = err
		return result, err
	}

	// 创建请求
	req, err := http.NewRequestWithContext(ctx, "POST", c.options.APIURL, nil)
	if err != nil {
		result.Success = false
		result.Error = err
		return result, err
	}

	req.Header.Set("Authorization", "Bearer "+c.options.APIKey)
	req.Header.Set("Content-Type", "application/json")
	req.Body = io.NopCloser(bytes.NewBuffer(jsonData))

	// 发送请求
	resp, err := c.httpClient.Do(req)
	if err != nil {
		result.Success = false
		result.Error = err
		return result, err
	}
	defer resp.Body.Close()

	// 解析响应
	var apiResp models.APIResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		result.Success = false
		result.Error = err
		return result, err
	}

	if apiResp.Error != nil {
		result.Success = false
		result.Error = fmt.Errorf("API error: %s", apiResp.Error.Message)
		return result, result.Error
	}

	// 处理响应内容
	if len(apiResp.Choices) > 0 {
		content := apiResp.Choices[0].Message.Content[0].Text
		result.Content = content

		// 提取图片URL
		re := regexp.MustCompile(`!\[.*?\]\((https?://[^\s]+)\)`)
		matches := re.FindAllStringSubmatch(content, -1)
		for _, match := range matches {
			if len(match) > 1 {
				result.ImageURLs = append(result.ImageURLs, match[1])
			}
		}
	}

	result.Success = true
	return result, nil
}

func (c *client) ProcessTasks(ctx context.Context, tasks []*models.Task) []*models.TaskResult {
	results := make([]*models.TaskResult, len(tasks))
	var wg sync.WaitGroup

	for i, task := range tasks {
		wg.Add(1)
		go func(idx int, t *models.Task) {
			defer wg.Done()
			result, _ := c.ProcessTask(ctx, t)
			results[idx] = result
		}(i, task)
	}

	wg.Wait()
	return results
}

func (c *client) ProcessTaskAsync(ctx context.Context, task *models.Task) <-chan *models.TaskResult {
	resultChan := make(chan *models.TaskResult, 1)

	go func() {
		defer close(resultChan)
		result, err := c.ProcessTask(ctx, task)
		if err != nil {
			result.Error = err
			resultChan <- result
			return
		}
		resultChan <- result
	}()

	return resultChan
}

func (c *client) ProcessTasksAsync(ctx context.Context, tasks []*models.Task) <-chan *models.TaskResult {
	resultChan := make(chan *models.TaskResult, len(tasks))

	go func() {
		defer close(resultChan)
		var wg sync.WaitGroup

		for _, task := range tasks {
			wg.Add(1)
			go func(t *models.Task) {
				defer wg.Done()
				result, err := c.ProcessTask(ctx, t)
				if err != nil {
					result.Error = err
					resultChan <- result
					return
				}
				resultChan <- result
			}(task)
		}

		wg.Wait()
	}()

	return resultChan
}

func (c *client) prepareImageData(imagePath string) (string, error) {
	// 支持相对路径和绝对路径
	var fullPath string
	if !filepath.IsAbs(imagePath) {
		fullPath = filepath.Join("input", "images", imagePath)
	} else {
		fullPath = imagePath
	}

	// 读取图片文件
	imgData, err := os.ReadFile(fullPath)
	if err != nil {
		return "", fmt.Errorf("failed to read image file: %w", err)
	}

	// 检测图片类型
	imgType := "png"
	if filepath.Ext(imagePath) == ".jpg" || filepath.Ext(imagePath) == ".jpeg" {
		imgType = "jpeg"
	}

	// 转换为base64
	encoded := base64.StdEncoding.EncodeToString(imgData)
	return fmt.Sprintf("data:image/%s;base64,%s", imgType, encoded), nil
}
