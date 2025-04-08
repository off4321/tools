from src.protocol_analyzer.base import ProtocolAnalyzer

class AnalyzeSctp(ProtocolAnalyzer):
    """
    SCTP内の詳細情報を解析するクラス
    """

    def analyze(self):
        if not hasattr(self.packet, 'sctp'):
            return None
        try:
            sctp_info = {}
            if hasattr(self.packet.sctp, 'srcport'):
                sctp_info['src_port'] = self.packet.sctp.srcport
            if hasattr(self.packet.sctp, 'dstport'):
                sctp_info['dst_port'] = self.packet.sctp.dstport
            if hasattr(self.packet.sctp, 'verification_tag'):
                sctp_info['verification_tag'] = self.packet.sctp.verification_tag
            return sctp_info
        except Exception as e:
            print(f"SCTP解析中にエラーが発生しました: {str(e)}")
            return None

    def get_display_info(self):
        info = self.analyze()
        if not info:
            return "SCTP情報なし"
        return f"SCTP {info.get('src_port','?')} -> {info.get('dst_port','?')} (VTAG={info.get('verification_tag','?')})"

    def get_summary(self):
        info = self.analyze()
        if not info:
            return "SCTP"
        return f"SCTP {info.get('src_port','?')} -> {info.get('dst_port','?')}"