### Github action
[release](.github/workflows/release.yaml) workflow actions:
- create new tag for the release
- building image and pushing it to the public dockerhub repo
- packaging helm chart and pushing it to OCI registry (same dockerhub, different repository)

### Helm chart installation
```bash
helm upgrade <my-release> --install -n pizza-store --create-namespace oci://registry-1.docker.io/len4i/pizza-store-helm --version <chart-version>
```
List of versions can be found in the [dockerhub repo](https://hub.docker.com/r/len4i/pizza-store-helm/tags)
Tag of the image and helm chart version are the same


### Logger
Using 2 loggers:
1. standard slog in json format for the app logs
2. slog based logger for http requests with middleware that enriches log with http request data
_http logger servers the only puprose to align loggin format for access logs, so that it can be collected same way as an app logs_

### HA
In current mode with local db (sqlite) there is no much point in HA  
However sqlite can be easily swapped with real mysql / postgres DB