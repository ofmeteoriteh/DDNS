# DDNS

<!-- Badges / 徽章 -->
[![License](https://img.shields.io/badge/license-Apache--2.0-blue.svg)](LICENSE)
[![Go](https://img.shields.io/badge/Go-1.26-00ADD8?logo=go&logoColor=white)](https://go.dev/)
[![CI](https://github.com/ofmeteoriteh/ddns/actions/workflows/ci.yml/badge.svg)](https://github.com/ofmeteoriteh/ddns/actions/workflows/ci.yml)

> *A minimal, self-hosted DDNS client — Cloudflare only, IPv4 and IPv6, zero third-party dependencies*
>
> 极简自托管 DDNS 客户端 — 只对接 Cloudflare，只做 IPv4 / IPv6，零第三方依赖

---

## What is this? / 这是什么？

DDNS is a lightweight command-line tool that automatically updates Cloudflare DNS A / AAAA records with your machine's current public IP address. It is deliberately minimal — one provider, two record types, standard library only.

DDNS 是一个轻量级命令行工具，自动将本机公网出口 IP 更新到 Cloudflare 的 DNS A / AAAA 记录。刻意精简 — 单一 provider、两种记录类型、纯标准库。

---

## Highlights / 技术亮点

- **Zero dependencies** — built entirely on the Go standard library; no third-party SDK, no Cloudflare client library

  **零依赖** — 完全基于 Go 标准库构建；无第三方 SDK，无 Cloudflare 客户端库

- **Dual-stack support** — fetches both IPv4 and IPv6 public IPs from multiple single-stack sources, updates A and AAAA records independently

  **双栈支持** — 从多个纯单栈源获取 IPv4 和 IPv6 公网 IP，独立更新 A 和 AAAA 记录

- **Skip on no-change** — compares current IP against existing DNS record content; skips the API call if nothing changed

  **无变化跳过** — 将当前 IP 与现有 DNS 记录比对；IP 未变则跳过 API 调用

- **Cross-platform** — pure Go with no CGO; compiles to Linux / macOS / Windows on amd64 / arm64 / arm / 386

  **跨平台** — 纯 Go 无 CGO；可编译至 Linux / macOS / Windows 的 amd64 / arm64 / arm / 386

---

## How it works / 工作流程

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

## Quick start / 快速开始

### 1. Environment variables / 环境变量

```bash
export DDNS_NAME="home.example.com"
export CLOUDFLARE_API_KEY="<YOUR_API_TOKEN>"
export CLOUDFLARE_ZONE_ID="<YOUR_ZONE_ID>"
```

Or copy `.env.example` to `.env` and fill in the values.

或复制 `.env.example` 为 `.env` 并填入对应值。

### 2. Run / 运行

```bash
go run .
```

### 3. Binary / 二进制

Download the appropriate binary for your platform from [Releases](https://github.com/ofmeteoriteh/ddns/releases).

从 [Releases](https://github.com/ofmeteoriteh/ddns/releases) 下载对应平台的二进制文件。

---

## Project structure / 项目结构

```text
.
├── main.go              # Entry point / 入口，读取环境变量，编排流程
├── getip/
│   └── client.go        # Fetch public IP from multiple sources / 从多个源获取公网 IP
├── cloudflare/
│   ├── verify.go        # API token verification / API Token 验证
│   └── dns.go           # DNS record CRUD / DNS 记录的列出、创建、覆盖
├── .env.example         # Environment variable template / 环境变量模板
└── go.mod
```

---

## Maintenance / 维护说明

> This is a hobby project maintained in my spare time. Quality issues and
> PRs are welcome, but **response is not guaranteed**. Please factor this
> in before depending on it.
>
> 这是一个出于个人兴趣的开源项目。我欢迎高质量的 issue 和 PR，
> 但 **不承诺响应时间** — 可能处理得慢，也可能不处理。
> 若你需要稳定维护的方案，请评估后再依赖本项目。

---

## License & Disclaimer / 许可与免责

This project is licensed under the [Apache License 2.0](LICENSE).

This project is developed and maintained by individual contributors on a voluntary, non-commercial basis.

This software is provided **"as is"**, without warranty of any kind. The author(s) accept no liability for any damages arising from the use of this software. See the [LICENSE](LICENSE) file for the full terms, including the disclaimer of warranty and limitation of liability.

Any commercial entity using this software is solely responsible for their own compliance with applicable laws and regulations, including but not limited to the EU Cyber Resilience Act (CRA) and any other regional requirements.
