# ユーザーガイド
このドキュメントは、packet_sequence ツールの `ユーザーガイド` を提供し、`インストール手順`、`使用例`、および `トラブルシューティングのヒント` を含んでいます。  
tcpdump ファイルからパケットシーケンスを解析し、シーケンスチャート形式で視覚化したいユーザーを対象としています。

## インストール
packet_sequence ツールをインストールするには、以下の手順に従ってください：
1. リポジトリをクローンします:
   ```bash
   git clone <URL>
   ```
2. プロジェクトディレクトリに移動します:
   ```bash
   cd packet_sequence
   cp -p conv_packet_sequence /bin
   ```

## 使用方法
packet_sequence ツールを使用するには、以下のコマンドを実行します:
```bash
packet_sequence <tcpdump_file.pcap> [options]
```

### オプション
- `-h`, `--help`: ヘルプメッセージを表示して終了します。
- `-o`, `--output`: 出力ファイル名を指定します。デフォルトは `output.mmd` です。
- `-f`, `--format`: 出力フォーマットを指定します。選択肢は `mermaid` または `text` です。デフォルトは `mermaid` です。
- `-v`, `--verbose`: 詳細な出力を有効にします。
- `-version`: ツールのバージョンを表示します。
- `-i`, `--input`: 入力ファイル名を指定します。デフォルトは `input.pcap` です。

## トラブルシューティング