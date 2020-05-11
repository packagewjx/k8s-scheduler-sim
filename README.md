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

采用`Kubernetes-v1.18.0`版本的调度器框架，根据该框架的设计扩展调度器。

### `Profile`提供的多调度器

本项目引入的是Kubernetes自带的`GenericScheduler`，因此只有一个调度器在运行。

`GenericScheduler`支持根据Pod的`spec.schedulerName`获取使用的调度器Profile的功能。每个Profile对应一套插件的配置。

### 插件支持

通过实现Kubernetes自带的插件接口，可以加入自己需要的插件，并在Profile中指定的生命周期函数设置该插件。

## 模拟器设计思想

### 时钟周期

时间的概念在这里简化为了一个个时钟周期，每个具有状态的对象，均需要实现`Tick`函数，该函数负责在每个时钟周期更新对象的状态。

`Tick`函数的参数没有时间，而目前一个Tick对应现实时间为多少暂时没有定下来，初步定位一秒钟。可能以后可以在模拟器初始化时
设置Tick的长度，以提高模拟的精度。由于其参数不带时间，因此更新的逻辑需要以近期的状态而定，而不是根据监控数据在此时的实际
数值来定。模拟的准确度依赖于Pod核心算法的设计。

由于时钟周期非现实时间，因此异步通知机制将会导致应该发生在某个时间时间的事件，发生在了后面的时间周期中。为了避免这个情况，
除了特别需要异步通知的组件外（如Informer），其余组件将采用同步设计。

### `Pod`

`Pod`被设计为一个`struct`，其保存了Kubernetes的`Pod`信息，以及集群内模拟所需要的通用信息。`Pod`的具体运行逻辑由`PodAlgorithm`
决定，该字段为接口，用户可以根据自己的模拟逻辑，设计并实现自己的`Pod`。

### `Node`

模拟节点，管理节点上的所有`Pod`，分配资源给所有的`Pod`运行。

#### CPU时间片分配

现实的节点中，进程时间片的分配由内核调度器做决定。为了提高模拟的精度，将调度的逻辑交给接口`CoreScheduler`处理。
用户可以实现属于自己的内核调度器，模拟现实中调度器调度进程的行为，可以使用不同的策略调度`Pod`。

#### 状态更新

每一轮的状态更新的主要流程如下：

1. 查询当前待运行和已停止的`Pod`，根据`Pod`的状态得知。
2. 通知已结束的`Pod`的`DeploymentController`其结束状态。
3. 使用`CoreScheduler`分配各个`Pod`的时间片。
4. 根据时间片以及`Pod`所需的内存，更新`Pod`的状态。
5. 根据`Pod`返回的负载信息，更新节点的负载情况。

## TODO List

- [ ] 数据读取接口的设计
  - [x] DeploymentController设计。用于在每个时钟周期向集群提交数据
  - [ ] 使用阿里巴巴2017年数据，实现符合该数据集的DeploymentController
- [ ] 模拟器设计
  - [x] Pod模拟器设计
    - [x] 初步可运行框架设计
    - [x] 整合Kubernetes API的Pod
  - [x] Node模拟器设计
    - [x] 模拟节点运行逻辑设计
    - [x] 整合Kubernetes API的Node
  - [x] 事件通知器接口`SharedInformerFactory`与`PodInformer`实现
  - [ ] client-go接口实现
    - [x] Pod增删改查与Watch，以及Bind接口
    - [x] Node增删改查与Watch
    - [ ] 根据模拟需求待定
  - [x] 调度Pod
    - [x] 将新的Pod放入调度队列
    - [x] 将Pod与Node绑定
  - [ ] 调度器插件编写
- [ ] 监控系统的设计
  - [x] 监控数据设计
  - [x] 节点监控数据采集
  - [ ] 监控数据统计
  
## 尚未计划实现的调度器功能

Kubernetes调度器拥有许多的Predicates和Priority插件，能够查看集群几乎所有资源的状态。由于本系统专注于测试调度器对集群资源
利用率的提升，以及性能的提升方面，因此与这些无关的资源接口将不被实现，以下列出部分调度器插件目前使用的接口

- 存储相关接口
  - StorageClass
  - PersistentVolume
  - PersistentVolumeClaim
- 部署相关
  - Deployment
  - ReplicaSet
  - Services
  - ReplicationController