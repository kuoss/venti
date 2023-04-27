package alerting

import (
	"fmt"
	"os"
	"testing"

	"github.com/kuoss/venti/pkg/model"
	"github.com/stretchr/testify/assert"
)

func init() {
	_ = os.Chdir("../../..")
}

func TestNew(t *testing.T) {
	testCases := []struct {
		file string
		want *AlertingStore
	}{
		{
			file: "",
			want: &AlertingStore{AlertingFile: &model.AlertingFile{Alertings: []model.Alerting{{Name: "alertmanager", Type: model.AlertingTypeAlertmanager, URL: "http://localhost:9093"}}}},
		},
		{
			file: "asdf",
			want: &AlertingStore{AlertingFile: &model.AlertingFile{Alertings: []model.Alerting{}}},
		},
		{
			file: "etc/alerting.yml",
			want: &AlertingStore{AlertingFile: &model.AlertingFile{Alertings: []model.Alerting{{Name: "alertmanager", Type: model.AlertingTypeAlertmanager, URL: "http://localhost:9093"}}}},
		},
		{
			file: "etc/alerting.yaml",
			want: &AlertingStore{AlertingFile: &model.AlertingFile{Alertings: []model.Alerting{}}},
		},
	}
	for i, tc := range testCases {
		t.Run(fmt.Sprintf("TESTCASE_#%d", i), func(t *testing.T) {
			got := New(tc.file)
			assert.Equal(t, tc.want, got)
		})
	}

}

func TestLoadAlertingFile(t *testing.T) {
	testCases := []struct {
		file      string
		want      *model.AlertingFile
		wantError string
	}{
		{
			"",
			&model.AlertingFile{Alertings: []model.Alerting{{Name: "alertmanager", Type: model.AlertingTypeAlertmanager, URL: "http://localhost:9093"}}},
			"",
		},
		{
			"asdfasdf",
			nil,
			"error on ReadFile: open asdfasdf: no such file or directory",
		},
		{
			"etc/alerting.yml",
			&model.AlertingFile{Alertings: []model.Alerting{{Name: "alertmanager", Type: model.AlertingTypeAlertmanager, URL: "http://localhost:9093"}}},
			"",
		},
		{
			"etc/alerting.yaml",
			nil,
			"error on ReadFile: open etc/alerting.yaml: no such file or directory",
		},
	}
	for i, tc := range testCases {
		t.Run(fmt.Sprintf("TESTCASE_#%d", i), func(t *testing.T) {
			got, err := loadAlertingFile(tc.file)
			if tc.wantError == "" {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, tc.wantError)
			}
			assert.Equal(t, tc.want, got)
		})
	}
}

func TestGetAlertmanagerURL(t *testing.T) {
	testCases := []struct {
		file string
		want string
	}{
		{
			"",
			"http://localhost:9093",
		},
		{
			"asdf",
			"",
		},
		{
			"etc/alerting.yml",
			"http://localhost:9093",
		},
		{
			"etc/alerting.yaml",
			"",
		},
	}
	for i, tc := range testCases {
		t.Run(fmt.Sprintf("TESTCASE_#%d", i), func(t *testing.T) {
			store := New(tc.file)
			assert.Equal(t, tc.want, store.GetAlertmanagerURL())
		})
	}
}
