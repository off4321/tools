# ユーザーガイド
このドキュメントは、pcapファイルからパケットシーケンスを分析し、シーケンスチャート形式で視覚化したいユーザー向けに、`インストール手順`、`使用例`、`トラブルシューティングのヒント`を含む`パケットシーケンスツールのユーザーガイド`を提供します。

## 事前準備
1. tsharkがインストールされていることを確認します
   - Wiresharkに付属しているCLIツールです
   - Windowsの場合は、`config/config.pkseq`ファイルにtsharkのパスを設定してください

## インストール
パケットシーケンストールをインストールするには、次の手順に従ってください：
1. リポジトリをクローンします：
   ```bash
   git clone <URL>
   ```
2. プロジェクトをビルドします：
   ```bash
   cd packet_sequence
   go build -o packet_sequence src/cmd/main.go
   ```

## 使用方法
パケットシーケンスツールを使用するには、次のコマンドを実行します：
```bash
./packet_sequence -file <pcapファイル.pcap> [オプション]
```

### オプション
- `-file <path>`: 解析するpcapファイルのパス（必須）
- `-out <path>`: 出力するマークダウンファイルのパス（デフォルト: output.md）
- `-max <n>`: 処理する最大パケット数（0=すべて）
- `-debug`: デバッグモードの有効化（詳細なログが表示されます）
- `-source <ip>`: 送信元IPアドレスでフィルタリング
- `-destination <ip>`: 宛先IPアドレスでフィルタリング
- `-protocol <name>`: プロトコル名でフィルタリング
- `-IP <ip>`: 送信元または宛先IPアドレスでフィルタリング
- `-startTime <time>`: フィルタリング開始時刻（例: 2023-10-01 12:00:00）
- `-endTime <time>`: フィルタリング終了時刻（例: 2023-10-01 12:00:00）
- `-version`: バージョン情報の表示

### 使用例
1. 基本的な使用方法：
   ```bash
   ./packet_sequence -file capture.pcap -out sequence.md
   ```

2. 特定のIPアドレスとプロトコルでフィルタリング：
   ```bash
   ./packet_sequence -file capture.pcap -IP 192.168.1.100 -protocol HTTP
   ```

3. 時間範囲でフィルタリング：
   ```bash
   ./packet_sequence -file capture.pcap -startTime "2023-10-01 12:00:00" -endTime "2023-10-01 12:05:00"
   ```

## 出力ファイル
出力されるMarkdownファイルには、Mermaid形式のシーケンス図が含まれています。
このファイルはGitHubやその他のMermaidレンダラーをサポートするMarkdownビューアで表示できます。

## トラブルシューティング
1. tsharkが見つからない場合：
   - Windowsの場合、`config/config.pkseq`ファイルに正しいパスが設定されているか確認してください
   - Linuxの場合、tsharkがインストールされているか確認してください: `which tshark`

2. パケット情報が取得できない場合：
   - `-debug`オプションを使用して詳細情報を確認してください
   - pcapファイルが正しいフォーマットか確認してください

3. 出力ファイルにシーケンス図が表示されない場合：
   - Mermaidをサポートするマークダウンビューアを使用しているか確認してください
   - 処理するパケットが存在するか確認してください
