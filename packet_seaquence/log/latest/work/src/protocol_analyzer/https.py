from src.protocol_analyzer.base import ProtocolAnalyzer

class AnalyzeHttps(ProtocolAnalyzer):
    """
    HTTPS内の詳細情報を解析するクラス
    """

    def analyze(self):
        # HTTPSは暗号化されるため、詳細フィールドを取得できない場合が多い
        if not hasattr(self.packet, 'ssl') and not hasattr(self.packet, 'tls'):
            return None
        try:
            https_info = {}
            if hasattr(self.packet, 'ssl'):
                https_info['ssl_record'] = True
            if hasattr(self.packet, 'tls'):
                https_info['tls_record'] = True
            return https_info
        except Exception as e:
            print(f"HTTPS解析中にエラーが発生しました: {str(e)}")
            return None

    def get_display_info(self):
        info = self.analyze()
        if not info:
            return "HTTPS情報なし"
        return "HTTPS セッション (暗号化)"

    def get_summary(self):
        info = self.analyze()
        if not info:
            return "HTTPS"
        return "HTTPS Session"