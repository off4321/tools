#!/usr/bin/env python3
"""
データ抽出クラスを提供するモジュール
"""
from src.models import PacketSequenceData
from src.discriminator.ethernet import DiscriminateEthernet
from src.discriminator.protocol import DiscriminateProtocol
from src.discriminator.mac import DiscriminateMac
from src.discriminator.ip import DiscriminateIp


class ExtractData:
    """
    PCAPファイルからデータを抽出するクラス
    """
    
    def __init__(self, packets, max_entries=50):
        """
        ExtractDataクラスの初期化
        
        Args:
            packets (list): pyshark.packet.Packetオブジェクトのリスト
            max_entries (int): 抽出するエントリの最大数
        """
        self.packets = packets
        self.max_entries = max_entries
        self.data = PacketSequenceData()
    
    def extract(self):
        """
        パケットからデータを抽出する
        
        パケットから送信元、宛先、プロトコルなどの情報を抽出し、
        PacketSequenceDataオブジェクトに格納する
        
        Returns:
            PacketSequenceData: パケットシーケンス情報を含むデータ構造
        """
        try:
            print("パケットデータ抽出開始...")
            
            # 最大エントリ数まで処理
            count = 0
            for packet in self.packets:
                if count >= self.max_entries:
                    break
                
                # パケットデータを抽出
                packet_info = self._extract_packet_data(packet)
                
                if packet_info:
                    self.data.add_packet(packet_info)
                    count += 1
            
            print(f"データ抽出完了: {count}エントリ")
            return self.data
            
        except Exception as e:
            print(f"データ抽出中にエラーが発生しました: {str(e)}")
            import traceback
            traceback.print_exc()
            raise

    def _extract_packet_data(self, packet):
        """
        1つのパケットからデータを抽出する
        
        Args:
            packet: pysharkパケットオブジェクト
        
        Returns:
            dict: 抽出したパケットデータ
        """
        # プロトコル判別
        discriminator = DiscriminateProtocol(packet)
        protocol_info = discriminator.discriminate()

        print(f"DEBUG_抽出前 -- プロトコル情報: {protocol_info}")

        # プロトコル情報が取得できなかった場合のフォールバック
        if not protocol_info:
            protocol_info = {"protocol_name": "Unknown"}
        
        # MACアドレス情報
        mac_discriminator = DiscriminateMac(packet)
        mac_info = mac_discriminator.discriminate()
        
        # IP情報
        ip_discriminator = DiscriminateIp(packet)
        ip_info = ip_discriminator.discriminate()
        
        # タイムスタンプ
        timestamp = self._extract_timestamp(packet)
        
        # プロトコル名の決定
        protocol = self._determine_protocol(packet, protocol_info, ip_info)  # ip_infoを追加
        
        # 送信元と宛先の取得
        src, dst = self._get_src_dst(ip_info, protocol_info)
        
        # パケットデータの作成
        packet_data = {
            'src': src,
            'dst': dst,
            'protocol': protocol,
            'time': timestamp,
            'info': {
                'mac_info': mac_info,
                'ip_info': ip_info,
                'protocol_info': protocol_info
            }
        }
        
        return packet_data

    def _determine_protocol(self, packet, protocol_info, ip_info):
        """
        パケットのプロトコルを決定する
        
        Args:
            packet: pysharkパケットオブジェクト
            protocol_info: プロトコル情報
            ip_info: IP情報
        
        Returns:
            str: プロトコル名
        """
        # protocol_name が存在すれば最優先で返す
        if 'protocol_name' in protocol_info:
            return protocol_info['protocol_name']

        if not protocol_info:
            return "UNKNOWN"
        
        # プロトコル判別ロジック
        if 'arp_info' in protocol_info:
            return "ARP"
        
        if 'ipv4_info' in protocol_info:
            # IPプロトコル番号に基づいて判別
            if ip_info and 'protocol' in ip_info:
                ip_proto = ip_info['protocol']
                if ip_proto == '1':
                    return "ICMP"
                elif ip_proto == '6':
                    return "TCP"
                elif ip_proto == '17':
                    return "UDP"
            return "IPv4"
        
        if 'tcp_info' in protocol_info:
            return "TCP"
        
        if 'udp_info' in protocol_info:
            return "UDP"
        
        if 'dns_info' in protocol_info:
            return "DNS"
        
        if 'http_info' in protocol_info:
            return "HTTP"
        
        if 'https_info' in protocol_info:
            return "HTTPS"
        
        # X.25プロトコルの判別を追加
        if 'x25_info' in protocol_info:
            return "X.25"
        
        # 不明なプロトコル
        return "UNKNOWN"

    def _extract_timestamp(self, packet):
        """
        パケットからタイムスタンプを抽出する
        
        Args:
            packet: pysharkパケットオブジェクト
            
        Returns:
            datetime: タイムスタンプ（日時）
        """
        try:
            from datetime import datetime
            
            # pysharkパケットからタイムスタンプを取得
            if hasattr(packet, 'sniff_time'):
                return packet.sniff_time
            
            # フォールバック: フレームタイムスタンプを使用
            if hasattr(packet, 'frame_info') and hasattr(packet.frame_info, 'time_epoch'):
                epoch_time = float(packet.frame_info.time_epoch)
                return datetime.fromtimestamp(epoch_time)
            
            # 現在時刻をフォールバックとして使用
            return datetime.now()
            
        except Exception as e:
            print(f"タイムスタンプ抽出中にエラーが発生しました: {str(e)}")
            # エラーが発生した場合は現在時刻を返す
            from datetime import datetime
            return datetime.now()

    def _get_src_dst(self, ip_info, protocol_info):
        """
        パケットの送信元と宛先を取得する
        
        Args:
            ip_info: IP情報
            protocol_info: プロトコル情報
        
        Returns:
            tuple: 送信元と宛先のタプル (src, dst)
        """
        src = None
        dst = None
        
        # IP情報から送信元と宛先を取得
        if ip_info:
            src = ip_info.get('src_ip')
            dst = ip_info.get('dst_ip')
        
        # IPが取得できない場合はARP情報を使用
        if (not src or not dst) and protocol_info and 'arp_info' in protocol_info:
            arp_info = protocol_info['arp_info']
            src = arp_info.get('src_ip')
            dst = arp_info.get('dst_ip')
        
        # それでも取得できない場合、MACアドレスを使用したり他の方法で取得を試みる
        # この部分は必要に応じて拡張
        
        return src or "Unknown", dst or "Unknown"
