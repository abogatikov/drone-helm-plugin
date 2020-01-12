package pkg

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"text/template"
)

const (
	cInstall = "install"
	cUpgrade = "upgrade"
	cDelete  = "delete"
	cLint    = "lint"

	cValuesFile = "/tmp/values.yaml"
)

//nolint:lll,maligned,staticcheck //reason:part of command line parser
type Config struct {
	HelmBinary         string   `long:"helm-binary" description:"Absolute path of the helm binary" required:"true" default:"/usr/bin/helm"`
	KubeConfigTemplate string   `long:"kube-config-template" description:"Absolute path of the kube config template" required:"true" default:"/opt/config"`
	APIServer          string   `long:"api-server" description:"Set api server address" env:"PLUGIN_API_SERVER" required:"true" default:"https://kubernetes.default"`
	Token              string   `long:"token" description:"Set client token" env:"PLUGIN_TOKEN" required:"true"`
	Certificate        string   `long:"certificate" description:"Set server certificate" env:"PLUGIN_CERTIFICATE" required:"true"`
	ServiceAccount     string   `long:"service-account" description:"Set service account name" env:"PLUGIN_SERVICE_ACCOUNT" required:"true"`
	KubeConfig         string   `long:"kube-config" description:"Absolute path of the kubeconfig file to be used" env:"PLUGIN_KUBE_CONFIG"`
	HelmCommand        string   `long:"helm-command" description:"The command Helm has to execute" env:"PLUGIN_HELM_COMMAND" required:"true" choice:"install" choice:"upgrade" choice:"delete" choice:"lint"`
	Namespace          string   `long:"namespace" description:"Namespace to install the release into." env:"PLUGIN_NAMESPACE"`
	TLSVerify          bool     `long:"tls-verify" description:"Enable TLS for request and verify remote" env:"PLUGIN_TLS_VERIFY"`
	Set                string   `long:"set" description:"Set values on the command line" env:"PLUGIN_SET"`
	SetString          string   `long:"set-string" description:"Set STRING values on the command line" env:"PLUGIN_SET_STRING"`
	Values             string   `long:"values" description:"Specify values in a YAML file or a URL" env:"PLUGIN_VALUES"`
	GetValues          bool     `long:"get-values" description:"Download the values file for a named release" env:"PLUGIN_GET_VALUES"`
	Release            string   `long:"release" description:"The release name. If unspecified, it will autogenerate one for you" env:"PLUGIN_RELEASE"`
	HelmRepos          []string `long:"helm-repos" description:"Add a chart repository" env:"PLUGIN_HELM_REPOS"`
	Chart              string   `long:"chart" description:"Chart reference" env:"PLUGIN_CHART" required:"true"`
	Version            string   `long:"chart-version" description:"Specify the exact chart version to install." env:"PLUGIN_CHART_VERSION"`
	Debug              bool     `long:"debug" description:"Enable verbose output" env:"PLUGIN_DEBUG"`
	DryRun             bool     `long:"dry-run" description:"Simulate an install" env:"PLUGIN_DRY-RUN"`
	TillerNamespace    string   `long:"tiller-namespace" description:"Namespace of Tiller" env:"PLUGIN_TILLER_NAMESPACE"`
	Wait               bool     `long:"wait" description:"If set, will wait until all parts a Deployment are in a ready state" env:"PLUGIN_WAIT"`
	RecreatePods       bool     `long:"recreate-pods" description:"Performs pods restart for the resource if applicable" env:"PLUGIN_RECREATE_PODS"`
	ReuseValues        bool     `long:"reuse-values" description:"When upgrading, reuse the last release's values and merge in any overrides from other commands" env:"PLUGIN_REUSE_VALUES"`
	Timeout            string   `long:"timeout" description:"Time in seconds to wait for any individual Kubernetes operation." env:"PLUGIN_TIMEOUT"`
	Force              bool     `long:"force" description:"Force resource update through delete/recreate if needed" env:"PLUGIN_FORCE"`
	UpdateDependencies bool     `long:"update-dependencies" description:"Update information of available charts locally from chart repositories" env:"PLUGIN_UPDATE_DEPENDENCIES"`
	Purge              bool     `long:"purge" description:"Remove the release from the store and make its name free for later use" env:"PLUGIN_PURGE"`
}

func (cfg *Config) Exec() error {
	if cfg.KubeConfig == "" {
		err := cfg.initKubeconfig()
		if err != nil {
			return err
		}
	}

	err := cfg.runCommand([]string{"init", "--client-only"})
	if err != nil {
		return err
	}

	err = cfg.addHelmRepo()
	if err != nil {
		return err
	}

	if cfg.UpdateDependencies {
		if err := cfg.runCommand([]string{"dependency", "update"}); err != nil {
			return err
		}
	}

	if cfg.GetValues {
		if err := cfg.runCommand([]string{"get", "values", cfg.Release, ">", cValuesFile}); err != nil {
			return err
		}
	}

	return cfg.runCommand(cfg.prepareHelmCommand())
}

func (cfg *Config) prepareHelmCommand() []string {
	switch cfg.HelmCommand {
	case cInstall:
		return cfg.setInstallCommand()
	case cUpgrade:
		return cfg.setUpgradeCommand()
	case cDelete:
		return cfg.setDeleteCommand()
	case cLint:
		return cfg.setLintCommand()
	default:
		return nil
	}
}

