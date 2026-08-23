#!/usr/bin/env python3

from __future__ import annotations

import re
import subprocess
from pathlib import Path

from config import load_config

SKIPPED_MODULES = (
    "qtwebengine",
    "qtquick3d",
    "qtgraphs",
    "qtquick3dphysics",
    "qtcanvaspainter",
    "qtdeclarative",
    "qtquicktimeline",
    "qtdoc",
    "qtlocation",
    "qtlottie",
    "qtmqtt",
    "qtopcua",
    "qtquickeffectmaker",
    "qtvirtualkeyboard",
    "qtwebview",
)


def run(command: list[str], *, cwd: Path | None = None) -> None:
    subprocess.run(command, cwd=cwd, check=True)


def enable_static_pkg_config(source_dir: Path) -> None:
    helper = source_dir / "qtbase/cmake/QtPkgConfigHelpers.cmake"
    contents = helper.read_text(encoding="utf-8")
    block = re.compile(
        r"^[ \t]*if\(NOT BUILD_SHARED_LIBS\)[ \t]*\n"
        r".*?"
        r"^[ \t]*endif\(\)[ \t]*\n",
        flags=re.MULTILINE | re.DOTALL,
    )
    updated, replacements = block.subn("", contents, count=1)
    if replacements != 1:
        raise RuntimeError(f"Could not patch static pkg-config support in {helper}")
    helper.write_text(updated, encoding="utf-8")


def main() -> None:
    config = load_config()
    config.qt_root.mkdir(parents=True, exist_ok=True)
    archive = config.qt_root / config.qt_source_archive

    run(
        [
            "curl",
            "--fail",
            "--location",
            "--retry",
            "3",
            "--output",
            str(archive),
            config.qt_source_url,
        ]
    )
    run(["tar", "-xf", str(archive), "-C", str(config.qt_root)])
    enable_static_pkg_config(config.qt_source_dir)

    configure = [
        "./configure",
        "-prefix",
        str(config.qt_prefix),
        "-release",
        "-static",
        "-static-runtime",
        "-opensource",
        "-confirm-license",
        "-nomake",
        "examples",
        "-nomake",
        "tests",
        "-no-opengl",
        "-no-dbus",
        "-no-gtk",
        "-no-feature-zstd",
    ]
    for module in SKIPPED_MODULES:
        configure.extend(("-skip", module))
    configure.extend(
        (
            "--",
            "-DCMAKE_POLICY_DEFAULT_CMP0177=NEW",
            "-DFEATURE_pkg_config=ON",
        )
    )
    run(configure, cwd=config.qt_source_dir)

    logical_cpus = subprocess.check_output(
        ["sysctl", "-n", "hw.logicalcpu"], text=True
    ).strip()
    run(
        ["cmake", "--build", ".", "--parallel", logical_cpus],
        cwd=config.qt_source_dir,
    )
    run(["cmake", "--install", "."], cwd=config.qt_source_dir)


if __name__ == "__main__":
    main()
