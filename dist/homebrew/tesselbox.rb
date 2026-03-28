class Tesselbox < Formula
  desc "Hexagon Sandbox Game"
  homepage "https://github.com/tesselstudio/TesselBox-game"
  url "https://github.com/tesselstudio/TesselBox-game/archive/refs/tags/v2.0.0.tar.gz"
  sha256 "sha256_placeholder"
  license "MIT"

  depends_on "go" => :build

  def install
    system "go", "run", "build/build.go", "-os=darwin", "-arch=#{Hardware::CPU.arch}", "-output=tesselbox", "-release"
    bin.install "tesselbox"
  end

  test do
    system "#{bin}/tesselbox", "--version"
  end
end
