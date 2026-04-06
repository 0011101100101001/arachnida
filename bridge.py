#!/usr/bin/env python3

import sys
import subprocess
from pathlib import Path

SPIDER_BIN = "./spider/spider"
SCORPION_BIN = "./scorpion/target/debug/scorpion"

DEFAULT = "\033[0m"
BOLD    = "\033[1m"
RED     = "\033[31m"
YELLOW  = "\033[33m"
BLUE    = "\033[34m"
WHITE   = "\033[37m"

IMAGE_EXTS = {".jpg", ".jpeg", ".png", ".gif", ".bmp"}

def get_download_dir(spider_args: list[str], spider_dir: str) -> Path:
    path = Path("./data")
    for i, arg in enumerate(spider_args):
        if arg == "-p" and i + 1 < len(spider_args):
            path = Path(spider_dir) / Path(spider_args[i + 1].rstrip("/"))
            break
    return path

def main() -> int:
    if len(sys.argv) < 2:
        print("Usage: ./bridge.py [spider options] URL", file=sys.stderr)
        return 2

    # SPIDER
    print(f"{BOLD}{BLUE}Bridge:{WHITE} launching spider...{DEFAULT}\n")
    spider_args = sys.argv[1:]
    spider_cmd = ["./" + str(Path(SPIDER_BIN).name)] + spider_args
    spider_dir = str(Path(SPIDER_BIN).parent)
    spider_res = subprocess.run(spider_cmd, cwd=spider_dir)

    if spider_res.returncode != 0:
        print(f"{BOLD}{RED}Error: spider failed with exit code {YELLOW}{spider_res.returncode}{DEFAULT}", file=sys.stderr)
        return spider_res.returncode


    download_dir = get_download_dir(spider_args, spider_dir)
    if not download_dir.exists():
        print(f"{BOLD}{RED}Error:{DEFAULT} can't find {download_dir}", file=sys.stderr)
        return 1

    image_paths = [
        str(p)
        for p in sorted(download_dir.iterdir())
        if p.is_file() and p.suffix.lower() in IMAGE_EXTS
    ]

    if not image_paths:
        print(f"{BOLD}{RED}Error:{DEFAULT} no images in {download_dir}", file=sys.stderr)
        return 0

    # SCORPION
    print(f"\n{BOLD}{BLUE}Bridge: {WHITE}launching scorpion...{DEFAULT}\n")
    scorpion_cmd = [SCORPION_BIN, *image_paths]
    return subprocess.call(scorpion_cmd)


if __name__ == "__main__":
    sys.exit(main())