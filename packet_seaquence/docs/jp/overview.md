# 概要
このファイルは、packet_sequence ツールの概要（目的および使用方法）を提供します。  
tcpdump ファイルからパケットシーケンスを解析し、それをシーケンスチャート形式で視覚化したいユーザーを対象としています。

## 目的
packet_sequence ツールは、tcpdump ファイルをシーケンスチャートに変換するために設計されています。  
このツールは、tcpdump ファイル内のパケットシーケンスを解析しやすくし、データフローの視覚的な表現を提供します。  
つまり、`tcpdump ファイル(.pcap)` を `mermaid.js` のシーケンスチャートに変換します。

## 技術詳細
 - **実装言語**: Python x.x
 - **依存関係**:
   - パケット操作および解析のための `scapy`
   - シーケンスチャート生成のための `mermaid`

## 使用方法
```bash
packet_sequence <tcpdump_file.pcap> [options]
```