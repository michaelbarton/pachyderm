package server

import (
	"bytes"
	"context"
	"crypto/tls"
	"fmt"
	"io"
	"net/http"
	"testing"
	"time"

	"github.com/minio/minio-go/v6"
	"github.com/pachyderm/pachyderm/v2/src/client"
	"github.com/pachyderm/pachyderm/v2/src/internal/minikubetestenv"
	"github.com/pachyderm/pachyderm/v2/src/internal/require"
	"github.com/pachyderm/pachyderm/v2/src/internal/testutil"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	applyappsv1 "k8s.io/client-go/applyconfigurations/apps/v1"
	applyv1 "k8s.io/client-go/applyconfigurations/core/v1"
	applymetav1 "k8s.io/client-go/applyconfigurations/meta/v1"
)

func mustGetOK(t *testing.T, url, description string) {
	t.Helper()
	hc := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}
	ctx, c := context.WithTimeout(context.Background(), 20*time.Second)
	defer c()

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	require.NoError(t, err, "should create an http request for %v ok (%v)", url, description)

	res, err := hc.Do(req)
	require.NoError(t, err, "should make http request for %v ok (%v)", url, description)

	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		t.Logf("warning: get %v (%v): problem reading body: %v", url, description, err)
	}
	t.Logf("get %v (%v): body:\n%s", url, description, body)

	require.Equal(t, http.StatusOK, res.StatusCode, "http status should be OK (200) for %v (%v)", url, description)
}

func proxyTest(t *testing.T, c *client.APIClient, secure bool) {
	t.Helper()
	httpPrefix := "http://"
	if secure {
		httpPrefix = "https://"
	}
	addr := fmt.Sprintf("%v:%d", c.GetAddress().Host, c.GetAddress().Port)

	// Test console.
	t.Run("TestConsole", func(t *testing.T) {
		mustGetOK(t, httpPrefix+addr+"/", "console")
	})

	// Test OIDC.
	t.Run("TestOIDC", func(t *testing.T) {
		mustGetOK(t, httpPrefix+addr+"/dex/.well-known/openid-configuration", "openid config")
	})

	testText := []byte("this is a test\n")
	// Test GRPC API.
	t.Run("TestGRPC", func(t *testing.T) {
		err := c.CreateRepo("test")
		require.NoError(t, err, "should be able to create repo")
		err = c.PutFile(client.NewRepo("test").NewCommit("master", ""), "test.txt", bytes.NewReader(testText))
		require.NoError(t, err, "should be able to put a file")
	})

	// Test S3 API.
	t.Run("TestS3", func(t *testing.T) {
		s3v4, err := minio.NewV4(addr, c.AuthToken(), c.AuthToken(), secure)
		require.NoError(t, err, "should create s3v4 client")
		s3v2, err := minio.NewV2(addr, c.AuthToken(), c.AuthToken(), secure)
		require.NoError(t, err, "should create s3v2 client")

		for name, client := range map[string]interface {
			GetObject(string, string, minio.GetObjectOptions) (*minio.Object, error)
		}{"v4": s3v4, "v2": s3v2} {
			obj, err := client.GetObject("master.test", "test.txt", minio.GetObjectOptions{})
			require.NoError(t, err, "should be able to get object using s3 protocol %v", name)
			content, err := io.ReadAll(obj)
			require.NoError(t, err, "should be able to read object using s3 protocol %v", name)
			require.Equal(t, testText, content, "should have read correct content using s3 protocol %v", name)
		}
	})
}

