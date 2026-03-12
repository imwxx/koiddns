# 1.koiddns

- **What**: Dynamic DNS Client scripts for Ali DNS, just support IPv4, config was suitable for OpenWRT
- **Cron**: The application runs every specified number of minutes according to the executionCycleMinutes parameter setting.

# 2.How

## 2.1.config

```
koiddns 使用说明：
--help                  显示帮助信息
--generate-config FILE  生成示例配置文件，参数为文件路径
--config FILE           指定配置文件，参数为文件路径
--daemon                以守护进程方式运行
--pidfile FILE          守护进程 PID 文件路径（默认 /var/run/koiddns.pid）
--logfile FILE          守护进程日志文件路径（默认 /var/log/koiddns.log）

配置文件示例（/etc/config/koiddns）：
main:
  executionCycleMinutes: 5                                     # 执行周期，单位为分钟

providers:
  - name: "aliyun"
    accessKeyId: "your_aliyun_access_key_id"                   # 阿里云访问密钥 ID
    accessKeySecret: "your_aliyun_access_key_secret"           # 阿里云访问密钥 Secret
  - name: "tencent"
    secretId: "your_tencent_secret_id"                         # 腾讯云访问密钥 ID
    secretKey: "your_tencent_secret_key"                       # 腾讯云访问密钥 Secret

domains:
  - subDomain: "api"
    primaryDomain: "example.com"
	value: "" # 留空则程序会自动获取或创建
    provider: "aliyun"
    recordType: "A"
    recordId: ""  # 留空则程序会自动获取或创建
    line: "default"
    priority: "10"
  - subDomain: "api"
    primaryDomain: "anotherdomain.com"
	value: "" # 留空则程序会自动获取或创建
    provider: "tencent"
    recordType: "CNAME"
    recordId: ""
    line: "unicom"
    priority: "20"

参数说明：
- accessKeyId / secretId: 云服务提供商分配的访问密钥 ID, 用于身份验证, 必须项
- accessKeySecret / secretKey: 云服务提供商分配的访问密钥 Secret, 用于身份验证, 必须项
- subDomain: 需要更新的域名api.example.com中, api即为子域名, 必须项
- primaryDomain: 需要更新的域名api.example.com中, example.com即为主域名, 主域名要与 解析服务商的解析服务配置对齐，必须项
- value: 记录值 有配置则获取使用，否则会以获取的出网NAT地址为值, 非必须项
- recordType: DNS 记录类型，如 A、CNAME 等, 必须项
- recordId: 云服务提供商分配的记录 ID，用于唯一标识一条 DNS 记录, 非必须项。
- line: DNS 解析线路，如 default、telecom、unicom 等, 必须项。
- priority: DNS 记录的优先级，数值越小优先级越高（适用于 MX 等记录类型）, 必须项。

```

## 2.2.run

前台运行：
```
/usr/bin/koiddns --config /etc/config/koiddns
```

守护进程运行（可指定 PID/日志路径）：
```
/usr/bin/koiddns --config /etc/config/koiddns --daemon
/usr/bin/koiddns --config /etc/config/koiddns --daemon --pidfile /var/run/koiddns.pid --logfile /var/log/koiddns.log
```

## 2.3 OpenWRT / luci-app-koiddns

在 OpenWRT 上可安装配套 **luci-app-koiddns**，提供：

- **配置管理**：直接读写 `/etc/config/koiddns`（YAML），在「General Settings」页展示并编辑；保存前自动做格式校验，也可通过「Validate format」单独校验当前文件。
- **服务进程管理**：在「Service」页查看运行状态，并执行启动、停止、重启。
