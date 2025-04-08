from src.protocol_analyzer.base import ProtocolAnalyzer

class AnalyzeHttp(ProtocolAnalyzer):
    """
    HTTP内の詳細情報を解析するクラス
    """

    def analyze(self):
        if not hasattr(self.packet, 'http'):
            return None
        try:
            http_info = {}
            if hasattr(self.packet.http, 'request_line'):
                http_info['request_line'] = self.packet.http.request_line
            if hasattr(self.packet.http, 'request_method'):
                http_info['method'] = self.packet.http.request_method
            if hasattr(self.packet.http, 'host'):
                http_info['host'] = self.packet.http.host
            if hasattr(self.packet.http, 'request_uri'):
                http_info['uri'] = self.packet.http.request_uri
            return http_info
        except Exception as e:
            print(f"HTTP解析中にエラーが発生しました: {str(e)}")
            return None

    def get_display_info(self):
        info = self.analyze()
        if not info:
            return "HTTP情報なし"
        line = info.get('request_line', '')
        method = info.get('method', '')
        uri = info.get('uri', '')
        if method and uri:
            return f"HTTP {method} {uri}"
        if line:
            return f"HTTP {line}"
        return "HTTP"

    def get_summary(self):
        info = self.analyze()
        if not info:
            return "HTTP"
        return f"HTTP {info.get('method','?')} {info.get('uri','?')}"