#!/usr/bin/env python3

from __future__ import annotations

import re
import subprocess
import sys
from pathlib import Path

HOMEBREW_PATH = re.compile(
	r"(?:/opt/homebrew|/usr/local/(?:Cellar|opt)|/home/linuxbrew|\.linuxbrew)(?:/|$)"
)


def homebrew_references(load_commands: str) -> list[str]:
	return sorted(
		{
			line.strip()
			for line in load_commands.splitlines()
			if HOMEBREW_PATH.search(line)
		}
	)


def main() -> None:
	if len(sys.argv) != 2:
		raise SystemExit(f"usage: {Path(sys.argv[0]).name} BINARY")

	binary = Path(sys.argv[1])
	if not binary.is_file():
		raise FileNotFoundError(binary)

	load_commands = subprocess.run(
		["otool", "-l", str(binary)],
		check=True,
		stdout=subprocess.PIPE,
		text=True,
	).stdout
	references = homebrew_references(load_commands)
	if references:
		details = "\n".join(f"  {reference}" for reference in references)
		raise RuntimeError(f"Binary contains Homebrew load paths:\n{details}")

	subprocess.run(["otool", "-L", str(binary)], check=True)


if __name__ == "__main__":
	main()
