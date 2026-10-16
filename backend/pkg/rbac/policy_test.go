package rbac

import (
	"testing"

	"github.com/casbin/casbin/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Проверяем боевые model.conf и policy.csv без базы: роли выдаются в рамках команды (домена).
func newTestEnforcer(t *testing.T) *casbin.Enforcer {
	t.Helper()

	e, err := casbin.NewEnforcer("../../casbin/model.conf", "../../casbin/policy.csv")
	require.NoError(t, err)

	_, err = e.AddGroupingPolicies([][]string{
		{"alice", "owner", "team-a"},
		{"rita", "recruiter", "team-a"},
		{"hank", "hiring_manager", "team-a"},
		{"bob", "recruiter", "team-b"},
	})
	require.NoError(t, err)

	return e
}

func TestPolicy(t *testing.T) {
	e := newTestEnforcer(t)

	tests := []struct {
		name               string
		sub, dom, obj, act string
		want               bool
	}{
		{"владелец может удалить вакансию", "alice", "team-a", "/api/v1/jobs/42", "DELETE", true},
		{"рекрутер читает что угодно", "rita", "team-a", "/api/v1/dashboard/stats", "GET", true},
		{"рекрутер управляет кандидатами", "rita", "team-a", "/api/v1/candidates/7", "PATCH", true},
		{"рекрутер получает список кандидатов вакансии", "rita", "team-a", "/api/v1/jobs/42/candidates/list", "POST", true},
		{"нанимающий менеджер читает", "hank", "team-a", "/api/v1/jobs/42", "GET", true},
		{"нанимающий менеджер получает списки", "hank", "team-a", "/api/v1/jobs/list", "POST", true},
		{"нанимающий менеджер не удаляет кандидатов", "hank", "team-a", "/api/v1/candidates/7", "DELETE", false},
		{"нанимающий менеджер не создаёт вакансии", "hank", "team-a", "/api/v1/jobs/42", "PUT", false},
		{"роль не действует в чужой команде", "bob", "team-a", "/api/v1/candidates/7", "GET", false},
		{"пользователь без роли не имеет доступа", "mallory", "team-a", "/api/v1/jobs/42", "GET", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ok, err := e.Enforce(tt.sub, tt.dom, tt.obj, tt.act)
			require.NoError(t, err)
			assert.Equal(t, tt.want, ok)
		})
	}
}
