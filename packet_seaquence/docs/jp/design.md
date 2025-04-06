# 設計
このドキュメントは、パケットシーケンスツールの設計概要（アーキテクチャ、コンポーネント、データフロー）を提供します。  
これは、ツールの内部動作を理解し、拡張または変更したい開発者を対象としています。

## アーキテクチャ
パケットシーケンスツールは、pcapファイルをシーケンスチャートに変換するように設計されています。

### コンポーネント
```mermaid
graph TD
    A[reader - pcapファイル読み込み] --> B[infogetter - パケット情報取得]
    B --> C[protocol - プロトコル解析]
    C --> D[writer - シーケンス図生成]
```

### パッケージ構成
- **cmd**: メインプログラムエントリーポイント
- **reader**: pcapファイルの読み込みを担当
- **infogetter**: tsharkを使ったパケット情報の取得
- **protocol**: 各プロトコル固有の解析ロジック
- **models**: データモデル定義
- **writer**: シーケンス図生成と出力
- **checker**: プロトコルサポート状況の確認
- **config**: 設定ファイル関連

### データフロー
1. **pcapファイルの読み込み**  
   `reader`パッケージはpcapファイルの存在確認を行います。

2. **パケット情報の抽出**  
   `infogetter`パッケージはtsharkを使用して、pcapファイルからパケット情報を抽出します。
   送信元、送信先、プロトコルなどの基本情報を取得します。

3. **プロトコル固有の解析**  
   `protocol`パッケージ内の各プロトコルアナライザーが、パケットの詳細情報を解析します。

4. **Mermaid図の生成**  
   `writer`パッケージがMermaid形式でシーケンス図を生成します。

5. **ファイル出力**  
   生成されたMermaidコンテンツがMarkdownファイルに出力されます。

```mermaid
sequenceDiagram
    participant M as main
    participant R as reader.PCAPReader
    participant I as infogetter.TsharkInfoGetter
    participant C as checker.ProtocolChecker
    participant P as protocol.Analyzer
    participant W as writer.MermaidWriter
    
    M->>R: 1. pcapファイルを読み込む
    M->>I: 2. パケット情報を取得
    I-->>M: パケットリスト
    M->>C: 3. プロトコルサポート確認
    M->>I: 4. 詳細情報を取得
    I->>P: プロトコル別の解析
    P-->>I: 解析結果
    I-->>M: 詳細情報付きパケット
    M->>W: 5. シーケンス図を生成
    W-->>M: 出力結果
```
