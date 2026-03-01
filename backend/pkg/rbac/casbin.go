package rbac

import (
	"backend/pkg/svc"
	"bufio"
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/casbin/casbin/v2"
	"github.com/jackc/pgx/v5"
	pgxadapter "github.com/pckhoi/casbin-pgx-adapter/v3"
	"go.uber.org/zap"
)

var _ svc.Service = (*CasbinClient)(nil)

type CasbinClient struct {
	log       *zap.Logger
	connConf  *pgx.ConnConfig
	modelPath string

	adapter  *pgxadapter.Adapter
	enforcer *casbin.Enforcer
}

func NewCasbinClient(log *zap.Logger, connConf *pgx.ConnConfig, modelPath string) *CasbinClient {
	return &CasbinClient{
		log:       log,
		connConf:  connConf,
		modelPath: modelPath,
	}
}

func (c *CasbinClient) Name() string {
	return "casbin"
}

func (c *CasbinClient) DependsOn() []string {
	return []string{"db"}
}

func (c *CasbinClient) Init(_ context.Context) error {
	adapter, err := pgxadapter.NewAdapter(c.connConf)
	if err != nil {
		return fmt.Errorf("failed to create casbin pgx adapter: %w", err)
	}

	c.adapter = adapter

	enforcer, err := casbin.NewEnforcer(c.modelPath, adapter)
	if err != nil {
		return fmt.Errorf("failed to create casbin enforcer: %w", err)
	}

	if err := enforcer.LoadPolicy(); err != nil {
		return fmt.Errorf("failed to load casbin policy: %w", err)
	}

	c.enforcer = enforcer

	// Синхронизируем политики из policy.csv при каждом запуске.
	// Метод AddPolicy идемпотентен: он добавит в БД только новые строки из CSV,
	// которых там ещё нет.
	if err := c.seedDefaultPolicies(); err != nil {
		return fmt.Errorf("failed to sync default policies: %w", err)
	}

	c.log.Debug("casbin initialized")

	return nil
}

func (c *CasbinClient) HealthCheck(_ context.Context) error {
	if c.enforcer == nil {
		return errors.New("casbin enforcer is not initialized")
	}

	c.log.Debug("casbin health check passed")
	return nil
}

func (c *CasbinClient) Run(_ context.Context) error {
	return nil
}

func (c *CasbinClient) Stop(_ context.Context) error {
	if c.adapter != nil {
		c.adapter.Close()
	}

	c.log.Debug("casbin adapter closed")
	return nil
}

// Enforce проверяет, имеет ли sub (userID) в домене dom (orgID)
// доступ к ресурсу obj с действием act.
func (c *CasbinClient) Enforce(sub, dom, obj, act string) (bool, error) {
	return c.enforcer.Enforce(sub, dom, obj, act)
}

// AddRoleForUserInDomain назначает роль пользователю в организации.
// g = user, role, domain
func (c *CasbinClient) AddRoleForUserInDomain(user, role, domain string) (bool, error) {
	ok, err := c.enforcer.AddGroupingPolicy(user, role, domain)
	if err != nil {
		return false, fmt.Errorf("add role %q for user %q in domain %q: %w", role, user, domain, err)
	}

	return ok, nil
}

// DeleteRoleForUserInDomain удаляет роль у пользователя в организации.
func (c *CasbinClient) DeleteRoleForUserInDomain(user, role, domain string) (bool, error) {
	ok, err := c.enforcer.RemoveGroupingPolicy(user, role, domain)
	if err != nil {
		return false, fmt.Errorf("delete role %q for user %q in domain %q: %w", role, user, domain, err)
	}

	return ok, nil
}

// GetRolesForUserInDomain возвращает список ролей пользователя в организации.
func (c *CasbinClient) GetRolesForUserInDomain(user, domain string) []string {
	return c.enforcer.GetRolesForUserInDomain(user, domain)
}

// GetUsersForRoleInDomain возвращает список пользователей с данной ролью в организации.
func (c *CasbinClient) GetUsersForRoleInDomain(role, domain string) []string {
	return c.enforcer.GetUsersForRoleInDomain(role, domain)
}

// DeleteAllRolesInDomain удаляет все role-assignments в указанном домене.
// Используется при удалении организации.
func (c *CasbinClient) DeleteAllRolesInDomain(domain string) (bool, error) {
	ok, err := c.enforcer.DeleteDomains(domain)
	if err != nil {
		return false, fmt.Errorf("delete all roles in domain %q: %w", domain, err)
	}

	return ok, nil
}

// seedDefaultPolicies загружает стартовые политики из policy.csv.
// Вызывается один раз при первом запуске, когда в БД нет ни одной политики.
func (c *CasbinClient) seedDefaultPolicies() error {
	policyPath := filepath.Join(filepath.Dir(c.modelPath), "policy.csv")

	f, err := os.Open(policyPath)
	if err != nil {
		return fmt.Errorf("open policy file %q: %w", policyPath, err)
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		parts := strings.Split(line, ",")
		for i := range parts {
			parts[i] = strings.TrimSpace(parts[i])
		}

		if len(parts) < 2 {
			continue
		}

		switch parts[0] {
		case "p":
			args := make([]interface{}, len(parts)-1)
			for i, v := range parts[1:] {
				args[i] = v
			}

			if _, err := c.enforcer.AddPolicy(args...); err != nil {
				return fmt.Errorf("add policy %v: %w", parts, err)
			}
		case "g":
			args := make([]interface{}, len(parts)-1)
			for i, v := range parts[1:] {
				args[i] = v
			}

			if _, err := c.enforcer.AddGroupingPolicy(args...); err != nil {
				return fmt.Errorf("add grouping policy %v: %w", parts, err)
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("scan policy file: %w", err)
	}

	c.log.Info("seeded default casbin policies", zap.String("file", policyPath))

	return nil
}
