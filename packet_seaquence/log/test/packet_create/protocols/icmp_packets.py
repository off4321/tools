#!/usr/bin/env python3
"""
ICMPパケットを作成するモジュール
"""
from scapy.all import Ether, IP, ICMP, Raw
from .utils import random_mac, random_ip, save_packets

def create_icmp_packets(output_file, count=10, verbose=False):
    """ICMPパケットの作成"""
    if verbose:
        print(f"ICMP パケット作成開始: {count}パケット")
    
    packets = []
    
    # ICMPエコー要求と応答のペアを作成
    for i in range(count // 2):
        # ソース/デスティネーションの情報
        src_mac = random_mac()
        dst_mac = random_mac()
        src_ip = random_ip("192.168.1")
        dst_ip = random_ip("192.168.2")
        
        # ICMPエコー要求 (ping)
        icmp_request = Ether(src=src_mac, dst=dst_mac) / IP(
            src=src_ip,
            dst=dst_ip
        ) / ICMP(
            type=8,  # エコー要求
            code=0,
            id=i+1000,
            seq=i+1
        ) / Raw(load=f"PING packet data #{i+1}")
        
        packets.append(icmp_request)
        
        # ICMPエコー応答 (ping応答)
        icmp_reply = Ether(src=dst_mac, dst=src_mac) / IP(
            src=dst_ip,
            dst=src_ip
        ) / ICMP(
            type=0,  # エコー応答
            code=0,
            id=i+1000,
            seq=i+1
        ) / Raw(load=f"PING packet data #{i+1}")
        
        packets.append(icmp_reply)
    
    # タイムエクシード、到達不能などのICMPエラーメッセージも追加
    if count > 4:
        # ICMPタイムエクシード (TTL超過)
        ttl_exceeded = Ether(src=random_mac(), dst=random_mac()) / IP(
            src="10.0.0.1",
            dst=src_ip
        ) / ICMP(
            type=11,  # タイムエクシード
            code=0    # トランジット中のTTL超過
        ) / IP(
            src=src_ip,
            dst="8.8.8.8"
        ) / ICMP(
            type=8,
            code=0
        )
        
        packets.append(ttl_exceeded)
        
        # ICMP宛先到達不能 (ポート到達不能)
        port_unreachable = Ether(src=random_mac(), dst=src_mac) / IP(
            src=dst_ip,
            dst=src_ip
        ) / ICMP(
            type=3,   # 宛先到達不能
            code=3    # ポート到達不能
        ) / IP(
            src=src_ip,
            dst=dst_ip
        ) / Raw(load=b'\x00' * 8)
        
        packets.append(port_unreachable)
    
    if verbose:
        print(f"ICMP パケット作成完了: {len(packets)}パケット")
    
    # パケットの保存
    save_packets(packets, output_file)
    
    return True

if __name__ == "__main__":
    # 単体テスト用
    create_icmp_packets("icmp_test.pcap", count=6, verbose=True)
