# beku
[![GoDoc](https://godoc.org/github.com/imroc/beku?status.svg)](https://godoc.org/github.com/yulibaozi/beku)
[![Go Report Card](https://goreportcard.com/badge/github.com/yulibaozi/beku)](https://goreportcard.com/badge/github.com/yulibaozi/beku)

Golang人性化Kubernetes资源对象创建库。极简，无额外心智负担。

### 安装

```
go get -u github.com/yulibaozi/beku
```

### 特性

- 极简的JSON和YAML输入/输出
- 自动判断资源对象的必要字段
- 人性化的资源对象关联发布
- 严谨的QOS等级设置
- 准确的字段自动填充
- 写意的链式调用

### 文档

- [中文](https://github.com/yulibaozi/beku/blob/master/doc/README-cn.md)
- [更多示例](https://github.com/yulibaozi/beku/blob/master/test/example_test.go)
- [Youtube:Deploy your application on Kubernetes with 3 LoC using Beku](xclearx)
- [Tencent Video:Beku--3行代码发布你的应用到Kubernetes](http://v.qq.com/x/page/d0783vtazs9.html)

### 概要

beku的诞生来源于对现实工作上的不满, 资源对象配置繁琐, 本身大量的字段已经让人难以忍受, 何况层级又带来了额外的心智负担, 还有多个资源对象重复劳动带来的低效能工作, 这三者的纠缠催生了beku, 这实在令人欢欣鼓舞。

beku的使用场景是对接Kubernetes的客户端(eg:client-go), 当然, 也可以使用beku生成json或者yaml供命令行创建。另外,beku库中存在Kubernetes一些源代码,这对于beku的实现提供了极大的帮助，这非常值得尊敬。

### beku的使用习惯

1. beku的使用习惯是以NewXXX()开始链式调用, 最终以调用Finish()表示调用结束, 从而得到完整的kubernetes资源对象配置信息。
2. 所有的填写都以SetXXX()开头, 所有的获取都以GetXXX()开头。
3. 在使用beku时, 尽量不使用强转的方式来满足函数所需要的变量类型, 这会引发未知错误。
4. 在使用beku时, 如果有函数的参数不知道填什么, 实现函数的注释中有相关阐述。
5. 在beku的应用场景中, Pod中的第一个container往往有至高地位, 拥有优先设置的权利, 随着序列的变大越显得平凡, 比如:第一次设置环境的时候, 只会为第一个container设置, 而不会为第二个设置, 当你第二次调用设置环境变量方法时, 才会设置第二container的环境变量, 以此类推。
5. 如果某结构体存在**union**字符时, 那么说明会同时创建两个Kubernetes资源对象, 例如:Deployment和Service的联合,PersistentVolume和PersistentVolumeClaim的联合

### 示例

如何快速创建一个Service(svc)?


正如你所见:

```go
func howToNewSvc() {
	svc, err := beku.NewSvc().SetNamespaceAndName("roc", "mysql-svc").
		SetSelector(map[string]string{"app": "mysql"}).SetServiceType(beku.ServiceTypeNodePort).
		SetPort(beku.ServicePort{Port: 3306, TargetPort: 3306}).Finish()
	if err != nil {
		panic(err)
	}
	yamlbyts, err := beku.ToYAML(svc)
	jsonbyts, err := beku.ToJSON(svc)
	if err != nil {
		panic(err)
	}
}
```

ToYAML

```yaml
apiVersion: v1
kind: Service
metadata:
  creationTimestamp: null
  name: mysql-svc
  namespace: roc
spec:
  ports:
  - port: 3306
    protocol: TCP
    targetPort: 3306
  selector:
    app: mysql
  type: NodePort
status:
  loadBalancer: {}
```

ToJSON

```json
{
    "kind":"Service",
    "apiVersion":"v1",
    "metadata":
    {
        "name":"mysql-svc",
        "namespace":"roc",
        "creationTimestamp":null
    },
    "spec":
    {
        "ports":
        [
            {
                "protocol":"TCP",
                "port":3306,
                "targetPort":3306
            }
        ],
        "selector":
        {
            "app":"mysql"
        },
        "type":"NodePort"
    },
    "status":
    {
        "loadBalancer":{}
    }
}
```
更多示例如: [Example.md](https://github.com/yulibaozi/beku/blob/master/doc/example.md)

### beku现阶段支持的资源对象列表

资源对象 | 简称|  版本
---|---|---
namespace   | ns| core/v1
service   | svc| core/v1
deployment | - | apps/v1
statefulset | sts | apps/v1
secret | - | core/v1
persistentVolumeClaim | pvc | core/v1
persistentVolume | pv | core/v1
daemonSet | ds | apps/v1
configMap | cm | core/v1

### beku的实现策略

1. **当某个资源对象存在多个版本时**, 只会实现其中一个版本, 版本的选择的首要条件是按照稳定性来选择的, 而不是功能多样性, 这可能导致实现的版本稍延后于最新流行的版本, 但, 没关系, 稳定可用才是第一前提, 另一方面, 在选择版本的优先级版本方面如下(优先级越高越在前面)。
    * core/v1
    * apps/v1
    ...

2. **当某资源对象缺少稳定级版本时**, 在实现kubernetes resouce object时, 不会实现Alpha, Beta级别的资源对象, 如果k8s某资源对象不存在stable版本, beku会选择放弃实现, 原因是其他不稳定版本或者较稳定版本存在更改的风险较大, 稳定性会降低。

3. 关于API版本的参考资料如下:
[Kubernetes API 概述](http://kubernetes.kansea.com/docs/api/)

### beku的构想

**在过去**的工作经历中,我发现kubernetes资源对象的配置甚似繁琐, 除去配置字段本身复杂外, 配置需要的额外心智负担也存在, 简单归纳如下:
 * 配置字段那些是必要, 那些是可需要, 这是个问题;
 * Kubernete资源对象层级较深, 配置字段该放在json/yaml的哪个位置, 这是个问题;
 * yaml配置时, 需要考虑缩进, 容易出错;
 * 实现多个资源对象, 上述的问题会N次重现, 头发掉一地谁买单, 这是个问题。

**现阶段**, 在实现上还有许多不足, 例如:某种资源对象的的字段有多种填写策略, 为了简便, 我们实现了其中常用的填写方式而隐藏了复杂的填写策略, 这可能会导致kubernetes资源对象本身就支持的多种填写策略被阉割, 如果您觉得是常用的, 必不可少的, 您可以提交pr或者issue来共同探讨。另一方面, 还有些功能还趋待完善, 我们可以共同成长。

**在未来**,期望能够消灭掉上述的三点不足和遵照以下目标前进
- [ ] 调用方提供必要的字段，不需要考虑那些字段需要还是不需要;
- [x] 调用方提供无层级的配置字段,免去这些额外负担;
- [x] 调用方提供必要的字段后，生成完整的json/yaml配置信息。
- [x] 调用方提供单一字段后,自动补充相关联的其他字段。
- [x] 当提供的必要字段不完整时,返回相关字段缺少，引导调用方逐步填写完整。
