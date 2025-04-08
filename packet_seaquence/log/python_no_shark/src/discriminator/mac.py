#!/usr/bin/env python3
"""
MACアドレス判別クラスを提供するモジュール
"""
from scapy.all import Ether, ARP


class DiscriminateMac:
    """
    MACを判別するクラス
    """
    
    def __init__(self, packet):
        """
        DiscriminateMacクラスの初期化
        
        Args:
            packet: scapy.packet.Packetオブジェクト
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
            if Ether in self.packet:
                eth = self.packet[Ether]
                # 送信元MACアドレスを取得
                mac_info['src_mac'] = eth.src
                
                # 宛先MACアドレスを取得
                mac_info['dst_mac'] = eth.dst
                
                # EtherType（上位プロトコル識別子）を取得
                mac_info['ether_type'] = hex(eth.type)
            
            # ARPパケットの場合、ARPレイヤーからもMACアドレスを取得
            if ARP in self.packet:
                arp = self.packet[ARP]
                # 送信元MACアドレス（ARPレイヤー）
                mac_info['arp_src_mac'] = arp.hwsrc
                
                # 宛先MACアドレス（ARPレイヤー）
                mac_info['arp_dst_mac'] = arp.hwdst
            
            return mac_info if mac_info else None
            
        except Exception as e:
            print(f"MAC判別中にエラーが発生しました: {str(e)}")
            import traceback
            traceback.print_exc()
            return None
