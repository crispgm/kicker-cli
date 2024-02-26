{
buildGoModule,
fetchFromGitHub,
stdenv,
lib,
...
}:



stdenv.mkDerivation buildGoModule rec {
    pname = "kicker-cli";
    version = "1.0.0"; #there are no versions?    
    src = fetchFromGitHub {
      owner = "crispgm";
      repo = "${pname}";
      rev = "18fb2d566fb62f68c90d19468bb5dc545ea47da7"; 
      sha256 = "0xwmpnznys6v7sl7rvnqhjdy78l17hw3bqgi8m12iyw7m1plj6nv";
    };
    vendorHash =  "sha256-sK+blqs5DDcF2Am6GNpkZhQ/AvoDqpxrWZa4rb3v5iE=";
    meta = with lib; {
      description = "Simple command-line snippet manager, written in Go";
      homepage = "https://github.com/knqyf263/pet";
      license = licenses.mit;
      maintainers = with maintainers; [ kalbasit ];
    };
}

