# Author: lmk17@tsinghua.org.cn
# Date: 2024-03-08
# Description: Deploy kasos in local minikube environment
#!/usr/bin/bash

KASOS_VERSION=v0.1

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

echo ">>> Start build and deploy kasos..."
start=$(date +%s)

# 切换docker环境
eval $(minikube -p minikube docker-env)

# 拉取golang基础镜像
docker pull golang:1.21.6
# 拉取python基础镜像
docker pull python:3.10.13
# 拉取mysql基础镜像
docker pull mysql:8.0
# 构建镜像
build_image ../server kasos-server $KASOS_VERSION
build_image ../trainer kasos-trainer $KASOS_VERSION
build_image ../infer-module kasos-infer-module $KASOS_VERSION
build_image ../hpa-executor kasos-hpa-executor $KASOS_VERSION
# 创建k8s资源
# 请务必注意创建顺序
# 先创建ServiceAccount，再创建ClusterRoleBinding
cmd_log "kubectl apply -f ./account.yaml"
cmd_log "kubectl apply -f ./role-binding.yaml"
# 创建存储类
cmd_log "kubectl apply -f ./local-storage.yaml"
# 创建持久化存储卷和卷声明
cmd_log "kubectl apply -f ./mysql-pv.yaml"
cmd_log "kubectl apply -f ./mysql-pvc.yaml"
cmd_log "kubectl apply -f ./public-pv.yaml"
cmd_log "kubectl apply -f ./public-pvc.yaml"
# 创建ConfigMap
cmd_log "kubectl apply -f ./configmap.yaml"
# 创建mysql Deployment and Service
cmd_log "kubectl apply -f ./mysql-deployment.yaml"
cmd_log "kubectl apply -f ./mysql-service.yaml"
# 创建server Deployment and Service
cmd_log "kubectl apply -f ./server-deployment.yaml"
cmd_log "kubectl apply -f ./server-service.yaml"
# 创建hpa-executor Deployment and Service
cmd_log "kubectl apply -f ./hpa-executor-deployment.yaml"
cmd_log "kubectl apply -f ./hpa-executor-service.yaml"
# 创建infer-module Deployment and Service
cmd_log "kubectl apply -f ./infer-module-deployment.yaml"
cmd_log "kubectl apply -f ./infer-module-service.yaml"
# 创建trainer CronJob
cmd_log "kubectl apply -f ./trainer-cronjob.yaml"

end=$(date +%s)
echo ">>> Build and deploy kasos successfully, cost $((end-start)) seconds."