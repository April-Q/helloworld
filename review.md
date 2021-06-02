# 试用期工作总结

## 自我介绍

我从2020年12月10日入职公司，到现在基本六个月了，这六个月感谢领导和同事的帮助和指导，使我能在工作上快速入手。

## 学习培训

前半个月主要是熟悉工作环境，熟悉团队工作内容，完成新员工线上与线下的培训，包括财务、IT、法务、人力、行政等5门基础制度课程并参与考试。

## 工作情况介绍

### 自动诊断系统项目

1. go profiler开发

   * 支持Profile,Heap and Goroutine的剖析类型
   * 内部实现使用 go tool pprof命令
   * 可配置 TLS，目前支持配置secret，通过指定secret中的ca.crt与token ，来获取 HTTPS 类型的剖析文件
   * 定时垃圾回收，且可配置
   * webhook监测go profiler参数
   * 默认安装可连接apiserver的serviceaccount，配置rbac，对集群apiserver的剖析实现开箱即用

1. 一些项目文档编写

   | Path | Description |
   |-|-|
   | docs/design/garbage-collection.md | 垃圾回收文档，介绍垃圾回收的设计，实现方式与用户配置方法。 |
   | docs/design/go-profiler.md | 介绍go 性能剖析器的设计实现，支持类型，实现方式与用户配置方法。 |
   | docs/website/concepts/go-profiler.md | 介绍go 性能剖析器的设计实现，支持类型，实现方式与用户配置方法。 |
   | docs/website/tutorials/how-to-profiler-apiserver.md | 介绍go 性能剖析器的设计实现，支持类型，实现方式与用户配置方法。 |
   | docs/website/tutorials/how-to-use-go-profiler-in-your-application.md | 介绍 go 性能剖析器的设计实现，支持类型，实现方式与用户配置方法。 |

1. e2e测试
   * 为自动诊断系统添加e2e框架
   * e2e test suite
   * a framework which automatically creates uniquely named namespace for each test，

1. prometheus metrics暴露

### 建行lb项目

1. snmp exporter开发
    支持通过SNMP等设备性能和状态采集方式，SNMP支持V2版本；(W1)
    SNMP适配，需要监听vs级别的信息
    并发（活跃连接数）吞吐（mbps）新建（cps），ssl的tps（连接的建立）
    确定一个prefix，例如1.3.6.1.4.1.163163，现在已经被用过的在这里可以看到：
    确定如何表达各VS，可以参考net-snmp自带的TCP的统计项，ipv4的tcp连接表达方式就是 .192.168.1.1.80 这样表达的
    确定各统计项的id，比方说1 = 并发连接数，2 = 吞吐。以此类推
    和envoy联调，用snmpwalk查询统计项
    导出mib用于交付
    查询结果支持数字格式（用起来最简单的net-snmp的shell脚本扩展不支持数字只支持字符串）
    能够动态增加/删除oid（因为listener是动态的）

1. lb client开发

目标组后端管理 server
支持服务器上线、下线功能：online
负载均衡集群管理；lb
负载均衡节点管理； lbn
虚拟服务器管理；vs
转发规则管理 fr
服务器证书管理 cc
支持HTTPS协议
lb集群配置备份、恢复

## 心得体会

积极沟通
多思考
从用户角度考虑

## 意见与建议
