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
      description = "A Foosball data aggregator, analyzers, and manager based on Kickertool.";
      homepage = "https://github.com/crispgm/kicker-cli/";
      license = licenses.mit;
      maintainers = with maintainers; [ crispgm ];
    };
}

