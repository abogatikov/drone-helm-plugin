package pkg

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfig_setInstallCommand(t *testing.T) {
	testCases := []struct {
		cfg      Config
		expected []string
	}{
		{cfg: Config{HelmCommand: cInstall, Chart: "test",}, expected: []string{"install", "test"}},
		{cfg: Config{HelmCommand: cInstall, Chart: "test", Version: "0.0.1-test"}, expected: []string{"install", "test", "--version", "0.0.1-test"}},
		{cfg: Config{HelmCommand: cInstall, Chart: "test", Release: "test"}, expected: []string{"install", "test", "--name", "test"}},
		{cfg: Config{HelmCommand: cInstall, Chart: "test", Set: "test=123"}, expected: []string{"install", "test", "--set", "test=123"}},
		{cfg: Config{HelmCommand: cInstall, Chart: "test", SetString: "test=123"}, expected: []string{"install", "test", "--set-string", "test=123"}},
		{cfg: Config{HelmCommand: cInstall, Chart: "test", Values: "/tmp/1.yml,/tmp/2.yml"}, expected: []string{"install", "test", "--values", "/tmp/1.yml", "--values", "/tmp/2.yml"}},
		{cfg: Config{HelmCommand: cInstall, Chart: "test", Namespace: "test"}, expected: []string{"install", "test", "--namespace", "test"}},
		{cfg: Config{HelmCommand: cInstall, Chart: "test", TillerNamespace: "test"}, expected: []string{"install", "test", "--tiller-namespace", "test"}},
		{cfg: Config{HelmCommand: cInstall, Chart: "test", DryRun: true}, expected: []string{"install", "test", "--dry-run"}},
		{cfg: Config{HelmCommand: cInstall, Chart: "test", Debug: true}, expected: []string{"install", "test", "--debug"}},
		{cfg: Config{HelmCommand: cInstall, Chart: "test", Wait: true}, expected: []string{"install", "test", "--wait"}},
		{cfg: Config{HelmCommand: cInstall, Chart: "test", RecreatePods: true}, expected: []string{"install", "test", "--recreate-pods"}},
		{cfg: Config{HelmCommand: cInstall, Chart: "test", ReuseValues: true}, expected: []string{"install", "test", "--reuse-values"}},
		{cfg: Config{HelmCommand: cInstall, Chart: "test", Timeout: "10s"}, expected: []string{"install", "test", "--timeout", "10s"}},
		{cfg: Config{HelmCommand: cInstall, Chart: "test", Force: true}, expected: []string{"install", "test", "--force"}},
	}

	for _, c := range testCases {
		assert.Equal(t, c.expected, c.cfg.setInstallCommand())
	}
}

