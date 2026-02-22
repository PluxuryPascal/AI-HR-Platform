package svc

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

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

	sorted, err := sortTopologically(services)
	if err != nil {
		return fmt.Errorf("topological sort failed: %w", err)
	}
	services = sorted

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

	// Ожидаем сигнала об остановке сервера
	<-ctx.Done()

	// Получаем новый контекст с таймаутом, чтобы не зависнуть при остановке
	stopCtx, stopCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer stopCancel()

	for i := len(services) - 1; i >= 0; i-- {
		s := services[i]

		if err := s.Stop(stopCtx); err != nil {
			return fmt.Errorf("[STOP] %s: %w", s.Name(), err)
		}

		fmt.Printf("[STOP] %s\n", s.Name())
	}

	if err := g.Wait(); err != nil {
		return fmt.Errorf("exit reason from service: %w", err)
	}

	return nil
}

// sortTopologically выполняет топологическую сортировку списка сервисов
// на основе графа зависимостей с помощью алгоритма Кана.
func sortTopologically(services []Service) ([]Service, error) {
	// serviceMap хранит сервисы по их именам для быстрого доступа.
	serviceMap := make(map[string]Service)
	// inDegree отслеживает количество входящих зависимостей для каждого сервиса
	// (сколько других сервисов должно запуститься ДО него).
	inDegree := make(map[string]int)
	// adjList (список смежности) хранит граф: кто зависит от данного сервиса.
	adjList := make(map[string][]string)

	// Шаг 1: Инициализируем мапы и проверяем на дубликаты имен.
	for _, s := range services {
		name := s.Name()
		if _, exists := serviceMap[name]; exists {
			return nil, fmt.Errorf("duplicate service name: %s", name)
		}
		serviceMap[name] = s
		inDegree[name] = 0
	}

	// Шаг 2: Строим граф зависимостей.
	for _, s := range services {
		name := s.Name()
		for _, dep := range s.DependsOn() {
			// Проверяем, существует ли сервис, от которого мы зависим.
			if _, exists := serviceMap[dep]; !exists {
				return nil, fmt.Errorf("service '%s' depends on unknown service '%s'", name, dep)
			}
			// Добавляем текущий сервис в список смежности его зависимости
			// (когда запустится dep, он освободит name).
			adjList[dep] = append(adjList[dep], name)
			// Увеличиваем счетчик заблокированности для текущего сервиса.
			inDegree[name]++
		}
	}

	// Шаг 3: Находим все сервисы, которые ни от кого не зависят (inDegree == 0).
	// Они первыми отправятся в очередь на запуск.
	var queue []string
	for name, degree := range inDegree {
		if degree == 0 {
			queue = append(queue, name)
		}
	}

	var sorted []Service

	// Шаг 4: Обходим граф в ширину (BFS), освобождая зависимости.
	for len(queue) > 0 {
		// Достаем сервис из очереди
		currName := queue[0]
		queue = queue[1:]

		// Раз он в очереди, значит его зависимости удовлетворены — добавляем в итог.
		sorted = append(sorted, serviceMap[currName])

		// Проходимся по всем сервисам, которые ждали текущий
		for _, dependentName := range adjList[currName] {
			// Уменьшаем их счетчик ожиданий на 1
			inDegree[dependentName]--
			// Если зависимостей больше не осталось, ставим в очередь
			if inDegree[dependentName] == 0 {
				queue = append(queue, dependentName)
			}
		}
	}

	// Шаг 5: Проверяем на наличие циклических зависимостей.
	// Если в отсортированном списке меньше элементов, чем в оригинальном,
	// значит какие-то сервисы застряли в бесконечном цикле (deadlock).
	if len(sorted) != len(services) {
		return nil, fmt.Errorf("cyclic dependency detected among services")
	}

	return sorted, nil
}
