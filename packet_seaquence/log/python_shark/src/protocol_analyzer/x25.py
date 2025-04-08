#!/usr/bin/env python3
"""
X.25プロトコル分析クラスを提供するモジュール
"""

from src.protocol_analyzer.base import ProtocolAnalyzer

class AnalyzeX25(ProtocolAnalyzer):
    """
    X.25プロトコルを分析するクラス
    """
    
    def __init__(self, packet):
        """
        AnalyzeX25クラスの初期化
        
        Args:
            packet: pysharkパケットオブジェクト
        """
        self.packet = packet
    
    def analyze(self):
        """
        X.25プロトコルを分析する
        
        Returns:
            dict: X.25情報を含む辞書
        """
        try:
            x25_info = {}
            
            if hasattr(self.packet, 'x25'):
                x25_layer = self.packet.x25
                
                # 基本情報の抽出
                if hasattr(x25_layer, 'lcn'):
                    x25_info['lcn'] = x25_layer.lcn  # 論理チャネル番号
                
                if hasattr(x25_layer, 'type'):
                    x25_info['packet_type'] = x25_layer.type
                    x25_info['packet_type_desc'] = self._get_packet_type_desc(x25_layer.type)
                
                # フラグと制御情報
                if hasattr(x25_layer, 'm'):
                    x25_info['more_data_flag'] = x25_layer.m == '1'  # 後続データフラグ
                
                if hasattr(x25_layer, 'p'):
                    x25_info['p_bit'] = x25_layer.p  # P(oll)ビット
                
                if hasattr(x25_layer, 'q'):
                    x25_info['q_bit'] = x25_layer.q  # 修飾子ビット
                
                if hasattr(x25_layer, 'd'):
                    x25_info['delivery_confirmation'] = x25_layer.d == '1'  # 配信確認
                
                # 送受信シーケンス番号
                if hasattr(x25_layer, 'ps'):
                    x25_info['send_seq'] = x25_layer.ps
                
                if hasattr(x25_layer, 'pr'):
                    x25_info['recv_seq'] = x25_layer.pr
                
                # ファシリティ情報（あれば）
                if hasattr(x25_layer, 'facilities'):
                    x25_info['facilities'] = x25_layer.facilities
                
                # 原因コードとダイアグノスティックコード（クリアやリセットの場合）
                if hasattr(x25_layer, 'cause'):
                    x25_info['cause_code'] = x25_layer.cause
                    x25_info['cause_desc'] = self._get_cause_desc(x25_layer.cause)
                
                if hasattr(x25_layer, 'diagnostic'):
                    x25_info['diagnostic_code'] = x25_layer.diagnostic
            
            return x25_info
            
        except Exception as e:
            print(f"X.25プロトコル分析中にエラーが発生しました: {str(e)}")
            return {}
    
    def _get_packet_type_desc(self, packet_type):
        """
        X.25パケットタイプの説明を取得
        
        Args:
            packet_type: パケットタイプ値
        
        Returns:
            str: パケットタイプの説明
        """
        packet_types = {
            '01': 'Call Request',
            '02': 'Call Accepted',
            '03': 'Clear Request',
            '04': 'Clear Confirmation',
            '05': 'Data',
            '06': 'RR (Receive Ready)',
            '07': 'RNR (Receive Not Ready)',
            '08': 'REJ (Reject)',
            '09': 'Reset Request',
            '0a': 'Reset Confirmation',
            '0b': 'Restart Request',
            '0c': 'Restart Confirmation',
            '0d': 'Interrupt',
            '0e': 'Interrupt Confirmation',
            '0f': 'Registration Request',
            '10': 'Registration Confirmation'
        }
        
        return packet_types.get(packet_type.lower(), 'Unknown')
    
    def _get_cause_desc(self, cause_code):
        """
        X.25原因コードの説明を取得
        
        Args:
            cause_code: 原因コード
        
        Returns:
            str: 原因コードの説明
        """
        cause_codes = {
            '00': 'DTE Originated',
            '01': 'Number Busy',
            '03': 'Invalid Facility Request',
            '05': 'Network Congestion',
            '09': 'Out of Order',
            '0b': 'Access Barred',
            '0d': 'Not Obtainable',
            '11': 'Remote Procedure Error',
            '13': 'Local Procedure Error',
            '15': 'RPOA Out of Order',
            '19': 'Reverse Charging Acceptance Not Subscribed',
            '21': 'Incompatible Destination',
            '29': 'Fast Select Acceptance Not Subscribed',
            '39': 'Ship Absent'
        }
        
        return cause_codes.get(cause_code.lower(), 'Unknown')
    
    def get_display_info(self):
        info = self.analyze()
        if not info:
            return "X.25情報なし"
        dlci = info.get('dlci', '?')
        ptype = info.get('packet_type', '?')
        return f"X.25 (DLCI={dlci}, Type={ptype})"

    def get_summary(self):
        info = self.analyze()
        if not info:
            return "X.25"
        return f"X.25 DLCI={info.get('dlci', '?')}"