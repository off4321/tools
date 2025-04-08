from src.protocol_analyzer.base import ProtocolAnalyzer

class AnalyzeUdp(ProtocolAnalyzer):
    """
    UDP内の詳細情報を解析するクラス
    """
    
    def analyze(self):
        if not hasattr(self.packet, 'udp'):
            return None
        try:
            udp_info = {}
            if hasattr(self.packet.udp, 'srcport'):
                udp_info['src_port'] = self.packet.udp.srcport
            if hasattr(self.packet.udp, 'dstport'):
                udp_info['dst_port'] = self.packet.udp.dstport
            if hasattr(self.packet.udp, 'length'):
                udp_info['length'] = self.packet.udp.length
            return udp_info
        except Exception as e:
            print(f"UDP解析中にエラーが発生しました: {str(e)}")
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
        return f"UDP {info.get('src_port','?')} -> {info.get('dst_port','?')}"