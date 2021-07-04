# 搬砖 Bot (brickbot)

一个用来润滑日常搬砖过程中团队协作的简单 bot。

目前支持以下代码托管方案：

- GitHub

目前支持以下 IM 软件：

- 企业微信

## License

```
// SPDX-License-Identifier: GPL-3.0-or-later
```

本软件的核心部分以 [GNU GPLv3 或更新版本](./LICENSE.GPLv3.md) 授权。

使用 GPL 协议是为了鼓励大家将改进回馈上游，这样所有用户都可从中受益；
不使用 AGPL 协议是为了防止万一有需求必须向软件的核心（不可插拔）部分加入不宜公开的信息，而对下游造成两难的境地。

## FAQ

### Why Go (and not `$YOUR_FAVORITE_LANGUAGE`)?

因为这个软件之后很可能会部署到七牛线上，也可能被其他七牛同学拿去用，所以选用 Go 对其他人来说门槛比较低。
