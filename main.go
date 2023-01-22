package main

import (
	"context"
	"fmt"
	"os"

	"dagger.io/dagger"
)

func main() {
	if err := build(context.Background()); err != nil {
		fmt.Println(err)
	}
}

func build(ctx context.Context) error {
	fmt.Println("Building with Dagger")

	// クライアントを初期化して Dagger Engine に接続する
	// dagger.WithLogOutput でログの出力先を指定できる
	client, err := dagger.Connect(ctx, dagger.WithLogOutput(os.Stdout))
	if err != nil {
		return err
	}
	defer client.Close()

	// Docker イメージを取得する
	golang := client.Container().From("golang:1.19")

	// カレントディレクトリをコンテナにマウントする
	// ワーキングディレクトリを指定する
	src := client.Host().Directory(".")
	golang = golang.
		WithMountedDirectory("/src", src).
		WithWorkdir("/src")

	// 実行するコマンドを設定する
	path := "build/"
	golang = golang.
		WithExec([]string{"go", "test", "-v", "./..."}).
		WithExec([]string{"go", "build", "-o", path})

	// バイナリを出力するディレクトリを取得します。
	output := golang.Directory(path)
	// バイナリを出力します。
	_, err = output.Export(ctx, path)

	return nil
}
