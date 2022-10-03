# makefile

sacloudプロダクトで共通利用するMakefile

- `go/`: Go言語向け

## Usage

利用するプロジェクト側で以下のように利用します。

#### リモートリポジトリの追加(初回のみ)

```bash
git remote add makefile https://github.com/sacloud/makefile.git
```

#### 追加(初回のみ)

```bash
git subtree add --prefix=includes --squash makefile v0.0.7
```

利用する側のプロジェクトではMakefileを以下のように記述します。

```makefile
# 必要に応じて変数定義
AUTHOR         ?= The sacloud/example Authors
COPYRIGHT_YEAR ?= 2022
BIN            ?= example
DEFAULT_GOALS  ?= fmt set-license go-licenses-check goimports lint test build

# 必要なファイルをインクルード
include includes/go/common.mk
include includes/go/simple.mk

# ゴールを追加
default: $(DEFAULT_GOALS)
tools: dev-tools # toolsゴールはsacloudプロダクト向け日次CIを行うプロジェクトでは必須
```

#### 更新

```bash
git subtree pull --prefix=includes --squash makefile v0.0.7
```

## License

`sacloud/makefile` Copyright (C) 2022 The sacloud/makefile Authors.

This project is published under [Apache 2.0 License](LICENSE).

