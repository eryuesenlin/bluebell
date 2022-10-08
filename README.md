# Bluebell

## Go Web进阶

### Gin框架源码解析



#### Gin框架路由详解

#### Gin框架中间件详解





### Go连接MySQL/Redis

#### database/sql及sqlx使用

##### database/sql

包sql提供了一个围绕SQL（或类似SQL）数据库的通用接口。

sql包必须与数据库驱动一起使用。

##### 使用MySQL驱动

Open打开一个dirverName指定的数据库，dataSourceName指定数据源

Open函数可能只是验证其参数格式是否正确，实际上并不创建与数据库的连接。如果要检查数据源的名称是否真实有效，应该调用Ping方法。

返回的DB对象可以安全地被多个goroutine并发使用，并且维护其自己的空闲连接池。因此，Open函数应该仅被调用一次，很少需要关闭这个DB对象。

接下来，我们定义一个全局变量`db`，用来保存数据库连接对象。将上面的示例代码拆分出一个独立的`initDB`函数，只需要在程序启动时调用一次该函数完成全局变量db的初始化，其他函数中就可以直接使用全局变量`db`了。

其中`sql.DB`是表示连接的数据库对象（结构体实例），它保存了连接数据库相关的所有信息。它内部维护着一个具有零到多个底层连接的连接池，它可以安全地被多个goroutine同时使用。

##### database/sqlx

#### go-redis库使用





### 搭建Go Web开发脚手架

#### zap日志库使用

#### viper配置管理库使用

#### 优雅关机与平滑关机

#### CLD代码分层

### 仿Reddit论坛项目



#### 分布式ID生成

不用数据库自增ID作为用户ID

原因：

别人注册一下，就可以知道你数据库里有多少用户。

当用户量很大时，对数据库分库分表时，不同库的用户ID可能会重复。

不用uid的原因：uid丢失了使用数组作为用户ID的特点，按照时间排序是递增的，uid是无序的。

特点：

全局唯一性：不能出现有重复的ID标识

递增性：确保生成ID对于用户或业务是递增的

高可用性：确保任何时候都能生成正确的ID

高性能性：在高并发的环境下依然表现良好

典型场景：电商促销，微博短时间大量转发、评论的消息。将数据插入数据库之前，我们需要给订单或者消息先分配一个唯一的ID，再存进数据库。对这个ID的要求是 希望带一些时间信息，到后端系统对消息进行分库分表的时候，也能以时间的顺序对数据进行排序。

