# はじめに
このドキュメントは、PacketSequenceの使い方を説明します。PacketSequenceは、パケットのシーケンスを管理するためのツールです。

## 概要
PacketSequenceは、pcapファイルを解析して、パケットのシーケンスをMermaid.jsを使ったシーケンス図として視覚化するツールです。ネットワーク通信の解析やトラブルシューティングに役立ちます。


# 事前準備
このツールを使用するには、以下のものが必要です：
- tshark (Wireshark CLI版)
  - Windowsの場合は`config/config.pkseq`ファイルにtsharkの実行パスを設定してください
  - 例: `tsharkDir=C:\Program Files\Wireshark\tshark.exe`
- Mermaidをサポートするマークダウンビューア（出力ファイルの表示用）

# インストール方法
```bash
git clone <リポジトリURL>
cd packet_seaquence
```

# 使い方
```bash
./packet_sequence -file <pcapファイル.pcap> [オプション]
```

## オプション
- `-file <パス>`: 解析するpcapファイルのパス（必須）
- `-out <パス>`: 出力するマークダウンファイルのパス（デフォルト: output.md）
- `-max <数値>`: 処理する最大パケット数（0=すべて）
- `-debug`: デバッグモードの有効化
- `-source <IP>`: 送信元IPアドレスでフィルタリング
- `-destination <IP>`: 宛先IPアドレスでフィルタリング
- `-protocol <名前>`: プロトコル名でフィルタリング
- `-IP <IP>`: 送信元または宛先IPアドレスでフィルタリング
- `-startTime <時刻>`: フィルタリング開始時刻（例: 2023-10-01 12:00:00）
- `-endTime <時刻>`: フィルタリング終了時刻（例: 2023-10-01 12:00:00）
- `-version`: バージョン情報の表示

## 使用例
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

# 詳細ドキュメント
より詳細な情報については、以下のドキュメントを参照してください：

- [概要](docs/jp/overview.md)
- [インストールと使用方法](docs/jp/User_guide.md)
- [設計](docs/jp/design.md)
- [API リファレンス](docs/jp/api.md)
- [システム要件](docs/jp/requirements.md)

# ライセンス
MIT License

# 著者
Junnosuke Horiuchi
