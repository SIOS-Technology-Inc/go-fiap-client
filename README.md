# go-fiap-client

## go-fiap-clientとは
go-fiap-clientは、IEEE1888プロトコルをGo言語で扱うためのクライアント実装です。現在はFETCHのみサポートしており、その他のクライアント メソッドは非対応です。サーバ実装はサポートされません。

IEEE1888 (UGCCNet, FIAPとも) は大量の時系列データをやりとりするための規格であり、BEMSやスマートグリッドでの利用を期待して開発されています。

### pkg
ライブラリとしてのFIAPクライアント実装です。
```golang
package main

import (
	"fmt"
	"github.com/SIOS-Technology-Inc/go-fiap-client/pkg/fiap"
)

func main() {
	var (
		cli fiap.Fetcher = &fiap.FetchClient{ConnectionURL: "http://example.jp/FIAPEndpoint"}
		id  string       = "sios/example/Temperature"
	)
	_, points, _, _ := cli.FetchLatest(nil, nil, id)
	fmt.Print(id + ":[")
	for i, value := range points[id] {
		if i > 0 {
			fmt.Print(", ")
		}
		fmt.Printf(`time:%d,value:"%s" `, value.Time.Unix(), value.Value)
	}
	fmt.Println("]")
}
```

```bash
go run main.go
# sios/example/Temperature:[{time:1722438000,value:"29.01"}]
```

### cmd
コマンドラインとしてのFIAPクライアント実装です。
```bash
go run main.go fetch --select max "http://example.jp/FIAPEndpoint" "sios/example/Temperature"
# {"points":{"sios/example/Temperature":[{"time":"2024-08-01T00:00:00+09:00","value":"29.01"}]}}
```

## how to use library
go.modに次の文を追加します。
```go
require github.com/SIOS-Technology-Inc/go-fiap-client v0.2.2
```
パッケージや関数の詳細は[ドキュメント](https://pkg.go.dev/github.com/SIOS-Technology-Inc/go-fiap-client@v0.2.2/pkg/fiap)を参照してください。

## how to use command line
次のコマンドでローカルにコマンドをインストールします。
```bash
go install github.com/SIOS-Technology-Inc/go-fiap-client v0.2.2
```
### コマンドラインの記法
#### Fetch
```bash
go-fiap-client fetch [flags] URL (POINT_ID | POINTSET_ID)
```
このコマンドは、指定した`URL`と`POINT_ID`または`POINTSET_ID`を用いて、FIAPサーバからデータをFetchし、JSON形式で出力します。
- `-h`, `--help`<br>オプション情報を含むコマンドのヘルプを表示します。
- `-d`, `--debug`<br>デバッグ用出力が表示されるようにします。
- `-o FILEPATH`, `--output FILEPATH`<br>Fetchの結果を指定したファイルに出力します。
- `-s TYPE`, `--select TYPE`<br>Fetchされるデータを変更するオプションです。`TYPE`は`max`, `min`, `none`を記述します。指定しない場合のデフォルトは`max`です。<br>FIAPのkeyクラスの`select`の、それぞれ`maximum`、`minimun`、指定なしに対応します。
- `--from DATETIME`
- `--until DATETIME`<br>指定した日付期間で取得するデータを絞り込みます。`DATETIME`には指定する日付日時をRFC3339形式の文字列で指定します。<br>FIAPのkeyクラスの`gteq`、`lteq`にそれぞれ対応します。
#### その他
```bash
go-fiap-client [flags]
```
- `-h`, `--help`<br>サブコマンドやオプションの情報を含むヘルプを表示します。
- `-v`, `--version`<br>go-fiap-clientのライブラリバージョンを表示します。


### how to develop
開発環境はVisual Studio Code Dev Containersで構築されています。
[Docker](https://www.docker.com/)環境も前提となります。

詳細な利用方法は[Dev Containers 公式ドキュメント](https://code.visualstudio.com/docs/devcontainers/containers)を確認してください。
