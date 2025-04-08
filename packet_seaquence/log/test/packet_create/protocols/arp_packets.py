#!/usr/bin/env python3
"""
ARPパケットを作成するモジュール
"""
from scapy.all import Ether, ARP, IP
from .utils import random_mac, random_ip, save_packets

def create_arp_packets(output_file, count=10, verbose=False):
    """ARPパケットの作成"""
    if verbose:
        print(f"ARP パケット作成開始: {count}パケット")
    
    packets = []
    
    # ARPリクエストとレスポンスのペアを作成
    for i in range(count // 2):
        src_mac = random_mac()
        dst_mac = "ff:ff:ff:ff:ff:ff"  # ブロードキャスト
        src_ip = random_ip("192.168.1")
        dst_ip = random_ip("192.168.1")
        target_mac = random_mac()
        
        # ARPリクエスト
        arp_request = Ether(src=src_mac, dst=dst_mac) / ARP(
            op=1,  # リクエスト
            hwsrc=src_mac,
            psrc=src_ip,
            hwdst="00:00:00:00:00:00",
            pdst=dst_ip
        )
        packets.append(arp_request)
        
        # ARPレスポンス
        arp_reply = Ether(src=target_mac, dst=src_mac) / ARP(
            op=2,  # レスポンス
            hwsrc=target_mac,
            psrc=dst_ip,
            hwdst=src_mac,
            pdst=src_ip
        )
        packets.append(arp_reply)
    
    save_packets(packets, output_file, verbose)
    return packets