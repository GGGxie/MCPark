# nginx 参考配置

线上实际生效的配置在 `/etc/nginx/`，这里是同步的参考副本，方便在新机器上重建。

| 文件 | 线上位置 |
|---|---|
| `00-edge-limits.conf` | `/etc/nginx/conf.d/00-edge-limits.conf` |
| `edge-security-headers.conf` | `/etc/nginx/snippets/edge-security-headers.conf` |
| 站点配置（本目录中的另一个 `.conf`） | `/etc/nginx/conf.d/<名称>.conf` |

要点：
- 站点经 Cloudflare Tunnel 进来，nginx 看到的来源都是 `127.0.0.1`，所以限流按 `CF-Connecting-IP`（没有时退回 `$remote_addr`）
- 通用 API 10 r/s（burst 30），`/api/auth/` 30 r/m（burst 5），超限返回 429；请求体上限 256KB
- 安全头放在 snippet 里：带自己 `add_header` 的 location 需要再 `include` 一次
- 未设置 CSP：登录使用 Google Identity Services，需要单独验证后再加
- 若站点端口（8083 / 8085）可被绕过 Cloudflare 直连，则 `CF-Connecting-IP` 可被伪造，应在云防火墙上只放行 Cloudflare 或回环

修改后先 `nginx -t` 再 `systemctl reload nginx`。
