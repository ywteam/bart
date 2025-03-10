# Stack ***yellowteam***
## Environment
 - **dev**
## Overview
This is a stack of services that are deployed using docker-compose.
## Variables
| Name | Required | Default | Alternate | Value |
| ---- | -------- | ------- | --------- | ----- |
| apr1 | false |  |  |  |
| CF_API_EMAIL | false |  |  |  |
| CF_API_KEY | false |  |  |  |
| nn6 | false |  |  |  |
| STACK_DEPLOY_RESTART_CONDITION | false | on-failure |  | on****e |
| STACK_DEPLOY_ROLLBACK_FAILURE_ACTION | false | rollback |  | ro****k |
| STACK_DOMAIN | false | yellowteam.dev |  | ye****z |
| STACK_ENV | false |  |  | de****v |
| STACK_GROUP_ID | false | 1000 |  | 10****0 |
| STACK_LOGGING_DRIVER | false | json-file |  | js****e |
| STACK_LOGGING_MAX_FILE | false | 5 |  | 5 |
| STACK_LOGGING_MAX_SIZE | false | 5m |  | 5m |
| STACK_LOGS_DIR | false | /var/log/yellow |  | /v****m |
| STACK_NAME | false |  |  | ye****m |
| STACK_NETWORK_BRIDGE_EXTERNAL | false | true |  | fa****e |
| STACK_NETWORK_BRIDGE_PREFIX | false | 10.0.0 |  | 10****0 |
| STACK_NETWORK_OVERLAY_EXTERNAL | false | false |  | fa****e |
| STACK_PREFIX | false |  |  | ye****v |
| STACK_TMP_DIR | false | /tmp/yellow |  | /t****m |
| STACK_TZ | false | UTC |  | UT****C |
| STACK_USER_ID | false | 1000 |  | 10****0 |
| STACK_VERSION | false | 0.0.1-alpha-0 |  | 0.****0 |
| SVC_AUTHENTIK_EMAIL__FROM | false |  |  | au****m |
| SVC_AUTHENTIK_EMAIL__HOST | false |  |  | sm****m |
| SVC_AUTHENTIK_EMAIL__PASSWORD | false |  |  | 1 |
| SVC_AUTHENTIK_EMAIL__PORT | false |  |  | 58****7 |
| SVC_AUTHENTIK_EMAIL__TIMEOUT | false |  |  | 5 |
| SVC_AUTHENTIK_EMAIL__USERNAME | false |  |  | 1 |
| SVC_AUTHENTIK_EMAIL__USE_SSL | false |  |  | tr****e |
| SVC_AUTHENTIK_EMAIL__USE_TLS | false |  |  | tr****e |
| SVC_AUTHENTIK_ERROR_REPORTING__ENABLED | false |  |  | fa****e |
| SVC_AUTHENTIK_PG_DB | false | authentik |  | au****k |
| SVC_AUTHENTIK_PG_PASS | false |  |  | gc****n |
| SVC_AUTHENTIK_PG_USER | false | authentik |  | au****k |
| SVC_AUTHENTIK_PORT_HTTP | false | 9000 |  | 90****0 |
| SVC_AUTHENTIK_PORT_HTTPS | false | 9443 |  | 94****3 |
| SVC_AUTHENTIK_SECRET_KEY | false |  |  | 0T****+ |
| SVC_AUTHENTIK_VERSION | false | 2025.2.1 |  | 20****1 |
| SVC_TRAEFIK_ACME_EMAIL | false |  |  | ra****m |
| SVC_TRAEFIK_AUTH_PASSWORD_HASH | false | '$apr1$yzXCOaAe$nn6/HUwjA/QwHLTaF6DYT1' |  | '$****' |
| SVC_TRAEFIK_AUTH_USER | false | admin |  | ad****n |
| SVC_TRAEFIK_DASHBOARD_CREDENTIALS | false |  |  |  |
| SVC_TRAEFIK_LOG_LEVEL | false | INFO |  | DE****G |
| SVC_TRAEFIK_VERSION | false | v3.3.4 |  | v3****4 |
| yzXCOaAe | false |  |  |  |
## Services
### authentik-postgresql
### authentik-proxy
### authentik-redis
### authentik-server
### authentik-worker
### error-pages
### traefik
### whoami
## Volumes
### authentik-certs
### authentik-database
### authentik-redis
### stack-data
### stack-logs
### stack-tmp
### stack-tmpfs
### traefik-acme
## Profiles
## Environment
<details>
<summary>STACK_ARCH</summary>
am****4
</details>
<details>
<summary>STACK_CURRENT_DIR</summary>
/d****r
</details>
<details>
<summary>STACK_DEPLOY_RESTART_CONDITION</summary>
on****e
</details>
<details>
<summary>STACK_DEPLOY_ROLLBACK_FAILURE_ACTION</summary>
ro****k
</details>
<details>
<summary>STACK_DOMAIN</summary>
ye****z
</details>
<details>
<summary>STACK_ENV</summary>
de****v
</details>
<details>
<summary>STACK_GROUP_ID</summary>
10****0
</details>
<details>
<summary>STACK_LOGGING_DRIVER</summary>
js****e
</details>
<details>
<summary>STACK_LOGGING_MAX_FILE</summary>
5
</details>
<details>
<summary>STACK_LOGGING_MAX_SIZE</summary>
5m
</details>
<details>
<summary>STACK_LOGS_DIR</summary>
/v****m
</details>
<details>
<summary>STACK_NAME</summary>
ye****m
</details>
<details>
<summary>STACK_NETWORK_BRIDGE_EXTERNAL</summary>
fa****e
</details>
<details>
<summary>STACK_NETWORK_BRIDGE_IP_RANGE</summary>
10****0
</details>
<details>
<summary>STACK_NETWORK_BRIDGE_PREFIX</summary>
10****0
</details>
<details>
<summary>STACK_NETWORK_NAME</summary>
ye****m
</details>
<details>
<summary>STACK_NETWORK_OVERLAY_EXTERNAL</summary>
fa****e
</details>
<details>
<summary>STACK_PLATFORM</summary>
li****x
</details>
<details>
<summary>STACK_PREFIX</summary>
ye****v
</details>
<details>
<summary>STACK_TMP_DIR</summary>
/t****m
</details>
<details>
<summary>STACK_TZ</summary>
UT****C
</details>
<details>
<summary>STACK_USER_ID</summary>
10****0
</details>
<details>
<summary>STACK_VERSION</summary>
0.****0
</details>
<details>
<summary>SVC_AUTHENTIK_EMAIL__FROM</summary>
au****m
</details>
<details>
<summary>SVC_AUTHENTIK_EMAIL__HOST</summary>
sm****m
</details>
<details>
<summary>SVC_AUTHENTIK_EMAIL__PASSWORD</summary>
1
</details>
<details>
<summary>SVC_AUTHENTIK_EMAIL__PORT</summary>
58****7
</details>
<details>
<summary>SVC_AUTHENTIK_EMAIL__TIMEOUT</summary>
5
</details>
<details>
<summary>SVC_AUTHENTIK_EMAIL__USERNAME</summary>
1
</details>
<details>
<summary>SVC_AUTHENTIK_EMAIL__USE_SSL</summary>
tr****e
</details>
<details>
<summary>SVC_AUTHENTIK_EMAIL__USE_TLS</summary>
tr****e
</details>
<details>
<summary>SVC_AUTHENTIK_ERROR_REPORTING__ENABLED</summary>
fa****e
</details>
<details>
<summary>SVC_AUTHENTIK_PG_DB</summary>
au****k
</details>
<details>
<summary>SVC_AUTHENTIK_PG_PASS</summary>
gc****n
</details>
<details>
<summary>SVC_AUTHENTIK_PG_USER</summary>
au****k
</details>
<details>
<summary>SVC_AUTHENTIK_PORT_HTTP</summary>
90****0
</details>
<details>
<summary>SVC_AUTHENTIK_PORT_HTTPS</summary>
94****3
</details>
<details>
<summary>SVC_AUTHENTIK_SECRET_KEY</summary>
0T****+
</details>
<details>
<summary>SVC_AUTHENTIK_TOKEN</summary>
JS****2
</details>
<details>
<summary>SVC_AUTHENTIK_VERSION</summary>
20****1
</details>
<details>
<summary>SVC_TRAEFIK_ACME_EMAIL</summary>
ra****m
</details>
<details>
<summary>SVC_TRAEFIK_AUTH_PASSWORD</summary>
05****d
</details>
<details>
<summary>SVC_TRAEFIK_AUTH_PASSWORD_HASH</summary>
'$****'
</details>
<details>
<summary>SVC_TRAEFIK_AUTH_USER</summary>
ad****n
</details>
<details>
<summary>SVC_TRAEFIK_CF_API_EMAIL</summary>
ra****m
</details>
<details>
<summary>SVC_TRAEFIK_CF_API_KEY</summary>
Dd****K
</details>
<details>
<summary>SVC_TRAEFIK_DASHBOARD_CREDENTIALS</summary>

</details>
<details>
<summary>SVC_TRAEFIK_LOG_LEVEL</summary>
DE****G
</details>
<details>
<summary>SVC_TRAEFIK_VERSION</summary>
v3****4
</details>

## Conclusion
This documentation provides an overview of the services and configurations in the stack.
