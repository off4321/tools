#!/usr/bin/env python3
"""
Ethernet判別クラスを提供するモジュール
"""
from scapy.all import Ether, Dot1Q
from src.discriminator.mac import DiscriminateMac
from src.discriminator.ip import DiscriminateIp
from src.discriminator.port import DiscriminatePort
from src.discriminator.protocol import DiscriminateProtocol


class DiscriminateEthernet:
    """
    Ethernetを判別するクラス (VLANタグを含む)
    """
    
    def __init__(self, packet):
        """
        DiscriminateEthernetクラスの初期化
        
        Args:
            packet: scapy.packet.Packetオブジェクト
        """
        self.packet = packet
    
    def discriminate(self):
        """
        パケットからEthernet情報を判別する
        
        各判別クラスを順に実行し、パケット情報を構築する
        
        Returns:
            dict: パケット情報を含む辞書
        """
        try:
            # パケット情報を初期化
            packet_info = {
                'time': self._extract_timestamp(),
                'src': '',
                'dst': '',
                'protocol': '',
                'info': {}
            }
            
            # MACアドレスの判別
            mac_discriminator = DiscriminateMac(self.packet)
            mac_info = mac_discriminator.discriminate()
            if mac_info:
                packet_info['src'] = mac_info.get('src_mac', '')
                packet_info['dst'] = mac_info.get('dst_mac', '')
                packet_info['info'].update({'mac_info': mac_info})
            
            # IPアドレスの判別
            ip_discriminator = DiscriminateIp(self.packet)
            ip_info = ip_discriminator.discriminate()
            if ip_info:
                # IPアドレスが利用可能ならMACアドレスより優先
                packet_info['src'] = ip_info.get('src_ip', packet_info['src'])
                packet_info['dst'] = ip_info.get('dst_ip', packet_info['dst'])
                packet_info['info'].update({'ip_info': ip_info})
            
            # ポート番号の判別
            port_discriminator = DiscriminatePort(self.packet)
            port_info = port_discriminator.discriminate()
            if port_info:
                packet_info['info'].update({'port_info': port_info})
            
            # プロトコルの判別
            protocol_discriminator = DiscriminateProtocol(self.packet)
            protocol_info = protocol_discriminator.discriminate()
            if protocol_info:
                packet_info['protocol'] = protocol_info.get('protocol', '')
                packet_info['info'].update({'protocol_info': protocol_info})
            
            # VLAN情報の判別
            vlan_info = self._extract_vlan_info()
            if vlan_info:
                packet_info['info'].update({'vlan_info': vlan_info})
            
            return packet_info
            
        except Exception as e:
            print(f"Ethernet判別中にエラーが発生しました: {str(e)}")
            import traceback
            traceback.print_exc()
            return None
    
    def _extract_timestamp(self):
        """
        パケットからタイムスタンプを抽出する
        
        Returns:
            datetime: タイムスタンプ（日時）
        """
        from datetime import datetime
        try:
            # Scapyパケットからタイムスタンプを取得
            if hasattr(self.packet, 'time'):
                return datetime.fromtimestamp(self.packet.time)
            
            # 現在時刻をフォールバックとして使用
            return datetime.now()
            
        except Exception as e:
            # エラーが発生した場合は現在時刻を返す
            return datetime.now()
    
    def _extract_vlan_info(self):
        """
        パケットからVLAN情報を抽出する
        
        Returns:
            dict: VLAN情報を含む辞書、または情報がない場合はNone
        """
        vlan_info = {}
        
        # 802.1QタグのあるVLAN情報を抽出
        if Dot1Q in self.packet:
            dot1q = self.packet[Dot1Q]
            vlan_info['id'] = str(dot1q.vlan)
            vlan_info['priority'] = str(dot1q.prio)
        
        return vlan_info if vlan_info else None
