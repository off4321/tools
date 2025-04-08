#!/usr/bin/env python3
"""
IPアドレス判別クラスを提供するモジュール
"""
from scapy.all import IP, IPv6, ARP


class DiscriminateIp:
    """
    IPを判別するクラス
    """
    
    def __init__(self, packet):
        """
        DiscriminateIpクラスの初期化
        
        Args:
            packet: scapy.packet.Packetオブジェクト
        """
        self.packet = packet
    
    def discriminate(self):
        """
        パケットからIPアドレス情報を判別する
        
        IPv4とIPv6に対応
        
        Returns:
            dict: IPアドレス情報を含む辞書、または情報がない場合はNone
        """
        try:
            ip_info = {}
            
            # IPv4情報の取得
            if IP in self.packet:
                ip_layer = self.packet[IP]
                # 送信元IPアドレス
                ip_info['src_ip'] = ip_layer.src
                
                # 宛先IPアドレス
                ip_info['dst_ip'] = ip_layer.dst
                
                # プロトコル番号
                ip_info['protocol'] = str(ip_layer.proto)
                
                # TTL
                ip_info['ttl'] = str(ip_layer.ttl)
            
            # IPv6情報の取得
            elif IPv6 in self.packet:
                ipv6_layer = self.packet[IPv6]
                # 送信元IPv6アドレス
                ip_info['src_ip'] = ipv6_layer.src
                
                # 宛先IPv6アドレス
                ip_info['dst_ip'] = ipv6_layer.dst
                
                # ホップリミット
                ip_info['hop_limit'] = str(ipv6_layer.hlim)
                
                # 次ヘッダ
                ip_info['next_header'] = str(ipv6_layer.nh)
            
            # ARPパケットからIPアドレスを取得
            if ARP in self.packet:
                arp_layer = self.packet[ARP]
                # 送信元IPアドレス
                ip_info['arp_src_ip'] = arp_layer.psrc
                # ARPの場合、送信元IPアドレスを通常の送信元IPアドレスとして設定
                if 'src_ip' not in ip_info:
                    ip_info['src_ip'] = arp_layer.psrc
                
                # 宛先IPアドレス
                ip_info['arp_dst_ip'] = arp_layer.pdst
                # ARPの場合、宛先IPアドレスを通常の宛先IPアドレスとして設定
                if 'dst_ip' not in ip_info:
                    ip_info['dst_ip'] = arp_layer.pdst
            
            return ip_info if ip_info else None
            
        except Exception as e:
            print(f"IP判別中にエラーが発生しました: {str(e)}")
            import traceback
            traceback.print_exc()
            return None
