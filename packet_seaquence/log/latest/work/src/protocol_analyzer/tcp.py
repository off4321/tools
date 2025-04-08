#!/usr/bin/env python3
"""
TCPプロトコル解析クラスを提供するモジュール
"""
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
        if not hasattr(self.packet, 'tcp'):
            return None
        
        try:
            tcp_info = {}
            
            # 送信元ポート
            if hasattr(self.packet.tcp, 'srcport'):
                tcp_info['src_port'] = self.packet.tcp.srcport
            
            # 宛先ポート
            if hasattr(self.packet.tcp, 'dstport'):
                tcp_info['dst_port'] = self.packet.tcp.dstport
            
            # シーケンス番号
            if hasattr(self.packet.tcp, 'seq'):
                tcp_info['seq'] = self.packet.tcp.seq
                if hasattr(self.packet.tcp, 'seq_raw'):
                    tcp_info['seq_raw'] = self.packet.tcp.seq_raw
            
            # 確認応答番号
            if hasattr(self.packet.tcp, 'ack'):
                tcp_info['ack'] = self.packet.tcp.ack
                if hasattr(self.packet.tcp, 'ack_raw'):
                    tcp_info['ack_raw'] = self.packet.tcp.ack_raw
            
            # フラグ
            if hasattr(self.packet.tcp, 'flags'):
                tcp_info['flags'] = self.packet.tcp.flags
                tcp_info['flags_desc'] = self._get_tcp_flags_desc(self.packet.tcp.flags)
            
            # ウィンドウサイズ
            if hasattr(self.packet.tcp, 'window_size'):
                tcp_info['window_size'] = self.packet.tcp.window_size
            elif hasattr(self.packet.tcp, 'window'):
                tcp_info['window_size'] = self.packet.tcp.window
            
            # チェックサム
            if hasattr(self.packet.tcp, 'checksum'):
                tcp_info['checksum'] = self.packet.tcp.checksum
            
            # 緊急ポインタ
            if hasattr(self.packet.tcp, 'urgent_pointer'):
                tcp_info['urgent_pointer'] = self.packet.tcp.urgent_pointer
            
            # オプション
            if hasattr(self.packet.tcp, 'options'):
                tcp_info['options'] = self.packet.tcp.options
            
            # ペイロード
            if hasattr(self.packet.tcp, 'payload'):
                tcp_info['payload'] = self.packet.tcp.payload
            elif hasattr(self.packet.tcp, 'payload_raw'):
                tcp_info['payload_raw'] = self.packet.tcp.payload_raw
            
            # ストリーム番号
            if hasattr(self.packet.tcp, 'stream'):
                tcp_info['stream'] = self.packet.tcp.stream
            
            return tcp_info
            
        except Exception as e:
            print(f"TCP解析中にエラーが発生しました: {str(e)}")
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
            flags (str): TCPフラグ
        
        Returns:
            str: TCPフラグの説明
        """
        descriptions = []
        if 'SYN' in flags:
            descriptions.append('SYN')
        if 'ACK' in flags:
            descriptions.append('ACK')
        if 'FIN' in flags:
            descriptions.append('FIN')
        if 'RST' in flags:
            descriptions.append('RST')
        if 'PSH' in flags:
            descriptions.append('PSH')
        if 'URG' in flags:
            descriptions.append('URG')
        if 'ECE' in flags:
            descriptions.append('ECE')
        if 'CWR' in flags:
            descriptions.append('CWR')
        if 'NS' in flags:
            descriptions.append('NS')
        
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
