#!/usr/bin/env python3
"""
HTTPパケットを作成するモジュール - より詳細な情報を含むHTTPパケットを生成
"""
from scapy.all import Ether, IP, TCP, Raw
from .utils import random_mac, random_ip, random_port, save_packets
import random

def create_http_packets(output_file, count=10, verbose=False):
    """HTTPパケットの作成 - より詳細なHTTPリクエスト/レスポンスを含む"""
    if verbose:
        print(f"HTTP パケット作成開始: {count}パケット")
    
    packets = []
    
    # HTTPリクエスト種別
    http_methods = ["GET", "POST", "PUT", "DELETE", "HEAD"]
    
    # URIパス
    uri_paths = [
        "/", 
        "/index.html", 
        "/api/users",
        "/products/123",
        "/images/logo.png",
        "/login",
        "/dashboard",
        "/search?q=example",
        "/blog/2025/03/latest-news",
        "/contact-us"
    ]
    
    # User-Agent
    user_agents = [
        "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36",
        "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/15.0 Safari/605.1.15",
        "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/92.0.4515.107 Safari/537.36",
        "Mozilla/5.0 (iPhone; CPU iPhone OS 14_6 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Mobile/15E148"
    ]
    
    # HTTPリクエストとレスポンスのペアを作成
    for i in range(count // 2):
        client_mac = random_mac()
        server_mac = random_mac()
        client_ip = random_ip("192.168.1")
        server_ip = random_ip("10.0.0") 
        client_port = random_port(1024, 65000)
        server_port = 80
        seq_base = 1000 * (i + 1)
        
        # ランダムなHTTPメソッドとURIを選択
        http_method = random.choice(http_methods)
        uri_path = random.choice(uri_paths)
        domain = f"example-{i}.com" if i % 3 != 0 else f"api.service-{i}.org"
        user_agent = random.choice(user_agents)
        
        # POSTリクエストの場合はデータを含める
        post_data = ""
        content_type = ""
        if http_method == "POST":
            content_type = "application/json" if random.random() > 0.5 else "application/x-www-form-urlencoded"
            if content_type == "application/json":
                post_data = f'{{"id": {i+1}, "name": "Test User {i}", "action": "update"}}'
            else:
                post_data = f"id={i+1}&name=Test+User+{i}&action=update"
        
        # HTTPリクエスト
        http_request = f"{http_method} {uri_path} HTTP/1.1\r\n"
        http_request += f"Host: {domain}\r\n"
        http_request += f"User-Agent: {user_agent}\r\n"
        http_request += "Accept: */*\r\n"
        
        if post_data:
            http_request += f"Content-Type: {content_type}\r\n"
            http_request += f"Content-Length: {len(post_data)}\r\n"
            http_request += "\r\n"
            http_request += post_data
        else:
            http_request += "\r\n"
        
        request_packet = Ether(src=client_mac, dst=server_mac) / IP(
            src=client_ip,
            dst=server_ip
        ) / TCP(
            sport=client_port,
            dport=server_port,
            flags="PA",
            seq=seq_base,
            ack=seq_base
        ) / Raw(load=http_request)
        
        packets.append(request_packet)
        
        # HTTPレスポンス
        status_code = 200
        status_text = "OK"
        
        # 一部のリクエストはエラーを返す
        if i % 5 == 0:
            status_code = random.choice([400, 401, 403, 404, 500])
            status_text = {
                400: "Bad Request",
                401: "Unauthorized",
                403: "Forbidden", 
                404: "Not Found",
                500: "Internal Server Error"
            }[status_code]
        
        # レスポンスの内容を設定
        content_type = "text/html"
        response_body = f"<html><body><h1>Hello from {domain}</h1></body></html>"
        
        if uri_path.endswith(".png"):
            content_type = "image/png"
            response_body = f"[バイナリデータ - {len(response_body)} バイト]"
        elif uri_path.startswith("/api/"):
            content_type = "application/json"
            response_body = f'{{"success": true, "message": "Data retrieved", "timestamp": "2025-03-29T10:55:35Z"}}'
        
        http_response = f"HTTP/1.1 {status_code} {status_text}\r\n"
        http_response += f"Content-Type: {content_type}\r\n"
        http_response += f"Content-Length: {len(response_body)}\r\n"
        http_response += "Server: Apache/2.4.41\r\n"
        http_response += "Connection: keep-alive\r\n"
        http_response += "\r\n"
        http_response += response_body
        
        response_packet = Ether(src=server_mac, dst=client_mac) / IP(
            src=server_ip,
            dst=client_ip
        ) / TCP(
            sport=server_port,
            dport=client_port,
            flags="PA",
            seq=seq_base,
            ack=seq_base + len(http_request)
        ) / Raw(load=http_response)
        
        packets.append(response_packet)
    
    # いくつかの追加のHTTPシナリオを追加
    if count > 4:
        # リダイレクト (301/302)
        client_mac = random_mac()
        server_mac = random_mac()
        client_ip = random_ip("192.168.1")
        server_ip = random_ip("10.0.0")
        client_port = random_port(1024, 65000)
        server_port = 80
        seq_base = 50000
        
        # リダイレクトリクエスト
        redirect_request = (
            "GET /old-page HTTP/1.1\r\n"
            "Host: redirect-example.com\r\n"
            "User-Agent: Mozilla/5.0\r\n"
            "\r\n"
        )
        
        redirect_req_packet = Ether(src=client_mac, dst=server_mac) / IP(
            src=client_ip, 
            dst=server_ip
        ) / TCP(
            sport=client_port,
            dport=server_port,
            flags="PA",
            seq=seq_base,
            ack=seq_base
        ) / Raw(load=redirect_request)
        
        packets.append(redirect_req_packet)
        
        # 301リダイレクトレスポンス
        redirect_response = (
            "HTTP/1.1 301 Moved Permanently\r\n"
            "Location: https://redirect-example.com/new-page\r\n"
            "Content-Length: 0\r\n"
            "Connection: close\r\n"
            "\r\n"
        )
        
        redirect_resp_packet = Ether(src=server_mac, dst=client_mac) / IP(
            src=server_ip,
            dst=client_ip
        ) / TCP(
            sport=server_port,
            dport=client_port,
            flags="PA",
            seq=seq_base,
            ack=seq_base + len(redirect_request)
        ) / Raw(load=redirect_response)
        
        packets.append(redirect_resp_packet)
    
    if verbose:
        print(f"HTTP パケット作成完了: {len(packets)}パケット")
    
    # パケットを保存
    save_packets(packets, output_file)
    
    return True

if __name__ == "__main__":
    # 単体テスト用
    create_http_packets("http_test.pcap", count=10, verbose=True)