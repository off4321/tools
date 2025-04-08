from src.protocol_analyzer.base import ProtocolAnalyzer

class AnalyzeIPv4(ProtocolAnalyzer):
    """
    IPv4内の詳細情報を解析するクラス
    """
    
    def analyze(self):
        if not hasattr(self.packet, 'ip'):
            return None
        try:
            ipv4_info = {}
            if hasattr(self.packet.ip, 'src'):
                ipv4_info['src_ip'] = self.packet.ip.src
            if hasattr(self.packet.ip, 'dst'):
                ipv4_info['dst_ip'] = self.packet.ip.dst
            if hasattr(self.packet.ip, 'ttl'):
                ipv4_info['ttl'] = self.packet.ip.ttl
            if hasattr(self.packet.ip, 'flags'):
                ipv4_info['flags'] = self.packet.ip.flags
            if hasattr(self.packet.ip, 'len'):
                ipv4_info['length'] = self.packet.ip.len
            return ipv4_info
        except Exception as e:
            print(f"IPv4解析中にエラーが発生しました: {str(e)}")
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