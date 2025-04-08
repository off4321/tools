#!/usr/bin/env python3
"""
X.25パケットを擬似的に生成するスクリプト
(scappy標準にはX.25レイヤーがないため、カスタムレイヤーを定義)
"""
from scapy.all import Ether, IP, UDP, Raw, wrpcap
from scapy.packet import Packet
from scapy.fields import ByteField, ShortField, BitField, XByteField, StrField
from .utils import random_mac, random_ip, save_packets

# X.25プロトコルのカスタムレイヤーを定義
class X25(Packet):
    name = "X.25"
    fields_desc = [
        BitField("gfi", 0, 4),              # General Format Identifier
        BitField("lcn", 1, 12),             # Logical Channel Number
        XByteField("type", 0x01),           # Packet Type
        ByteField("p_r", 0),                # Receive Sequence Number
        BitField("more", 0, 1),             # More data
        BitField("p_s", 0, 7),              # Send Sequence Number
    ]

# X.25呼設定/切断用
class X25Control(Packet):
    name = "X25Control"
    fields_desc = [
        ByteField("cause", 0),              # Cause code
        ByteField("diagnostic", 0),         # Diagnostic code
        StrField("facilities", "")          # Facilities
    ]

# X.25データ
class X25Data(Packet):
    name = "X25Data"
    fields_desc = [
        StrField("data", "X.25 payload")
    ]

def create_x25_packets(output_file, count=10, verbose=False):
    """
    擬似的なX.25パケットを生成し、PCAPファイルに保存
    
    Args:
        output_file (str): 保存するファイルパス
        count (int): 生成するパケット数
        verbose (bool): 詳細ログ出力するかどうか
    """
    if verbose:
        print(f"X.25 パケット作成開始: {count}パケット")
    
    packets = []
    
    # 接続端点の情報
    dte_mac = random_mac()
    dce_mac = random_mac()
    dte_ip = "192.168.0.10"
    dce_ip = "192.168.0.1"
    
    # 基本パケット構造 (トンネリング風)
    base_packet = Ether(src=dte_mac, dst=dce_mac) / IP(src=dte_ip, dst=dce_ip) / UDP(sport=1998, dport=1998)
    
    # X.25通信シーケンスの作成 (典型的なセッション)
    # 1. 呼設定 (Call Request)
    call_req = base_packet / Raw(load=bytes(X25(type=0x01, lcn=1) / X25Control(cause=0x00)))
    packets.append(call_req)
    
    if verbose:
        print("X.25 Call Request パケット作成")
    
    # 2. 呼設定受付 (Call Accepted)
    call_acc = Ether(src=dce_mac, dst=dte_mac) / IP(src=dce_ip, dst=dte_ip) / UDP(sport=1998, dport=1998) / \
               Raw(load=bytes(X25(type=0x02, lcn=1) / X25Control(cause=0x00)))
    packets.append(call_acc)
    
    if verbose:
        print("X.25 Call Accepted パケット作成")
    
    # 3. データパケットの交換 (Data)
    data_count = max(2, (count - 6) // 2)  # 少なくとも2つのデータパケット
    
    for i in range(data_count):
        # DTEからDCEへのデータ
        data_dte_to_dce = base_packet / Raw(load=bytes(X25(type=0x05, lcn=1, p_s=i % 8) / X25Data(data=f"Data packet {i+1} from DTE to DCE")))
        packets.append(data_dte_to_dce)
        
        # 受信確認 (RR)
        rr_dce_to_dte = Ether(src=dce_mac, dst=dte_mac) / IP(src=dce_ip, dst=dte_ip) / UDP(sport=1998, dport=1998) / \
                         Raw(load=bytes(X25(type=0x06, lcn=1, p_r=(i+1) % 8)))
        packets.append(rr_dce_to_dte)
        
        # DCEからDTEへのデータ
        data_dce_to_dte = Ether(src=dce_mac, dst=dte_mac) / IP(src=dce_ip, dst=dte_ip) / UDP(sport=1998, dport=1998) / \
                          Raw(load=bytes(X25(type=0x05, lcn=1, p_s=i % 8) / X25Data(data=f"Data packet {i+1} from DCE to DTE")))
        packets.append(data_dce_to_dte)
        
        # 受信確認 (RR)
        rr_dte_to_dce = base_packet / Raw(load=bytes(X25(type=0x06, lcn=1, p_r=(i+1) % 8)))
        packets.append(rr_dte_to_dce)
    
    if verbose:
        print(f"X.25 データパケット {data_count * 2} 組作成")
    
    # 4. クリア要求 (Clear Request)
    clear_req = base_packet / Raw(load=bytes(X25(type=0x03, lcn=1) / X25Control(cause=0x00)))
    packets.append(clear_req)
    
    # 5. クリア確認 (Clear Confirmation)
    clear_conf = Ether(src=dce_mac, dst=dte_mac) / IP(src=dce_ip, dst=dte_ip) / UDP(sport=1998, dport=1998) / \
                 Raw(load=bytes(X25(type=0x04, lcn=1)))
    packets.append(clear_conf)
    
    # エラーシナリオ - リセット要求と確認
    if count > 10:
        reset_req = base_packet / Raw(load=bytes(X25(type=0x09, lcn=1) / X25Control(cause=0x09, diagnostic=0x01)))
        packets.append(reset_req)
        
        reset_conf = Ether(src=dce_mac, dst=dte_mac) / IP(src=dce_ip, dst=dte_ip) / UDP(sport=1998, dport=1998) / \
                     Raw(load=bytes(X25(type=0x0A, lcn=1))))
        packets.append(reset_conf)
    
    if verbose:
        print(f"X.25 パケット作成完了: 合計 {len(packets)} パケット")
    
    # パケットの保存
    save_packets(packets, output_file)
    
    return True

if __name__ == "__main__":
    create_x25_packets("x25_packets.pcap", count=20, verbose=True)
