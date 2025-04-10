# mermaid_output.xlsmについて
`mermaid_output.xlsm`は、Mermaid形式のシーケンス図をExcelで表示するためのマクロ有効なExcelファイルです。

情報設定シートのD3からD12にフィルタしたい情報を入力します。
## フィルタ情報
- D3: 対象ファイル名(captureファイル)
- D4: 出力ファイル名(参照からフォルダを選択すると、output.mdがデフォルトで選択されます)
- D5: 最大パケット数(全数取得する場合は空白)
- D6: 送信元IPアドレス
- D7: 送信先IPアドレス
- D8: フィルタプロトコル名
- D9: 送信元または送信先IPアドレス
- D10: フィルタ適用開始時刻
- D11: フィルタ適用終了時刻
- D12: デバックモード(1:デバックモード、0:通常モード)

## 注意事項
1. xlsmファイルとバイナリファイル(packet_seaquence.exe)のフォルダに配置位置に注意してください。
```
ファイル構成
│ packet_sequence.exe
├─execel_output
│  ├─exec_command.bas
│  ├─convert_mermaid.bas
│  └─mermaid_output.xlsm(本ファイル)
└config
   └─config.pkseq
```

2. 大量のパケットを取得する場合、Excelのメモリ制限により、正常に動作しない場合があります。
フィルタ条件を指定することを推奨します。

2025-4-10 ver0.0.1