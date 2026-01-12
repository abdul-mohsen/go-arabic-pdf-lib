# PowerShell script to download Arabic fonts

$fontsDir = "fonts"
New-Item -ItemType Directory -Force -Path $fontsDir | Out-Null

Write-Host "Downloading Amiri font..."

$amiriUrl = "https://github.com/aliftype/amiri/releases/download/1.000/Amiri-1.000.zip"
$amiriZip = "amiri.zip"

Invoke-WebRequest -Uri $amiriUrl -OutFile $amiriZip
Expand-Archive -Path $amiriZip -DestinationPath "temp_amiri" -Force

Copy-Item "temp_amiri\Amiri-1.000\Amiri-Regular.ttf" -Destination $fontsDir
Copy-Item "temp_amiri\Amiri-1.000\Amiri-Bold.ttf" -Destination $fontsDir

Remove-Item $amiriZip -Force
Remove-Item "temp_amiri" -Recurse -Force

Write-Host "Fonts downloaded to .\fonts\"
Get-ChildItem $fontsDir
