
golang.org/x/term の調査
========================

+ ReadLine 関数があるが、よい example が go.dev に載っていない
    + → 検索：[Example of interactive terminal in Go](https://gist.github.com/artyom/a59e2707976124f387f5)
        + 古いバージョンの golang.org/x/crypto/ssh/terminal をベースとしている
            + golang.org/x/term に書き換える
        + 標準入力が 0 , 標準出力が 1 と決め打ちしているので、Windows では動かない
            + 0 → int(os.Stdin.Fd()) , 1 → int(os.Stdout.Fd()) に書き換える
        + 修正してうごくようにしたバージョン：[./chat.go](./chat.go)

+ Windows では Ctrl-F/B など、ASCIIコードのあるキーしか認識しないが、Linux(WSLで検証) では矢印キーもきちんと入力する（カーソルがちゃんと移動する）

+ term.MakeRaw ~ term.Restore の間、.Read([]byte) でキー入力をすればよいようだが、具体的にどういう値が返ってくるか？
    + → 調査コード： [./main.go](./main.go)

+ Linux(WSL)だと、↑ が `\x1B[A` になるので、ほぼ想定通り。Windows では ASCII コードを持っていないキーでは何も返ってこない（が、ASCIIコードを持っているキーであれば期待どおりになった）
