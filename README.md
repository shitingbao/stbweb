# stbweb
>自己的web框架，封装各个模块，api拦截，异步端口监听，多个httpserver同时服务，后台日志收集服务,websocekt为基础的消息功能，nsq消息队列，grpc外部服务，redis地理位置处理，ants和cron结合的任务机制等

## 1.api拦截实现
>功能核心，采用类多态的性质，引入工作元素的概念（以url和对应header进行划分），每一个工作元素将对应不同的业务类，实现功能划分。权限上将api分为内部调用和外部调用，来约束调用过程
将每个对应功能都分配一个业务类，对应业务类实现get，post方法，在init的时候，将所有业务类注册到对应全局注册中心，等待，将拦截到的请求，与全局注册中的类匹配成功后，指向该类的方法内部执行。
在进行拦截时，获取到所有url请求，进行url处理，获取到标识，对应到相应的工作元素。这里将保存request和ResponseWriter，工作元素信息，用户信息等。然后，获取控制器的标识与全局注册中心匹配，成功后将进入api调用，否之，则无效请求。

eg：@/modules/common/api_get_ecample.go  

## 2.用户管理模块
>用户操作通过redis进行状态保存。注册将用户信息保存至数据库中，期间，将获取部分用户信息，结合加密算法生成唯一salt,使用md5结合salt生成password密文,采用二进制存储。登录时，使用相同加密方式进行验证，即使代码加密方式泄密，第三方也无法获取唯一salt，保证密码的安全性。登录后，将用户信息注册入redis用户组，同时保存用户状态，反馈api请求token和权限信息。  

eg：@/modules/common/register.go和login  

## 3.日志收集模块
>使用了第三方logrus包，将信息输出分成多个收集，定向不同的文件句柄中。logrus中将定义一个每日log文件，记录基本的信息。而系统中的日记分两种，一个panic的异常信息，另一个为普通日志信息，这里调用了系统的dll方式，不同操作系统使用不用的编译文件，将panic的错误信息重定向至指定文件内，防止控制台因为行数的约束未能记录获取信息。

eg：@core/global中logrul的启动  

## 4.聊天室及消息传递模块
>聊天室以及消息功能，是以基于websocket为基础，使用改进了聊天室的框架搭建的。每个连接使用Sec-WebSocket-Protocol标记用户，也可以在发送的信息内标记用户，达到一对一和多对多的双向信息传输过程。  

eg：@/lib/ws和/load/websocket 

## 5.任务处理以及过程结果记录
>内部调用了cron和ants两个包结合使用，达到了定时功能和线程池化的处理，同时，在处理过程中，使用redis标记任务，防止同类型同一用户多次重复提交，并执行后记录至数据库中。使用过程中，任务为定义的一个单独的对象，内部包含任务基本信息（定时信息及用户信息等）以及待处理逻辑，触发ants定时系统后，将该任务提交至ants协程池中实际处理。

eg：@/lib/task

## 6.缓存及数据保持
>使用redis记录实时信息及保值频繁的重复信息，生成验证码和状态标识。使用了reidsGeo的位置处理，反馈各个成员之间的位置信息和筛选功能。

eg：@/lib/rediser

## 7.grpc为基础的外部服务
>引入grpc机制，在proto中写入需要外部调用的服务，编译生成对应go服务文件。在load中新启端口建立tcp监听。  

eg: @/lib/external_service

## 8.nsq为基础的消息队列
>封装了原始的nsq方法，快速构建生产者或者消费，达到消息通讯的目的  

eg: @/lib/snsq---@/modules/common/nsq_send---@/modules/common/nsq_customer

## 项目架构
### 目录介绍
    builds 构建目录
    core    核心功能目录
    lib     功能列表目录
    loader  启动目录
    modules/common  功能实现和展示目录
#### --builds\common  
    main.go 主函数入口  
    assets静态资源存放  
    log logrus每日日志记录（暂剔除，待定）  
    config 配置文件  
    err.txt panic定向文件  
    log.txt普通日志收集文件  
#### --core  
    核心功能模块，定义了主要的功能内容和全局引用，以及各个引用类型，或者引用功能列表实现项目逻辑（内部任务等）  
#### --lib  
    功能列表，单独分开独立的功能工具  
    现有功能：  
    1.聊天室架构（基于websocket实现）  
    2.信息加密过程  
    3.提取图片文字（基于baiduAPI实现）  
    4.task任务模块（ants+cron）
    5.外服服务模块定制(grpc)
    6.相关文件处理（excel、csv等）
    7.nsq消息队列的构建
    8.redis地理位置操作
    9.redis信息保存和用户状态
#### --loader  
    启动服务设定，开启项目基本依赖（数据库，redis连接，webocekt监听以及日志记录功能等），开启服务监听，拦截api请求等  
#### --modules  
    实际业务实现模块  

## sql及nosql  
    后天数据库使用了mysql,缓存使用redis。
