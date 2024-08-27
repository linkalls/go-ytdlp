# rcloneでアップロード後自動削除するスクリプト

## 変更場所
main関数内の**dir**,**errorDir**,**dropboxDir**を変更してね
#### rcloneはpathを通しておいてね
``` go
func main() {
    dir := `E:\minato` 
    errorDir := `E:\error`
    dropboxDir := "dropbox:youtube"

```

## yt-dlpコマンド
```py
yt-dlp -f "bestvideo[ext=mp4]+bestaudio[ext=m4a]/best" --merge-output-format mp4 -o "%(title)s.%(ext)s" --match-filter "!is_live & !was_live" "https://www.youtube.com/[ここにアカウント名]"
```
**dir**で指定したフォルダ内でこのコマンドを実行してね
ffmpegとyt-dlpをインストールしといてね