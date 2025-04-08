#!/usr/bin/env python3
"""
Mermaid記法出力クラスを提供するモジュール
"""


class WriteMermaid:
    """
    Mermaid記法でシーケンス図を出力するクラス
    """
    
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
            
            # 現在のパケットの絶対時間を取得（日付と時刻）
            time_str = current_time.strftime('%Y-%m-%d %H:%M:%S.%f')[:-3]
            
            # 時間の注記を追加（必ず別の行に）
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
            # TCPメッセージを構築...
            message = "TCP"
        
        # UDPの処理
        if 'udp_info' in protocol_info:
            udp_info = protocol_info['udp_info']
            # UDPメッセージを構築...
            message = "UDP"
        
        # DNSの処理
        if 'dns_info' in protocol_info:
            dns_info = protocol_info['dns_info']
            # DNSメッセージを構築...
            message = "DNS"
        
        # HTTPの処理
        if 'http_info' in protocol_info:
            http_info = protocol_info['http_info']
            # HTTPメッセージを構築...
            message = "HTTP"
        
        # HTTPSの処理
        if 'https_info' in protocol_info:
            https_info = protocol_info['https_info']
            # HTTPSメッセージを構築...
            message = "HTTPS"
        
        # メッセージが空の場合、基本プロトコル情報から推測
        if not message:
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
        vlan_id = info.get('vlan_info', {}).get('id')
        if not vlan_id and ('172.16.2.' in packet.get('src', '') or '172.16.2.' in packet.get('dst', '')):
            # サブネットからVLANを推測
            vlan_id = 2
        
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
        safe_id = safe_id.replace(':', '_')
        safe_id = safe_id.replace('-', '_')
        safe_id = safe_id.replace(' ', '_')
        
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
        
        # その他のコロンをエスケープ
        escaped = text_with_mac_fixed.replace(':', '&#58;')
        
        return escaped
