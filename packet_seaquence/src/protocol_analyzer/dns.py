from src.protocol_analyzer.base import ProtocolAnalyzer

class AnalyzeDns(ProtocolAnalyzer):
    """
    DNS内の詳細情報を解析するクラス
    """

    def analyze(self):
        if not hasattr(self.packet, 'dns'):
            return None
        try:
            dns_info = {}
            if hasattr(self.packet.dns, 'qry_name'):
                dns_info['query_name'] = self.packet.dns.qry_name
            if hasattr(self.packet.dns, 'resp_name'):
                dns_info['response_name'] = self.packet.dns.resp_name
            if hasattr(self.packet.dns, 'qry_type'):
                dns_info['query_type'] = self.packet.dns.qry_type
            return dns_info
        except Exception as e:
            print(f"DNS解析中にエラーが発生しました: {str(e)}")
            return None

    def get_display_info(self):
        info = self.analyze()
        if not info:
            return "DNS情報なし"
        qname = info.get('query_name', '?')
        rname = info.get('response_name', '')
        if rname:
            return f"DNS {qname} -> {rname}"
        return f"DNS Query: {qname}"

    def get_summary(self):
        info = self.analyze()
        if not info:
            return "DNS"
        return f"DNS {info.get('query_name', '?')}"