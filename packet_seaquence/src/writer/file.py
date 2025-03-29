#!/usr/bin/env python3
"""
ファイル出力クラスを提供するモジュール
"""
import os


class WriteFile:
    """
    ファイルに出力するクラス
    複数のファイルに分割して出力する機能を持つ
    """
    
    def __init__(self, output_file):
        """
        WriteFileクラスの初期化
        
        Args:
            output_file (str): 出力ファイルのパス
        """
        self.output_file = output_file
    
    def write(self, content):
        """
        コンテンツをファイルに書き込む
        
        Args:
            content (str): 書き込むコンテンツ
        
        Returns:
            bool: 書き込みが成功したかどうか
        """
        try:
            print(f"ファイルに書き込み中: {self.output_file}")
            
            # 出力ディレクトリが存在しない場合は作成
            output_dir = os.path.dirname(self.output_file)
            if output_dir and not os.path.exists(output_dir):
                os.makedirs(output_dir)
            
            # ファイルに書き込み
            with open(self.output_file, 'w', encoding='utf-8') as f:
                f.write(content)
            
            print(f"ファイルへの書き込みが完了しました: {self.output_file}")
            return True
            
        except Exception as e:
            print(f"ファイルへの書き込み中にエラーが発生しました: {str(e)}")
            import traceback
            traceback.print_exc()
            return False
    
    def write_split(self, contents_list):
        """
        複数のコンテンツを分割して複数のファイルに書き込む
        
        Args:
            contents_list (list): 書き込むコンテンツのリスト
        
        Returns:
            bool: 書き込みが成功したかどうか
        """
        try:
            if not contents_list:
                print("書き込むコンテンツがありません")
                return False
            
            # ファイル名とパスの分離
            output_dir = os.path.dirname(self.output_file)
            filename, ext = os.path.splitext(os.path.basename(self.output_file))
            
            # 出力ディレクトリが存在しない場合は作成
            if output_dir and not os.path.exists(output_dir):
                os.makedirs(output_dir)
            
            # 各コンテンツをファイルに書き込み
            success = True
            for i, content in enumerate(contents_list):
                # ファイル名を生成（例: output_1.md, output_2.md, ...）
                output_file = os.path.join(output_dir, f"{filename}_{i+1}{ext}")
                
                print(f"ファイルに書き込み中: {output_file}")
                
                # ファイルに書き込み
                with open(output_file, 'w', encoding='utf-8') as f:
                    f.write(content)
                
                print(f"ファイルへの書き込みが完了しました: {output_file}")
            
            return success
            
        except Exception as e:
            print(f"ファイルへの書き込み中にエラーが発生しました: {str(e)}")
            import traceback
            traceback.print_exc()
            return False
