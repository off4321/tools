#!/usr/bin/env python3
"""
packet_sequenceのセットアップスクリプト
"""
from setuptools import setup, find_packages

with open("README.md", "r", encoding="utf-8") as fh:
    long_description = fh.read()

setup(
    name="packet_sequence",
    version="0.1.0",
    author="Your Name",
    author_email="your.email@example.com",
    description="PCAPファイルからシーケンス図を生成するツール",
    long_description=long_description,
    long_description_content_type="text/markdown",
    url="https://github.com/yourusername/packet_sequence",
    packages=find_packages(),
    classifiers=[
        "Programming Language :: Python :: 3",
        "License :: OSI Approved :: MIT License",
        "Operating System :: OS Independent",
    ],
    python_requires='>=3.7',
    install_requires=[
        "pyshark",
    ],
    entry_points={
        "console_scripts": [
            "packet_sequence=src.main:main",
        ],
    },
)
