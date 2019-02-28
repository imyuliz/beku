ROADMAP
-
This document defines the roadmap for Beku development.

## v0.1.0: (2018-12-11)

1. 支持Apply方法，存在即更新，不存在即部署；
2. 修复PVCMount() 带来的异常；
3. 支持对Deployment对Pod优先级的设置,保证进一步保障应用的稳定性；
4. 支持资源对象PriorityClass,保证高优先级、重量级的应用优先被调度和部署成功，提升应用的稳定性；
5. 增加Deployment,StatefulSet,DaemonSet 对PreStop的支持，此功能支持在容器容器删除前，做一些清理工作，支持Shell/HTTP；
6. 增加Deployment,StatefulSet,DaemonSet 对PostStart的支持，此功能支持在容器启动前做一些初始化操作，支持Shell/HTTP；
7. apiServer证书支持base64编码之后注册；


## v0.1.1: (2019-01-22)

1. 支持ServiceAccount
2. 支持ClusterRole
3. 支持ClusterRoleBinding
4. 支持RoleBinding


## v0.1.2: (2019-02-18)

1. 支持SA(serviceAccount)应用身份认证
2. 支持ClusterRole访问权限控制
3. 支持User,Group,SA三种方式到ClusterRole的绑定---ClusterRoleBinding


## v0.1.3: (2019-02-18)

1. 支持Deployment的NodeAffinity
2. 支持StatefulSet的NodeAffinity


## v0.1.3 (2019-02-28)

1. 支持Deployment的污染容忍Toleration
2. 支持Node资源对象