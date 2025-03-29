#!/usr/bin/env python3
"""
使用可能なプロトコル判別クラスを提供するモジュール
"""


class DiscriminateAvailable:
    """
    使用可能なプロトコルを判別するクラス
    """
    
    def __init__(self, packet):
        """
        DiscriminateAvailableクラスの初期化
        
        Args:
            packet: pyshark.packet.Packetオブジェクト
        """
        self.packet = packet
    
    def discriminate(self):
        """
        パケットから使用可能なレイヤー情報を判別する
        
        Returns:
            list: 使用可能なレイヤー名のリスト
        """
        try:
            # 利用可能なレイヤー名のリストを取得
            available_layers = [layer.layer_name for layer in self.packet.layers]
            
            # レイヤーの詳細情報を収集
            layer_info = {}
            for layer_name in available_layers:
                try:
                    # 各レイヤーのフィールド名を収集
                    layer = getattr(self.packet, layer_name.lower(), None)
                    if layer:
                        fields = dir(layer)
                        # プライベートフィールドや関数を除外
                        fields = [f for f in fields if not f.startswith('_') and not callable(getattr(layer, f))]
                        layer_info[layer_name] = fields
                except AttributeError:
                    continue
            
            return {
                'layer_names': available_layers,
                'layer_details': layer_info
            }
            
        except Exception as e:
            print(f"使用可能なレイヤー判別中にエラーが発生しました: {str(e)}")
            return None
