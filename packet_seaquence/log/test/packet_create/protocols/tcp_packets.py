#!/usr/bin/env python3
"""
TCPパケットを作成するモジュール
"""
from scapy.all import Ether, IP, TCP, Raw
from .utils import random_mac, random_ip, random_port, save_packets

def create_tcp_packets(output_file, count=10, verbose=False):
    """TCPパケットの作成"""
    if verbose:
        print(f"TCP パケット作成開始: {count}パケット")
    
    packets = []
    
    # TCP接続のシーケンスを作成
    for i in range(max(1, count // 10)):
        client_mac = random_mac()
        server_mac = random_mac()
        client_ip = random_ip("192.168.1")
        server_ip = random_ip("10.0.0")
        client_port = random_port(10000, 60000)
        server_port = 80
        
        # TCP 3-wayハンドシェイク
        # 1. SYN
        seq_num = 1000 * i
        syn = Ether(src=client_mac, dst=server_mac) / IP(
            src=client_ip,
            dst=server_ip
        ) / TCP(
            sport=client_port,
            dport=server_port,
            flags="S",  # SYN
            seq=seq_num
        )
        
        packets.append(syn)
        
        # 2. SYN-ACK
        server_seq = 2000 * i
        syn_ack = Ether(src=server_mac, dst=client_mac) / IP(
            src=server_ip,
            dst=client_ip
        ) / TCP(
            sport=server_port,
            dport=client_port,
            flags="SA",  # SYN-ACK
            seq=server_seq,
            ack=seq_num + 1
        )
        
        packets.append(syn_ack)
        
        # 3. ACK
        ack = Ether(src=client_mac, dst=server_mac) / IP(
            src=client_ip,
            dst=server_ip
        ) / TCP(
            sport=client_port,
            dport=server_port,
            flags="A",  # ACK
            seq=seq_num + 1,
            ack=server_seq + 1
        )
        
        packets.append(ack)
        
        # データ送信 (クライアント→サーバー)
        data1 = "Hello, server! This is packet " + str(i)
        push_ack1 = Ether(src=client_mac, dst=server_mac) / IP(
            src=client_ip,
            dst=server_ip
        ) / TCP(
            sport=client_port,
            dport=server_port,
            flags="PA",  # PSH-ACK
            seq=seq_num + 1,
            ack=server_seq + 1
        ) / Raw(load=data1)
        
        packets.append(push_ack1)
        
        # サーバからのACK
        ack2 = Ether(src=server_mac, dst=client_mac) / IP(
            src=server_ip,
            dst=client_ip
        ) / TCP(
            sport=server_port,
            dport=client_port,
            flags="A",  # ACK
            seq=server_seq + 1,
            ack=seq_num + 1 + len(data1)
        )
        
        packets.append(ack2)
        
        # データ送信 (サーバー→クライアント)
        data2 = "Hello, client! This is response " + str(i)
        push_ack2 = Ether(src=server_mac, dst=client_mac) / IP(
            src=server_ip,
            dst=client_ip
        ) / TCP(
            sport=server_port,
            dport=client_port,
            flags="PA",  # PSH-ACK
            seq=server_seq + 1,
            ack=seq_num + 1 + len(data1)
        ) / Raw(load=data2)
        
        packets.append(push_ack2)
        
        # クライアントからのACK
        ack3 = Ether(src=client_mac, dst=server_mac) / IP(
            src=client_ip,
            dst=server_ip
        ) / TCP(
            sport=client_port,
            dport=server_port,
            flags="A",  # ACK
            seq=seq_num + 1 + len(data1),
            ack=server_seq + 1 + len(data2)
        )
        
        packets.append(ack3)
        
        # TCP接続終了
        # 1. FIN-ACK (クライアント→サーバー)
        fin_ack1 = Ether(src=client_mac, dst=server_mac) / IP(
            src=client_ip,
            dst=server_ip
        ) / TCP(
            sport=client_port,
            dport=server_port,
            flags="FA",  # FIN-ACK
            seq=seq_num + 1 + len(data1),
            ack=server_seq + 1 + len(data2)
        )
        
        packets.append(fin_ack1)
        
        # 2. ACK (サーバー→クライアント)
        ack4 = Ether(src=server_mac, dst=client_mac) / IP(
            src=server_ip,
            dst=client_ip
        ) / TCP(
            sport=server_port,
            dport=client_port,
            flags="A",  # ACK
            seq=server_seq + 1 + len(data2),
            ack=seq_num + 1 + len(data1) + 1
        )
        
        packets.append(ack4)
        
        # 3. FIN-ACK (サーバー→クライアント)
        fin_ack2 = Ether(src=server_mac, dst=client_mac) / IP(
            src=server_ip,
            dst=client_ip
        ) / TCP(
            sport=server_port,
            dport=client_port,
            flags="FA",  # FIN-ACK
            seq=server_seq + 1 + len(data2),
            ack=seq_num + 1 + len(data1) + 1
        )
        
        packets.append(fin_ack2)
        
        # 4. ACK (クライアント→サーバー)
        ack5 = Ether(src=client_mac, dst=server_mac) / IP(
            src=client_ip,
            dst=server_ip
        ) / TCP(
            sport=client_port,
            dport=server_port,
            flags="A",  # ACK
            seq=seq_num + 1 + len(data1) + 1,
            ack=server_seq + 1 + len(data2) + 1
        )
        
        packets.append(ack5)
    
    if verbose:
        print(f"TCP パケット作成完了: {len(packets)}パケット")
    
    # パケットの保存
    save_packets(packets, output_file)
    
    return True

if __name__ == "__main__":
    # 単体テスト用
    create_tcp_packets("tcp_test.pcap", count=20, verbose=True)