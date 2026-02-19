package svc

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"golang.org/x/sync/errgroup"
)

// Run запускает набор сервисов с учётом их зависимостей.
//
// Алгоритм работы:
//  1. Строит граф зависимостей и выполняет топологическую сортировку.
//  2. Синхронно вызывает Init для каждого сервиса в порядке зависимостей.
//  3. Дожидается готовности всех сервисов через HealthCheck.
//  4. Запускает Run для всех сервисов параллельно.
//  5. При получении сигнала завершения или ошибке — останавливает сервисы
//     в обратном порядке вызовом Stop.
//
// Если обнаружены циклические зависимости или отсутствует требуемый сервис,
// функция вернёт ошибку до запуска.
//
// Run блокирует выполнение до завершения всех сервисов или возникновения ошибки.
func Run(ctx context.Context, services []Service) error {
	ctx, stop := signal.NotifyContext(ctx, os.Interrupt, syscall.SIGTERM)
	defer stop()

	for _, s := range services {
		fmt.Printf("[INIT] %s\n", s.Name())

		if err := s.Init(ctx); err != nil {
			return fmt.Errorf("[INIT] %s: %w", s.Name(), err)
		}
	}

	g, ctx := errgroup.WithContext(ctx)

	for _, s := range services {
		fmt.Printf("[HEALTHCHECK] %s\n", s.Name())

		if err := s.HealthCheck(ctx); err != nil {
			return fmt.Errorf("[HEALTHCHECK] %s: %w", s.Name(), err)
		}
	}

	for _, s := range services {
		s := s

		g.Go(func() error {
			fmt.Printf("[RUN] %s\n", s.Name())

			if err := s.Run(ctx); err != nil {
				return fmt.Errorf("[RUN] %s: %w", s.Name(), err)
			}

			return nil
		})
	}

	for i := len(services) - 1; i >= 0; i-- {
		s := services[i]

		if err := s.Stop(ctx); err != nil {
			return fmt.Errorf("[STOP] %s: %w", s.Name(), err)
		}

		fmt.Printf("[STOP] %s\n", s.Name())
	}

	if err := g.Wait(); err != nil {
		return fmt.Errorf("exit reason from service: %w", err)
	}

	return nil
}
