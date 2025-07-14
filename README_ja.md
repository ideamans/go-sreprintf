# go-sreprintf

[English README](README.md)

sprintf形式の文字列からパラメータを逆エンジニアリングし、別のテンプレート（翻訳など）に再適用するGoライブラリです。

## インストール

```bash
go get github.com/miyanaga/go-sreprintf
```

## 使い方

```go
package main

import (
    "fmt"
    "github.com/miyanaga/go-sreprintf"
)

func main() {
    template := "Hello %s, you have %d new messages"
    message := "Hello Alice, you have 5 new messages."
    translation := "こんにちは %sさん、あなたには%d件の新しいメッセージがあります"
    
    result, err := sreprintf.Sreprintf(template, message, translation)
    if err != nil {
        panic(err)
    }
    
    fmt.Println(result)
    // 出力: こんにちは Aliceさん、あなたには5件の新しいメッセージがあります
}
```

## 対応フォーマット指定子

以下のフォーマット指定子に対応しています：

- `%s` - 文字列
- `%d` - 整数
- `%f` - 浮動小数点数
- `%v` - 汎用フォーマット（Go特有）
- `%t` - bool型（true/false）
- `%%` - リテラルの%記号

### フォーマットオプション

以下のフォーマットオプションも正しく処理します：

- 幅指定: `%10s`, `%5d`
- 左寄せ: `%-10s`
- 符号表示: `%+d`
- 代替形式: `%#x`, `%#o`
- 精度指定: `%.2f`, `%.5s`
- ゼロパディング: `%05d`
- 複合指定: `%+10.2f`, `%-20.10s`

パース時はフォーマットオプションを無視して値を抽出し、再適用時は翻訳テンプレートのフォーマットオプションをそのまま使用します。

## 動作原理

1. **テンプレート解析**: テンプレート文字列を解析してすべてのフォーマット指定子を識別
2. **値の抽出**: 正規表現パターンを使用してフォーマット済みメッセージから値を抽出
3. **型変換**: 抽出した文字列値を適切な型（int、float、bool）に変換
4. **再適用**: 抽出した値を`fmt.Sprintf`を使用して翻訳テンプレートに適用

## エラーハンドリング

極力エラーを出さない方針に従っています：

- メッセージがテンプレート構造と一致しない場合のみエラーを返します
- 型の不一致は優雅に処理されます - 変換に失敗した場合は文字列として保持
- プレースホルダの過不足は空文字列でパディングまたは余分な値を無視することで処理

## テスト

以下のコマンドでテストを実行：

```bash
go test -v ./...
```

## ライセンス

MIT License