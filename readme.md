# KASOS

> kasos 是一个基于 kubernetes 的水平自动扩缩容系统, 用基于数据的深度学习方法改善K8s水平自动扩缩容的效果
>

## 项目说明

### /server

**系统管控模块**用于管理系统的元信息并提供对外交互接口

### /kube-prometheus

**数据收集模块**采用了开源的kube-prometheus解决方案，用于灵活增删待监控的服务并配置Grafana看板

### /trainer

**模型训练模块**自动从Prometheus拉取任务并触发模型训练和更新

### /infer-module

**模型应用模块**用于周期性地推理模型预测服务在下一周期的QPS值, 并提供了训练和推理脚本检查的功能

### /hpa-executor

**扩缩容执行模块**用于收到QPS值后计算服务下一周期的Pod逾期数并做出调整, 提供***临界QPS***的测试功能

### /build

**K8s资源构建工具**收录了所有本模块在K8s集群中需要构建出的资源定义

### /bootstrap.sh

**一键启动及构建脚本**可一键启动minikube集群并在集群中部署本系统, 执行前请确保运行环境的docker守护进程在后台运行

### /examples/measure

**测试服务**用于测试本项目的功能和效果的简单web服务

### /examples/product

**产物与数据**本项目的测试数据与模型的对比结果收录

### /examples/scripts

**脚本**本项目的测试脚本, 数据分析脚本
