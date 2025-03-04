#!/usr/bin/env bash
#shellcheck disable=SC1091,SC2015,SC2317
set -e -o pipefail
trap:exit() {
    {
        docker stop opa
        docker rm opa
        docker stop busybox
        docker rm busybox
        docker network rm opa
    } &>/dev/null

}
trap "trap:exit" EXIT INT TERM
{
    docker network rm opa 
    docker network create opa
} &>/dev/null
docker run --name opa -d \
    -u root \
    -v "$(pwd)/policies":/policies \
    --network opa \
    -p 8181:8181 openpolicyagent/opa run --server --addr :8181 &>/dev/null

until curl -s http://localhost:8181/health &>/dev/null; do
    echo "Waiting for OPA server to start"
    sleep 1
done
echo "OPA server started"
# ! docker exec opa opa build -b /policies -t rego && {
! docker exec opa opa build --bundle /policies --output /policies/bundle.tar.gz && {
    echo "Failed to compile policies"
    exit 1
} || {
    echo "Policies compiled"
}

# run opa test
! docker exec opa opa test /policies -v && {
    echo "Failed to run OPA test"
    exit 1
} || {
    echo "OPA test run"
}
# create policies
opaPolicies=(
    "${PWD}/policies/public.rego"
    "${PWD}/policies/zerotrust.rego"
)
for policy in "${opaPolicies[@]}"; do
    curl -i -X PUT --data-binary "@$policy" http://localhost:8181/v1/policies/"$(basename "$policy")" &>/dev/null
done
# curl -s http://localhost:8181/v1/policies | jq .

# use k6 docker image grafana/k6
# docker run -u root  --rm -i -v $PWD:/app -w /app grafana/k6 new
! docker run --rm -it --name k6 -v "$(pwd)":/tests --network opa grafana/k6 run /tests/script.js && {
    echo "Failed to run k6"
    exit 1
} || {
    echo "K6 run"
}

# # show opa logs
# docker logs opa | grep -i error && {
#     echo "OPA error"
#     exit 1
# } || {
#     echo "OPA logs"
# }
exit 0

# # start busybox
# docker run -u root --name busybox -d node:20-alpine tail -f /dev/null
# # install curl
# ! docker exec busybox sh -c "apk add --no-cache curl bash" && {
#     echo "Failed to install curl"
#     exit 1
# } || {
#     echo "Curl installed"
# }
# # install k6
# ! docker exec busybox sh -c "npm install -g k6" && {
#     echo "Failed to install k6"
#     exit 1
# } || {
#     echo "K6 installed"
# }
# # run k6
