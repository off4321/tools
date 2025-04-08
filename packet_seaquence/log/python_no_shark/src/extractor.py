#!/usr/bin/env python3
"""
データ抽出クラスを提供するモジュール
"""
from datetime import datetime
from src.models import PacketSequenceData
from src.discriminator.ethernet import DiscriminateEthernet
from src.discriminator.protocol import DiscriminateProtocol
from src.discriminator.mac import DiscriminateMac
from src.discriminator.ip import DiscriminateIp
from scapy.all import conf, Ether, Raw


class ExtractData:
    """
    PCAPファイルからデータを抽出するクラス
    """
    
    def __init__(self, packets, max_entries=50):
        """
        ExtractDataクラスの初期化
        
        Args:
            packets (list): scapy.packet.Packetオブジェクトのリスト
            max_entries (int): 抽出するエントリの最大数
        """
        # 未知のリンク層タイプをEtherにマッピング
        conf.l2types.register(196, Ether)  # 196は未知のリンク層タイプ
        self.packets = packets  # RAWパケットの変換をやめ、そのまま保持
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
            packet: scapyパケットオブジェクト
        
        Returns:
            dict: 抽出したパケットデータ
        """
        # DEBUG: パケットの詳細情報を表示
        print("DEBUG - パケット詳細情報:")
        packet.show()  # Scapyの詳細表示

        # Rawパケットのデコードを試行
        if packet.haslayer(Raw):
            raw_data = packet[Raw].load
            try:
                decoded_data = raw_data.decode('utf-8')  # UTF-8でデコードを試行
                print(f"DEBUG - Rawパケットデコード結果 (UTF-8): {decoded_data}")
            except UnicodeDecodeError:
                try:
                    decoded_data = raw_data.decode('ascii')  # ASCIIでデコードを試行
                    print(f"DEBUG - Rawパケットデコード結果 (ASCII): {decoded_data}")
                except UnicodeDecodeError:
                    # デコード失敗時に16進数表記でデータを表示
                    hex_data = raw_data.hex()
                    print(f"DEBUG - Rawパケットデコード失敗: {hex_data}")
                    # プロトコル固有のデコードを試行
                    self._decode_custom_protocol(raw_data)

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
        protocol = self._determine_protocol(packet, protocol_info, ip_info)
        
        # 送信元と宛先の取得
        src, dst = self._get_src_dst(packet, ip_info, protocol_info, mac_info)
        
        # パケットデータの作成
        packet_data = {
            'src': src,
            'dst': dst,
            'protocol': protocol,
            'protocol_name': protocol,  # プロトコル名を明示的に設定
            'time': timestamp,
            'info': {
                'mac_info': mac_info,
                'ip_info': ip_info,
                'protocol_info': protocol_info
            }
        }
        
        return packet_data

    def _decode_custom_protocol(self, raw_data):
        """
        特定のプロトコルのデコードを試行
        
        Args:
            raw_data: バイナリデータ
        """
        try:
            # 例: X.25プロトコルのデコード処理
            if raw_data.startswith(b'\x00\x1f'):  # X.25のヘッダー例
                print("DEBUG - X.25プロトコルデコード開始")
                # ヘッダー解析例
                header = raw_data[:4]
                payload = raw_data[4:]
                print(f"DEBUG - X.25ヘッダー: {header.hex()}")
                print(f"DEBUG - X.25ペイロード: {payload.hex()}")
            else:
                print("DEBUG - 未知のプロトコル形式")
        except Exception as e:
            print(f"DEBUG - カスタムプロトコルデコード中にエラー: {str(e)}")

    def _determine_protocol(self, packet, protocol_info, ip_info):
        """
        パケットのプロトコルを決定する
        
        Args:
            packet: scapyパケットオブジェクト
            protocol_info: プロトコル情報
            ip_info: IP情報
        
        Returns:
            str: プロトコル名
        """
        if not protocol_info:
            return "UNKNOWN"
        
        # アプリケーションプロトコルを最優先で判別
        if 'dns_info' in protocol_info:
            return "DNS"
        
        if 'http_info' in protocol_info:
            return "HTTP"
        
        if 'https_info' in protocol_info:
            return "HTTPS"
        
        # X.25プロトコルの判別
        if 'x25_info' in protocol_info:
            return "X.25"

        # SCTPプロトコルの判別（優先度を上げる）
        if 'sctp_info' in protocol_info:
            return "SCTP"
        
        # ARPプロトコルの判別
        if 'arp_info' in protocol_info:
            return "ARP"
        
        # TCPトラフィックでポート番号に基づくプロトコル判別
        if 'tcp_info' in protocol_info:
            tcp_info = protocol_info['tcp_info']
            if tcp_info and ('dst_port' in tcp_info or 'src_port' in tcp_info):
                dst_port = tcp_info.get('dst_port', '0')
                src_port = tcp_info.get('src_port', '0')
                if dst_port == '23' or src_port == '23':
                    return "TELNET"
                if dst_port == '21' or src_port == '21':
                    return "FTP"
                if dst_port == '22' or src_port == '22':
                    return "SSH"
                if dst_port == '25' or src_port == '25':
                    return "SMTP"
            return "TCP"
        
        # UDPプロトコルの判別
        if 'udp_info' in protocol_info:
            return "UDP"
        
        # IPプロトコル番号に基づいて判別
        if 'ipv4_info' in protocol_info:
            if ip_info and 'protocol' in ip_info:
                ip_proto = ip_info['protocol']
                if ip_proto == '1':
                    return "ICMP"
                elif ip_proto == '6':  # TCP
                    return "TCP"
                elif ip_proto == '17':  # UDP
                    return "UDP"
                elif ip_proto == '132':  # SCTP
                    return "SCTP"
            return "IPv4"
        
        # 汎用的な未対応プロトコル判定
        # 明示的な判定を削除し、より汎用的なアプローチに変更
        pname = protocol_info.get('protocol_name')
        if pname and pname != "Unknown":
            return f"{pname}[未サポート]"
            
        # パケットのレイヤー情報を取得して判定
        if hasattr(packet, "haslayer") and hasattr(packet, "layers"):
            try:
                # パケットの各レイヤーを調査
                for layer in packet.layers():
                    layer_name = layer.__name__
                    if layer_name not in ["Padding", "Raw", "Ethernet", "IP"]:
                        return f"{layer_name}[未サポート]"
            except Exception as e:
                print(f"レイヤー判定中にエラー: {str(e)}")
        
                    
        return "UNKNOWN"
        
    def _get_packet_layers(self, packet):
        """
        パケット内の全レイヤーのクラス名をリストで返すヘルパー
        
        Args:
            packet: scapyパケットオブジェクト
            
        Returns:
            list: レイヤー名のリスト
        """
        layers = []
        current = packet
        while current:
            layers.append(current.__class__.__name__)
            # payloadがない場合は終了
            if not hasattr(current, "payload") or current.payload is None or current.payload == current:
                break
            current = current.payload
        return layers

    def _extract_timestamp(self, packet):
        """
        パケットからタイムスタンプを抽出する
        
        Args:
            packet: scapyパケットオブジェクト
            
        Returns:
            datetime: タイムスタンプ（日時）
        """
        try:
            # Scapyパケットからタイムスタンプを取得
            if hasattr(packet, 'time'):
                return datetime.fromtimestamp(float(packet.time))  # 修正: floatに変換
            
            # 現在時刻をフォールバックとして使用
            return datetime.now()
            
        except Exception as e:
            print(f"タイムスタンプ抽出中にエラーが発生しました: {str(e)}")
            # エラーが発生した場合は現在時刻を返す
            return datetime.now()

    def _get_src_dst(self, packet, ip_info, protocol_info, mac_info):
        """
        パケットの送信元と宛先を取得する
        
        Args:
            packet: scapyパケットオブジェクト
            ip_info: IP情報
            protocol_info: プロトコル情報
            mac_info: MAC情報
            
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
            
        # MACアドレス情報を使用
        if (not src or not dst) and mac_info:
            src = mac_info.get('src_mac')
            dst = mac_info.get('dst_mac')
            
        # それでも取得できない場合は、パケットから直接取得を試みる
        if (not src or not dst) and hasattr(packet, "src") and hasattr(packet, "dst"):
            src = getattr(packet, "src", None)
            dst = getattr(packet, "dst", None)
            
        return src or "Unknown", dst or "Unknown"

    def _decode_packet(self, packet):
        """
        パケットを適切にデコードする
        
        Args:
            packet: scapyパケットオブジェクト
        
        Returns:
            scapyパケットオブジェクト
        """
        # RAWパケットの変換を削除
        return packet
