package user

import (
	"os"
	"testing"

	"github.com/kuoss/venti/pkg/model"
	"github.com/stretchr/testify/require"
)

func init() {
	err := os.Chdir("../../..")
	if err != nil {
		panic(err)
	}
}

func TestNew(t *testing.T) {
	userService, err := New("./data/venti.sqlite3", model.UserConfig{})
	require.NoError(t, err)
	require.NotEmpty(t, userService)
}

func TestSetEtcUsers(t *testing.T) {

}

func TestFindByUsername(t *testing.T) {

}

func TestFindByUserIdAndToken(t *testing.T) {

}

func TestSave(t *testing.T) {

}
