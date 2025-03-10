#!/bin/bash
# shellcheck disable=all
set -e
clear
declare -gxI STACK_CURRENT_DIR="$(dirname "$(realpath "${BASH_SOURCE[0]:-$0}")")"
declare -gxIA STACK_VARIABLES=()
declare -gxIA STACK_SERVICES=()
declare -gxIA STACK_VOLUMES=()
declare -gxIA STACK_NETWORKS=()
declare -gxIA STACK_PROFILES=()
declare -gxA STACK_ENV=()
declare -gxA STACK_DOTENVS=()
dotenv:load() {

    ! test -f "$1" && return 0
    [[ -n "${STACK_DOTENVS["$1"]}" ]] && return 0
    local prefix=${2:-STACK} && [[ "$prefix" =~ _$ ]] && prefix=${prefix%_}    
    echo "Loading $1 into $prefix"
    while read -r line || [[ -n $line ]]; do
        [[ -z "$line" ]] && continue
        [[ "$line" =~ ^#.*$ ]] && continue
        read -r key value <<<"$(echo "$line" | awk -F'=' '{print $1, $2}')"
        [[ -z "$key" ]] && continue
        [[ "$key" =~ ^${prefix^^}_ ]] && key=${key#*_}
        [[ "$value" =~ \$\{.*\} ]] && value=$(eval echo "$value")
        declare -gxI "${prefix^^}_${key}"="$value"
        # dotenv:upinset .env.example "${prefix^^}_${key}" "$value"
        STACK_ENV["${prefix^^}_${key}"]="$value"
        unset key value
    done <"$1"
    STACK_DOTENVS["$1"]="$prefix"
    unset prefix
    return 0
}
dotenv:upinset() {
    ! test -f "$1" && return 0
    local file=$1
    local key=$2
    local value=$3
    local delimiter=${4:-=}
    local pattern="^$key$delimiter"
    if grep -q "$pattern" "$file"; then
        # ignore value
        sed -i -r "s|^$key$delimiter.*$|$key$delimiter$value|" "$file"
        # sed -i "s|$pattern.*|$key$delimiter$value|g" "$file"
    else
        echo "$key$delimiter$value" >>"$file"
    fi
    unset file key value delimiter pattern
}
dotenv:svc() {
    local keys=$(echo "${!STACK_SERVICES[@]}" | tr ' ' '\n' | sort -u)
    for svc in $keys; do
        local svcDir=${STACK_CURRENT_DIR}/apps/${svc,,}
        test -f "${svcDir}/.env" &&  dotenv:load "${svcDir}/.env" "SVC_${svc^^}" && continue
        # when service is not found, try get first part of svc splited by -
        local svcName=${svc%%-*}
        local svcDir=${STACK_CURRENT_DIR}/apps/${svcName,,}
        test -f "${svcDir}/.env" && dotenv:load "${svcDir}/.env" "SVC_${svcName^^}" && continue
        # when service is not found, try get first part of svc splited by _
        local svcName=${svc%%_*}
        local svcDir=${STACK_CURRENT_DIR}/apps/${svcName,,}
        test -f "${svcDir}/.env" && dotenv:load "${svcDir}/.env" "SVC_${svcName^^}" && continue
        unset svcDir svcName        
    done
    unset keys svc
    return 0
}
dotenv:export(){
    {
        local keys=$(echo "${!STACK_ENV[@]}" | tr ' ' '\n' | sort -u)
        for key in $keys; do
            echo "$key=\"${STACK_ENV[$key]}\""
        done
        unset keys key
    } > ".env.${STACK_NAME}.${STACK_ENV}"    
}
compose:inspect() {
    for dockerFlag in {"services","volumes","profiles"}; do
        while read -r line || [[ -n $line ]]; do
            case $dockerFlag in
            "services")
                STACK_SERVICES["${line}"]="${line}"
                ;;
            "volumes")
                STACK_VOLUMES["${line}"]="${line}"
                ;;
            "profiles")
                STACK_PROFILES["${line}"]="${line}"
                ;;
            esac
        done < <(docker compose config --$dockerFlag 2>/dev/null)
        unset dockerFlag
    done
    return 0
}
compose:doc:variables(){
    echo "## Variables"
    echo "| Name | Required | Default | Alternate | Value |"
    echo "| ---- | -------- | ------- | --------- | ----- |"
    while read -r varLine || [[ -n $varLine ]]; do
        var=$(echo "$varLine" | awk '{print $1}')
        required=$(echo "$varLine" | awk '{print $2}')
        default=$(echo "$varLine" | awk '{print $3}')
        alternate=$(echo "$varLine" | awk '{print $4}')
        STACK_ENV["${var}"]="${!var:-${default:-""}}"
        varValue=$(echo "${!var}" | sed -r 's/(.{2}).*(.{1})/\1****\2/')
        echo "| $var | $required | $default | $alternate | ${varValue} |"
        STACK_VARIABLES["${var}"]="${default:-""}"
        unset var required default alternate varValue
    done < <(docker compose config --variables | tail -n +2 | sort -u)
}
compose:doc:services(){
    echo "## Services"
    local keys=$(echo "${!STACK_SERVICES[@]}" | tr ' ' '\n' | sort -u)
    for svc in $keys; do
        echo "### $svc"
    done
}
compose:doc:volumes(){
    echo "## Volumes"
    keys=$(echo "${!STACK_VOLUMES[@]}" | tr ' ' '\n' | sort -u)
    for vol in $keys; do
        echo "### $vol"
    done
}
compose:doc:profiles(){
    echo "## Profiles"
    keys=$(echo "${!STACK_PROFILES[@]}" | tr ' ' '\n' | sort -u)
    for prof in $keys; do
        echo "### $prof"
    done
}
compose:doc:environment(){
    echo "## Environment"
    while read -r env || [[ -n $env ]]; do
        read -r varName varValue <<<"$(echo "$env" | awk -F'=' '{print $1, $2}')"
        echo "<details>"
        echo "<summary>$varName</summary>"
        # mask value showing only first 3 and last 3 characters
        varValue=$(echo "$varValue" | sed -r 's/(.{2}).*(.{1})/\1****\2/')
        echo "$varValue"
        echo "</details>"
        unset env varName varValue
    done < <(docker compose config --environment | grep -E "STACK_|SVC_" | sort -u)
}
compose:doc() {
    {
        echo "# Stack ***${STACK_NAME}***"
        echo "## Environment"
        echo " - **${STACK_ENV}**"
        echo "## Overview"
        echo "This is a stack of services that are deployed using docker-compose."
        compose:doc:variables
        compose:doc:services
        compose:doc:volumes
        compose:doc:profiles
        compose:doc:environment
        echo 
        echo "## Conclusion"
        echo "This documentation provides an overview of the services and configurations in the stack."
        # printenv | grep -E "STACK_" | sort
    } >"${STACK_NAME}.${STACK_ENV}".md
    return 0
}
compose:generate() {
    # --no-interpolate --resolve-image-digests --no-consistency --no-normalize --no-path-resolution
    docker compose config >"${STACK_NAME}.${STACK_ENV}.compose.yaml"
    # envsubst <"compose.${STACK_NAME}.${STACK_ENV}.yaml" #>"compose.${STACK_NAME}.${STACK_ENV}.yaml.tmp"
    #mv "compose.${STACK_NAME}.${STACK_ENV}.yaml.tmp" "compose.${STACK_NAME}.${STACK_ENV}.yaml"
}
compose:deploy(){
    docker compose -f "${STACK_NAME}.${STACK_ENV}.compose.yaml" down --remove-orphans --volumes
    docker compose -f "${STACK_NAME}.${STACK_ENV}.compose.yaml" up --detach --remove-orphans --renew-anon-volumes --wait # --force-recreate
    docker compose -f "${STACK_NAME}.${STACK_ENV}.compose.yaml" ps
}
swarm:generate() {
    cp "${STACK_NAME}.${STACK_ENV}.compose.yaml" "${STACK_NAME}.${STACK_ENV}.stack.yaml"
    yq eval '
        del(.name) |
        del(.networks[] | .ipam) |
        del(.networks[] | .enable_ipv4) |
        del(.networks[] | .enable_ipv6) |
        del(.networks[] | .driver_opts) |
        del(.services[].depends_on) |
        with(
            .services[];         
            .deploy.resources.limits.cpus |= (select(.) | tostring) |
            .deploy.resources.limits.memory |= (select(.) | tostring) |
            .deploy.resources.reservations.cpus |= (select(.) | tostring) |
            .deploy.resources.reservations.memory |= (select(.) | tostring)
        ) |
        .services[] |= (
            with(select(.deploy.resources == null);
                del(.deploy.resources)
            ) |
            with(select(.deploy.resources.limits == null);
                del(.deploy.resources.limits)
            ) |
            with(select(.deploy.resources.limits.memory == null);
                del(.deploy.resources.limits.memory)
            ) |
            with(select(.deploy.resources.limits.cpus == null);
                del(.deploy.resources.limits.cpus)
            ) |
            with(select(.deploy.resources.reservations == null);
                del(.deploy.resources.reservations)
            ) |
            with(select(.deploy.resources.reservations.memory == null);
                del(.deploy.resources.reservations.memory)
            ) |
            with(select(.deploy.resources.reservations.cpus == null);
                del(.deploy.resources.reservations.cpus)
            )                
        ) |
        .services[] |= (
            with(select(.ports != null);
                .ports[] |= (
                    .published |= (select(.) | tonumber) |
                    .target |= (select(.) | tonumber) |
                    del(.host_ip)
                )
            )
        )
    ' "${STACK_NAME}.${STACK_ENV}.stack.yaml" -i        
    # |
    # .mode = "host"
}
swarm:deploy(){
    docker stack deploy --detach=true --compose-file "${STACK_NAME}.${STACK_ENV}.stack.yaml" "${STACK_PREFIX}"
}
stack:setup(){
    return 0
    test -n "${STACK_TMP_DIR}" && ! test -d "${STACK_TMP_DIR}" && mkdir -p "${STACK_TMP_DIR}"
    test -n "${STACK_LOGS_DIR}" && ! test -d "${STACK_LOGS_DIR}" && mkdir -p "${STACK_LOGS_DIR}"
}
docker:events(){
    {
        while read -r line || [[ -n $line ]]; do
            echo "$line"
        done < <(docker events --since 1m --until 1m)
    } & DOCKER_EVENTS_PID=$!
    trap 'kill '"$DOCKER_EVENTS_PID"'' EXIT INT TERM
    wait $DOCKER_EVENTS_PID
}
dotenv:load .env
compose:inspect
dotenv:svc
compose:doc
dotenv:export
stack:setup
compose:generate
swarm:generate
# docker:events # & DOCKER_EVENTS_PID=$!
swarm:deploy



