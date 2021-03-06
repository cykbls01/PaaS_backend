# PaaS功能点整理

（2019.11.07）

## 项目目标

- 完成一个简洁、易用、较稳定、较安全的、功能较完善、工作量足够的PaaS平台
  - 简洁易用：降低使用PaaS平台的难度，简化创建容器的过程
  - 稳定安全：权限设置比较完善，服务不会因为小错误就挂掉
  - 功能完善、工作量足够：当做毕设要有足够的工作量（比如数据库只有一个表就不太行）

## 项目参考

- [新云平台线上测试环境](vlab.beihangsoft.cn)
- 测试账号（密码均为12345-abcde）
  - 管理员：admin
  - 教师：99999
  - 学生：16219999
- 界面、交互逻辑可以参考（至少不能比这个丑）
- 以后PaaS可能会合并入云平台(.Net)
  - 考虑API的兼容性，目前设想是使用一个类似于超级管理员的账号进行鉴权，然后调用Go服务提供的API
  - 可能要在API中显式加入请求的ID信息
    - 比如云平台可以向PaaS服务发送请求`api/podlist/limituser?id=16219999`、`api/podlist/limituser/16219999`
    - POST创建请求的时候也会加上ID
- 仅供参考，并不是“指导”

## 使用方面

很多用户可能之前没有接触过docker，我们的平台要让没接触过docker的用户使用没有障碍。大部分用户是只会用一个数据库、部署一个网站，而不会使用其他更高阶的服务。

### 单独docker、共享docker

- 现在是每个人按照一定的配置独享一个docker
- 更多情况下，不需要docker，只需要一个服务
  - 比如MySQL，只需要一个临时的数据库，同一个docker中的MySQL可以创建多个database和user来提供给多个人数据库
  - 比如tomcat，同一个docker中的tomcat可以提供给多个人部署不同端口号应用
    - 可能需要tomcat自动部署的接口
- 可选共享docker
  - 共享的MySQL，直接返回一个数据库访问地址、数据库名、数据库用户的账号密码
  - 共享的Tomcat，上传一个打包好的网站应用，直接返回一个访问地址

### 创建容器

- 自动填充名字可以删除了（因为肯定会重复）

### 简化创建流程

- 选择端口、配置环境变量可以默认省略
  - 原因是这些镜像提供的服务端口都比较固定，一般也不会用到环境变量。可以把这些设置默认隐藏，放到“高级设置”之类的里面。

- 选择配置可以跳跃式选择
  - 可以参考，阿里云/腾讯云等。实际的PaaS可能会使用高中低等配置，或者内存增强型、计算增强型等等，可以去各个云网站查看各自的文字描述。尽量不要使用微内核（因为云平台那边使用的1 CPU就是1核，这边1000CPU=1核会造成误解，没使用过docker的用户也不会了解微内核的概念）。  
    - 可以提供“自定义设置”来自己指定微内核数、存储大小和内存大小。
    - 可能要重新设计一下配额设置问题。
      - 参考云平台的逻辑：没有固定的配额（比如限定4U8G200G之类的），而是每创建一个虚拟机都需要向管理员申请
      - PaaS平台也可以加入申请-审批的流程

### 表格部分

- 镜像名规范化（MySQL 5.7）
- 端口的映射显示更加清楚一点
- 上边的搜索框可以更灵活一点
- 状态改为中文

### 容器日志

- 日志显示不太清楚，当日志过长时比较不友好
- 可以考虑

## 管理员的管理功能很不完善

### 用户管理

- 用户的增删查（包括用户的信息查看、密码重置）
- 无法得知用户的配额大小，重新设定时也无法得知
- 用户列表增加“已创建容器数”
- 用户表加入学号/工号、姓名等展示
- 批量管理用户
  - 比如，批量导入一门课程的所有选课学生的学号，创建账号，并统一分配额度
    - 带来的问题：额度收回的批量操作

### 日志管理

- 此处的日志不是网站服务后台的错误日志，而是网站操作日志
  - 比如修改密码、创建容器、登录等敏感操作，操作的内容、操作者、登录IP等记录到数据库中
- 管理员增加查看日志功能

### 容器管理

- 管理员点进用户管理-管理容器后，并不能看到当前看的是哪个用户，容易产生误操作
- 按类统计每大类docker各创建了多少，可以与监控页面合并

### 节点监控

- 未显示单位
- 未显示存储数

## 自动配置、自动部署

- k8s新增节点的自动部署，增加可用资源，扩展集群
  - 克隆的k8s虚拟机自动加入集群中
    - Linux：编写.sh脚本，开机自动执行配置过程（可以加入一个类似于flag的文件/变量来标记已经执行过自动配置过程了，不需要再次执行），参考[Github上的一键配置脚本](https://github.com/luckman666/kkitdeploy_install)。
    - Windows：暂时未想到
- 应用的自动部署
  - Tomcat、Nginx、.Net、PHP等镜像，往往需要上传文件，但是无论是通过ftp还是scp，使用命令行传递文件都比较麻烦（WinSCP等客户端需要额外配置端口，也需要额外下载客户端）
  - 考虑一种自动化部署Web应用的程序/脚本，用户可以在创建容器后，直接通过前端上传部署包（或者发布文件夹等）的方式来部署自己的应用

## 已部署的应用列表、资源列表、发布内容

- 在首页展示一下在PaaS中已经部署的所有应用
  - 列表？云图？其他可视化方法？

## 镜像管理

### 管理员可以从前端添加镜像

- 直接pull?

### 增加镜像种类（包括版本种类）

- Redis
- .Net
  - 未找到Linux的.Net Framework的可用镜像
  - 考虑使用win下的k8s，以支持.net系列应用
- Tomcat多种版本
- MangoDB
- PHP
- ...

### SaaS

- 参考[UZER.me](Uzer.me)

## 代码规范

- 代码的可读性、关键函数的功能注释
  - 这个项目是要实际上线的，以后可能会有人继续做
- 及时更新代码，防止因版本不统一而出现的错误
- 每完成一个功能点commit一次，便于以后的debug
- commit message要显示出更改内容
  - feature_增加xxx功能
  - bugfix_修复xxxbug
    - description可以写一下bug的原因，如果这个bug的原因值得写的话
  - optimazation_优化xxx过程
  - ...
- 错误日志的记录
  - 错误堆栈、请求内容
  - 与INFO级别的最好分开