func TestConfig_setUpgradeCommand(t *testing.T) {
	testCases := []struct {
		cfg      Config
		expected []string
	}{
		{cfg: Config{HelmCommand: cUpgrade, Chart: "test", Release: "test_release"}, expected: []string{"upgrade", "--install", "test_release", "test"}},
		{cfg: Config{HelmCommand: cUpgrade, Chart: "test", Release: "test_release", Version: "0.0.1-test"}, expected: []string{"upgrade", "--install", "test_release", "test", "--version", "0.0.1-test"}},
		{cfg: Config{HelmCommand: cUpgrade, Chart: "test", Release: "test_release", Set: "test=123"}, expected: []string{"upgrade", "--install", "test_release", "test", "--set", "test=123"}},
		{cfg: Config{HelmCommand: cUpgrade, Chart: "test", Release: "test_release", SetString: "test=123"}, expected: []string{"upgrade", "--install", "test_release", "test", "--set-string", "test=123"}},
		{cfg: Config{HelmCommand: cUpgrade, Chart: "test", Release: "test_release", Values: "/tmp/1.yml,/tmp/2.yml"}, expected: []string{"upgrade", "--install", "test_release", "test", "--values", "/tmp/1.yml", "--values", "/tmp/2.yml"}},
		{cfg: Config{HelmCommand: cUpgrade, Chart: "test", Release: "test_release", Namespace: "test"}, expected: []string{"upgrade", "--install", "test_release", "test", "--namespace", "test"}},
		{cfg: Config{HelmCommand: cUpgrade, Chart: "test", Release: "test_release", TillerNamespace: "test"}, expected: []string{"upgrade", "--install", "test_release", "test", "--tiller-namespace", "test"}},
		{cfg: Config{HelmCommand: cUpgrade, Chart: "test", Release: "test_release", DryRun: true}, expected: []string{"upgrade", "--install", "test_release", "test", "--dry-run"}},
		{cfg: Config{HelmCommand: cUpgrade, Chart: "test", Release: "test_release", Debug: true}, expected: []string{"upgrade", "--install", "test_release", "test", "--debug"}},
		{cfg: Config{HelmCommand: cUpgrade, Chart: "test", Release: "test_release", Wait: true}, expected: []string{"upgrade", "--install", "test_release", "test", "--wait"}},
		{cfg: Config{HelmCommand: cUpgrade, Chart: "test", Release: "test_release", RecreatePods: true}, expected: []string{"upgrade", "--install", "test_release", "test", "--recreate-pods"}},
		{cfg: Config{HelmCommand: cUpgrade, Chart: "test", Release: "test_release", ReuseValues: true}, expected: []string{"upgrade", "--install", "test_release", "test", "--reuse-values"}},
		{cfg: Config{HelmCommand: cUpgrade, Chart: "test", Release: "test_release", Timeout: "10s"}, expected: []string{"upgrade", "--install", "test_release", "test", "--timeout", "10s"}},
		{cfg: Config{HelmCommand: cUpgrade, Chart: "test", Release: "test_release", Force: true}, expected: []string{"upgrade", "--install", "test_release", "test", "--force"}},
	}

	for _, c := range testCases {
		assert.Equal(t, c.expected, c.cfg.setUpgradeCommand())
	}
}

func TestConfig_setDeleteCommand(t *testing.T) {
	testCases := []struct {
		cfg      Config
		expected []string
	}{
		{cfg: Config{HelmCommand: cDelete, Release: "test", TillerNamespace: "test"}, expected: []string{"delete", "test", "--tiller-namespace", "test"}},
		{cfg: Config{HelmCommand: cDelete, Release: "test", DryRun: true}, expected: []string{"delete", "test", "--dry-run"}},
		{cfg: Config{HelmCommand: cDelete, Release: "test", Purge: true}, expected: []string{"delete", "test", "--purge"}},
	}

	for _, c := range testCases {
		assert.Equal(t, c.expected, c.cfg.setDeleteCommand())
	}
}

func TestConfig_setLintCommand(t *testing.T) {
	testCases := []struct {
		cfg      Config
		expected []string
	}{
		{cfg: Config{HelmCommand: cLint, Chart: "test", TillerNamespace: "test"}, expected: []string{"lint", "test", "--tiller-namespace", "test"}},
		{cfg: Config{HelmCommand: cLint, Chart: "test", Namespace: "test"}, expected: []string{"lint", "test", "--namespace", "test"}},
		{cfg: Config{HelmCommand: cLint, Chart: "test", Debug: true}, expected: []string{"lint", "test", "--debug"}},
		{cfg: Config{HelmCommand: cLint, Chart: "test", Set: "test=123"}, expected: []string{"lint", "test", "--set", "test=123"}},
		{cfg: Config{HelmCommand: cLint, Chart: "test", SetString: "test=123"}, expected: []string{"lint", "test", "--set-string", "test=123"}},
		{cfg: Config{HelmCommand: cLint, Chart: "test", Values: "/tmp/1.yml,/tmp/2.yml"}, expected: []string{"lint", "test", "--values", "/tmp/1.yml", "--values", "/tmp/2.yml"}},
	}

	for _, c := range testCases {
		assert.Equal(t, c.expected, c.cfg.setLintCommand())
	}
}
