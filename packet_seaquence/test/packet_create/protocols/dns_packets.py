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
    
    # DNS用のドメイン名リスト
    domains = ["example.com", "google.com", "github.com", "microsoft.com", "wikipedia.org"]
    
    # DNSクエリとレスポンスのペアを作成
    for i in range(min(count // 2, len(domains))):
        client_mac = random_mac()
        server_mac = random_mac()
        client_ip = random_ip("192.168.1")
        server_ip = random_ip("10.0.0")  # DNSサーバーIP
        client_port = random_port()
        server_port = 53  # DNS
        domain = domains[i]
        
        # DNSクエリ
        dns_query = Ether(src=client_mac, dst=server_mac) / IP(
            src=client_ip,
            dst=server_ip
        ) / UDP(
            sport=client_port,
            dport=server_port
        ) / DNS(
            rd=1,  # recursion desired
            qd=DNSQR(qname=domain)
        )
        
        packets.append(dns_query)
        
        # DNSレスポンス
        dns_response = Ether(src=server_mac, dst=client_mac) / IP(
            src=server_ip,
            dst=client_ip
        ) / UDP(
            sport=server_port,
            dport=dns_query[UDP].sport
        ) / DNS(
            id=dns_query[DNS].id,
            qr=1,  # response
            rd=1,  # recursion desired
            ra=1,  # recursion available
            qd=dns_query[DNS].qd,
            an=DNSRR(
                rrname=domain,
                ttl=86400,
                type='A',
                rdata=random_ip("104.16")
            )
        )
        
        packets.append(dns_response)
    
    # NXDOMAINケースなど他のDNSレスポンスタイプも追加
    if count > len(domains) * 2:
        client_mac = random_mac()
        server_mac = random_mac()
        client_ip = random_ip("192.168.1")
        server_ip = random_ip("10.0.0")
        client_port = random_port()
        server_port = 53
        non_existent = f"nonexistent-{random.randint(1000, 9999)}.example.org"
        
        # 存在しないドメインへのDNSクエリ
        nxdomain_query = Ether(src=client_mac, dst=server_mac) / IP(
            src=client_ip,
            dst=server_ip
        ) / UDP(
            sport=client_port,
            dport=server_port
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
    
    # パケットの保存（正しい引数で呼び出し）
    save_packets(packets, output_file)
    
    return True

if __name__ == "__main__":
    # 単体テスト用
    create_dns_packets("dns_test.pcap", count=10, verbose=True)