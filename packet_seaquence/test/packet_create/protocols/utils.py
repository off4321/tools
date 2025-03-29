#!/usr/bin/env python3
"""
パケット生成に使用するユーティリティ関数
"""
import os
import random
import socket
from scapy.all import wrpcap

def random_mac():
    """ランダムなMACアドレスを生成"""
    mac = [random.randint(0x00, 0xff) for _ in range(6)]
    # 最初のバイトを偶数にして、ユニキャストアドレスにする
    mac[0] = mac[0] & 0xfe
    return ':'.join('{:02x}'.format(b) for b in mac)

def random_ip(prefix="192.168"):
    """
    ランダムなIPアドレスを生成
    
    Args:
        prefix: IPアドレスのプレフィックス (例: "192.168" -> "192.168.x.y")
                または単一IPアドレス (例: "10.0.0.1")
    
    Returns:
        str: 有効なIPアドレス文字列
    """
    # プレフィックスの構造を確認
    parts = prefix.split('.')
    
    if len(parts) == 4:
        # すでに完全なIPアドレス
        return prefix
    elif len(parts) == 3:
        # クラスCネットワーク (例: "192.168.1")
        return f"{prefix}.{random.randint(1, 254)}"
    elif len(parts) == 2:
        # クラスBネットワーク (例: "192.168")
        return f"{prefix}.{random.randint(1, 254)}.{random.randint(1, 254)}"
    elif len(parts) == 1:
        # クラスAネットワーク (例: "10")
        return f"{prefix}.{random.randint(1, 254)}.{random.randint(1, 254)}.{random.randint(1, 254)}"
    else:
        # デフォルト
        return f"192.168.{random.randint(1, 254)}.{random.randint(1, 254)}"

def random_port(min_port=1024, max_port=65535):
    """ランダムなポート番号を生成"""
    return random.randint(min_port, max_port)

def save_packets(packets, output_file):
    """パケットをPCAPファイルに保存"""
    # 出力ディレクトリの確認と作成
    output_dir = os.path.dirname(output_file)
    if output_dir and not os.path.exists(output_dir):
        os.makedirs(output_dir)
    
    # パケットの保存
    try:
        wrpcap(output_file, packets)
        print(f"パケットを {output_file} に保存しました ({len(packets)} パケット)")
        return True
    except Exception as e:
        print(f"パケットの保存中にエラーが発生しました: {str(e)}")
        return False

def validate_ip(ip):
    """IPアドレスが有効かどうかを検証"""
    try:
        socket.inet_aton(ip)
        return True
    except:
        return False