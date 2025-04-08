#!/bin/bash
# パケット生成ツールの実行スクリプト

# 出力ディレクトリの作成
mkdir -p output

# パケット生成コマンドの実行例
echo "===== テスト用パケット生成開始 ====="

# HTTP パケット作成
echo "HTTPパケットを作成中..."
python3 create_packets.py --protocol http --count 20 --output output/http_packets.pcap --verbose

# ARP パケット作成
echo "ARPパケットを作成中..."
python3 create_packets.py --protocol arp --count 10 --output output/arp_packets.pcap --verbose

# ICMP パケット作成
echo "ICMPパケットを作成中..."
python3 create_packets.py --protocol icmp --count 10 --output output/icmp_packets.pcap --verbose

# DNS パケット作成
echo "DNSパケットを作成中..."
python3 create_packets.py --protocol dns --count 10 --output output/dns_packets.pcap --verbose

# TCP パケット作成
echo "TCPパケットを作成中..."
python3 create_packets.py --protocol tcp --count 20 --output output/tcp_packets.pcap --verbose

# UDP パケット作成
echo "UDPパケットを作成中..."
python3 create_packets.py --protocol udp --count 10 --output output/udp_packets.pcap --verbose

# X.25 パケット作成
echo "X.25パケットを作成中..."
python3 create_packets.py --protocol x25 --count 15 --output output/x25_packets.pcap --verbose

# 全プロトコル一括生成
echo "すべてのプロトコルを一括生成中..."
python3 create_packets.py --protocol all --count 5 --output output/all_protocols.pcap --verbose

echo "===== テスト用パケット生成完了 ====="