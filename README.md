### beku
A golang kubernetes deploy library for humans.


### About beku

beku 是kubernetes的一个资源对象工具，它以提供简单的，方便的，稳定的服务为宗旨。

### 实现资源对象的策略如下

1. **当某个资源对象存在多个版本时**,只会实现其中一个版本,版本的选择的首要条件是按照稳定性来选择的，而不是功能多样性，这可能导致实现的版本稍延后于最新流行的版本,但是,没关系,稳定可用才是第一前提,另一方面，在选择版本的优先级版本方面如下（优先级越高越在前面)。
    * core/v1
    * apps/v1
    ...

2. **当某资源对象缺少稳定级版本时**,在实现kubernetes resouce object时,不会实现Alpha,Beta级别的资源对象，如果k8s某资源对象不存在stable版本的,会选择放弃实现,原因是其他不稳定版本或者较稳定版本存在更改的风险较大,稳定性会降低。

3. 关于API版本的参考资料如下:
[Kubernetes API 概述](http://kubernetes.kansea.com/docs/api/)

### 现阶段实现的资源对象列表
资源对象 | 简称|  版本
---|---|---
service   | svc| core/v1
deployment | - | apps/v1
statefulset | sts | apps/v1
secret | - | core/v1
persistentVolumeClaim | pvc | core/v1
persistentVolume | pv | core/v1
daemonSet | ds | apps/v1
configMap | cm | core/v1

### 关于项目的构想

**在过去**的工作经历中,我发现kubernetes资源对象的配置甚似繁琐，除去配置字段本身复杂外，配置需要的额外心智负担也存在,简单归纳如下：
 * 配置字段那些是必要，那些是可需要,我不知道；
 * 配置字段该放在json/yaml的哪个位置，我不知道，因为kubernetes的资源对象配置层级足够深；
 * yaml配置时，需要考虑缩进,容易出错。

**现阶段**,在实现上还有许多不足,例如:某种资源对象的的字段有多种填写策略，为了简便，我们实现了其中常用的填写方式而隐藏了复杂的填写策略，这可能会导致kubernetes资源对象本身就支持的多种填写策略被阉割，如果您觉得是常用的，必不可少的，您可以提交pr或者issue来共同探讨。另一方面,还有些功能还趋待完善，我们可以共同成长。

**在未来**,期望能够消灭掉上述的三点不足和遵照以下目标前进
 * 调用方提供必要的字段，不需要考虑那些字段需要还是不需要;
 * 调用方提供无层级的配置字段,免去这些额外负担;
 * 调用方提供必要的json字段后，生成完整的json/yaml配置信息。
 * 调用方提供单一字段后,自动补充相关联的其他字段。
 * 当提供的必要字段不完整时,返回相关字段缺少，引导调用方逐步填写完整。


### 关于beku的使用习惯

1. beku的使用习惯是以NewXXX()开始链式调用,最终以调用Finish()表示调用结束,从而得到完整的kubernetes资源对象配置信息。
2. 所有的填写都以SetXXX()开头,所有的获取都以GetXXX()开头。
3. 在使用beku时,尽量不使用强转的方式来满足函数所需要的变量类型,这会引发未知错误。
4. 在使用beku时,如果有函数的参数不知道填什么,实现函数的注释中有相关阐述。
5. 在beku的应用场景中,Pod中的第一个container往往有至高地位,拥有优先设置的权利,而第二，第三个容器越来越显得平凡,比如:第一次设置环境的时候,只会为第一个container设置,而不会为第二个设置,当你第二次调用设置环境变量方法时，才会设置第二container的环境变量，以此类推
5. 如果某结构体存在**union**字符时,那么说明会同时创建两个Kubernetes资源对象,例如:Deployment和Service的联合,PersistentVolume和PersistentVolumeClaim的联合

### 新特性

* 支持设置QOS等级
* 支持自动填充Pod的Label和Selector