exit 0
[[ "${STACK_NETWORK_BRIDGE_EXTERNAL}" == "true" ]] && {
    # --ip-range="${STACK_NETWORK_BRIDGE_IP_RANGE}.0/24" \
    echo "Creating network stack-${STACK_NETWORK_NAME} with prefix ${STACK_NETWORK_BRIDGE_PREFIX}"
    docker network create --driver=bridge --attachable \
        --subnet=10.0.0.0/16 \
        --gateway=10.0.0.1 \
        --opt "com.docker.network.bridge.name=stack-${STACK_NETWORK_NAME}" \
        "stack-${STACK_NETWORK_NAME}" >/dev/null || true
}
[[ "${STACK_NETWORK_OVERLAY_EXTERNAL}" == true ]] && {
    docker network create --driver=overlay --attachable "stack-${STACK_NETWORK_NAME}-edge" >/dev/null || true
    docker network create --driver=overlay --attachable "stack-${STACK_NETWORK_NAME}-ipv6" >/dev/null || true
}



# ensure volumes local paths has been created
docker compose -f "compose.${STACK_NAME}.${STACK_ENV}.yaml" down --remove-orphans --volumes
docker compose -f "compose.${STACK_NAME}.${STACK_ENV}.yaml" up --detach --remove-orphans --renew-anon-volumes --wait # --force-recreate
docker compose -f "compose.${STACK_NAME}.${STACK_ENV}.yaml" ps

