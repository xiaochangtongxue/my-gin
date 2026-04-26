# API Contract

本文档约定后端和前端之间的基础接口契约。Swagger/OpenAPI 以 `api/` 目录生成文件为准。

## Base URL

- 本地后端：`http://localhost:8080`
- API 前缀：`/api/v1`
- 认证前缀：`/api/v1/auth`

## Response Envelope

所有业务接口使用统一 JSON 包装：

```json
{
  "code": 0,
  "message": "success",
  "data": {},
  "timestamp": 1700000000,
  "request_id": "optional-request-id"
}
```

分页数据统一放在 `data` 内：

```json
{
  "items": [],
  "total": 0,
  "page": 1,
  "page_size": 20,
  "total_page": 0
}
```

## Auth APIs

- `GET /api/v1/auth/captcha`
- `POST /api/v1/auth/register`
- `POST /api/v1/auth/login`
- `POST /api/v1/auth/refresh`
- `POST /api/v1/auth/logout`

前端登录成功后保存：

- `access_token`：放入后续请求头 `Authorization: Bearer <token>`
- `refresh_token`：仅用于刷新 token 和登出
- `expires_in`：Access Token 剩余秒数

## Error Codes

| code | HTTP | 含义 |
|---:|---:|---|
| `0` | `200` | 成功 |
| `10001` | `400` | 参数错误 |
| `10002` | `401` | 未授权 |
| `10003` | `403` | 禁止访问 |
| `10004` | `404` | 资源不存在 |
| `10005` | `500` | 服务器错误 |
| `10006` | `500` | 数据库错误 |
| `10007` | `500` | Redis 错误 |
| `10008` | `401` | Token 过期 |
| `10009` | `401` | Token 无效 |
| `10303` | `403` | 权限不足 |
| `10400` | `400` | 验证码错误或已过期 |
| `10401` | `401` | 用户名或密码错误 |
| `10402` | `202` | 需要验证码 |
| `10423` | `403` | 账号已锁定 |
| `10429` | `429` | 请求过于频繁 |

## Axios Policy

前端请求拦截器：

- 自动带上 `Authorization: Bearer <access_token>`
- 自动带上 `X-Request-ID`，便于排查日志
- 登录、注册、验证码、刷新 token 不强制要求 Access Token

前端响应拦截器：

- `code === 0` 返回 `data`
- `10008` 或 `10009` 尝试使用 refresh token 调用 `/api/v1/auth/refresh`
- refresh 失败时清理本地 token 并跳转登录页
- `10402` 展示验证码输入
- `10429` 展示限流提示，不自动重试高频写请求

## CORS

开发环境允许跨域；生产环境应把 `middleware.cors.allow_origins` 改为明确的前端域名，避免继续使用 `*`。
