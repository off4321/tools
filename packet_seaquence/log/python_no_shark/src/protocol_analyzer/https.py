from scapy.all import TCP, Raw
from src.protocol_analyzer.base import ProtocolAnalyzer

# Try to import TLS, but don't fail if it's not available
TLS = None
try:
    from scapy.all import TLS
except ImportError:
    pass

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
        
        # TLSレイヤーがあるかどうかをチェック
        has_tls = False
        if TLS is not None:
            try:
                has_tls = TLS in self.packet
            except:
                pass
            
        if not (is_https_port and (has_payload or has_tls)):
            return None
            
        try:
            https_info = {
                'secure': True,
                'port': tcp_layer.dport if tcp_layer.dport in (443, 8443) else tcp_layer.sport
            }
            
            # TLSバージョン情報を取得（可能な場合）
            if has_tls and TLS is not None:
                tls_layer = self.packet[TLS]
                if hasattr(tls_layer, 'version'):
                    https_info['tls_version'] = self._get_tls_version(tls_layer.version)
                if hasattr(tls_layer, 'type'):
                    https_info['tls_type'] = self._get_tls_type(tls_layer.type)
            
            # If TLS is not available, try to detect TLS from raw data
            elif has_payload:
                raw_data = self.packet[Raw].load
                if len(raw_data) > 5 and raw_data[0] == 0x16:  # Handshake message type
                    if raw_data[1] == 0x03:  # SSL/TLS version (major)
                        if raw_data[2] == 0x01:
                            https_info['tls_version'] = "TLS 1.0"
                        elif raw_data[2] == 0x02:
                            https_info['tls_version'] = "TLS 1.1"
                        elif raw_data[2] == 0x03:
                            https_info['tls_version'] = "TLS 1.2"
                        elif raw_data[2] == 0x04:
                            https_info['tls_version'] = "TLS 1.3"
                        else:
                            https_info['tls_version'] = f"TLS/SSL (0x03{raw_data[2]:02x})"
                        
                        # Record type
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
        
    def _get_tls_version(self, version):
        """TLSバージョンを人間が読める形式に変換"""
        versions = {
            0x0301: "TLS 1.0",
            0x0302: "TLS 1.1",
            0x0303: "TLS 1.2", 
            0x0304: "TLS 1.3"
        }
        return versions.get(version, f"TLS/SSL (0x{version:04x})")
        
    def _get_tls_type(self, type_val):
        """TLSレコードタイプを人間が読める形式に変換"""
        types = {
            20: "Change Cipher Spec",
            21: "Alert",
            22: "Handshake",
            23: "Application Data"
        }
        return types.get(type_val, f"Type {type_val}")