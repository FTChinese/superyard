param (
    $stage = 'build'
)

$version = git tag -l --sort=-v:refname | Select-Object -First 1
$build_at = Get-Date -Format "yyyy-MM-ddTHH:mm:ssK"
$commit = git log --max-count=1 --pretty=format:%aI_%h

$ldflags = "-w -s -X main.version=$version -X main.build=$build_at -X main.commit=$commit"

$app_name = 'superyard.exe'
$out_dir = './build'

$exec = "$out_dir/$app_name"

switch ($stage)
{
    'build' { go build -o $exec -ldflags $ldflags -tags production -v . }
    'run' { $exec | Invoke-Expression }
    'version' {
        New-Item -Path "$out_dir/version" -ItemType "file" -Value $version -Force

        New-Item -Path "$out_dir/build_time" -ItemType "file" -Value $build_at -Force

        New-Item -Path "$out_dir/commit" -ItemType "file" -Value $commit -Force

        Copy-Item -Path '~/config/env.dev.toml' -Destination "$out_dir/api.toml"
    }
}
