package jobs

import (
	"context"
	"log"

	"github.com/hibiken/asynq"
	"github.com/mnizarzr/dot-test/config"
)

type JobManager struct {
	client *asynq.Client
	server *asynq.Server
	mux    *asynq.ServeMux
	config *config.Config
}

// NewJobManager creates a new job manager instance
func NewJobManager(cfg *config.Config) *JobManager {
	redisOpt := asynq.RedisClientOpt{
		Addr:     cfg.RedisAddress,
		Password: cfg.RedisPassword,
	}

	client := asynq.NewClient(redisOpt)

	server := asynq.NewServer(redisOpt, asynq.Config{
		Concurrency: 10,
		Queues: map[string]int{
			"critical": 6,
			"default":  3,
			"low":      1,
		},
	})

	mux := asynq.NewServeMux()

	return &JobManager{
		client: client,
		server: server,
		mux:    mux,
		config: cfg,
	}
}

// RegisterHandlers registers all job handlers
func (jm *JobManager) RegisterHandlers() {
	emailJobHandler := NewEmailJobHandler(jm.config)
	jm.mux.HandleFunc(TypeEmailWelcome, emailJobHandler.HandleWelcomeEmail)
}

// EnqueueJob enqueues a job for processing
func (jm *JobManager) EnqueueJob(ctx context.Context, task *asynq.Task, opts ...asynq.Option) (*asynq.TaskInfo, error) {
	return jm.client.Enqueue(task, opts...)
}

// Start starts the job processing server
func (jm *JobManager) Start() error {
	jm.RegisterHandlers()
	log.Println("Starting job processing server...")
	return jm.server.Run(jm.mux)
}

// GetClient returns the asynq client for enqueueing jobs
func (jm *JobManager) GetClient() *asynq.Client {
	return jm.client
}
