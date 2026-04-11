#!/usr/bin/env python3

import sys
import subprocess
from pathlib import Path

SPIDER_DIR = Path("spider")
SCORPION_DIR = Path("scorpion")
SPIDER_BIN = SPIDER_DIR / "spider"
SCORPION_BIN = SCORPION_DIR / "target/debug/scorpion"

DEFAULT = "\033[0m"
BOLD = "\033[1m"
RED = "\033[31m"
YELLOW = "\033[33m"
BLUE = "\033[34m"
WHITE = "\033[37m"

IMAGE_EXTS = {".jpg", ".jpeg", ".png", ".gif", ".bmp"}


def check_compilation(
    name: str, path: str, binary: str, command: list[str]
) -> None:
    if not binary.exists():
        print(f"{BOLD}{BLUE}Bridge:{WHITE} compiling {name}...{DEFAULT}")
    try:
        subprocess.run(
            command, cwd=path, check=True, capture_output=True, text=True
        )

    except subprocess.CalledProcessError as e:
        raise RuntimeError(
            f"{BOLD}{RED}Bridge:{DEFAULT} {name} failed to compile"
        ) from e


def get_download_dir(spider_args: list[str], spider_dir: str) -> Path:
    path = Path(f"{spider_dir}/data")
    for i, arg in enumerate(spider_args):
        if arg == "-p" and i + 1 < len(spider_args):
            path = Path(spider_dir) / Path(spider_args[i + 1].rstrip("/"))
            break
    return path


def run_spider(spider_args: list[str]) -> Path:
    try:
        print(f"{BOLD}{BLUE}Bridge:{WHITE} launching spider...{DEFAULT}\n")
        cmd = ["./" + str(Path(SPIDER_BIN).name)] + spider_args
        subprocess.run(cmd, cwd=SPIDER_DIR, check=True)

    except subprocess.CalledProcessError as e:
        raise RuntimeError(
            f"{BOLD}{RED}Bridge:{DEFAULT} spider failed with exit code "
            f"{BOLD}{YELLOW}{e.returncode}{DEFAULT}"
        ) from e

    return get_download_dir(spider_args, str(SPIDER_BIN.parent))


def run_scorpion(image_paths: list[str]) -> None:
    try:
        print(f"\n{BOLD}{BLUE}Bridge: {WHITE}launching scorpion...{DEFAULT}\n")
        cmd = [SCORPION_BIN, *image_paths]
        subprocess.call(cmd)

    except subprocess.CalledProcessError as e:
        raise RuntimeError(
            f"{BOLD}{RED}Bridge:{DEFAULT} scorpion failed with exit code "
            f"{BOLD}{YELLOW}{e.returncode}{DEFAULT}"
        ) from e


def main() -> int:
    if len(sys.argv) < 2:
        print(
            f"{BOLD}{WHITE}Usage:{DEFAULT} ./bridge.py [spider options] URL",
            file=sys.stderr,
        )
        return 2

    try:
        check_compilation(
            "spider", SPIDER_DIR, SPIDER_BIN, ["go", "build", "-o", "spider"]
        )
        check_compilation(
            "scorpion", SCORPION_DIR, SCORPION_BIN, ["cargo", "build"]
        )

        download_dir = run_spider(sys.argv[1:])
        if not download_dir.exists():
            print(
                f"{BOLD}{RED}Bridge:{DEFAULT} can't find",
                f"{BOLD}{WHITE}{download_dir}{DEFAULT}",
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
            return 1

        run_scorpion(image_paths)

    except RuntimeError as e:
        print(f"{e}", file=sys.stderr)
        return 1

    except KeyboardInterrupt:
        print(f"\n{BOLD}{RED}Bridge:{DEFAULT} abording...", file=sys.stderr)
        return 1

    return 0


if __name__ == "__main__":
    sys.exit(main())
