# Security Policy / 安全策略

Reports are accepted in **English or Chinese**.

报告接受**英文或中文**。

## Supported versions / 支持的版本

Only the latest release is supported.

只支持最新发布版本。

## Reporting a vulnerability / 报告漏洞

**Do not open a public issue for security vulnerabilities.**

**不要为安全漏洞开公开 issue。**

Use GitHub's private vulnerability reporting:

改用 GitHub 的私密漏洞报告功能：

1. Open the **Security** tab of this repository. / 进入本仓库的 **Security** 标签页。
2. Click **Report a vulnerability**. / 点击 **Report a vulnerability**。

This keeps the report private until a fix is available.

这样在修复发布前，报告将保持私密。

> **Note:** This is a hobby project maintained in spare time. While security reports are taken seriously, **response time is not guaranteed**.

> **注意：** 这是一个利用业余时间维护的兴趣项目。安全报告会被认真对待，但**不承诺响应时间**。

## Scope / 范围

DDNS stores a Cloudflare API token and updates your DNS records over the network.
Security-sensitive areas include, but are not limited to:

DDNS 保存 Cloudflare API 令牌，并通过网络更新你的 DNS 记录。
安全敏感面包括但不限于：

- **Token storage** — the API token is kept in plaintext in `config.json` (file mode `0600`).

  令牌存储 — API 令牌以明文保存在 `config.json`（文件权限 `0600`）。
- **DNS updates** — record creation and updates against `api.cloudflare.com`; a compromise here redirects or takes down hostnames.

  DNS 更新 — 针对 `api.cloudflare.com` 的记录创建与更新；此处被攻破会劫持或下线主机名。
- **Network calls** — HTTPS requests to Cloudflare and public IP echo services.

  网络请求 — 到 Cloudflare 与公共 IP 回显服务的 HTTPS 请求。
- **Config file handling** — reading, writing and relocating `config.json`.

  配置文件处理 — `config.json` 的读取、写入与迁移。

## Out of scope / 不在范围内

- Vulnerabilities in Cloudflare or the IP echo services (report upstream).

  Cloudflare 或 IP 回显服务的漏洞（向上游报告）。
- Issues requiring physical access to the host. / 需要物理访问主机的问题。
- Running DDNS with an over-scoped token on a host you do not control.

  在你不受控的主机上用超出最小权限的令牌运行 DDNS。

## Disclaimer / 免责声明

This software is provided "as is", without warranty of any kind.
Any commercial entity using this software is solely responsible for its own compliance with applicable laws and regulations, including the EU Cyber Resilience Act (CRA).

本软件按"原样"提供，不附带任何形式的担保。
任何使用本软件的商业实体，须自行负责遵守适用的法律法规，包括欧盟《网络弹性法案》（CRA）。
