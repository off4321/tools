#!/usr/bin/env python3
"""
Mermaid記法出力クラスを提供するモジュール
"""


class WriteMermaid:
    """
    Mermaid記法でシーケンス図を出力するクラス
    """
    
    # サポートしているプロトコルのリスト
    supported_protocols = ['X.25', 'IPv4', 'ARP', 'TCP', 'UDP', 'DNS', 'HTTP', 'HTTPS', 'ICMP']
    
    def __init__(self, packet_data):
        """
        WriteMermaidクラスの初期化
        
        Args:
            packet_data (PacketSequenceData): パケットシーケンス情報を含むデータ構造
        """
        self.packet_data = packet_data
    
    def generate(self):
        """
        Mermaid記法のシーケンス図を生成する
        
        Returns:
            str: Mermaid記法のシーケンス図
        """
        try:
            print("Mermaid記法のシーケンス図を生成中...")
            
            # シーケンス図のヘッダー
            mermaid_content = "```mermaid\nsequenceDiagram\n"
            
            # 参加者の定義
            participants = self._generate_participants()
            mermaid_content += participants
            
            # シーケンスの生成
            sequences = self._generate_sequences()
            mermaid_content += sequences
            
            # シーケンス図のフッター
            mermaid_content += "```\n"
            
            print("Mermaid記法のシーケンス図の生成完了")
            return mermaid_content
            
        except Exception as e:
            print(f"Mermaid記法のシーケンス図の生成中にエラーが発生しました: {str(e)}")
            import traceback
            traceback.print_exc()
            return "```mermaid\nsequenceDiagram\n  Note over Error: シーケンス図の生成中にエラーが発生しました\n```\n"
    
    def _generate_participants(self):
        """
        Mermaid記法の参加者（アクター）定義を生成する
        
        Returns:
            str: 参加者定義の文字列
        """
        participants = ""
        
        # 時間列用の特別な参加者を最初に追加
        participants += "    participant Time as \"🕒 Time\"\n"
        
        # ノードのセットから参加者を生成
        nodes = sorted(list(self.packet_data.get_nodes()))
        
        # 各ノードを参加者として定義
        for node in nodes:
            # ノード名が長すぎる場合は短縮する
            safe_id = self._sanitize_id(node)
            if len(node) > 20:
                short_name = node[:8] + "..." + node[-8:] if len(node) > 20 else node
                participants += f"    participant {safe_id} as \"{node}\"\n"
            else:
                participants += f"    participant {safe_id} as \"{node}\"\n"
        
        return participants

    def _generate_sequences(self):
        """
        Mermaid記法のシーケンスを生成する。各パケットに時間表示を追加。
        
        Returns:
            str: シーケンスの文字列
        """
        sequences = ""
        
        # パケット情報のリストからシーケンスを生成
        packets = self.packet_data.get_packets()
        
        for i, packet in enumerate(packets):
            # 送信元、宛先、プロトコル情報を取得
            src = packet.get('src', '')
            dst = packet.get('dst', '')
            protocol = packet.get('protocol', 'UNKNOWN')
            info = packet.get('info', {})
            current_time = packet.get('time')
            
            # タイムゾーン変換をせず、パケットの時刻をそのまま表示
            time_str = current_time.strftime('%Y-%m-%d %H:%M:%S.%f')
            
            sequences += f"    Note right of Time: {time_str}\n"
            sequences += f"    Time-->>Time: 📍\n"
            
            # メッセージの詳細を構築
            message = self._build_message(packet)
            
            # シーケンス矢印を追加（メッセージと時間の注記を分離）
            if src and dst:
                src_id = self._sanitize_id(src)
                dst_id = self._sanitize_id(dst)
                
                # メッセージ本体（プロトコル情報やVLAN情報を含む）
                message_text = self._escape(message) if message else ""
                sequences += f"    {src_id}->>+{dst_id}: {message_text}\n"
    
        return sequences
    
    def _build_message(self, packet):
        """
        パケットの詳細メッセージを構築する
        
        Args:
            packet: パケット情報の辞書
        
        Returns:
            str: メッセージ文字列
        """
        protocol = packet.get('protocol', 'UNKNOWN')
        info = packet.get('info', {})
        protocol_info = info.get('protocol_info', {})
        message = ""

        print(f"DEBUG_buildmessage--プロトコル: {protocol}, プロトコル情報: {protocol_info}")

        # もしprotocolがUNKNOWNでprotocol_infoにプロトコル名があれば上書き
        if protocol != 'protocol_name' in protocol_info:
            protocol = protocol_info['protocol_name']

        # もし上書き後のprotocolがサポート外なら [未サポートプロトコル] を付加して返す
        if protocol != 'UNKNOWN' and protocol not in self.supported_protocols:
            return f"{protocol} [未サポートプロトコル]"
        
        # X.25の処理を追加
        if 'x25_info' in protocol_info:
            x25_info = protocol_info['x25_info']
            packet_type_desc = x25_info.get('packet_type_desc', '')
            lcn = x25_info.get('lcn', '')
            
            if packet_type_desc:
                message = f"X.25 {packet_type_desc}"
                
                # 論理チャネル番号があれば追加
                if lcn:
                    message += f" (LCN:{lcn})"
                
                # 送受信シーケンス番号があれば追加（データパケットの場合）
                if 'send_seq' in x25_info and 'recv_seq' in x25_info:
                    message += f" S:{x25_info['send_seq']}/R:{x25_info['recv_seq']}"
                
                # クリアやリセットの原因があれば追加
                if 'cause_desc' in x25_info:
                    message += f": {x25_info['cause_desc']}"
            else:
                message = "X.25"
        
        # その他の既存のプロトコル処理（変更なし）
        if 'ipv4_info' in protocol_info:
            ipv4_info = protocol_info['ipv4_info']
            ttl = ipv4_info.get('ttl', '')
            flags = ipv4_info.get('flags', '')
            
            # ICMPの処理 (IPプロトコルナンバー1はICMP)
            if info.get('ip_info', {}).get('protocol') == '1':  # 修正箇所: packetからinfoに変更
                message = "ICMP"
                if int(ttl) == 64:
                    # TTLと他の特徴からPingを推測
                    if packet['src'] == '172.16.1.253' or packet['src'] == '172.16.2.254':
                        message = "ICMP Echo Request (Ping)"
                    else:
                        message = "ICMP Echo Reply (Ping)"
            else:
                message = "IPv4"
        
        # ARPの処理
        if 'arp_info' in protocol_info:
            arp_info = protocol_info['arp_info']
            operation = arp_info.get('operation', '')
            
            if operation == 'REQUEST':
                message = f"ARP Request: Who has {arp_info.get('dst_ip', '')}?"
            elif operation == 'REPLY':
                message = f"ARP Reply: {arp_info.get('src_ip', '')} is at {arp_info.get('src_mac', '')}"
            else:
                message = f"ARP {operation}"
        
        # TCPの処理
        if 'tcp_info' in protocol_info:
            tcp_info = protocol_info['tcp_info']
            src_port = tcp_info.get('src_port', '')
            dst_port = tcp_info.get('dst_port', '')
            flags_hex = tcp_info.get('flags', '0x0000')
            
            # フラグの簡易デコード
            flags_map = {
                0x0002: 'SYN',
                0x0012: 'SYN,ACK',
                0x0010: 'ACK',
                0x0011: 'FIN,ACK',
                0x0004: 'RST',
                # ...必要に応じて追加...
            }
            try:
                flags_int = int(flags_hex, 16)
                decoded_flags = flags_map.get(flags_int, flags_hex)
            except ValueError:
                decoded_flags = flags_hex

            options = tcp_info.get('options', '')
            # オプションが空でない場合のみ表示
            if options:
                message = f"TCP {src_port}->{dst_port} (Flags: {decoded_flags}, Options: {options})"
            else:
                message = f"TCP {src_port}->{dst_port} (Flags: {decoded_flags})"
        
        # UDPの処理
        if 'udp_info' in protocol_info:
            udp_info = protocol_info['udp_info']
            src_port = udp_info.get('src_port', '')
            dst_port = udp_info.get('dst_port', '')
            message = f"UDP {src_port}->{dst_port}"
        
        # DNSの処理
        if 'dns_info' in protocol_info:
            dns_info = protocol_info['dns_info']
            domain = dns_info.get('domain', '')
            record_type = dns_info.get('type', '')
            resolved_ips = dns_info.get('resolved_ips', [])
            is_query = dns_info.get('is_query', False)
            error_code = dns_info.get('error_code', 0)
            error_message = dns_info.get('error_message', '')

            if is_query:
                message = f"DNS Query: {domain} (Type: {record_type})"
            else:
                if error_code != 0:  # エラーがある場合
                    message = f"DNS Reply Error: {domain} (Type: {record_type}) - Code: {error_code}, {error_message}"
                else:
                    message = f"DNS Reply: {domain} (Type: {record_type}) => {', '.join(resolved_ips)}"
        
        # HTTPの処理
        if 'http_info' in protocol_info:
            http_info = protocol_info['http_info']
            is_request = http_info.get('is_request', True)  # デフォルトはリクエストと仮定
            version = http_info.get('version', '')
            
            if is_request:
                method = http_info.get('method', '')
                url = http_info.get('url', '')
                message = f"HTTP {method} {url} ({version})"
            else:
                status_code = http_info.get('status_code', '')
                status_message = http_info.get('status_message', '')
                message = f"HTTP Response: {status_code} {status_message} ({version})"
        
        # HTTPSの処理
        if 'https_info' in protocol_info:
            https_info = protocol_info['https_info']
            tls_version = https_info.get('tls_version', '')
            handshake_type = https_info.get('handshake_type', '')
            server_name = https_info.get('server_name', '')  # SNI情報
            
            if handshake_type:
                if server_name:
                    message = f"TLS {handshake_type} ({tls_version}) - SNI: {server_name}"
                else:
                    message = f"TLS {handshake_type} ({tls_version})"
            else:
                message = f"HTTPS/TLS ({tls_version})"
        
        # SCTPの処理
        if 'sctp_info' in protocol_info:
            sctp_info = protocol_info['sctp_info']
            src_port = sctp_info.get('src_port', '')
            dst_port = sctp_info.get('dst_port', '')
            chunks = sctp_info.get('chunks', [])
            chunk_descriptions = []
            for chunk in chunks:
                chunk_type = chunk.get('type', 'Unknown')
                chunk_flags = chunk.get('flags', '0x00')
                chunk_len = chunk.get('len', '0')
                chunk_descriptions.append(f"{chunk_type}({chunk_flags}, len={chunk_len})")
            
            options = sctp_info.get('options', '')  # SCTP Optionsの取得
            message = f"SCTP {src_port}->{dst_port} Chunks: {', '.join(chunk_descriptions)} Options: {options}"
        
        # メッセージが空の場合、プロトコル名を取得
        if not message:
            # 1. プロトコル情報から直接名前を取得
            if 'protocol_name' in protocol_info:
                proto_name = protocol_info['protocol_name']
                # サポート外ならタグ追加
                if proto_name not in self.supported_protocols:
                    message = f"{proto_name} [未サポートプロトコル]"
                else:
                    message = proto_name
            # 2. パケットのプロトコルフィールドを使用
            elif protocol != 'UNKNOWN':
                message = f"{protocol}"
            # 3. 既存のフォールバック処理
            else:
                ether_type = packet.get('mac_info', {}).get('ether_type', '')
                ip_proto = packet.get('ip_info', {}).get('protocol', '')
                
                if ether_type == '0x0800':  # IPv4
                    if ip_proto == '1':  # ICMP
                        message = "ICMP"
                    elif ip_proto == '6':  # TCP
                        message = "TCP"
                    elif ip_proto == '17':  # UDP
                        message = "UDP"
                    else:
                        message = f"IP Protocol {ip_proto}"
                elif ether_type == '0x0806':  # ARP
                    message = "ARP"
                elif ether_type == '0x86dd':  # IPv6
                    message = "IPv6"
                else:
                    message = f"EtherType {ether_type}"
        
        # VLAN情報があれば追加
        vlan_id = None
        # 明示的なVLAN情報があるか確認
        if 'vlan_info' in info and info['vlan_info'].get('id'):
            vlan_id = info['vlan_info'].get('id')
        # mac_infoからVLAN情報を取得
        if 'vlan_info' in info and info['vlan_info'].get('id'):
            vlan_id = info['vlan_info'].get('id')
        # mac_infoからVLAN情報を取得
        elif 'mac_info' in info and info['mac_info'] is not None and info['mac_info'].get('vlan_id'):
            vlan_id = info['mac_info'].get('vlan_id')
        
        if vlan_id:
            message += f" (VLAN {vlan_id})"
        
        return message
    
    def _sanitize_id(self, text):
        """
        Mermaid記法で使用できる識別子に変換する
        
        Args:
            text (str): 変換する文字列
        
        Returns:
            str: 安全な識別子
        """
        # 特殊文字を置換
        safe_id = text.replace('.', '_')
        safe_id = text.replace(':', '_')
        safe_id = text.replace('-', '_')
        safe_id = text.replace(' ', '_')
        
        # 数字から始まる場合、先頭に'n'を追加
        if safe_id and safe_id[0].isdigit():
            safe_id = 'n' + safe_id
        
        return safe_id
    
    def _escape(self, text):
        """
        Mermaid記法でエスケープが必要な文字をエスケープする
        MACアドレスのコロンはハイフンに変換して見やすく
        
        Args:
            text (str): エスケープする文字列
        
        Returns:
            str: エスケープされた文字列
        """
        # MACアドレスパターン
        import re
        mac_pattern = re.compile(r'([0-9a-fA-F]{2}:){5}[0-9a-fA-F]{2}')
        
        # MACアドレス内のコロンをハイフンに変換
        def replace_mac(match):
            return match.group(0).replace(':', '-')
        
        text_with_mac_fixed = re.sub(mac_pattern, replace_mac, text)
        
        return text_with_mac_fixed
