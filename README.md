# auto-github-hosts
自动获取github的地址，并更新本地HOSTS文件

### 架构、技术栈说明：

架构：本程序采用传统的分层架构，包括界面层、应用层、领域层、基础设施层四个层次。其中界面层为状态栏图标，应用层为程序的主逻辑，领域层为对hosts配置文件及ip地址的处理逻辑，基础设施层为对操作系统的接口及I/O操作等。
技术栈：本程序采用Golang编写，使用了第三方库Viper实现对YAML配置文件的读取，使用了第三方库Hostsfile实现对hosts文件的读写。同时，程序通过调用系统API获取IP地址，并使用了第三方库trayhost来实现跨平台状态栏图标的绘制。


### Golang 工程目录：
```

    ├── cmd                     # 包含程序入口main.go
    │   └── main.go 
    ├── config                  # 包含配置文件auto-github-hosts.toml和读取配置文件所需的代码
    │   ├── auto-github-hosts.toml
    │   └── config.go
    ├── domain                  # 领域层，包括主要的hosts文件与IP地址的处理逻辑
    │   ├── hosts.go
    │   ├── ip.go
    │   └── service.go
    ├── infrastructure          # 基础设施层，包括对系统API的调用与对hosts文件的读写
    │   ├── system.go
    │   └── hostsfile.go
    └── ui                      # 界面层，包括对状态栏图标的绘制和程序主逻辑的响应
        ├── icon.go
        ├── menu.go
        ├── status.go
        └── service.go
```

需要创建的源代码文件，名称及说明：
* ./main.go：程序入口
* ./config/config.go：读取配置文件及初始化配置
* ./domain/hosts.go：处理hosts文件的逻辑，包括读、写、更新hosts文件等
* ./domain/ip.go：处理IP地址的逻辑，包括从url获取IP地址等
* ./domain/service.go：领域层服务接口定义
* ./infrastructure/system.go：操作系统API调用
* ./infrastructure/hostsfile.go：对hosts文件的操作（读、写、更新）实现
* ./ui/icon.go：状态栏图标实现
* ./ui/menu.go：状态栏菜单实现
* ./ui/status.go：状态栏实现
* ./ui/service.go：界面层服务接口定义及响应实现