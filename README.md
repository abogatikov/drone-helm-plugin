# Helm plugin for drone.io

This plugin allows to deploy a [Helm](https://github.com/kubernetes/helm) chart into a [Kubernetes](https://github.com/kubernetes/kubernetes) cluster.

* Current `helm` version: 2.16.1
* Current `kubectl` version: 1.16.2

## Usage

```yaml
    - name: helm_command
      image: docker.pkg.github.com/abogatikov/drone-helm-plugin/drone-helm-plugin
      settings:
        api_server: "https://kubernetes.default"
        token: YZ.....
        certificate: YZ.....
        service_account: kubeszer
        kube_config: ./.kube/config
        helm_command: upgrade
        namespace: sre
        tls_verify: true
        set: alice=alice_value,bob=bob_value
        set_string: alice_string=alice_string_value,bob_string=bob_string_value
        values: /drone/src/values1.yml,/drone/src/values2.yml
        get_values: true
        helm_repos: [https://charts.jetstack.io, https://charts.gitlab.io/]
        release: app-release-name
        chart: app
        chart_version: 1.0.0
        debug: false
        dry-run: false
        tiller_namespace: kube-system
        wait: true
        recreate_pods: true
        reuse_values: true
        timeout: 1s
        force: true
        update_dependencies: true
```

## License

[MIT](LICENSE) © [Alex Bogatikov](https://github.com/abogatikov) 