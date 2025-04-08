from scapy.all import IP
from src.protocol_analyzer.base import ProtocolAnalyzer

class AnalyzeIPv4(ProtocolAnalyzer):
    """
    IPv4内の詳細情報を解析するクラス
    """
    
    def analyze(self):
        if IP not in self.packet:
            return None
        try:
            ipv4_info = {}
            ip_layer = self.packet[IP]
            
            ipv4_info['src_ip'] = ip_layer.src
            ipv4_info['dst_ip'] = ip_layer.dst
            ipv4_info['ttl'] = str(ip_layer.ttl)
            ipv4_info['flags'] = str(ip_layer.flags)
            ipv4_info['length'] = str(ip_layer.len)
            ipv4_info['protocol'] = str(ip_layer.proto)
            
            return ipv4_info
        except Exception as e:
            print(f"IPv4解析中にエラーが発生しました: {str(e)}")
            import traceback
            traceback.print_exc()
            return None
    
    def get_display_info(self):
        info = self.analyze()
        if not info:
            return "IPv4情報なし"
        return f"IPv4 {info.get('src_ip','?')} -> {info.get('dst_ip','?')} (TTL={info.get('ttl','?')})"
    
    def get_summary(self):
        info = self.analyze()
        if not info:
            return "IPv4"
        return f"IPv4 {info.get('src_ip','?')} -> {info.get('dst_ip','?')}"