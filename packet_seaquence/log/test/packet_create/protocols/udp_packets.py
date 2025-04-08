#!/usr/bin/env python3
"""
UDPパケットを作成するモジュール
"""
from scapy.all import Ether, IP, UDP, Raw
from .utils import random_mac, random_ip, random_port, save_packets

def create_udp_packets(output_file, count=10, verbose=False):
    """UDPパケットの作成"""
    if verbose:
        print(f"UDP パケット作成開始: {count}パケット")
    
    packets = []
    
    # 単純なUDPパケット交換を作成
    for i in range(count // 2):
        client_mac = random_mac()
        server_mac = random_mac()
        client_ip = random_ip("192.168.1")
        server_ip = random_ip("203.0.113")
        client_port = random_port()
        server_port = 53  # DNS
        
        # クライアントからサーバーへのUDPパケット
        client_data = f"クライアントからのUDPデータ #{i+1}"
        client_to_server = Ether(src=client_mac, dst=server_mac) / IP(
            src=client_ip,
            dst=server_ip
        ) / UDP(
            sport=client_port,
            dport=server_port
        ) / Raw(load=client_data)
        
        packets.append(client_to_server)
        
        # サーバーからクライアントへのUDPパケット
        server_data = f"サーバーからのUDPデータ #{i+1}"
        server_to_client = Ether(src=server_mac, dst=client_mac) / IP(
            src=server_ip,
            dst=client_ip
        ) / UDP(
            sport=server_port,
            dport=client_port
        ) / Raw(load=server_data)
        
        packets.append(server_to_client)
    
    # 追加の特殊UDPパケットも含める
    if count > 2:
        # UDPブロードキャスト
        broadcast_mac = "ff:ff:ff:ff:ff:ff"
        broadcast_ip = "255.255.255.255"
        
        broadcast_packet = Ether(src=random_mac(), dst=broadcast_mac) / IP(
            src=random_ip("192.168.1"),
            dst=broadcast_ip
        ) / UDP(
            sport=random_port(),
            dport=67  # DHCP
        ) / Raw(load="UDPブロードキャストパケット")
        
        packets.append(broadcast_packet)
    
    if verbose:
        print(f"UDP パケット作成完了: {len(packets)}パケット")
    
    # パケットの保存
    save_packets(packets, output_file)
    
    return True

if __name__ == "__main__":
    # 単体テスト用
    create_udp_packets("udp_test.pcap", count=10, verbose=True)