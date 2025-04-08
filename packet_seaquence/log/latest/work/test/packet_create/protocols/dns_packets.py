#!/usr/bin/env python3
"""
DNSパケットを作成するモジュール
"""
from scapy.all import Ether, IP, UDP, DNS, DNSQR, DNSRR
from .utils import random_mac, random_ip, random_port, save_packets

def create_dns_packets(output_file, count=10, verbose=False):
    """DNSパケットの作成"""
    if verbose:
        print(f"DNS パケット作成開始: {count}パケット")
    
    packets = []
    domains = [
        "example.com", "google.com", "github.com", "amazon.com", 
        "microsoft.com", "apple.com", "facebook.com", "twitter.com"
    ]
    
    # DNSクエリとレスポンスを作成
    for i in range(min(count, len(domains))):
        client_mac = random_mac()
        server_mac = random_mac()
        client_ip = random_ip("192.168.1")
        server_ip = "8.8.8.8"  # DNSサーバーのIP
        domain = domains[i]
        
        # DNSクエリ
        dns_query = Ether(src=client_mac, dst=server_mac) / IP(
            src=client_ip,
            dst=server_ip
        ) / UDP(
            sport=random_port(),
            dport=53
        ) / DNS(
            rd=1,  # 再帰的クエリを要求
            qd=DNSQR(qname=domain)
        )
        
        packets.append(dns_query)
        
        # DNSレスポンス
        dns_response = Ether(src=server_mac, dst=client_mac) / IP(
            src=server_ip,
            dst=client_ip
        ) / UDP(
            sport=53,
            dport=dns_query[UDP].sport
        ) / DNS(
            id=dns_query[DNS].id,
            qr=1,  # レスポンス
            rd=1,  # 再帰的クエリが要求された
            ra=1,  # 再帰的クエリが利用可能
            qd=DNSQR(qname=domain),
            an=DNSRR(
                rrname=domain,
                type='A',
                ttl=3600,
                rdata=random_ip("203.0.113")
            )
        )
        
        packets.append(dns_response)
    
    # 追加のDNSエラーケースやCNAMEレコードなど
    if count > len(domains) * 2:
        # NXDOMAIN (存在しないドメイン)
        non_existent = "nonexistent-domain-example.com"
        nxdomain_query = Ether(src=client_mac, dst=server_mac) / IP(
            src=client_ip,
            dst=server_ip
        ) / UDP(
            sport=random_port(),
            dport=53
        ) / DNS(
            rd=1,
            qd=DNSQR(qname=non_existent)
        )
        
        packets.append(nxdomain_query)
        
        # NXDOMAINレスポンス
        nxdomain_response = Ether(src=server_mac, dst=client_mac) / IP(
            src=server_ip,
            dst=client_ip
        ) / UDP(
            sport=53,
            dport=nxdomain_query[UDP].sport
        ) / DNS(
            id=nxdomain_query[DNS].id,
            qr=1,
            rd=1,
            ra=1,
            rcode=3,  # NXDOMAIN
            qd=DNSQR(qname=non_existent)
        )
        
        packets.append(nxdomain_response)
    
    if verbose:
        print(f"DNS パケット作成完了: {len(packets)}パケット")
    
    # パケットの保存
    save_packets(packets, output_file)
    
    return True

if __name__ == "__main__":
    # 単体テスト用
    create_dns_packets("dns_test.pcap", count=10, verbose=True)