func (cfg *Config) setInstallCommand() []string {
	command := []string{cInstall, cfg.Chart}
	if cfg.Release != "" {
		command = append(command, "--name", cfg.Release)
	}

	if cfg.Version != "" {
		command = append(command, "--version", cfg.Version)
	}

	if cfg.Set != "" {
		command = append(command, "--set", unQuote(cfg.Set))
	}

	if cfg.SetString != "" {
		command = append(command, "--set-string", unQuote(cfg.SetString))
	}

	if cfg.Values != "" {
		for _, valuesFile := range strings.Split(cfg.Values, ",") {
			command = append(command, "--values", valuesFile)
		}
	}

	if cfg.Namespace != "" {
		command = append(command, "--namespace", cfg.Namespace)
	}

	if cfg.TillerNamespace != "" {
		command = append(command, "--tiller-namespace", cfg.TillerNamespace)
	}

	if cfg.DryRun {
		command = append(command, "--dry-run")
	}

	if cfg.Debug {
		command = append(command, "--debug")
	}

	if cfg.Wait {
		command = append(command, "--wait")
	}

	if cfg.RecreatePods {
		command = append(command, "--recreate-pods")
	}

	if cfg.ReuseValues {
		command = append(command, "--reuse-values")
	}

	if cfg.Timeout != "" {
		command = append(command, "--timeout", cfg.Timeout)
	}

	if cfg.Force {
		command = append(command, "--force")
	}

	return command
}

func (cfg *Config) setUpgradeCommand() []string {
	command := []string{cUpgrade, "--install", cfg.Release, cfg.Chart}
	if cfg.Version != "" {
		command = append(command, "--version", cfg.Version)
	}

	if cfg.Set != "" {
		command = append(command, "--set", unQuote(cfg.Set))
	}

	if cfg.SetString != "" {
		command = append(command, "--set-string", unQuote(cfg.SetString))
	}

	if cfg.Values != "" {
		for _, valuesFile := range strings.Split(cfg.Values, ",") {
			command = append(command, "--values", valuesFile)
		}
	}

	if cfg.GetValues {
		command = append(command, "--values", cValuesFile)
	}

	if cfg.Namespace != "" {
		command = append(command, "--namespace", cfg.Namespace)
	}

	if cfg.TillerNamespace != "" {
		command = append(command, "--tiller-namespace", cfg.TillerNamespace)
	}

	if cfg.DryRun {
		command = append(command, "--dry-run")
	}

	if cfg.Debug {
		command = append(command, "--debug")
	}

	if cfg.Wait {
		command = append(command, "--wait")
	}

	if cfg.RecreatePods {
		command = append(command, "--recreate-pods")
	}

	if cfg.ReuseValues {
		command = append(command, "--reuse-values")
	}

	if cfg.Timeout != "" {
		command = append(command, "--timeout", cfg.Timeout)
	}

	if cfg.Force {
		command = append(command, "--force")
	}

	return command
}

func (cfg *Config) setDeleteCommand() []string {
	command := make([]string, 2)
	command[0] = cDelete
	command[1] = cfg.Release

	if cfg.TillerNamespace != "" {
		command = append(command, "--tiller-namespace", cfg.TillerNamespace)
	}

	if cfg.DryRun {
		command = append(command, "--dry-run")
	}

	if cfg.Purge {
		command = append(command, "--purge")
	}

	return command
}

func (cfg *Config) setLintCommand() []string {
	command := make([]string, 2)
	command[0] = cLint
	command[1] = cfg.Chart

	if cfg.Set != "" {
		command = append(command, "--set", unQuote(cfg.Set))
	}

	if cfg.SetString != "" {
		command = append(command, "--set-string", unQuote(cfg.SetString))
	}

	if cfg.Values != "" {
		for _, valuesFile := range strings.Split(cfg.Values, ",") {
			command = append(command, "--values", valuesFile)
		}
	}

	if cfg.Namespace != "" {
		command = append(command, "--namespace", cfg.Namespace)
	}

	if cfg.TillerNamespace != "" {
		command = append(command, "--tiller-namespace", cfg.TillerNamespace)
	}

	if cfg.Debug {
		command = append(command, "--debug")
	}

	return command
}

func (cfg *Config) addHelmRepo() error {
	for _, v := range cfg.HelmRepos {
		err := cfg.runCommand([]string{"repo", "add", v})
		if err != nil {
			return err
		}
	}

	return nil
}

func (cfg *Config) initKubeconfig() (err error) {
	userHomeDir, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	f, err := os.Create(fmt.Sprintf("%s/.kube/config", userHomeDir))
	if err != nil {
		return err
	}

	defer func() {
		err1 := f.Close()
		if err1 != nil {
			err = fmt.Errorf("%v: %w", err1.Error(), err)
		}
	}()

	t, err := template.ParseFiles(cfg.KubeConfigTemplate)
	if err != nil {
		return err
	}

	return t.Execute(f, cfg)
}

func (cfg *Config) runCommand(params []string) error {
	cmd := new(exec.Cmd)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd = exec.Command(cfg.HelmBinary, params...) //nolint:gosec //reason:this app is command line builder

	return cmd.Run()
}

func unQuote(s string) string {
	unquoted, err := strconv.Unquote(s)
	if err != nil {
		return s
	}

	return unquoted
}
