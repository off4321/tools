#!/usr/bin/env python3
"""
X.25プロトコル分析クラスを提供するモジュール
"""

from src.protocol_analyzer.base import ProtocolAnalyzer

class AnalyzeX25(ProtocolAnalyzer):
    """
    X.25プロトコルを分析するクラス
    注意: ScapyではX.25の直接サポートが限定的なため、簡易実装
    """
    
    def analyze(self):
        """
        X.25プロトコルを分析する
        
        Returns:
            dict: X.25情報を含む辞書（Scapyでは限定的な情報のみ）
        """
        try:
            # Scapyには標準でX.25の詳細なサポートがないため、
            # 下位レイヤーの情報から推測するか、RAWデータから解析する必要がある
            # ここでは単純にプロトコル名情報のみ返す
            x25_info = {
                "protocol_name": "X.25",
                "note": "Limited support in Scapy"
            }
            
            return x25_info
            
        except Exception as e:
            print(f"X.25プロトコル分析中にエラーが発生しました: {str(e)}")
            import traceback
            traceback.print_exc()
            return {}
    
    def get_display_info(self):
        """
        表示用のX.25情報を取得する
        
        Returns:
            str: 表示用のX.25情報
        """
        return "X.25 プロトコル (詳細情報なし)"

    def get_summary(self):
        """
        X.25情報のサマリーを取得する
        
        Returns:
            str: X.25情報のサマリー
        """
        return "X.25"