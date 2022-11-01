#!/bin/bash
PROJECT_NAME="aiapp"
DOCKER_REGISTRY_SERVER="registry.aiedge.ndsl-lab.cn"
CPU_ARCH=$(dpkg --print-architecture)
TAG=v2.0
if [[ "$1" != "" ]]; then
    TAG=$1
fi

function build-and-push () {
    function log-info() { echo -e "\033[32m[info] $1\033[0m"; }
    function log-error() { echo -e "\033[31m[error] $1\033[0m"; }
    local project_name=$1
    local build_arch=$2
    local tag=$3
    local dist_image=aiedge/${project_name}:${tag}
    log-info "build for ${build_arch}"
    case ${build_arch} in
    "amd64")
        # dist_image=aiedge/${project_name}:${tag}
        log-info "building image ${dist_image}"
        docker build -t ${dist_image} .
        ;;
    "arm64")
        dist_image=aiedge/${project_name}-${build_arch}:${tag}
        log-info "building image ${dist_image}"
        if [ -f "${build_arch}.dockerfile" ]; then
            docker build -f ${build_arch}.dockerfile -t ${dist_image} .
        else
            log-info "no ${build_arch}.dockerfile, use dockerfile"
            docker build -t ${dist_image} .
        fi
        ;;
    *)
        log-error "$build_arch arch is not supported"
        exit 1
        ;;
    esac
    log-info "pushing image ${DOCKER_REGISTRY_SERVER}/${dist_image}"
    docker tag ${dist_image} ${DOCKER_REGISTRY_SERVER}/${dist_image}
    docker push ${DOCKER_REGISTRY_SERVER}/${dist_image}
}

build-and-push $PROJECT_NAME $CPU_ARCH $TAG