# deploy as a stack
cp "compose.${STACK_NAME}.${STACK_ENV}.yaml" "compose.${STACK_NAME}.${STACK_ENV}.stack.yaml"
# fix: (root) Additional property name is not allowed
# fix: networks.*.ipam, delete ipam
# fix: networks.stack.ipam.config.0 Additional property aux_addresses is not allowed
# fix: networks.stack.ipam.config.0 Additional property ip_range is not allowed
# fix: networks.* Additional property enable_ipv4 is not allowed
# fix: networks.* Additional property enable_ipv6 is not allowed
yq eval '
    del(.name) |
    del(.networks[] | .ipam) |
    del(.networks[] | .enable_ipv4) |
    del(.networks[] | .enable_ipv6) |
    del(.networks[] | .driver_opts) |
    del(.services[].depends_on) |
    with(
        .services[];         
        .deploy.resources.limits.cpus |= (select(.) | tostring) |
        .deploy.resources.limits.memory |= (select(.) | tostring) |
        .deploy.resources.reservations.cpus |= (select(.) | tostring) |
        .deploy.resources.reservations.memory |= (select(.) | tostring)
    ) |
    .services[] |= (
        with(select(.deploy.resources == null);
            del(.deploy.resources)
        ) |
        with(select(.deploy.resources.limits == null);
            del(.deploy.resources.limits)
        ) |
        with(select(.deploy.resources.limits.memory == null);
            del(.deploy.resources.limits.memory)
        ) |
        with(select(.deploy.resources.limits.cpus == null);
            del(.deploy.resources.limits.cpus)
        ) |
        with(select(.deploy.resources.reservations == null);
            del(.deploy.resources.reservations)
        ) |
        with(select(.deploy.resources.reservations.memory == null);
            del(.deploy.resources.reservations.memory)
        ) |
        with(select(.deploy.resources.reservations.cpus == null);
            del(.deploy.resources.reservations.cpus)
        )                
    )
' "compose.${STACK_NAME}.${STACK_ENV}.stack.yaml" -i

docker stack deploy --detach=false --compose-file "compose.${STACK_NAME}.${STACK_ENV}.stack.yaml" "${STACK_NAME}"
