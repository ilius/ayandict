#!/usr/bin/env python3

from __future__ import annotations

import tomllib
from dataclasses import dataclass
from pathlib import Path

MAC_CI_DIR = Path(__file__).resolve().parent
MAC_CONFIG_FILE = MAC_CI_DIR / "config.toml"


@dataclass(frozen=True)
class MacConfig:
	qt_version: str
	qt_source_url: str
	qt_root: Path
	qt_cache_revision: str
	go_version: str
	go_cache_root: Path

	@property
	def qt_source_archive(self) -> str:
		return self.qt_source_url.rsplit("/", 1)[-1]

	@property
	def qt_source_dir(self) -> Path:
		archive = self.qt_source_archive
		source_name = archive.removesuffix(".tar.xz")
		return self.qt_root / source_name

	@property
	def qt_prefix(self) -> Path:
		return self.qt_root / f"{self.qt_version}-static"

	@property
	def qt_cache_key(self) -> str:
		return f"qt{self.qt_version}-static-macos-{self.qt_cache_revision}"

	@property
	def go_cache_dir(self) -> Path:
		return self.go_cache_root / f"qt-{self.qt_version}-static"


def _required_string(values: dict[str, object], key: str) -> str:
	value = values.get(key)
	if not isinstance(value, str) or not value:
		raise ValueError(f"'{key}' must be a non-empty string in {MAC_CONFIG_FILE}")
	return value


def load_config() -> MacConfig:
	with MAC_CONFIG_FILE.open("rb") as config_file:
		values = tomllib.load(config_file)

	return MacConfig(
		qt_version=_required_string(values, "qt_version"),
		qt_source_url=_required_string(values, "qt_source_url"),
		qt_root=Path(_required_string(values, "qt_root")),
		qt_cache_revision=_required_string(values, "qt_cache_revision"),
		go_version=_required_string(values, "go_version"),
		go_cache_root=Path(_required_string(values, "go_cache_root")),
	)
