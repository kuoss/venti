package reloader

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	v1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
)

var clientset1 *kubernetes.Clientset

func init() {
	var userAgent string
	ts := httptest.NewUnstartedServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userAgent = r.Header.Get("User-Agent")
		if userAgent != "test-agent" {
			panic("not test-agent")
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte("{}"))
	}))
	ts.Start()
	defer ts.Close()

	gv := v1.SchemeGroupVersion
	config := &rest.Config{
		Host: ts.URL,
	}
	config.GroupVersion = &gv
	config.NegotiatedSerializer = scheme.Codecs.WithoutConversion()
	config.UserAgent = userAgent
	config.ContentType = "application/json"

	var err error
	clientset1, err = kubernetes.NewForConfig(config)
	if err != nil {
		panic(err)
	}
}

func TestNew(t *testing.T) {
	reloader, err := New(clientset1)
	require.NoError(t, err)
	require.NotEmpty(t, reloader)
}

func TestStart(t *testing.T) {
	reloader, err := New(clientset1)
	require.NoError(t, err)
	require.NotEmpty(t, reloader)
	reloader.Start()
	time.Sleep(100 * time.Millisecond)
}
