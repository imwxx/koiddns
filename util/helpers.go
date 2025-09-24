package util

import "fmt"

func ShowHelp() {
	fmt.Print(`koiddns 使用说明：
--help                  显示帮助信息
--generate-config FILE  生成示例配置文件，参数为文件路径
--config FILE           指定配置文件，参数为文件路径

配置文件示例（/etc/config/koiddns）：

main:
  executionCycleMinutes: 5                                # 执行周期，单位为分钟

providers:
  - name: "aliyun"
    accessKeyId: "your_aliyun_access_key_id"              # 阿里云访问密钥 ID
    accessKeySecret: "your_aliyun_access_key_secret"      # 阿里云访问密钥 Secret
  - name: "tencent"
    secretId: "your_tencent_secret_id"                    # 腾讯云访问密钥 ID
    secretKey: "your_tencent_secret_key"                  # 腾讯云访问密钥 Secret

domains:
  - subDomain: "api"
    primaryDomain: "example.com"
    value: ""          # 留空则程序会自动获取或创建
    provider: "aliyun"
    recordType: "A"
    recordId: ""       # 留空则程序会自动获取或创建
    line: "default"
    priority: "10"
  - subDomain: "api"
    primaryDomain: "anotherdomain.com"
    value: ""          # 留空则程序会自动获取或创建
    provider: "tencent"
    recordType: "CNAME"
    recordId: ""
    line: "unicom"
    priority: "20"

参数说明：
- accessKeyId / secretId: 云服务提供商分配的访问密钥 ID，用于身份验证，必须项
- accessKeySecret / secretKey: 云服务提供商分配的访问密钥 Secret，用于身份验证，必须项
- subDomain: 需要更新的域名 api.example.com 中，api 即为子域名，必须项
- primaryDomain: 需要更新的域名 api.example.com 中，example.com 即为主域名，主域名要与解析服务商的解析服务配置对齐，必须项
- value: 记录值，有配置则使用配置值，否则会以获取的出网 NAT 地址为值，非必须项
- recordType: DNS 记录类型，如 A、CNAME 等，必须项
- recordId: 云服务提供商分配的记录 ID，用于唯一标识一条 DNS 记录，非必须项
- line: DNS 解析线路，如 default、telecom、unicom、mobile、oversea、edu 等，必须项
- priority: DNS 记录的优先级，数值越小优先级越高（适用于 MX 等记录类型），必须项
`)
}
