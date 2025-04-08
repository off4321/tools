#!/usr/bin/env python3
"""
HTTPパケットを作成するモジュール
"""
from scapy.all import Ether, IP, TCP, Raw
from .utils import random_mac, random_ip, random_port, save_packets

def create_http_packets(output_file, count=10, verbose=False):
    """HTTPパケットの作成"""
    if verbose:
        print(f"HTTP パケット作成開始: {count}パケット")
    
    packets = []
    
    # HTTPリクエストとレスポンスのペアを作成
    for i in range(count // 2):
        client_mac = random_mac()
        server_mac = random_mac()
        client_ip = random_ip("192.168.1")
        server_ip = random_ip("10.0.0")  # 安全なIPプレフィックスを使用
        client_port = 12345 + i
        server_port = 80
        seq_base = 1000 * (i + 1)
        
        # HTTPリクエスト
        http_request = (
            "GET / HTTP/1.1\r\n"
            f"Host: example-{i}.com\r\n"
            "User-Agent: Mozilla/5.0\r\n"
            "Accept: */*\r\n"
            "\r\n"
        )
        
        request_packet = Ether(src=client_mac, dst=server_mac) / IP(
            src=client_ip,
            dst=server_ip
        ) / TCP(
            sport=client_port,
            dport=server_port,
            flags="PA",
            seq=seq_base,
            ack=seq_base
        ) / Raw(load=http_request)
        
        packets.append(request_packet)
        
        # HTTPレスポンス
        http_response = (
            "HTTP/1.1 200 OK\r\n"
            "Content-Type: text/html\r\n"
            "Content-Length: 52\r\n"
            "\r\n"
            f"<html><body><h1>Hello from example-{i}.com</h1></body></html>"
        )
        
        response_packet = Ether(src=server_mac, dst=client_mac) / IP(
            src=server_ip,
            dst=client_ip
        ) / TCP(
            sport=server_port,
            dport=client_port,
            flags="PA",
            seq=seq_base,
            ack=seq_base + len(http_request)
        ) / Raw(load=http_response)
        
        packets.append(response_packet)
    
    if verbose:
        print(f"HTTP パケット作成完了: {len(packets)}パケット")
    
    # パケットを保存（verbose引数を削除）
    save_packets(packets, output_file)
    
    return True

if __name__ == "__main__":
    # 単体テスト用
    create_http_packets("http_test.pcap", count=10, verbose=True)