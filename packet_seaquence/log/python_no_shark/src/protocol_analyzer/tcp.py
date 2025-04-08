#!/usr/bin/env python3
"""
TCPプロトコル解析クラスを提供するモジュール
"""
from scapy.all import TCP
from src.protocol_analyzer.base import ProtocolAnalyzer


class AnalyzeTcp(ProtocolAnalyzer):
    """
    TCP内の詳細情報を解析するクラス
    """
    
    def analyze(self):
        """
        パケットからTCP情報を解析する
        
        Returns:
            dict: TCP情報を含む辞書、または情報がない場合はNone
        """
        if TCP not in self.packet:
            return None
        
        try:
            tcp_info = {}
            tcp_layer = self.packet[TCP]
            
            # 送信元ポート
            tcp_info['src_port'] = str(tcp_layer.sport)
            
            # 宛先ポート
            tcp_info['dst_port'] = str(tcp_layer.dport)
            
            # シーケンス番号
            tcp_info['seq'] = str(tcp_layer.seq)
            tcp_info['seq_raw'] = str(tcp_layer.seq)
            
            # 確認応答番号
            tcp_info['ack'] = str(tcp_layer.ack)
            tcp_info['ack_raw'] = str(tcp_layer.ack)
            
            # フラグ
            flags_value = tcp_layer.flags
            tcp_info['flags'] = str(flags_value)
            tcp_info['flags_desc'] = self._get_tcp_flags_desc(flags_value)
            
            # ウィンドウサイズ
            tcp_info['window_size'] = str(tcp_layer.window)
            
            # チェックサム
            tcp_info['checksum'] = hex(tcp_layer.chksum) if tcp_layer.chksum is not None else "None"
            
            # 緊急ポインタ
            tcp_info['urgent_pointer'] = str(tcp_layer.urgptr)
            
            # オプション
            if hasattr(tcp_layer, 'options') and tcp_layer.options:
                tcp_info['options'] = str(tcp_layer.options)
            
            # ペイロード
            if hasattr(tcp_layer, 'payload') and tcp_layer.payload:
                tcp_info['payload'] = str(tcp_layer.payload)
            
            return tcp_info
            
        except Exception as e:
            print(f"TCP解析中にエラーが発生しました: {str(e)}")
            import traceback
            traceback.print_exc()
            return None
    
    def get_display_info(self):
        """
        表示用のTCP情報を取得する
        
        Returns:
            str: 表示用のTCP情報
        """
        tcp_info = self.analyze()
        if not tcp_info:
            return "TCP情報なし"
        
        src_port = tcp_info.get('src_port', '不明')
        dst_port = tcp_info.get('dst_port', '不明')
        flags = tcp_info.get('flags_desc', '不明')
        seq = tcp_info.get('seq', 'N/A')
        ack = tcp_info.get('ack', 'N/A')
        
        port_info = f"{src_port} -> {dst_port}"
        flag_info = f"[{flags}]"
        
        if 'SYN' in flags and 'ACK' not in flags:
            return f"TCP {port_info} {flag_info} 接続開始"
        elif 'SYN' in flags and 'ACK' in flags:
            return f"TCP {port_info} {flag_info} 接続確認"
        elif 'FIN' in flags:
            return f"TCP {port_info} {flag_info} 接続終了"
        elif 'RST' in flags:
            return f"TCP {port_info} {flag_info} 接続リセット"
        elif 'ACK' in flags:
            return f"TCP {port_info} {flag_info} 確認応答 (SEQ={seq}, ACK={ack})"
        else:
            return f"TCP {port_info} {flag_info}"
    
    def get_summary(self):
        """
        TCP情報のサマリーを取得する
        
        Returns:
            str: TCP情報のサマリー
        """
        tcp_info = self.analyze()
        if not tcp_info:
            return "TCP"
        
        src_port = tcp_info.get('src_port', '')
        dst_port = tcp_info.get('dst_port', '')
        flags = tcp_info.get('flags_desc', '')
        
        service_info = self._get_service_name_from_port(dst_port)
        
        if service_info:
            return f"TCP {src_port} -> {dst_port} ({service_info}) {flags}"
        else:
            return f"TCP {src_port} -> {dst_port} {flags}"
    
    def _get_tcp_flags_desc(self, flags):
        """
        TCPフラグの説明を取得する
        
        Args:
            flags: TCPフラグ値
        
        Returns:
            str: TCPフラグの説明
        """
        descriptions = []
        # Scapyのフラグは文字列ではなく整数値なので、各フラグのビットをチェック
        flag_bits = {
            'F': 0x01,  # FIN
            'S': 0x02,  # SYN
            'R': 0x04,  # RST
            'P': 0x08,  # PSH
            'A': 0x10,  # ACK
            'U': 0x20,  # URG
            'E': 0x40,  # ECE
            'C': 0x80   # CWR
        }
        
        if isinstance(flags, str):
            # 文字列の場合の処理（互換性のため）
            if 'F' in flags: descriptions.append('FIN')
            if 'S' in flags: descriptions.append('SYN')
            if 'R' in flags: descriptions.append('RST')
            if 'P' in flags: descriptions.append('PSH')
            if 'A' in flags: descriptions.append('ACK')
            if 'U' in flags: descriptions.append('URG')
            if 'E' in flags: descriptions.append('ECE')
            if 'C' in flags: descriptions.append('CWR')
        else:
            # 整数値やその他の場合（Scapyのデフォルト）
            try:
                flag_int = int(flags) if not isinstance(flags, int) else flags
                
                if flag_int & 0x01: descriptions.append('FIN')
                if flag_int & 0x02: descriptions.append('SYN')
                if flag_int & 0x04: descriptions.append('RST')
                if flag_int & 0x08: descriptions.append('PSH')
                if flag_int & 0x10: descriptions.append('ACK')
                if flag_int & 0x20: descriptions.append('URG')
                if flag_int & 0x40: descriptions.append('ECE')
                if flag_int & 0x80: descriptions.append('CWR')
            except (ValueError, TypeError):
                return str(flags)  # 変換できない場合はそのまま返す
        
        return ', '.join(descriptions) if descriptions else 'None'
    
    def _get_service_name_from_port(self, port):
        """
        ポート番号からサービス名を取得する
        
        Args:
            port (str): ポート番号
        
        Returns:
            str: サービス名、または空文字列
        """
        try:
            port_int = int(port)
            well_known_ports = {
                20: 'FTP-DATA',
                21: 'FTP',
                22: 'SSH',
                23: 'Telnet',
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
            return well_known_ports.get(port_int, '')
        except (ValueError, TypeError):
            return ''
