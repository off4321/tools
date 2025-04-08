#!/usr/bin/env python3
import pyshark
import argparse
import traceback

def read_pcap_file(pcap_file_path):
    """
    PySharkを使用してPCAPファイルを読み込み、パケット情報を表示する
    L2プロトコル（ARPなど）も含めて解析する
    """
    print(f"Reading PCAP file: {pcap_file_path}")
    
    try:
        # PCAPファイルを読み込む
        capture = pyshark.FileCapture(pcap_file_path)
        
        packet_count = 0
        for packet in capture:
            packet_count += 1
            print(f"\nPacket {packet_count}:")
            
            # 通信時間
            print(f"  Time: {packet.sniff_time}")
            
            # 最高位のプロトコル名を表示
            highest_layer = packet.highest_layer
            print(f"  Highest Layer Protocol: {highest_layer}")
            
            # イーサネットレイヤー情報
            if hasattr(packet, 'eth'):
                print(f"  Ethernet: {packet.eth.src} -> {packet.eth.dst}")
                print(f"  Ether Type: 0x{packet.eth.type}")
            
            # ARPプロトコル情報
            if hasattr(packet, 'arp'):
                print(f"  ARP:")
                # ARPフィールドを安全に取得する
                try:
                    # 多くの場合、操作タイプは operation または op_code にある
                    if hasattr(packet.arp, 'opcode'):
                        print(f"    Operation: {packet.arp.opcode}")
                    elif hasattr(packet.arp, 'operation'):
                        print(f"    Operation: {packet.arp.operation}")
                    
                    # 送信元MACアドレス
                    if hasattr(packet.arp, 'src_hw_mac'):
                        print(f"    Sender MAC: {packet.arp.src_hw_mac}")
                    elif hasattr(packet.arp, 'src.hw_mac'):
                        print(f"    Sender MAC: {packet.arp.src.hw_mac}")
                    
                    # 送信元IPアドレス
                    if hasattr(packet.arp, 'src_proto_ipv4'):
                        print(f"    Sender IP: {packet.arp.src_proto_ipv4}")
                    elif hasattr(packet.arp, 'src.proto_ipv4'):
                        print(f"    Sender IP: {packet.arp.src.proto_ipv4}")
                    
                    # 宛先MACアドレス
                    if hasattr(packet.arp, 'dst_hw_mac'):
                        print(f"    Target MAC: {packet.arp.dst_hw_mac}")
                    elif hasattr(packet.arp, 'dst.hw_mac'):
                        print(f"    Target MAC: {packet.arp.dst.hw_mac}")
                    
                    # 宛先IPアドレス
                    if hasattr(packet.arp, 'dst_proto_ipv4'):
                        print(f"    Target IP: {packet.arp.dst_proto_ipv4}")
                    elif hasattr(packet.arp, 'dst.proto_ipv4'):
                        print(f"    Target IP: {packet.arp.dst.proto_ipv4}")
                    
                except AttributeError:
                    # ARPフィールドの詳細を取得できない場合、利用可能なすべてのフィールドを表示
                    print(f"    Available ARP fields: {dir(packet.arp)}")
            
            # ICMPプロトコル情報
            if hasattr(packet, 'icmp'):
                print(f"  ICMP:")
                try:
                    print(f"    Type: {packet.icmp.type}")
                except AttributeError:
                    pass
                
                if hasattr(packet.icmp, 'code'):
                    print(f"    Code: {packet.icmp.code}")
                if hasattr(packet.icmp, 'seq'):
                    print(f"    Sequence: {packet.icmp.seq}")
            
            # IPv6関連
            if hasattr(packet, 'ipv6'):
                print(f"  IPv6: {packet.ipv6.src} -> {packet.ipv6.dst}")
            
            # IPv6のICMPv6
            if hasattr(packet, 'icmpv6'):
                print(f"  ICMPv6:")
                try:
                    print(f"    Type: {packet.icmpv6.type}")
                except AttributeError:
                    pass
                
                if hasattr(packet.icmpv6, 'code'):
                    print(f"    Code: {packet.icmpv6.code}")
            
            # IPレイヤー情報
            if hasattr(packet, 'ip'):
                print(f"  IP: {packet.ip.src} -> {packet.ip.dst}")
                print(f"  Protocol: {packet.ip.proto}")
            
            # TCPレイヤー情報
            if hasattr(packet, 'tcp'):
                print(f"  TCP: Port {packet.tcp.srcport} -> {packet.tcp.dstport}")
                print(f"  Flags: {packet.tcp.flags}")
                
                if hasattr(packet.tcp, 'payload'):
                    print(f"  Payload: {packet.tcp.payload}")
            
            # UDPレイヤー情報
            if hasattr(packet, 'udp'):
                print(f"  UDP: Port {packet.udp.srcport} -> {packet.udp.dstport}")
                
                if hasattr(packet.udp, 'payload'):
                    print(f"  Payload: {packet.udp.payload}")
            
            # 利用可能なすべてのレイヤーを表示
            print(f"  Available Layers: {[layer.layer_name for layer in packet.layers]}")
        
        print(f"\nTotal packets: {packet_count}")
        
    except Exception as e:
        print(f"Error reading PCAP file: {e}")
        traceback.print_exc()  # 例外のスタックトレースを表示
        return None
    
    return capture

def main():
    parser = argparse.ArgumentParser(description='Read and analyze PCAP files using PyShark')
    parser.add_argument('pcap_file', help='Path to the PCAP file to analyze')
    parser.add_argument('-v', '--verbose', action='store_true', help='Display detailed packet information')
    args = parser.parse_args()
    
    read_pcap_file(args.pcap_file)

if __name__ == "__main__":
    main()