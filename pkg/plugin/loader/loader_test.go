package loader

import (
	"path"
	"reflect"
	"sort"
	"testing"

	"github.com/heptio/sonobuoy/pkg/plugin"
	"github.com/heptio/sonobuoy/pkg/plugin/driver/daemonset"
	"github.com/heptio/sonobuoy/pkg/plugin/driver/job"
	"github.com/pkg/errors"
	corev1 "k8s.io/api/core/v1"
)

func TestFindPlugins(t *testing.T) {
	testdir := path.Join("testdata", "plugin.d")
	plugins, err := findPlugins(testdir)
	if err != nil {
		t.Fatalf("unexpected err %v", err)
	}

	expected := []string{
		"testdata/plugin.d/daemonset.yaml",
		"testdata/plugin.d/invalid.yml",
		"testdata/plugin.d/job.yml",
	}
	sort.Strings(plugins)

	if !reflect.DeepEqual(expected, plugins) {
		t.Errorf("expected %v, got %v", expected, plugins)
	}
}

func TestLoadNonexistentPlugin(t *testing.T) {
	_, err := loadDefinition("non/existent/path")
	if errors.Cause(err).Error() != "open non/existent/path: no such file or directory" {
		t.Errorf("Expected ErrNotExist, got %v", errors.Cause(err))
	}
}

func TestLoadValidPlugin(t *testing.T) {
	jobDefFile := "testdata/plugin.d/job.yml"
	jobDef, err := loadDefinition(jobDefFile)
	if err != nil {
		t.Fatalf("Unexpected error creating job plugin: %v", err)
	}

	if jobDef.SonobuoyConfig.Driver != "Job" {
		t.Errorf("expected driver Job, got %q", jobDef.SonobuoyConfig.Driver)
	}
	if jobDef.SonobuoyConfig.PluginName != "test-job-plugin" {
		t.Errorf("expected name test-job-plugin, got %q", jobDef.SonobuoyConfig.PluginName)
	}

	if jobDef.Spec.Image != "gcr.io/heptio-images/heptio-e2e:master" {
		t.Errorf("expected name gcr.io/heptio-images/heptio-e2e:master, got %q", jobDef.Spec.Image)
	}

	daemonDefFile := "testdata/plugin.d/daemonset.yaml"
	daemonDef, err := loadDefinition(daemonDefFile)
	if err != nil {
		t.Fatalf("Unexpected error creating job plugin: %v", err)
	}

	if daemonDef.SonobuoyConfig.Driver != "DaemonSet" {
		t.Errorf("expected driver DaemonSet, got %q", daemonDef.SonobuoyConfig.Driver)
	}
	if daemonDef.SonobuoyConfig.PluginName != "test-daemon-set-plugin" {
		t.Errorf("expected name test-daemon-set-plugin, got %q", daemonDef.SonobuoyConfig.PluginName)
	}
	if daemonDef.Spec.Image != "gcr.io/heptio-images/heptio-e2e:master" {
		t.Errorf("expected name gcr.io/heptio-images/heptio-e2e:master, got %q", jobDef.Spec.Image)
	}
}

func TestLoadJobPlugin(t *testing.T) {
	namespace := "loader_test"
	jobDef := &pluginDefinition{
		SonobuoyConfig: sonobuoyConfig{
			Driver:     "Job",
			PluginName: "test-job-plugin",
		},
		Spec: corev1.Container{
			Image: "gcr.io/heptio-images/heptio-e2e:master",
		},
	}

	pluginIface, err := loadPlugin(jobDef, namespace)
	if err != nil {
		t.Fatalf("unexpected error loading plugin: %v", err)
	}

	jobPlugin, ok := pluginIface.(*job.Plugin)

	if !ok {
		t.Fatalf("loaded plugin not a job.Plugin")
	}

	if jobPlugin.Definition.Name != "test-job-plugin" {
		t.Errorf("expected plugin name 'test-job-plugin', got '%v'", jobPlugin.Definition.Name)
	}
	if jobPlugin.Definition.Spec.Image != "gcr.io/heptio-images/heptio-e2e:master" {
		t.Errorf("expected plugin name 'gcr.io/heptio-images/heptio-e2e:master', got '%v'", jobPlugin.Definition.Spec.Image)
	}
	if jobPlugin.Namespace != namespace {
		t.Errorf("expected plugin name '%q', got '%v'", namespace, jobPlugin.Namespace)
	}

}

func TestLoadDaemonSet(t *testing.T) {
	namespace := "loader_test"
	daemonDef := &pluginDefinition{
		SonobuoyConfig: sonobuoyConfig{
			Driver:     "DaemonSet",
			PluginName: "test-daemon-set-plugin",
		},
		Spec: corev1.Container{
			Image: "gcr.io/heptio-images/heptio-e2e:master",
		},
	}

	pluginIface, err := loadPlugin(daemonDef, namespace)
	if err != nil {
		t.Fatalf("unexpected error loading plugin: %v", err)
	}

	daemonPlugin, ok := pluginIface.(*daemonset.Plugin)

	if !ok {
		t.Fatalf("loaded plugin not a daemon.Plugin")
	}

	if daemonPlugin.Definition.Name != "test-daemon-set-plugin" {
		t.Errorf("expected plugin name 'test-daemon-set-plugin', got '%v'", daemonPlugin.Definition.Name)
	}
	if daemonPlugin.Definition.Spec.Image != "gcr.io/heptio-images/heptio-e2e:master" {
		t.Errorf("expected plugin name 'gcr.io/heptio-images/heptio-e2e:master', got '%v'", daemonPlugin.Definition.Spec.Image)
	}
	if daemonPlugin.Namespace != namespace {
		t.Errorf("expected plugin name '%q', got '%v'", namespace, daemonPlugin.Namespace)
	}
}

func TestFilterList(t *testing.T) {
	definitions := []*pluginDefinition{
		{SonobuoyConfig: sonobuoyConfig{PluginName: "test1"}},
		{SonobuoyConfig: sonobuoyConfig{PluginName: "test2"}},
	}

	selections := []plugin.Selection{
		{Name: "test1"},
		{Name: "test3"},
	}

	expected := []*pluginDefinition{definitions[0]}
	filtered := filterPluginDef(definitions, selections)
	if !reflect.DeepEqual(filtered, expected) {
		t.Errorf("expected %+#v, got %+#v", expected, filtered)
	}
}
