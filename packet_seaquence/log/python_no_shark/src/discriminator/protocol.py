#!/usr/bin/env python3
"""
プロトコル判別クラスを提供するモジュール
"""
from scapy.all import TCP, UDP, IP, IPv6, ARP, ICMP, DNS, SCTP
from src.discriminator.available import DiscriminateAvailable

# 循環参照を避けるため、アナライザーを直接インポート
from src.protocol_analyzer.base import ProtocolAnalyzer
from src.protocol_analyzer.arp import AnalyzeArp
from src.protocol_analyzer.tcp import AnalyzeTcp
from src.protocol_analyzer.udp import AnalyzeUdp
from src.protocol_analyzer.sctp import AnalyzeSctp
from src.protocol_analyzer.ipv4 import AnalyzeIPv4
from src.protocol_analyzer.http import AnalyzeHttp
from src.protocol_analyzer.dns import AnalyzeDns
from src.protocol_analyzer.https import AnalyzeHttps
from src.protocol_analyzer.x25 import AnalyzeX25
from src.protocol_analyzer.unSupportedProtocol import AnalyzeUnsupportedProtocol

class DiscriminateProtocol:
    """
    プロトコルを判別するクラス
    """
    
    def __init__(self, packet):
        """
        DiscriminateProtocolクラスの初期化
        
        Args:
            packet: scapy.packet.Packetオブジェクト
        """
        self.packet = packet
    
    def discriminate(self):
        """
        パケットからプロトコル情報を判別する
        
        パケットの最高位レイヤーと各レイヤー情報を使用して
        プロトコルを判別する
        
        Returns:
            dict: プロトコル情報を含む辞書、または情報がない場合はNone
        """
        try:
            available = DiscriminateAvailable(self.packet).discriminate()
            if not available:
                return None
            
            protocol_info = {}
            self._extract_protocol_details(protocol_info, available)
            print(f"DEBUG - Protocol info before extraction: {protocol_info}")
            return protocol_info
            
        except Exception as e:
            print(f"プロトコル判別中にエラーが発生しました: {str(e)}")
            import traceback
            traceback.print_exc()
            return None
    
    def _extract_protocol_details(self, protocol_info, available):
        highest_layer = available.get('highest_layer', 'Unknown')
        # protocol_analyzer各クラスで詳細を分析させる
        layers = available.get('layer_names', [])
        
        if 'ip' in layers:
            analyzer = AnalyzeIPv4(self.packet)
            protocol_info['ipv4_info'] = analyzer.analyze()
        
        if 'arp' in layers:
            analyzer = AnalyzeArp(self.packet)
            protocol_info['arp_info'] = analyzer.analyze()
        
        if 'tcp' in layers:
            analyzer = AnalyzeTcp(self.packet)
            protocol_info['tcp_info'] = analyzer.analyze()
        
        if 'udp' in layers:
            analyzer = AnalyzeUdp(self.packet)
            protocol_info['udp_info'] = analyzer.analyze()
        
        if 'sctp' in layers:
            analyzer = AnalyzeSctp(self.packet)
            protocol_info['sctp_info'] = analyzer.analyze()
        
        if 'dns' in layers:
            analyzer = AnalyzeDns(self.packet)
            protocol_info['dns_info'] = analyzer.analyze()
        
        if 'http' in layers:
            analyzer = AnalyzeHttp(self.packet)
            protocol_info['http_info'] = analyzer.analyze()
        
        if 'ssl' in layers or 'tls' in layers:
            analyzer = AnalyzeHttps(self.packet)
            protocol_info['https_info'] = analyzer.analyze()
        
        # X.25の検出と分析を追加
        if 'x25' in layers:
            analyzer = AnalyzeX25(self.packet)
            protocol_info['x25_info'] = analyzer.analyze()
        
        # 認識できなかった場合にprotocol_nameを格納する
        if not protocol_info:
            protocol_info["protocol_name"] = highest_layer

        print(f"DEBUG - Protocol info after extraction: {protocol_info}")
        return protocol_info
    
    def _get_arp_operation(self, code):
        # ARP操作コード
        operations = {
            '1': 'REQUEST',
            '2': 'REPLY'
        }
        return operations.get(code, 'UNKNOWN')
    
    def _get_icmp_type(self, code):
        # ICMPタイプコード
        types = {
            '0': 'Echo Reply',
            '8': 'Echo Request',
            '3': 'Destination Unreachable'
        }
        return types.get(code, 'UNKNOWN')
    
    def _get_icmpv6_type(self, code):
        # ICMPv6タイプコード
        types = {
            '128': 'Echo Request',
            '129': 'Echo Reply',
            '1': 'Destination Unreachable'
        }
        return types.get(code, 'UNKNOWN')
    
    def _get_tcp_flags_desc(self, flags):
        # TCPフラグの説明
        flag_descs = []
        if 'F' in flags:
            flag_descs.append('FIN')
        if 'S' in flags:
            flag_descs.append('SYN')
        if 'R' in flags:
            flag_descs.append('RST')
        if 'P' in flags:
            flag_descs.append('PSH')
        if 'A' in flags:
            flag_descs.append('ACK')
        if 'U' in flags:
            flag_descs.append('URG')
        return ', '.join(flag_descs) if flag_descs else 'None'
