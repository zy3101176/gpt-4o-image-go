package models

import (
	"context"
)

// IClient 定义GPT-4O客户端接口
type IClient interface {
	// ProcessTask 同步处理单个任务
	ProcessTask(ctx context.Context, task *Task) (*TaskResult, error)

	// ProcessTasks 同步处理多个任务
	ProcessTasks(ctx context.Context, tasks []*Task) []*TaskResult

	// ProcessTaskAsync 异步处理单个任务
	ProcessTaskAsync(ctx context.Context, task *Task) <-chan *TaskResult

	// ProcessTasksAsync 异步处理多个任务
	ProcessTasksAsync(ctx context.Context, tasks []*Task) <-chan *TaskResult

	// DownloadImages 下载任务结果中的图片
	DownloadImages(ctx context.Context, result *TaskResult, options *DownloadOptions) error
}
