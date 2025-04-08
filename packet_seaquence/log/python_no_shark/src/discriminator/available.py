#!/usr/bin/env python3
"""
使用可能なプロトコル判別クラスを提供するモジュール
"""
from scapy.all import Ether, IP, IPv6, TCP, UDP, ARP, ICMP, DNS, Raw
try:
    from scapy.all import X25
except ImportError:
    X25 = None


class DiscriminateAvailable:
    """
    使用可能なプロトコルを判別するクラス
    """
    
    def __init__(self, packet):
        """
        DiscriminateAvailableクラスの初期化
        
        Args:
            packet: scapy.packet.Packetオブジェクト
        """
        self.packet = packet
    
    def discriminate(self):
        """
        パケットから使用可能なレイヤー情報を判別する
        
        Returns:
            dict: 使用可能なレイヤー情報を含む辞書
        """
        # DEBUG - 取得したパケット情報を表示
        print(f"DEBUG - Received packet summary: {self.packet.summary()}")
        self.packet.show()  # Scapyの詳細表示

        try:
            # 利用可能なレイヤー名のリストを取得
            available_layers = []
            layer_info = {}
            highest_layer = "Unknown"
            
            # イーサネットレイヤーの確認
            if Ether in self.packet:
                available_layers.append('ether')
                layer_info['ether'] = list(self.packet[Ether].fields.keys())
                highest_layer = 'ETHERNET'
            
            # ARPの確認
            if ARP in self.packet:
                available_layers.append('arp')
                layer_info['arp'] = list(self.packet[ARP].fields.keys())
                highest_layer = 'ARP'
            
            # IPレイヤーの確認
            if IP in self.packet:
                available_layers.append('ip')
                layer_info['ip'] = list(self.packet[IP].fields.keys())
                highest_layer = 'IP'
            
            # IPv6レイヤーの確認
            if IPv6 in self.packet:
                available_layers.append('ipv6')
                layer_info['ipv6'] = list(self.packet[IPv6].fields.keys())
                highest_layer = 'IPv6'
            
            # TCPレイヤーの確認
            if TCP in self.packet:
                available_layers.append('tcp')
                layer_info['tcp'] = list(self.packet[TCP].fields.keys())
                highest_layer = 'TCP'
            
            # UDPレイヤーの確認
            if UDP in self.packet:
                available_layers.append('udp')
                layer_info['udp'] = list(self.packet[UDP].fields.keys())
                highest_layer = 'UDP'
            
            # ICMPレイヤーの確認
            if ICMP in self.packet:
                available_layers.append('icmp')
                layer_info['icmp'] = list(self.packet[ICMP].fields.keys())
                highest_layer = 'ICMP'
            
            # DNSレイヤーの確認
            if DNS in self.packet:
                available_layers.append('dns')
                layer_info['dns'] = list(self.packet[DNS].fields.keys())
                highest_layer = 'DNS'
            
            # HTTPとHTTPSの識別（アプリケーションレイヤー）
            if TCP in self.packet:
                tcp_layer = self.packet[TCP]
                if tcp_layer.dport == 80 or tcp_layer.sport == 80:
                    available_layers.append('http')
                    highest_layer = 'HTTP'
                if tcp_layer.dport == 443 or tcp_layer.sport == 443:
                    available_layers.append('ssl')  # HTTPSはSSL/TLSを使用
                    highest_layer = 'HTTPS'
            
            # X.25レイヤーの確認：scapyにX.25がない場合はパケットsummaryから検出する
            if X25 is None:
                if 'x.25' in self.packet.summary().lower() or 'x25' in self.packet.summary().lower():
                    available_layers.append('x25')
                    layer_info['x25'] = ['x.25']
                    highest_layer = 'X.25'
            else:
                if X25 in self.packet:
                    available_layers.append('x25')
                    layer_info['x25'] = list(self.packet[X25].fields.keys())
                    highest_layer = 'X.25'
            
            # Rawパケットの処理
            if self.packet.haslayer(Raw):
                available_layers.append('raw')
                layer_info['raw'] = ['Raw payload']
                highest_layer = 'RAW'
            
            # ペイロードの処理
            if hasattr(self.packet, 'payload') and self.packet.payload:
                # ここでは簡易的な処理。実際には深い解析が必要かも
                pass
            
            # サポート済みプロトコルリストの用意
            supported_packets = ['ETHERNET', 'ARP', 'IP', 'IPv6', 'TCP', 'UDP', 'ICMP',
                                  'DNS', 'HTTP', 'HTTPS', 'X.25', 'RAW']
                              
            # サポート済みプロトコルに含まれない場合、パケットのsummaryからプロトコル名を取得してhighest_layerに設定
            if highest_layer not in supported_packets:
                highest_layer = self.packet.summary().split()[0]  # プロトコル名を取得
            
            available = {
                'layer_names': available_layers,
                'layer_details': layer_info,
                'highest_layer': highest_layer
            }
            
            return available
            
        except Exception as e:
            print(f"使用可能なレイヤー判別中にエラーが発生しました: {str(e)}")
            import traceback
            traceback.print_exc()
            return None
