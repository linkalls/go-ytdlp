package main

import (
    "fmt"
    "os"
    "os/exec"
    "path/filepath"
    "regexp"
    "strings"
    "time"
)

func main() {
    dir := `E:\minato` 
    errorDir := `E:\error`
    dropboxDir := "dropbox:youtube"

    // エラーフォルダが存在しない場合は作成
    if _, err := os.Stat(errorDir); os.IsNotExist(err) {
        err := os.Mkdir(errorDir, os.ModePerm)
        if err != nil {
            fmt.Println("エラーディレクトリの作成エラー:", err)
            return
        }
    }

    // 特殊文字（記号や絵文字など）を削除する正規表現
    re := regexp.MustCompile(`[^\w\s\p{Han}\p{Hiragana}\p{Katakana}#\[\]【】.-]`)
    // 全角スペースをアンダースコアに置換する正規表現
    spaceRe := regexp.MustCompile(`\p{Zs}`)

    for {
        files, err := os.ReadDir(dir)
        if err != nil {
            fmt.Println("ディレクトリの読み取りエラー:", err)
            return
        }

        found := false

        for _, file := range files {
            if !file.IsDir() && strings.HasSuffix(file.Name(), ".mp4") && !strings.Contains(file.Name(), ".f") && !strings.HasSuffix(file.Name(), ".temp.mp4") {
                found = true
                oldName := file.Name()
                // 特殊文字を削除
                newName := re.ReplaceAllString(oldName, "")
                // 全角スペースをアンダースコアに置換
                newName = spaceRe.ReplaceAllString(newName, "_")
                // 全角の「【」と「】」を半角の「[」と「]」に置換
                newName = strings.ReplaceAll(newName, "【", "[")
                newName = strings.ReplaceAll(newName, "】", "]")
                oldPath := filepath.Join(dir, oldName)
                newPath := filepath.Join(dir, newName)

                // ファイルが使用中でないことを確認
                if !isFileInUse(oldPath) {
                    if oldName != newName {
                        err := os.Rename(oldPath, newPath)
                        if err != nil {
                            fmt.Println("ファイル名の変更エラー:", err)
                            continue
                        }
                    }

                    fmt.Printf("アップロード中: %s\n", newPath)
                    cmd := exec.Command("rclone", "copy", newPath, dropboxDir, "--progress")
                    cmd.Stdout = os.Stdout
                    cmd.Stderr = os.Stderr
                    err = cmd.Run()
                    if err != nil {
                        fmt.Println("ファイルのアップロードエラー:", err)
                        // エラーファイルをエラーフォルダに移動
                        errorPath := filepath.Join(errorDir, newName)
                        moveErr := os.Rename(newPath, errorPath)
                        if moveErr != nil {
                            fmt.Println("エラーディレクトリへのファイル移動エラー:", moveErr)
                        } else {
                            fmt.Printf("エラーディレクトリに移動: %s\n", errorPath)
                        }
                        continue
                    }

                    fmt.Printf("アップロード完了: %s\n", newPath)
                    err = os.Remove(newPath)
                    if err != nil {
                        fmt.Println("ファイルの削除エラー:", err)
                    } else {
                        fmt.Printf("削除完了: %s\n", newPath)
                    }
                } else {
                    fmt.Printf("ファイルが使用中のためスキップ: %s\n", oldPath)
                }
            }
        }

        if !found {
            fmt.Println("動画ファイルが見つかりません。2秒後に再試行します...")
            time.Sleep(2 * time.Second)
        }
    }
}

// ファイルが使用中かどうかを確認する関数
func isFileInUse(filePath string) bool {
    file, err := os.Open(filePath)
    if err != nil {
        return true
    }
    file.Close()
    return false
}