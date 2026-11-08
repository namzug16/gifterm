#!/bin/sh
# Prints the Homebrew formula for a released version.
# Usage: homebrew-formula.sh VERSION DIST_DIR
# DIST_DIR holds the gifterm-<os>-<arch>.tar.gz archives built by `make release-artifacts`.
set -eu

version="$1"
dist_dir="$2"

archive_sha256() {
  sha256sum "$dist_dir/gifterm-$1.tar.gz" | cut -d ' ' -f 1
}

release_url="https://github.com/namzug16/gifterm/releases/download/v${version}"

cat <<FORMULA
class Gifterm < Formula
  desc "Play GIFs in the terminal as ASCII art"
  homepage "https://github.com/namzug16/gifterm"
  version "${version}"

  on_macos do
    on_arm do
      url "${release_url}/gifterm-darwin-arm64.tar.gz"
      sha256 "$(archive_sha256 darwin-arm64)"
    end
    on_intel do
      url "${release_url}/gifterm-darwin-amd64.tar.gz"
      sha256 "$(archive_sha256 darwin-amd64)"
    end
  end

  on_linux do
    on_arm do
      url "${release_url}/gifterm-linux-arm64.tar.gz"
      sha256 "$(archive_sha256 linux-arm64)"
    end
    on_intel do
      url "${release_url}/gifterm-linux-amd64.tar.gz"
      sha256 "$(archive_sha256 linux-amd64)"
    end
  end

  def install
    bin.install "gifterm"
  end

  test do
    assert_match version.to_s, shell_output("#{bin}/gifterm -version")
  end
end
FORMULA
