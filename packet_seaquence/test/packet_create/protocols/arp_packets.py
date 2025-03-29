#!/usr/bin/env python3
"""
ARPパケットを作成するモジュール
"""
from scapy.all import Ether, ARP
from .utils import random_mac, random_ip, save_packets

def create_arp_packets(output_file, count=10, verbose=False):
    """ARPパケットの作成"""
    if verbose:
        print(f"ARP パケット作成開始: {count}パケット")
    
    packets = []
    
    # ARPリクエストとレスポンスのペアを作成
    for i in range(count // 2):
        src_mac = random_mac()
        dst_mac = random_mac()
        src_ip = random_ip("192.168.1")
        dst_ip = random_ip("192.168.1")
        
        # ARPリクエスト
        arp_request = Ether(src=src_mac, dst="ff:ff:ff:ff:ff:ff") / ARP(
            op=1,  # who-has (request)
            hwsrc=src_mac,
            psrc=src_ip,
            hwdst="00:00:00:00:00:00",
            pdst=dst_ip
        )
        
        packets.append(arp_request)
        
        # ARPレスポンス
        arp_response = Ether(src=dst_mac, dst=src_mac) / ARP(
            op=2,  # is-at (response)
            hwsrc=dst_mac,
            psrc=dst_ip,
            hwdst=src_mac,
            pdst=src_ip
        )
        
        packets.append(arp_response)
    
    # 追加のARP Gratuitousなども生成
    if count > 2:
        # Gratuitous ARP
        grat_mac = random_mac()
        grat_ip = random_ip("192.168.1")
        
        gratuitous_arp = Ether(src=grat_mac, dst="ff:ff:ff:ff:ff:ff") / ARP(
            op=1,  # who-has (request)
            hwsrc=grat_mac,
            psrc=grat_ip,
            hwdst="ff:ff:ff:ff:ff:ff",
            pdst=grat_ip
        )
        
        packets.append(gratuitous_arp)
    
    if verbose:
        print(f"ARP パケット作成完了: {len(packets)}パケット")
    
    # パケットの保存（引数2つに修正）
    save_packets(packets, output_file)
    
    return True

if __name__ == "__main__":
    # 単体テスト用
    create_arp_packets("arp_test.pcap", count=10, verbose=True)