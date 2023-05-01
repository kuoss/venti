package store

import (
	"testing"

	"github.com/kuoss/venti/pkg/model"
	"github.com/stretchr/testify/assert"
)

func TestNewStores(t *testing.T) {
	stores, err := NewStores(&model.Config{})
	assert.NoError(t, err)
	assert.NotEmpty(t, stores)
	assert.NotEmpty(t, stores.AlertRuleStore)
	assert.NotEmpty(t, stores.AlertingStore)
	assert.NotEmpty(t, stores.DashboardStore)
	assert.NotEmpty(t, stores.DatasourceStore)
	assert.NotEmpty(t, stores.RemoteStore)
	assert.NotEmpty(t, stores.StatusStore)
	assert.NotEmpty(t, stores.UserStore)
}
