#!/usr/bin/env python3
"""
パケット生成に使用するユーティリティ関数
"""
import os
import random
import struct
from datetime import datetime, timedelta
from scapy.all import wrpcap

def random_mac():
    """ランダムなMACアドレスを生成"""
    mac = [random.randint(0x00, 0xff) for _ in range(6)]
    # 最初のバイトを偶数にして、ユニキャストアドレスにする
    mac[0] = mac[0] & 0xfe
    return ':'.join('{:02x}'.format(b) for b in mac)

def random_ip(prefix="192.168"):
    """ランダムなIPアドレスを生成"""
    return f"{prefix}.{random.randint(1, 254)}.{random.randint(1, 254)}"

def random_port(min_port=1024, max_port=65535):
    """ランダムなポート番号を生成"""
    return random.randint(min_port, max_port)

def random_time_delta(base_time=None, max_seconds=5):
    """ランダムな時間間隔を生成"""
    if base_time is None:
        base_time = datetime.now()
    delta_seconds = random.uniform(0.001, max_seconds)
    return base_time + timedelta(seconds=delta_seconds)

def save_packets(packets, output_file):
    """パケットをPCAPファイルに保存"""
    # 出力ディレクトリの確認と作成
    output_dir = os.path.dirname(output_file)
    if (output_dir and not os.path.exists(output_dir)):
        os.makedirs(output_dir)
    
    # パケットの保存
    try:
        wrpcap(output_file, packets)
        print(f"パケットを {output_file} に保存しました ({len(packets)} パケット)")
        return True
    except Exception as e:
        print(f"パケットの保存中にエラーが発生しました: {str(e)}")
        return False

def generate_hex_dump(data, bytes_per_line=16):
    """データのHEXダンプを生成"""
    result = []
    for i in range(0, len(data), bytes_per_line):
        chunk = data[i:i + bytes_per_line]
        hex_part = ' '.join(f'{b:02x}' for b in chunk)
        ascii_part = ''.join(chr(b) if 32 <= b <= 126 else '.' for b in chunk)
        result.append(f"{i:04x}:  {hex_part:<{bytes_per_line * 3}}  {ascii_part}")
    return '\n'.join(result)