# 要件
このドキュメントは、パケットシーケンスツールの`依存関係`と`システム要件`を含む`要件リスト`を提供します。

## システム要件

### OS
- Linux
- Windows

MacOSは基本的にサポートされていますが、十分なテストが行われていません。

### 依存関係
- `tshark` (Wireshark CLI版): パケット解析に使用
- ブラウザ: Mermaid形式のシーケンス図を表示

### 設定
- Windowsの場合は`config/config.pkseq`ファイルにtsharkの実行パスを設定する必要があります
  ```
  tsharkDir=C:\Program Files\Wireshark\tshark.exe
  ```

## 開発要件

### 開発言語
- Go言語

### ビルド要件
- Go 1.24.2以上
