package main

import (
	"backend/internal/notification"
	"backend/pkg/config"
	"backend/pkg/email"
	"backend/pkg/logger"
	"backend/pkg/mq"
	"backend/pkg/svc"
	"context"
	"fmt"
	"log"
)

func main() {
	ctx := context.Background()

	if err := run(ctx); err != nil {
		log.Fatalf("application error: %v", err)
	}

	log.Println("success shutdown")
}

func run(ctx context.Context) error {
	conf, err := config.LoadConfig("config.yaml")
	if err != nil {
		return fmt.Errorf("load config error: %w", err)
	}

	zapLog, err := logger.New(
		logger.WithFile(conf.Logger.File.Path),
		logger.WithLevel(conf.Logger.Level),
		logger.WithStdOut(conf.Logger.StdOut),
	)
	if err != nil {
		return fmt.Errorf("create logger error: %w", err)
	}

	rabbitMQ := mq.NewRabbitMQ(zapLog.Log, &conf.RabbitMQ)
	emailClient := email.NewClient(zapLog.Log, &conf.SMTP)

	inviteConsumer := mq.NewMQConsumer(
		zapLog.Log,
		rabbitMQ.Conn,
		notification.InviteCreatedHandler(emailClient, zapLog.Log, conf.Invite),
		mq.WithQueueName("notification.invite.created"),
		mq.WithConsumerExchange("hiring.events"),
		mq.WithConsumerRoutingKey("invite.created"),
		mq.WithQuorumQueue(),
		mq.WithPrefetchCount(10),
		mq.WithConcurrency(5),
	)

	if err := svc.Run(ctx, zapLog.Log, []svc.Service{
		zapLog,
		rabbitMQ,
		emailClient,
		inviteConsumer,
	}); err != nil {
		return fmt.Errorf("run service error: %w", err)
	}

	return nil
}
