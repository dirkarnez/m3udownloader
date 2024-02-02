export PATH="/D/Softwares/ffmpeg-2021-10-28-git-e84c83ef98-full_build/bin:/P/Downloads/ffmpeg-2021-10-28-git-e84c83ef98-full_build/bin"



read -p "Enter Web URL: " webURL
echo "You entered $webURL"

read -p "Enter output file name / full path: " filename
echo "You entered $filename"



ffmpeg.exe -i "$(./m3udownloader.exe --url=$webURL)" -c:v libx264 -c:a aac "$filename" 

read -p "Done"
