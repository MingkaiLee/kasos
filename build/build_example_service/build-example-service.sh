# Author: lmk17@tsinghua.org.cn
# Date: 2024-03-10
# Description: Deploy service whose name is measure in example directory
#!/usr/bin/bash

build_image() {
    local image_path=$1
    local image_name=$2
    local image_tag=$3
    local work_dir=$(pwd)
    cd $image_path
    docker build -t $image_name:$image_tag .
    cd $work_dir
}

cmd_log() {
    local cmd=$1
    echo ">>> $cmd"
    eval $cmd
}

# 切换docker环境
eval $(minikube -p minikube docker-env)

build_image "../../examples/measure" "example-measure" "v0.1"

cmd_log "kubectl apply -f ./measure-deployment.yaml"
cmd_log "kubectl apply -f ./measure-service.yaml"