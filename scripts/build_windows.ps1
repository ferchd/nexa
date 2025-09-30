param(
    [string]$Version = (git describe --tags --abbrev=0)
)

$OutputDir = "dist\windows"
New-Item -ItemType Directory -Force -Path $OutputDir

Write-Host "Building NetCheck for Windows - Version: $Version"

$Architectures = @("amd64", "386")

foreach ($arch in $Architectures) {
    Write-Host "Building for $arch..."
    
    $env:GOOS = "windows"
    $env:GOARCH = $arch
    
    $BinaryName = "netcheck-windows-$arch.exe"
    $OutputPath = Join-Path $OutputDir $BinaryName
    
    go build -ldflags="-s -w -X main.version=$Version" `
             -o $OutputPath `
             ./cmd/netcheck
    
    $ZipName = "netcheck-windows-$arch.zip"
    $ZipPath = Join-Path $OutputDir $ZipName
    
    Compress-Archive -Path $OutputPath -DestinationPath $ZipPath -Force
    
    Write-Host "✅ Built $ZipPath"
}

Set-Location $OutputDir
Get-FileHash *.zip | Format-List > checksums.txt
Write-Host "✅ Created checksums"

Write-Host "Windows builds completed successfully!"