#!/usr/bin/env python3
"""
save_packets関数呼び出しを修正するスクリプト
"""
import os
import re
import sys

def fix_file(file_path):
    """save_packets呼び出しを修正"""
    try:
        with open(file_path, 'r') as f:
            content = f.read()
        
        # save_packets(packets, output_file, verbose) を save_packets(packets, output_file) に置き換え
        pattern = r'save_packets\s*\(\s*packets\s*,\s*output_file\s*,\s*verbose\s*\)'
        replacement = 'save_packets(packets, output_file)'
        
        new_content = re.sub(pattern, replacement, content)
        
        # 変更があった場合のみ書き込み
        if new_content != content:
            with open(file_path, 'w') as f:
                f.write(new_content)
            return True
        
        return False
    except Exception as e:
        print(f"エラー: {file_path} の処理中にエラーが発生しました: {str(e)}")
        return False

def main():
    protocols_dir = "protocols"
    if not os.path.exists(protocols_dir):
        print(f"エラー: {protocols_dir} ディレクトリが見つかりません")
        return
    
    fixed_count = 0
    
    # プロトコルファイルを探して処理
    for filename in os.listdir(protocols_dir):
        if filename.endswith("_packets.py"):
            file_path = os.path.join(protocols_dir, filename)
            print(f"ファイル {filename} を処理中...")
            
            if fix_file(file_path):
                print(f"  ✅ {filename} を修正しました")
                fixed_count += 1
            else:
                print(f"  ⏭️ {filename} に変更はありませんでした")
    
    print(f"\n処理完了: {fixed_count} ファイルを修正しました")

if __name__ == "__main__":
    main()