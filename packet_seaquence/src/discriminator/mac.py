#!/usr/bin/env python3
"""
MACアドレス判別クラスを提供するモジュール
"""


class DiscriminateMac:
    """
    MACを判別するクラス
    """
    
    def __init__(self, packet):
        """
        DiscriminateMacクラスの初期化
        
        Args:
            packet: pyshark.packet.Packetオブジェクト
        """
        self.packet = packet
    
    def discriminate(self):
        """
        パケットからMACアドレス情報を判別する
        
        Returns:
            dict: MACアドレス情報を含む辞書、または情報がない場合はNone
        """
        try:
            mac_info = {}
            
            # Ethernetレイヤーがあるか確認
            if hasattr(self.packet, 'eth'):
                # 送信元MACアドレスを取得
                if hasattr(self.packet.eth, 'src'):
                    mac_info['src_mac'] = self.packet.eth.src
                
                # 宛先MACアドレスを取得
                if hasattr(self.packet.eth, 'dst'):
                    mac_info['dst_mac'] = self.packet.eth.dst
                
                # EtherType（上位プロトコル識別子）を取得
                if hasattr(self.packet.eth, 'type'):
                    mac_info['ether_type'] = self.packet.eth.type
            
            # ARPパケットの場合、ARPレイヤーからもMACアドレスを取得
            if hasattr(self.packet, 'arp'):
                try:
                    # 送信元MACアドレス（ARPレイヤー）
                    if hasattr(self.packet.arp, 'src_hw_mac'):
                        mac_info['arp_src_mac'] = self.packet.arp.src_hw_mac
                    elif hasattr(self.packet.arp, 'src.hw_mac'):
                        mac_info['arp_src_mac'] = self.packet.arp.src.hw_mac
                    
                    # 宛先MACアドレス（ARPレイヤー）
                    if hasattr(self.packet.arp, 'dst_hw_mac'):
                        mac_info['arp_dst_mac'] = self.packet.arp.dst_hw_mac
                    elif hasattr(self.packet.arp, 'dst.hw_mac'):
                        mac_info['arp_dst_mac'] = self.packet.arp.dst.hw_mac
                except AttributeError:
                    # ARPフィールドの取得に失敗した場合は無視
                    pass
            
            return mac_info if mac_info else None
            
        except Exception as e:
            print(f"MAC判別中にエラーが発生しました: {str(e)}")
            return None
