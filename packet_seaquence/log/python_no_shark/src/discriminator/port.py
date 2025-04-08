#!/usr/bin/env python3
"""
ポート番号判別クラスを提供するモジュール
"""
from scapy.all import TCP, UDP, SCTP


class DiscriminatePort:
    """
    ポート番号を判別するクラス
    """
    
    def __init__(self, packet):
        """
        DiscriminatePortクラスの初期化
        
        Args:
            packet: scapy.packet.Packetオブジェクト
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
            if TCP in self.packet:
                tcp_layer = self.packet[TCP]
                # 送信元ポート番号
                port_info['src_port'] = str(tcp_layer.sport)
                
                # 宛先ポート番号
                port_info['dst_port'] = str(tcp_layer.dport)
                
                port_info['protocol'] = 'TCP'
            
            # UDPポート番号の取得
            elif UDP in self.packet:
                udp_layer = self.packet[UDP]
                # 送信元ポート番号
                port_info['src_port'] = str(udp_layer.sport)
                
                # 宛先ポート番号
                port_info['dst_port'] = str(udp_layer.dport)
                
                port_info['protocol'] = 'UDP'
            
            # SCTPポート番号の取得
            elif SCTP in self.packet:
                sctp_layer = self.packet[SCTP]
                # 送信元ポート番号
                port_info['src_port'] = str(sctp_layer.sport)
                
                # 宛先ポート番号
                port_info['dst_port'] = str(sctp_layer.dport)
                
                port_info['protocol'] = 'SCTP'
            
            # よく知られているポート番号に基づくプロトコルの推測
            if 'dst_port' in port_info:
                try:
                    dst_port = int(port_info['dst_port'])
                    port_info['service'] = self._get_service_name(dst_port)
                except (ValueError, TypeError):
                    pass
            
            return port_info if port_info else None
            
        except Exception as e:
            print(f"ポート判別中にエラーが発生しました: {str(e)}")
            import traceback
            traceback.print_exc()
            return None
            
    def _get_service_name(self, port):
        """
        ポート番号からサービス名を取得する
        
        Args:
            port (int): ポート番号
            
        Returns:
            str: サービス名または空文字列
        """
        well_known_ports = {
            20: 'FTP-DATA',
            21: 'FTP',
            22: 'SSH',
            23: 'TELNET',
            25: 'SMTP',
            53: 'DNS',
            67: 'DHCP-Server',
            68: 'DHCP-Client',
            80: 'HTTP',
            110: 'POP3',
            143: 'IMAP',
            443: 'HTTPS',
            465: 'SMTPS',
            993: 'IMAPS',
            995: 'POP3S',
            3389: 'RDP'
        }
        return well_known_ports.get(port, '')
