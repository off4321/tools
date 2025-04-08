from scapy.all import TCP, Raw
from src.protocol_analyzer.base import ProtocolAnalyzer

class AnalyzeHttp(ProtocolAnalyzer):
    """
    HTTP内の詳細情報を解析するクラス
    """

    def analyze(self):
        if TCP not in self.packet:
            return None
            
        tcp_layer = self.packet[TCP]
        # HTTP通信の一般的なポートをチェック
        is_http_port = (tcp_layer.dport == 80 or tcp_layer.sport == 80 or
                         tcp_layer.dport == 8080 or tcp_layer.sport == 8080)
                          
        # ペイロードがあるかどうかをチェック
        has_payload = Raw in self.packet
            
        if not (is_http_port and has_payload):
            return None
            
        try:
            http_info = {
                'secure': False,
                'port': tcp_layer.dport if tcp_layer.dport in (80, 8080) else tcp_layer.sport
            }
            
            # HTTPヘッダーの解析を試みる
            if has_payload:
                raw_data = self.packet[Raw].load
                try:
                    # バイナリデータをテキストに変換
                    payload_text = raw_data.decode('utf-8', errors='ignore')
                    
                    # HTTPリクエスト/レスポンスを検出
                    if payload_text.startswith('GET ') or payload_text.startswith('POST ') or \
                       payload_text.startswith('PUT ') or payload_text.startswith('DELETE ') or \
                       payload_text.startswith('HEAD ') or payload_text.startswith('OPTIONS '):
                        http_info['type'] = 'Request'
                        
                        # メソッドを抽出
                        method_end = payload_text.find(' ')
                        if method_end > 0:
                            http_info['method'] = payload_text[:method_end]
                        
                        # URIを抽出
                        uri_start = method_end + 1
                        uri_end = payload_text.find(' ', uri_start)
                        if uri_end > uri_start:
                            http_info['uri'] = payload_text[uri_start:uri_end]
                        
                    elif payload_text.startswith('HTTP/'):
                        http_info['type'] = 'Response'
                        
                        # ステータスコードを抽出
                        status_start = payload_text.find(' ') + 1
                        status_end = payload_text.find(' ', status_start)
                        if status_end > status_start:
                            http_info['status'] = payload_text[status_start:status_end]
                        
                except:
                    # デコードできない場合はバイナリデータとして扱う
                    http_info['binary_data'] = True
            
            return http_info
            
        except Exception as e:
            print(f"HTTP解析中にエラーが発生しました: {str(e)}")
            import traceback
            traceback.print_exc()
            return None

    def get_display_info(self):
        info = self.analyze()
        if not info:
            return "HTTP情報なし"
            
        if info.get('type') == 'Request':
            method = info.get('method', '')
            uri = info.get('uri', '')
            if method and uri:
                return f"HTTP {method} {uri}"
            return "HTTP Request"
            
        if info.get('type') == 'Response':
            status = info.get('status', '')
            if status:
                return f"HTTP Response (Status: {status})"
            return "HTTP Response"
            
        return "HTTP Session"

    def get_summary(self):
        info = self.analyze()
        if not info:
            return "HTTP"
        
        if info.get('type'):
            return f"HTTP {info.get('type')}"
        return "HTTP Session"

# この名前は混乱を避けるために別の場所に移動すべきですが、
# 互換性のために一時的にここに保持します
class AnalyzeHttps(ProtocolAnalyzer):
    """
    HTTPS内の詳細情報を解析するクラス
    """

    def analyze(self):
        if TCP not in self.packet:
            return None
            
        tcp_layer = self.packet[TCP]
        # HTTPS通信のポートをチェック (443または8443が一般的)
        is_https_port = (tcp_layer.dport == 443 or tcp_layer.sport == 443 or 
                          tcp_layer.dport == 8443 or tcp_layer.sport == 8443)
                          
        # 暗号化されたペイロードがあるかどうかをチェック
        has_payload = Raw in self.packet
        
        # TLSレイヤーがあるかどうかをチェック - TLSモジュールがない環境でも動作するように変更
        has_tls = False
            
        if not (is_https_port and has_payload):
            return None
            
        try:
            https_info = {
                'secure': True,
                'port': tcp_layer.dport if tcp_layer.dport in (443, 8443) else tcp_layer.sport
            }
            
            # TLSがない場合はRawペイロードの最初の数バイトを調査してTLSの特徴を判別
            if has_payload:
                raw_data = self.packet[Raw].load
                # TLSハンドシェイクのContent Type (0x16)とバージョン(0x0301, 0x0302, 0x0303)を確認
                if len(raw_data) > 5 and raw_data[0] == 0x16:
                    if raw_data[1] == 0x03:
                        if raw_data[2] == 0x01:
                            https_info['tls_version'] = "TLS 1.0"
                        elif raw_data[2] == 0x02:
                            https_info['tls_version'] = "TLS 1.1"
                        elif raw_data[2] == 0x03:
                            https_info['tls_version'] = "TLS 1.2"
                        else:
                            https_info['tls_version'] = f"TLS/SSL (0x03{raw_data[2]:02x})"
                            
                        # TLSレコードタイプを判別
                        if raw_data[0] == 0x16:
                            https_info['tls_type'] = "Handshake"
                        elif raw_data[0] == 0x14:
                            https_info['tls_type'] = "Change Cipher Spec"
                        elif raw_data[0] == 0x15:
                            https_info['tls_type'] = "Alert"
                        elif raw_data[0] == 0x17:
                            https_info['tls_type'] = "Application Data"
            
            return https_info
            
        except Exception as e:
            print(f"HTTPS解析中にエラーが発生しました: {str(e)}")
            import traceback
            traceback.print_exc()
            return None

    def get_display_info(self):
        info = self.analyze()
        if not info:
            return "HTTPS情報なし"
            
        version = info.get('tls_version', '')
        if version:
            return f"HTTPS セッション ({version}, 暗号化)"
        return "HTTPS セッション (暗号化)"

    def get_summary(self):
        info = self.analyze()
        if not info:
            return "HTTPS"
        
        tls_type = info.get('tls_type', '')
        if tls_type:
            return f"HTTPS {tls_type}"
        return "HTTPS Session"