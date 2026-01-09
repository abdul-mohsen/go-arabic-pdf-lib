#!/bin/bash

# Create fonts directory
mkdir -p fonts

# Download Amiri font (excellent Arabic support)
echo "Downloading Amiri font..."
curl -L "https://github.com/aliftype/amiri/releases/download/1.000/Amiri-1.000.zip" -o amiri.zip
unzip -o amiri.zip -d temp_amiri
cp temp_amiri/Amiri-1.000/Amiri-Regular.ttf fonts/
cp temp_amiri/Amiri-1.000/Amiri-Bold.ttf fonts/
rm -rf amiri.zip temp_amiri

echo "Fonts downloaded to ./fonts/"
ls -la fonts/
