# API

## 目录结构

```
api
├── product_name           // 产品名称
│   └── app_name           // 应用名称
│       └── v1             // 版本号
│           └── rpc.proto  // 具体接口
├── buf.gen.yaml           // Buf 代码生成插件配置
└── buf.yaml               // Buf Deps、Lint、Break 等配置
```
