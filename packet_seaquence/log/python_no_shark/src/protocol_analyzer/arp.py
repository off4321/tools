#!/usr/bin/env python3
"""
ARPプロトコル解析クラスを提供するモジュール
"""
from scapy.all import ARP
from src.protocol_analyzer.base import ProtocolAnalyzer


class AnalyzeArp(ProtocolAnalyzer):
    """
    ARP内の詳細情報を解析するクラス
    """
    
    def analyze(self):
        """
        パケットからARP情報を解析する
        
        Returns:
            dict: ARP情報を含む辞書、または情報がない場合はNone
        """
        if ARP not in self.packet:
            return None
        
        try:
            arp_info = {}
            arp_layer = self.packet[ARP]
            
            # ARPオペレーションコード
            arp_code = int(arp_layer.op)
            arp_info['operation_code'] = arp_code
            arp_info['operation'] = self._get_arp_operation(arp_code)
            
            # 送信元MACアドレス
            arp_info['src_mac'] = arp_layer.hwsrc
            
            # 送信元IPアドレス
            arp_info['src_ip'] = arp_layer.psrc
            
            # 宛先MACアドレス
            arp_info['dst_mac'] = arp_layer.hwdst
            
            # 宛先IPアドレス
            arp_info['dst_ip'] = arp_layer.pdst
            
            return arp_info
            
        except Exception as e:
            print(f"ARP解析中にエラーが発生しました: {str(e)}")
            import traceback
            traceback.print_exc()
            return None
    
    def get_display_info(self):
        """
        表示用のARP情報を取得する
        
        Returns:
            str: 表示用のARP情報
        """
        arp_info = self.analyze()
        if not arp_info:
            return "ARP情報なし"
        
        operation = arp_info.get('operation', '不明')
        src_mac = arp_info.get('src_mac', '不明')
        src_ip = arp_info.get('src_ip', '不明')
        dst_mac = arp_info.get('dst_mac', '不明')
        dst_ip = arp_info.get('dst_ip', '不明')
        
        if 'REQUEST' in operation:
            return f"ARP Request: {src_ip} ({src_mac}) は {dst_ip} のMACアドレスを要求"
        elif 'REPLY' in operation:
            return f"ARP Reply: {src_ip} は {src_mac} です"
        else:
            return f"ARP {operation}: {src_ip} ({src_mac}) -> {dst_ip} ({dst_mac})"
    
    def get_summary(self):
        """
        ARP情報のサマリーを取得する
        
        Returns:
            str: ARP情報のサマリー
        """
        arp_info = self.analyze()
        if not arp_info:
            return "ARP"
        
        operation = arp_info.get('operation', '').split('-')[0]  # REQUEST, REPLY などの部分だけ取得
        src_ip = arp_info.get('src_ip', '')
        dst_ip = arp_info.get('dst_ip', '')
        
        return f"ARP {operation} {src_ip} -> {dst_ip}"
    
    def _get_arp_operation(self, code):
        """
        ARPオペレーションコードの説明を取得する
        
        Args:
            code (int): ARPオペレーションコード
        
        Returns:
            str: ARPオペレーションの説明
        """
        arp_operations = {
            1: 'REQUEST',
            2: 'REPLY',
            3: 'RARP-REQUEST',
            4: 'RARP-REPLY',
            5: 'DRARP-REQUEST',
            6: 'DRARP-REPLY',
            7: 'DRARP-ERROR',
            8: 'InARP-REQUEST',
            9: 'InARP-REPLY'
        }
        return arp_operations.get(code, f'UNKNOWN ({code})')
