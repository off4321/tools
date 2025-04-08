from scapy.all import UDP
from src.protocol_analyzer.base import ProtocolAnalyzer

class AnalyzeUdp(ProtocolAnalyzer):
    """
    UDP内の詳細情報を解析するクラス
    """
    
    def analyze(self):
        if UDP not in self.packet:
            return None
        try:
            udp_info = {}
            udp_layer = self.packet[UDP]
            
            # 送信元ポート
            udp_info['src_port'] = str(udp_layer.sport)
            
            # 宛先ポート
            udp_info['dst_port'] = str(udp_layer.dport)
            
            # UDPデータグラム長
            udp_info['length'] = str(udp_layer.len) if hasattr(udp_layer, 'len') else str(len(udp_layer))
            
            # チェックサム
            udp_info['checksum'] = hex(udp_layer.chksum) if hasattr(udp_layer, 'chksum') and udp_layer.chksum is not None else "None"
            
            return udp_info
        except Exception as e:
            print(f"UDP解析中にエラーが発生しました: {str(e)}")
            import traceback
            traceback.print_exc()
            return None
    
    def get_display_info(self):
        info = self.analyze()
        if not info:
            return "UDP情報なし"
        return f"UDP {info.get('src_port','?')} -> {info.get('dst_port','?')} (len={info.get('length','?')})"
    
    def get_summary(self):
        info = self.analyze()
        if not info:
            return "UDP"
        
        src_port = info.get('src_port','?')
        dst_port = info.get('dst_port','?')
        
        # 既知のUDPサービスの判別
        service = self._get_service_name_from_port(dst_port)
        if service:
            return f"UDP {src_port} -> {dst_port} ({service})"
        return f"UDP {src_port} -> {dst_port}"
    
    def _get_service_name_from_port(self, port):
        """
        ポート番号からUDPサービス名を取得する
        
        Args:
            port (str): ポート番号
            
        Returns:
            str: サービス名、または空文字列
        """
        try:
            port_int = int(port)
            udp_services = {
                53: 'DNS',
                67: 'DHCP-Server',
                68: 'DHCP-Client',
                69: 'TFTP',
                123: 'NTP',
                137: 'NetBIOS-NS',
                138: 'NetBIOS-DGM',
                161: 'SNMP',
                162: 'SNMP-Trap',
                514: 'Syslog',
                1900: 'SSDP',
                5353: 'mDNS'
            }
            return udp_services.get(port_int, '')
        except (ValueError, TypeError):
            return ''