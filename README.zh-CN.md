[English](README.md) | 中文

# DDNS

[![License](https://img.shields.io/badge/license-Apache--2.0-lightgrey?style=flat-square)](LICENSE)
[![Go](https://img.shields.io/badge/Go--9ca3af?style=flat-square)](https://go.dev/)
[![Status](https://img.shields.io/badge/status-stable-blue?style=flat-square)]()

> 极简自托管 DDNS 客户端 — 只对接 Cloudflare，只做 IPv4 / IPv6。

---

## 这是什么？

DDNS 是一个轻量级命令行工具，自动将本机公网出口 IP 更新到 Cloudflare 的 DNS A / AAAA 记录。
刻意精简：单一 provider、两种记录类型、纯标准库。

---

## 技术亮点

- **精简依赖** — Cloudflare API 交互纯标准库；交互式配置由 [survey](https://github.com/AlecAivazis/survey) 驱动
- **双栈支持** — 从多个纯单栈源获取 IPv4 和 IPv6 公网 IP，独立更新 A 和 AAAA 记录
- **无变化跳过** — 将当前 IP 与现有 DNS 记录比对；IP 未变则跳过 API 调用
- **跨平台** — 纯 Go 无 CGO；可编译至 Linux / macOS / Windows 的 amd64 / arm64 / arm / 386

---

## 工作流程

```text
getip                cloudflare
  │                      │
  ├─ GetIPv4IP() ──────► DDNS("A", ipv4)
  │                        ├─ GET  /dns_records?name=...&type=A
  │                        ├─ IP changed?
  │                        │    ├─ no  → skip
  │                        │    ├─ yes → PUT /dns_records/{id}
  │                        │    └─ empty → POST /dns_records
  │
  └─ GetIPv6IP() ──────► DDNS("AAAA", ipv6)
                           └─ (same flow)
```

---

## 快速开始

### 下载

从 [Releases](https://github.com/ofmeteoriteh/DDNS/releases) 下载对应平台的二进制文件。

### 配置

```bash
ddns setup
```

交互式向导会提示输入 API Key、Zone、DNS 条目，并可选择生成 systemd service + timer 文件。
配置保存在 `config.json`。

### 运行

```bash
ddns
```

---

## 命令行参数

```text
ddns [flags]          运行 DDNS 更新
ddns setup            交互式配置向导

Flags:
  -v, --version       打印版本号
  -h, --help          打印帮助
  -c, --config PATH   配置文件路径（默认 config.json）
  --dry-run           模拟运行，不实际调 API
```

---

## 各平台使用方法

### Linux (systemd)

```bash
# 安装二进制
cp ddns /opt/ddns/ddns
chmod +x /opt/ddns/ddns

# 交互式配置（会生成 service + timer 文件）
cd /opt/ddns && ./ddns setup

# 安装 service + timer
cp ddns.service ddns.timer /etc/systemd/system/
systemctl daemon-reload
systemctl enable --now ddns.timer

# 手动触发一次
systemctl start ddns
```

### macOS (crontab)

```bash
# 安装二进制
cp ddns /usr/local/bin/ddns
chmod +x /usr/local/bin/ddns

# 交互式配置
cd /usr/local/bin && ddns setup

# 定时运行（crontab）
crontab -e
# 添加：*/5 * * * * cd /usr/local/bin && ./ddns
```

### Windows (Task Scheduler)

```powershell
# 放置二进制
# 将 ddns-windows-amd64.exe 放到 C:\ddns\ddns.exe

# 交互式配置
cd C:\ddns
.\ddns.exe setup

# 定时任务（Task Scheduler）
# 创建基本任务，触发器设为每 5 分钟，操作为启动 C:\ddns\ddns.exe
```

---

## 从源码构建

需要 Go 1.26+。

```bash
git clone https://github.com/ofmeteoriteh/DDNS.git
cd DDNS
go build -o ddns .
```

---

## 配置示例

`config.json` 由 `ddns setup` 自动生成，结构如下：

```json
{
  "keys": [
    { "name": "main", "key": "<YOUR_CLOUDFLARE_API_TOKEN>" }
  ],
  "zones": [
    { "name": "main-zone", "domain": "<YOUR_DOMAIN>", "zone_id": "<YOUR_ZONE_ID>" }
  ],
  "entries": [
    {
      "name": "<YOUR_RECORD_NAME>.<YOUR_DOMAIN>",
      "zone_id": "<YOUR_ZONE_ID>",
      "key": "<YOUR_CLOUDFLARE_API_TOKEN>",
      "types": ["A", "AAAA"],
      "proxied": false
    }
  ]
}
```

---

## 安全说明

DDNS 处理敏感数据，并跨越本地机器与 Cloudflare 之间的信任边界。

- **凭证** — Cloudflare API token 以明文存储在 `config.json` 中。应为此文件设置适当的文件系统权限。
- **网络端点** — 客户端通过 HTTPS 调用 Cloudflare API（`api.cloudflare.com`）和公网 IP 检测端点。
- **持久化状态** — `config.json` 保存凭证、zone ID 和记录名；它必须可写，并在多次运行之间持久存在。
- **信任边界** — API token 必须具有编辑目标 zone DNS 记录的权限。应使用最小权限范围，仅限所需的 zone 和记录。

---

## 项目结构

```text
.
├── main.go              # 入口，子命令路由
├── setup.go             # 交互式配置向导
├── config/
│   └── config.go        # 配置加载与保存
├── getip/
│   └── client.go        # 从多个源获取公网 IP
├── cloudflare/
│   ├── verify.go        # API Token 验证
│   └── dns.go           # DNS 记录的列出、创建、覆盖
└── go.mod
```

---

## 文档

| 文档 | 用途 |
| ----------------- | ---------------- |
| [`.github/PULL_REQUEST_TEMPLATE.md`](.github/PULL_REQUEST_TEMPLATE.md) | PR 模板 |
| [`LICENSE`](LICENSE) | Apache 2.0 许可 |

---

## 维护说明

> **注意：** 这是一个出于个人兴趣的开源项目，在业余时间维护。欢迎高质量的 issue 和 PR，但不承诺响应时间。在依赖本项目前，请将这一点纳入考量。

---

## 诚实局限

- 仅支持 Cloudflare，不支持其他 DNS 提供商。
- 仅支持 IPv4 / IPv6 的 A / AAAA 记录。不支持其他记录类型、泛域名管理，也不自动调整 TTL。
- IP 检测失败时没有 webhook、通知或多提供商降级方案。
- 需要可写的 `config.json` 持久化在磁盘；未针对无卷挂载的临时容器环境设计。

---

## 许可与免责

本项目采用 [Apache License 2.0](LICENSE) 许可。

本项目由个人贡献者在自愿、非商业基础上开发和维护。

本软件按 **"as is"** 提供，不含任何形式的担保。作者不对因使用本软件而造成的任何损害承担责任。详见 [LICENSE](LICENSE) 文件中的完整条款，包括免责声明和责任限制。

任何商业实体使用本软件时，须自行负责遵守适用的法律法规，包括但不限于欧盟《网络弹性法案》（CRA）及其他地区性要求。
