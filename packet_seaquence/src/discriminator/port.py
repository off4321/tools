#!/usr/bin/env python3
"""
ポート番号判別クラスを提供するモジュール
"""


class DiscriminatePort:
    """
    ポート番号を判別するクラス
    """
    
    def __init__(self, packet):
        """
        DiscriminatePortクラスの初期化
        
        Args:
            packet: pyshark.packet.Packetオブジェクト
        """
        self.packet = packet
    
    def discriminate(self):
        """
        パケットからポート番号情報を判別する
        
        TCP、UDP、SCTPプロトコルに対応
        
        Returns:
            dict: ポート番号情報を含む辞書、または情報がない場合はNone
        """
        try:
            port_info = {}
            
            # TCPポート番号の取得
            if hasattr(self.packet, 'tcp'):
                # 送信元ポート番号
                if hasattr(self.packet.tcp, 'srcport'):
                    port_info['src_port'] = self.packet.tcp.srcport
                
                # 宛先ポート番号
                if hasattr(self.packet.tcp, 'dstport'):
                    port_info['dst_port'] = self.packet.tcp.dstport
                
                port_info['protocol'] = 'TCP'
            
            # UDPポート番号の取得
            elif hasattr(self.packet, 'udp'):
                # 送信元ポート番号
                if hasattr(self.packet.udp, 'srcport'):
                    port_info['src_port'] = self.packet.udp.srcport
                
                # 宛先ポート番号
                if hasattr(self.packet.udp, 'dstport'):
                    port_info['dst_port'] = self.packet.udp.dstport
                
                port_info['protocol'] = 'UDP'
            
            # SCTPポート番号の取得
            elif hasattr(self.packet, 'sctp'):
                # 送信元ポート番号
                if hasattr(self.packet.sctp, 'srcport'):
                    port_info['src_port'] = self.packet.sctp.srcport
                
                # 宛先ポート番号
                if hasattr(self.packet.sctp, 'dstport'):
                    port_info['dst_port'] = self.packet.sctp.dstport
                
                port_info['protocol'] = 'SCTP'
            
            # よく知られているポート番号に基づくプロトコルの推測
            if 'dst_port' in port_info:
                dst_port = int(port_info['dst_port'])
                if dst_port == 80:
                    port_info['service'] = 'HTTP'
                elif dst_port == 443:
                    port_info['service'] = 'HTTPS'
                elif dst_port == 25:
                    port_info['service'] = 'SMTP'
                elif dst_port == 53:
                    port_info['service'] = 'DNS'
                elif dst_port == 22:
                    port_info['service'] = 'SSH'
                elif dst_port == 21:
                    port_info['service'] = 'FTP'
                # 他のポート番号に基づくサービスも追加可能
            
            return port_info if port_info else None
            
        except Exception as e:
            print(f"ポート判別中にエラーが発生しました: {str(e)}")
            return None
