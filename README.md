# Kubernetes Scheduler Simulator

本项目用于测试Kubernetes调度器的性能。

## 主要功能

### 模拟集群构建

采用真实数据进行集群的构建和模拟，比如[阿里云的监控数据](https://github.com/alibaba/clusterdata)。

若要提供新的数据，只需要实现接口即可。

### 模拟

每一时刻，计算集群内节点的负载情况，并将新的Pod提供给配置的调度器进行调度。

### 监控数据采集

用于衡量调度器的性能。

## 调度器扩展

采用Kubernetes-v1.18.0版本的调度器框架，根据该框架的设计扩展调度器。

### Profile提供的多调度器

本项目引入的是Kubernetes自带的`GenericScheduler`，因此只有一个调度器在运行。

`GenericScheduler`支持根据Pod的`spec.schedulerName`获取使用的调度器Profile的功能。每个Profile对应一套插件的配置。

### 插件支持

通过实现Kubernetes自带的插件接口，可以加入自己需要的插件，并在Profile中指定的生命周期函数设置该插件。

## TODO List

- [ ] 数据读取接口的设计
- [ ] 模拟器设计
  - [ ] Pod模拟器设计
  - [ ] Node模拟器设计
- [ ] 监控系统的设计