# scripts/analyzer.spec
# -*- mode: python ; coding: utf-8 -*-

import sys
import os

block_cipher = None

a = Analysis(
    ['../analyzer/__main__.py'],
    pathex=['.'],
    binaries=[],
    datas=[('../analyzer', 'analyzer')],
    hiddenimports=[
        'chromadb',
        'chromadb.telemetry.posthog',
        'sentence_transformers',
        'pydantic',
        'pydantic_core',
        'rich',
        'pyvis',
        'networkx',
        'requests',
        'jinja2'
    ],
    hookspath=[],
    hooksconfig={},
    runtime_hooks=[],
    excludes=['tkinter', 'matplotlib', 'notebook'],
    win_no_prefer_redirects=False,
    win_private_assemblies=False,
    cipher=block_cipher,
    noarchive=False,
)

pyz = PYZ(a.pure, a.zipped_data, cipher=block_cipher)

exe = EXE(
    pyz,
    a.scripts,
    [],
    exclude_binaries=True,
    name='spectre-analyzer',
    debug=False,
    bootloader_ignore_signals=False,
    strip=False,
    upx=True,
    console=True,
)

coll = COLLECT(
    exe,
    a.binaries,
    a.zipfiles,
    a.datas,
    strip=False,
    upx=True,
    upx_exclude=[],
    name='spectre-analyzer',
)