func deployFakeConsole(t *testing.T, ns string) {
	t.Helper()
	ctx := context.Background()
	c := testutil.GetKubeClient(t)
	lbl := map[string]string{
		"app":   "console",
		"suite": "pachyderm",
	}
	aopts := metav1.ApplyOptions{
		FieldManager: "proxy_test",
	}

	cm := applyv1.ConfigMap("fake-console-html", ns)
	cm.Data = map[string]string{
		"index.html": "<html><head><title>Hi</title></head><body><p>Hello from console!</p></body></html>",
	}
	_, err := c.CoreV1().ConfigMaps(ns).Apply(ctx, cm, aopts)
	require.NoError(t, err, "should create a configmap for the fake console HTML files")

	cm = applyv1.ConfigMap("fake-console-config", ns)
	cm.Data = map[string]string{
		"default.conf": `server {
			listen       4000;
			server_name  localhost;

			access_log  /dev/stdout  main;

			location / {
			    root   /usr/share/nginx/html;
			    index  index.html;
			}
		    }`,
	}
	_, err = c.CoreV1().ConfigMaps(ns).Apply(ctx, cm, aopts)
	require.NoError(t, err, "should create a configmap for the fake console config files")

	dep := applyappsv1.Deployment("console", ns)
	dep.ObjectMetaApplyConfiguration.Labels = lbl
	dep.Spec = applyappsv1.DeploymentSpec().WithReplicas(1).WithSelector(applymetav1.LabelSelector().WithMatchLabels(lbl))
	dep.Spec.Template = applyv1.PodTemplateSpec()
	dep.Spec.Template.ObjectMetaApplyConfiguration = applymetav1.ObjectMeta().WithName("console").WithNamespace(ns).WithLabels(lbl)
	dep.Spec.Template.Spec = applyv1.PodSpec().
		WithContainers(applyv1.Container().
			WithName("console").
			WithImage("nginx:latest").
			WithPorts(applyv1.ContainerPort().WithContainerPort(4000).WithName("console-http").WithProtocol(v1.ProtocolTCP)).
			WithVolumeMounts(
				applyv1.VolumeMount().WithName("html").WithMountPath("/usr/share/nginx/html"),
				applyv1.VolumeMount().WithName("config").WithMountPath("/etc/nginx/conf.d"),
			).
			WithReadinessProbe(applyv1.Probe().WithHTTPGet(
				applyv1.HTTPGetAction().WithPort(intstr.FromInt(4000)).WithPath("/")))).
		WithVolumes(
			applyv1.Volume().WithName("html").WithConfigMap(applyv1.ConfigMapVolumeSource().WithName("fake-console-html")),
			applyv1.Volume().WithName("config").WithConfigMap(applyv1.ConfigMapVolumeSource().WithName("fake-console-config")),
		)
	_, err = c.AppsV1().Deployments(ns).Apply(ctx, dep, aopts)
	require.NoError(t, err, "should create a 'console' deployment")

	t.Log("waiting for the console backend service to have endpoints")
	for i := 0; i < 20; i++ {
		ep, err := c.CoreV1().Endpoints(ns).Get(ctx, "console-proxy-backend", metav1.GetOptions{})
		if err != nil {
			t.Logf("get console service endpoints: %v", err)
			time.Sleep(time.Second)
			continue
		}
		if ss := ep.Subsets; len(ss) > 0 {
			t.Logf("console endpoints: %v", ss)
			break
		}
		t.Log("no endpoints yet")
		time.Sleep(time.Second)
	}
}

func TestTrafficThroughProxy(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}
	t.Parallel()
	c, ns := minikubetestenv.AcquireCluster(t)
	deployFakeConsole(t, ns)
	testutil.ActivateAuthClient(t, c)
	testutil.ConfigureOIDCProvider(t, c)
	proxyTest(t, c, false)
}

// func TestTrafficThroughTLSProxy(t *testing.T) {
// 	if testing.Short() {
// 		t.Skip("Skipping integration test in short mode")
// 	}
// 	t.Parallel()
// 	c, ns := minikubetestenv.AcquireCluster(t)
// 	deployFakeConsole(t, ns)
// 	c.GetAddress().Secured = true
// 	testutil.ActivateAuthClient(t, c)
// 	proxyTest(t, c, true)
// }
