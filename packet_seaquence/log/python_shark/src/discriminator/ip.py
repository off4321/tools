#!/usr/bin/env python3
"""
IPアドレス判別クラスを提供するモジュール
"""


class DiscriminateIp:
    """
    IPを判別するクラス
    """
    
    def __init__(self, packet):
        """
        DiscriminateIpクラスの初期化
        
        Args:
            packet: pyshark.packet.Packetオブジェクト
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
            if hasattr(self.packet, 'ip'):
                # 送信元IPアドレス
                if hasattr(self.packet.ip, 'src'):
                    ip_info['src_ip'] = self.packet.ip.src
                
                # 宛先IPアドレス
                if hasattr(self.packet.ip, 'dst'):
                    ip_info['dst_ip'] = self.packet.ip.dst
                
                # プロトコル番号
                if hasattr(self.packet.ip, 'proto'):
                    ip_info['protocol'] = self.packet.ip.proto
                
                # TTL
                if hasattr(self.packet.ip, 'ttl'):
                    ip_info['ttl'] = self.packet.ip.ttl
            
            # IPv6情報の取得
            elif hasattr(self.packet, 'ipv6'):
                # 送信元IPv6アドレス
                if hasattr(self.packet.ipv6, 'src'):
                    ip_info['src_ip'] = self.packet.ipv6.src
                
                # 宛先IPv6アドレス
                if hasattr(self.packet.ipv6, 'dst'):
                    ip_info['dst_ip'] = self.packet.ipv6.dst
                
                # ホップリミット
                if hasattr(self.packet.ipv6, 'hlim'):
                    ip_info['hop_limit'] = self.packet.ipv6.hlim
                
                # 次ヘッダ
                if hasattr(self.packet.ipv6, 'nxt'):
                    ip_info['next_header'] = self.packet.ipv6.nxt
            
            # ARPパケットからIPアドレスを取得
            if hasattr(self.packet, 'arp'):
                try:
                    # 送信元IPアドレス
                    if hasattr(self.packet.arp, 'src_proto_ipv4'):
                        ip_info['arp_src_ip'] = self.packet.arp.src_proto_ipv4
                        # ARPの場合、送信元IPアドレスを通常の送信元IPアドレスとして設定
                        if 'src_ip' not in ip_info:
                            ip_info['src_ip'] = self.packet.arp.src_proto_ipv4
                    elif hasattr(self.packet.arp, 'src.proto_ipv4'):
                        ip_info['arp_src_ip'] = self.packet.arp.src.proto_ipv4
                        if 'src_ip' not in ip_info:
                            ip_info['src_ip'] = self.packet.arp.src.proto_ipv4
                    
                    # 宛先IPアドレス
                    if hasattr(self.packet.arp, 'dst_proto_ipv4'):
                        ip_info['arp_dst_ip'] = self.packet.arp.dst_proto_ipv4
                        # ARPの場合、宛先IPアドレスを通常の宛先IPアドレスとして設定
                        if 'dst_ip' not in ip_info:
                            ip_info['dst_ip'] = self.packet.arp.dst_proto_ipv4
                    elif hasattr(self.packet.arp, 'dst.proto_ipv4'):
                        ip_info['arp_dst_ip'] = self.packet.arp.dst.proto_ipv4
                        if 'dst_ip' not in ip_info:
                            ip_info['dst_ip'] = self.packet.arp.dst.proto_ipv4
                except AttributeError:
                    # ARPフィールドの取得に失敗した場合は無視
                    pass
            
            return ip_info if ip_info else None
            
        except Exception as e:
            print(f"IP判別中にエラーが発生しました: {str(e)}")
            return None
