from scapy.all import DNS, DNSQR, DNSRR
from src.protocol_analyzer.base import ProtocolAnalyzer

class AnalyzeDns(ProtocolAnalyzer):
    """
    DNS内の詳細情報を解析するクラス
    """

    def analyze(self):
        if DNS not in self.packet:
            return None
        try:
            dns_info = {}
            dns_layer = self.packet[DNS]
            
            # DNS ID
            dns_info['id'] = str(dns_layer.id)
            
            # DNSフラグ
            dns_info['qr'] = "Response" if dns_layer.qr else "Query"
            
            # クエリ情報
            if dns_layer.qd and len(dns_layer.qd) > 0:
                query = dns_layer.qd[0]
                dns_info['query_name'] = query.qname.decode('utf-8', errors='replace').rstrip('.')
                dns_info['query_type'] = self._get_query_type(query.qtype)
            
            # 応答情報
            if dns_layer.an and len(dns_layer.an) > 0:
                answers = []
                for rr in dns_layer.an:
                    answer = {}
                    answer['name'] = rr.rrname.decode('utf-8', errors='replace').rstrip('.')
                    answer['type'] = self._get_query_type(rr.type)
                    if rr.type == 1:  # A record
                        answer['data'] = rr.rdata
                    elif rr.type == 5:  # CNAME
                        answer['data'] = rr.cname.decode('utf-8', errors='replace').rstrip('.')
                    elif rr.type == 28:  # AAAA record
                        answer['data'] = rr.rdata
                    else:
                        answer['data'] = str(rr.rdata)
                    answers.append(answer)
                
                if answers:
                    dns_info['answers'] = answers
                    dns_info['response_name'] = answers[0]['name']
            
            return dns_info
        except Exception as e:
            print(f"DNS解析中にエラーが発生しました: {str(e)}")
            import traceback
            traceback.print_exc()
            return None

    def get_display_info(self):
        info = self.analyze()
        if not info:
            return "DNS情報なし"
        qname = info.get('query_name', '?')
        rname = info.get('response_name', '')
        qr = info.get('qr', '')
        
        if qr == "Response" and rname:
            answers = info.get('answers', [])
            if answers:
                answer_str = ", ".join([f"{a.get('data', '?')}" for a in answers[:2]])
                if len(answers) > 2:
                    answer_str += f" (+ {len(answers)-2} more)"
                return f"DNS Response: {qname} -> {answer_str}"
        
        return f"DNS {qr}: {qname}"

    def get_summary(self):
        info = self.analyze()
        if not info:
            return "DNS"
        qr = info.get('qr', '')
        qname = info.get('query_name', '?')
        qtype = info.get('query_type', '')
        
        return f"DNS {qr} {qname} {qtype}"
    
    def _get_query_type(self, qtype):
        """
        DNSクエリタイプを文字列に変換
        
        Args:
            qtype (int): クエリタイプ番号
        
        Returns:
            str: クエリタイプ名
        """
        dns_types = {
            1: 'A',
            2: 'NS',
            5: 'CNAME',
            6: 'SOA',
            12: 'PTR',
            15: 'MX',
            16: 'TXT',
            28: 'AAAA',
            33: 'SRV',
            41: 'OPT',
            43: 'DS',
            46: 'RRSIG',
            47: 'NSEC',
            48: 'DNSKEY',
            50: 'NSEC3'
        }
        return dns_types.get(qtype, str(qtype))