# gitlab-prometheus-exporter
A prom exporter for the gitlab installations. This exporter talks directly to
your Gitlab API.

## Docker

```
docker run -p 9115:9115 -e GITLAB_API=http://web.gitlab.svc/api -e GITLAB_TOKEN=some-token -e HTTP_LISTENADDR=":9115" -it --rm deviavir/gitlab-prometheus-exporter:latest
```
