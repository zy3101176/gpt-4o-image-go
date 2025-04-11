package models

// Option 定义客户端配置选项
type Option[T any] func(*T)

// Task 表示一个处理任务
type Task struct {
	Name   string   `json:"name"`
	Prompt string   `json:"prompt"`
	Images []string `json:"images"`
	Model  string   `json:"model"`
}

// TaskResult 表示任务处理结果
type TaskResult struct {
	TaskID    string   `json:"task_id"`
	TaskName  string   `json:"task_name"`
	Success   bool     `json:"success"`
	Content   string   `json:"content,omitempty"`
	ImageURLs []string `json:"image_urls,omitempty"`
	Error     error    `json:"error,omitempty"`
	// 新增字段，保存下载的图片路径
	ImagePaths []string `json:"image_paths,omitempty"`
}

// Message 表示API请求消息
type Message struct {
	Role    string        `json:"role"`
	Content []ContentItem `json:"content"`
}

// ContentItem 表示消息内容项
type ContentItem struct {
	Type     string    `json:"type"`
	Text     string    `json:"text,omitempty"`
	ImageURL *ImageURL `json:"image_url,omitempty"`
}

// ImageURL 表示图片URL
type ImageURL struct {
	URL string `json:"url"`
}

// APIResponse 表示API响应
type APIResponse struct {
	ID      string    `json:"id"`
	Choices []Choice  `json:"choices"`
	Error   *APIError `json:"error,omitempty"`
}

// Choice 表示API响应选项
type Choice struct {
	Index   int     `json:"index"`
	Message Message `json:"message"`
}

// APIError 表示API错误
type APIError struct {
	Message string `json:"message"`
	Type    string `json:"type"`
	Code    string `json:"code,omitempty"`
}

// DownloadOptions 下载选项
type DownloadOptions struct {
	// OutputDir 输出目录
	OutputDir string
	// Timeout 下载超时
	Timeout int
}
