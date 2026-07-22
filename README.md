# CPA 插件仓库

[CLIProxyAPI (CPA)](https://github.com/router-for-me/CLIProxyAPI) 插件集合。当前提供 **WorkBuddy / CodeBuddy** OAuth Provider。

## 插件

| ID | 说明 | 源码 |
|---|---|---|
| `workbuddy` | Tencent CodeBuddy OAuth、动态模型、executor、签到、积分面板、可选积分调度 | [workbuddy/](workbuddy/) |

## 多架构 Release（回应 [issue #1](https://github.com/Sliverkiss/cpa-plugin/issues/1)）

每个版本 GitHub Release 提供 CPA 插件商店标准产物：

```text
workbuddy_<version>_linux_amd64.zip     # zip 根目录: workbuddy.so
workbuddy_<version>_linux_arm64.zip
workbuddy_<version>_darwin_amd64.zip    # workbuddy.dylib
workbuddy_<version>_darwin_arm64.zip
workbuddy_<version>_windows_amd64.zip   # workbuddy.dll
workbuddy_<version>_windows_arm64.zip
workbuddy_<version>_freebsd_amd64.zip
checksums.txt
```

命名规则与官方一致：`ArchiveName(id, version, goos, goarch) = {id}_{version}_{goos}_{goarch}.zip`  
（见 CLIProxyAPI `internal/pluginstore`）。

CI：push / tag `v*` / PR 触发 `.github/workflows/build.yml`。

## 安装（linux/amd64 示例）

```bash
# 从 Release 下载
unzip workbuddy_0.4.1_linux_amd64.zip
# 扁平 plugins 目录（常见 docker 挂载）
cp workbuddy.so /path/to/cliproxyapi/plugins/workbuddy.so
# 或平台子目录布局
# mkdir -p plugins/linux/amd64 && cp workbuddy.so plugins/linux/amd64/
```

`config.yaml`：

```yaml
plugins:
  enabled: true
  dir: "plugins"
  configs:
    workbuddy:
      enabled: true
      # checkin_auto: true
      # scheduler_mode: off   # or credits
```

重启 CPA 后：

1. 管理端 **OAuth / 登录** 选择 WorkBuddy（或面板导入凭证 JSON）
2. 凭证默认写入 CPA `auths/workbuddy-<uid>.json`
3. `/v1/models` 会出现 workbuddy 动态模型（如 `deepseek-v4-flash`、`glm-5.2`、`kimi-k2.7`、`hy3` 等，以账号权限为准）
4. 面板：`/v0/resource/plugins/workbuddy/panel`

更细说明见 [workbuddy/README.md](workbuddy/README.md)。

## 远程更新（插件源 registry）

仓库根目录提供 CPA 插件商店可消费的源：

```text
https://raw.githubusercontent.com/Sliverkiss/cpa-plugin/main/registry.json
```

在 CPA / CPAMP **插件商店 → 来源** 中添加该 URL，即可搜索安装 / 更新 `workbuddy`（`install` 类型为默认 `github-release`，从本仓库 Releases 拉 zip）。

校验：

```bash
python3 scripts/validate-registry.py registry.json
```

## 本地构建

```bash
cd workbuddy
make test && make vet && make build VERSION=$(cat VERSION)
# dist/workbuddy.so （当前主机架构）
```

## 许可证

各插件目录内 LICENSE（WorkBuddy：MIT）。
