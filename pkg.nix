{
buildGoModule,
fetchFromGithub,
...
}:



buildGoModule {
    pname = "kicker-cli";
    version = "1.0.0"; #there are no versions?
    
    src = fetchFromGithub {
      owner = "crispgm";
      repo = "kicker-cli";
      #rev = latest 
      sha256 = "sha256-YsR2KU5Np6xQHkjM8KAoDp/XZ/9DkwBlMbu2IX5OQlk=";
    };
}

