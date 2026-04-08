#!/usr/bin/env python3

import sys
import subprocess
from pathlib import Path

SPIDER_BIN = "./spider/spider"
SCORPION_BIN = "./scorpion/target/debug/scorpion"

DEFAULT = "\033[0m"
BOLD = "\033[1m"
RED = "\033[31m"
YELLOW = "\033[33m"
BLUE = "\033[34m"
WHITE = "\033[37m"

IMAGE_EXTS = {".jpg", ".jpeg", ".png", ".gif", ".bmp"}

def get_download_dir(spider_args: list[str], spider_dir: str) -> Path:
    path = Path("./data")
    for i, arg in enumerate(spider_args):
        if arg == "-p" and i + 1 < len(spider_args):
            path = Path(spider_dir) / Path(spider_args[i + 1].rstrip("/"))
            break
    return path


def check_compilation() -> int:
    # Compile Spider
    binary_path = Path(SPIDER_BIN)
    if not binary_path.exists():
        spider_build = subprocess.run(
            ["go", "build", "-o", "spider"], cwd="spider"
        )
        if spider_build.returncode != 0:
            print(
                f"{BOLD}{RED}Bridge:{DEFAULT} failed to build spider",
                file=sys.stderr,
            )
            return spider_build.returncode

    # Compile Scorpion
    binary_path = Path(SCORPION_BIN)
    if not binary_path.exists():
        scorpion_build = subprocess.run(["cargo", "build"], cwd="scorpion")
        if scorpion_build.returncode != 0:
            print(
                f"{BOLD}{RED}Bridge:{DEFAULT} failed to build scorpion",
                file=sys.stderr,
            )
            return scorpion_build.returncode

    return 0


def main() -> int:
    if len(sys.argv) < 2:
        print("Usage: ./bridge.py [spider options] URL", file=sys.stderr)
        return 2

    if (code := check_compilation()) != 0:
        return code

    # Run Spider
    print(f"{BOLD}{BLUE}Bridge:{WHITE} launching spider...{DEFAULT}\n")
    spider_args = sys.argv[1:]
    spider_cmd = ["./" + str(Path(SPIDER_BIN).name)] + spider_args
    spider_dir = str(Path(SPIDER_BIN).parent)
    spider_res = subprocess.run(spider_cmd, cwd=spider_dir)

    if spider_res.returncode != 0:
        print(
            f"{BOLD}{RED}Bridge: spider failed with exit code "
            f"{YELLOW}{spider_res.returncode}{DEFAULT}",
            file=sys.stderr,
        )
        return spider_res.returncode

    download_dir = get_download_dir(spider_args, spider_dir)
    if not download_dir.exists():
        print(
            f"{BOLD}{RED}Bridge:{DEFAULT} can't find {download_dir}",
            file=sys.stderr,
        )
        return 1

    image_paths = [
        str(path)
        for path in sorted(download_dir.iterdir())
        if path.is_file() and path.suffix.lower() in IMAGE_EXTS
    ]

    if not image_paths:
        print(
            f"{BOLD}{RED}Bridge:{DEFAULT} no images in {download_dir}",
            file=sys.stderr,
        )
        return 0

    # Run Scorpion
    print(f"\n{BOLD}{BLUE}Bridge: {WHITE}launching scorpion...{DEFAULT}\n")
    scorpion_cmd = [SCORPION_BIN, *image_paths]
    return subprocess.call(scorpion_cmd)


if __name__ == "__main__":
    sys.exit(main())
