#!/usr/bin/env python3
"""
テスト用パケット生成ツール

使用方法:
  python create_packets.py --protocol http --count 10 --output http_packets.pcap
  python create_packets.py --protocol all --output all_protocols.pcap
"""
import argparse
import os
import sys
from datetime import datetime

def main():
    """メイン関数"""
    parser = argparse.ArgumentParser(description='テスト用パケット生成ツール')
    parser.add_argument('--protocol', type=str, default='all',
                       choices=['http', 'arp', 'icmp', 'dns', 'tcp', 'udp', 'x25', 'all'],
                       help='生成するプロトコル')
    parser.add_argument('--count', type=int, default=10,
                       help='生成するパケット数 (プロトコルごと)')
    parser.add_argument('--output', type=str, default='output/packets.pcap',
                       help='出力ファイル名 (デフォルト: output/packets.pcap)')
    parser.add_argument('--verbose', '-v', action='store_true',
                       help='詳細な出力を表示')
    
    args = parser.parse_args()
    
    # 出力ディレクトリの確認
    output_dir = os.path.dirname(args.output)
    if output_dir and not os.path.exists(output_dir):
        os.makedirs(output_dir)
    
    # タイムスタンプをファイル名に追加（オプション）
    timestamp = datetime.now().strftime("%Y%m%d%H%M%S")
    
    # 各プロトコルのパケット生成を試行（個別にエラーハンドリング）
    if args.protocol == 'http' or args.protocol == 'all':
        try:
            from protocols.http_packets import create_http_packets
            http_output = args.output.replace('.pcap', f'_http_{timestamp}.pcap') if args.protocol == 'all' else args.output
            create_http_packets(http_output, args.count, args.verbose)
        except ImportError as e:
            print(f"HTTPパケット生成モジュールのインポートに失敗しました: {str(e)}")
    
    if args.protocol == 'arp' or args.protocol == 'all':
        try:
            from protocols.arp_packets import create_arp_packets
            arp_output = args.output.replace('.pcap', f'_arp_{timestamp}.pcap') if args.protocol == 'all' else args.output
            create_arp_packets(arp_output, args.count, args.verbose)
        except ImportError as e:
            print(f"ARPパケット生成モジュールのインポートに失敗しました: {str(e)}")
    
    if args.protocol == 'icmp' or args.protocol == 'all':
        try:
            from protocols.icmp_packets import create_icmp_packets
            icmp_output = args.output.replace('.pcap', f'_icmp_{timestamp}.pcap') if args.protocol == 'all' else args.output
            create_icmp_packets(icmp_output, args.count, args.verbose)
        except ImportError as e:
            print(f"ICMPパケット生成モジュールのインポートに失敗しました: {str(e)}")
    
    if args.protocol == 'dns' or args.protocol == 'all':
        try:
            from protocols.dns_packets import create_dns_packets
            dns_output = args.output.replace('.pcap', f'_dns_{timestamp}.pcap') if args.protocol == 'all' else args.output
            create_dns_packets(dns_output, args.count, args.verbose)
        except ImportError as e:
            print(f"DNSパケット生成モジュールのインポートに失敗しました: {str(e)}")
    
    if args.protocol == 'tcp' or args.protocol == 'all':
        try:
            from protocols.tcp_packets import create_tcp_packets
            tcp_output = args.output.replace('.pcap', f'_tcp_{timestamp}.pcap') if args.protocol == 'all' else args.output
            create_tcp_packets(tcp_output, args.count, args.verbose)
        except ImportError as e:
            print(f"TCPパケット生成モジュールのインポートに失敗しました: {str(e)}")
    
    if args.protocol == 'udp' or args.protocol == 'all':
        try:
            from protocols.udp_packets import create_udp_packets
            udp_output = args.output.replace('.pcap', f'_udp_{timestamp}.pcap') if args.protocol == 'all' else args.output
            create_udp_packets(udp_output, args.count, args.verbose)
        except ImportError as e:
            print(f"UDPパケット生成モジュールのインポートに失敗しました: {str(e)}")
    
    if args.protocol == 'x25' or args.protocol == 'all':
        try:
            from protocols.x25_packets import create_x25_packets
            x25_output = args.output.replace('.pcap', f'_x25_{timestamp}.pcap') if args.protocol == 'all' else args.output
            create_x25_packets(x25_output, args.count, args.verbose)
        except ImportError as e:
            print(f"X.25パケット生成モジュールのインポートに失敗しました: {str(e)}")
    
    print("パケット生成処理が完了しました。")

if __name__ == "__main__":
    main()