snowflake[算法](https://github.com/bwmarrin/snowflake)



#### JWT认证

##### cookie-sesion模式

`Cookie-Session`模式实现用户认证:

1. 客户端使用用户名、密码进行认证
2. 服务端校验用户名、密码正确后生成并存储session，将sessionid通过cookie返回给客户端
3. 客户端访问需要认证的接口时在cookie中携带sessionid
4. 服务端通过sessionid查找session并进行鉴权，返回给客户端需要的数据

基于session的方式存在多种问题：

1. 服务端需要存储session，并且session需要经常快速查找，通常存储在内存或者数据库中，同时在线用户较多时需要占用大量的服务器资源。
2. 当许需要扩展时，创建session的服务器可能不是验证session的服务器，所以还需要将所有的session单独存储并共享。
3. 由于客户端使用cookie存储sessionid，在跨域场景下需要进行兼容性处理，同时这种方式也难以防范CSRF攻击。

Token认证模式：

1. 客户端使用用户名、密码进行认证
2. 服务端验证用户名、密码正确后生成Token返回给客户端
3. 客户端保存Token，访问需要认证的接口时在URL参数或HTTP Hander中加入Token
4. 服务端通过解码Token进行鉴权，返回客户端需要的数据。

##### jwt是什么？

JWT全称JSON Web Token是一种跨域认证解决方案。JWT就是一种基于Token的轻量级认证模式，服务端认证通过后，会生成一个JSON对象，经过签名后得到一个Token（令牌）再发回给用户，用户后续请求只需要带上这个Token，服务端解密之后就能获取该用户的相关信息了。

JWT Token：

它是由`.`分割的三部分组成，这三个部分依次是：

- 头部（Header）
- 负载（Payloa）
- 签名（Signature）

头部和负载以JSON形式存在，三部分的内容都分别单独经过了Base64编码，以`.`拼接成一个JWT Token。

Sifnature 部分是是对前两个部分的签名，防止数据纂改。

##### jwt优缺点

JWT拥有基于Token的会话管理方式所拥有的一切优势，不依赖cookie，使得其可以防止CSRF攻击，也能在禁用cookie的浏览器环境中运行。

而JWT的最大优势是服务端不再需要存储session，使得服务端认证鉴权业务可以方便扩展，避免存储session所需要引入的Redis等组件，降低了系统架构复杂度。**但这也是JWT最大的劣势**：由于有效期存储在Token中，JWT Token一旦签发，就会在有效期内一直可用，无法在服务端废止，当用户进行登出操作，只能依赖客户端删除掉本地存储的JWT Token，如果需要禁用用户，单纯使用JWT就无法做到了。

##### 基于jwt实现认证实践

前⾯讲的 Token，都是 Access Token，也就是访问资源接⼝时所需要的 Token，还有另外⼀种 Token，Refresh Token，通常情况下，Refresh Token 的有效期会⽐较⻓，⽽ Access Token 的有效期 ⽐较短，当 Access Token 由于过期⽽失效时，使⽤ Refresh Token 就可以获取到新的 Access Token， 如果 Refresh Token 也失效了，⽤户就只能重新登录了。

 在 JWT 的实践中，引⼊ Refresh Token，将会话管理流程改进如下。 

- 客户端使⽤⽤户名密码进⾏认证
- 服务端⽣成有效时间较短的 Access Token（例如 10 分钟），和有效时间较⻓的 Refresh Token（例如 7 天）
- 客户端访问需要认证的接⼝时，携带 Access Token
- 如果 Access Token 没有过期，服务端鉴权后返回给客户端需要的数据
- 如果携带 Access Token 访问需要认证的接⼝时鉴权失败（例如返回 401 错误），则客户端使⽤ Refresh Token 向刷新接⼝申请新的 Access Token
- 如果 Refresh Token 没有过期，服务端向客户端下发新的 Access Token
- 客户端使⽤新的 Access Token 访问需要认证的接⼝

##### 生成jwt和解析jwt

使用`jwt-go`这个库





#### Makefile（mac或linux）

##### make

`make`是一个构建自动化工具，会在当前目录下寻找`Makefile`或`makefile`文件。如果存在相应的文件，它就会依据其中定义好的规则完成构建任务

##### Makefile介绍

可以把`Makefile`简单理解为它定义了一个项目文件的编译规则。借助`Makefile`我们在编译过程中不再需要每次手动输入编译的命令和编译的参数，可以极大简化项目编译过程。同时使用`Makefile`也可以在项目中确定具体的编译规则和流程，很多开源项目中都会定义`Makefile`文件。





#### 基于MySQL实现主业务

#### 基于Redis实现投票业务

#### 基于Docker搭建开发环境

#### 代码发布与项目部署

#### 多到写不下的实战经验和技巧

## viper --配置管理



### 什么是Viper?

Viper是适用于Go应用程序（包括`Twelve-Factor App`）的**完整配置解决方案**。它被设计用于在应用程序中工作，并且可以处理所有类型的配置需求和格式。它支持以下特性：

- 设置默认值
- 从`JSON`、`TOML`、`YAML`、`HCL`、`envfile`和`Java properties`格式的配置文件读取配置信息
- 实时监控和重新读取配置文件（可选）
- 从环境变量中读取
- 从远程配置系统（etcd或Consul）读取并监控配置变化
- 从命令行参数读取配置
- 从buffer读取配置
- 显式配置值



### 为什么选择Viper？

在构建现代应用程序时，你无需担心配置文件格式；你想要专注于构建出色的软件。Viper的出现就是为了在这方面帮助你的。

Viper能够为你执行下列操作：

1. 查找、加载和反序列化`JSON`、`TOML`、`YAML`、`HCL`、`INI`、`envfile`和`Java properties`格式的配置文件。
2. 提供一种机制为你的不同配置选项设置默认值。
3. 提供一种机制来通过命令行参数覆盖指定选项的值。
4. 提供别名系统，以便在不破坏现有代码的情况下轻松重命名参数。
5. 当用户提供了与默认值相同的命令行或配置文件时，可以很容易地分辨出它们之间的区别。

Viper会按照下面的优先级。每个项目的优先级都高于它下面的项目:

- 显示调用`Set`设置值
- 命令行参数（flag）
- 环境变量
- 配置文件
- key/value存储
- 默认值

**重要：** 目前Viper配置的键（Key）是大小写不敏感的。目前正在讨论是否将这一选项设为可选。



### 把值存入Viper



### 从Viper获取值



### 使用单个还是多个Viper实例？



## zap--日志库

### Choosing a Logger

Zap提供了两种类型的日志记录器—`Sugared Logger`和`Logger`。

在性能很好但不是很关键的上下文中，使用`SugaredLogger`。它比其他结构化日志记录包快4-10倍，并且支持结构化和printf风格的日志记录。

在每一微秒和每一次内存分配都很重要的上下文中，使用`Logger`。它甚至比`SugaredLogger`更快，内存分配次数也更少，但它只支持强类型的结构化日志记录。

在Logger和SugaredLogger之间的选择不需要是整个应用的决定：在两者之间的转换很简单，而且成本很低。

### Configuring Zap

#### Logger

- 通过调用`zap.NewProduction()`/`zap.NewDevelopment()`或者`zap.Example()`创建一个Logger。
- 上面的每一个函数都将创建一个logger。唯一的区别在于它将记录的信息不同。例如production logger默认记录调用函数信息、日期和时间等。
- 通过Logger调用Info/Error等。
- 默认情况下日志都会打印到应用程序的console界面。

#### SugaredLogger

- 大部分的实现基本都相同。
- 惟一的区别是，我们通过调用主logger的`. Sugar()`方法来获取一个`SugaredLogger`。
- 然后使用`SugaredLogger`以`printf`格式记录语句

### Extending Zap

- 使用`zap.New(…)`方法来手动传递所有配置

```go
func New(core zapcore.Core, options ...Option) *Logger
```

- `zapcore.Core`需要三个配置——`Encoder`，`WriteSyncer`，`LevelEnabler`。

**Encoder**:编码器(如何写入日志)。

**WriteSyncer**:指定日志将写到哪里去。

**LevelEnabler**:哪种级别的日志将被写入。



### Lumberjack--日志分割

Lumberjack旨在成为日志基础设施的一个部分。它不是一个全能的解决方案，而是在日志堆栈的底部的一个可插拔的组件，它只是控制日志被写入的文件。



## 优雅关机和平滑重启

### 优雅地关机

优雅关机就是服务端关机命令发出后不是立即关机，而是等待当前还在处理的请求全部处理完毕后再退出程序，是一种对客户端友好的关机方式。而执行`Ctrl+C`关闭服务端时，会强制结束进程导致正在访问的请求出现问题。

Go 1.8版本之后， http.Server 内置的 [Shutdown()](https://golang.org/pkg/net/http/#Server.Shutdown) 方法就支持优雅地关机。



## docker



## 传给前端数字id失真问题(go语言json技巧)

go中int64类型的值在有些场景下有可能超过前端js能够表示的最大的数，导致数字失真。

解决方法：在json序列化的时候转化成string类型。 

例： 

```go
type User struct{
	ID int64 `json:"id,string"`
}
```







## swagger生成接口文档

## 编写单元测试
