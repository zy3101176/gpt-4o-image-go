# 🚀 GPT-4O Go SDK

一个用于调用GPT-4O图像处理API的Go语言SDK，支持同步和异步处理模式。

## ✨ 特性

- 🔄 支持同步和异步处理模式
- 🚦 内置请求限流机制
- 🚀 支持并发处理多个任务
- 🖼️ 支持图片处理（jpg和png格式）
- ⚠️ 完整的错误处理机制
- ⚙️ 可配置的客户端选项
- 📥 支持图片下载功能

## 📦 安装

```bash
go get github.com/zy3101176/gpt-4o-image-go
```

## 🚀 快速开始

### 基础用法

```go
package main

import (
    "context"
    "fmt"
    "log"
    "time"

    "github.com/zy3101176/gpt-4o-image-go/pkg/models"
)

func main() {
    // 创建客户端
    client, err := gpt4o.NewClient(
        gpt4o.WithAPIKey("your-api-key"),
        gpt4o.WithRateLimit(500*time.Millisecond),
    )
    if err != nil {
        log.Fatal(err)
    }

    // 创建任务
    task := &models.Task{
        Name:   "测试任务",
        Prompt: "描述这张图片",
        Images: []string{"test.jpg"},
        Model:  "gpt-4o-image-vip",
    }

    // 同步处理
    ctx := context.Background()
    result, err := client.ProcessTask(ctx, task)
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("处理结果: %+v\n", result)
}
```

### 图片下载

```go
// 定义下载选项
downloadOptions := &models.DownloadOptions{
    OutputDir: "output",
    Timeout:   60, // 秒
}

// 处理任务
result, err := client.ProcessTask(ctx, task)
if err != nil {
    log.Fatal(err)
}

// 下载结果中的图片
if result.Success {
    err := client.DownloadImages(ctx, result, downloadOptions)
    if err != nil {
        log.Printf("下载图片失败: %v", err)
    } else {
        // 打印下载的图片路径
        for i, path := range result.ImagePaths {
            fmt.Printf("图片 %d 保存在: %s\n", i+1, path)
        }
    }
}
```

### 批量处理

```go
// 准备多个任务
tasks := []*models.Task{task1, task2, task3}

// 批量处理任务
results := client.ProcessTasks(ctx, tasks)

// 批量下载结果中的图片
for _, result := range results {
    if result.Success {
        client.DownloadImages(ctx, result, downloadOptions)
    }
}
```

### 异步处理

```go
// 异步处理单个任务
resultChan := client.ProcessTaskAsync(ctx, task)
for result := range resultChan {
    if result.Error != nil {
        log.Printf("错误: %v\n", result.Error)
        continue
    }
    fmt.Printf("处理结果: %+v\n", result)
    
    // 下载结果中的图片
    if result.Success {
        client.DownloadImages(ctx, result, downloadOptions)
    }
}

// 异步处理多个任务
tasks := []*models.Task{task1, task2, task3}
resultChan := client.ProcessTasksAsync(ctx, tasks)
for result := range resultChan {
    // 处理结果并下载图片
    if result.Success {
        client.DownloadImages(ctx, result, downloadOptions)
    }
}
```

## ⚙️ 配置选项

| 选项 | 说明 | 默认值 |
|------|------|--------|
| `WithAPIKey(apiKey string)` | 设置API密钥 | 必填 |
| `WithAPIURL(url string)` | 设置API URL | 默认API地址 |
| `WithRateLimit(rateLimit time.Duration)` | 设置请求速率限制 | 500ms |
| `WithMaxWorkers(maxWorkers int)` | 设置最大工作协程数 | 5 |

## 🛠️ 错误处理

SDK提供了完整的错误处理机制：

- API调用错误
- 图片处理错误
- 请求超时
- 上下文取消
- 图片下载错误

## ⚠️ 注意事项

1. 确保API密钥正确设置
2. 图片文件必须存在且可访问
3. 建议使用context控制超时
4. 注意请求速率限制
5. 并发处理时注意内存使用
6. 图片下载功能需要网络连接且目标URL必须可访问

## 📝 许可证

MIT 