golang.org/x/term の調査
========================

ReadLine 関数があるが、よい example が go.dev に載っていない
--------------------

+ 検索で見つかった example：[Example of interactive terminal in Go](https://gist.github.com/artyom/a59e2707976124f387f5)
+ 古いバージョンの golang.org/x/crypto/ssh/terminal をベースとしている
    + golang.org/x/term に書き換える
+ 標準入力が 0 , 標準出力が 1 と決め打ちしているので、Windows では動かない
    + 0 → int(os.Stdin.Fd()) , 1 → int(os.Stdout.Fd()) に書き換える
+ 修正してうごくようにしたバージョン：[./chat.go](https://github.com/hymkor/study-go-readchar/blob/038ef2b0e371c205842fd2184b08ba5ddb04ddf1/chat.go)
+ Windows では Ctrl-F/B など、ASCIIコードのあるキーしか認識しないが、Linux(WSLで検証) では矢印キーもきちんと入力する（カーソルがちゃんと移動する）

getch に相当する処理ではどういうコードが得られるか？
--------------------------

+ term.MakeRaw ~ term.Restore の間、.Read([]byte) でキー入力をすればよいようだが、具体的にどういう値が返ってくるか？
    + → 調査コード： [./main.go](https://github.com/hymkor/study-go-readchar/blob/038ef2b0e371c205842fd2184b08ba5ddb04ddf1/main.go)

+ Linux(WSL)だと、↑ が `\x1B[A` になるので、ほぼ想定通り。Windows では ASCII コードを持っていないキーでは何も返ってこない（が、ASCIIコードを持っているキーであれば期待どおりになった）

[SetConsoleMode] で ENABLE_VIRTUAL_TERMINAL_INPUT(0x0200) をセットすればいけるのでは？
----------------------

いけた！(new [./main.go](https://github.com/hymkor/study-go-readchar/blob/0296afe3d7a1903842d8c7c329e36085fe3edfff/main.go) and [./main_windows.go](https://github.com/hymkor/study-go-readchar/blob/0296afe3d7a1903842d8c7c329e36085fe3edfff/main_windows.go)) Windows でも Linux 同様に ↑ で `\x1B[A` が得られるようになった

結論
----

+ 準標準ライブラリ(golang.org/x/term and x/sys/windows)だけでも、readline 的なことはできる
+ ただし、あいかわらず倍角文字の入力ではカーソル移動位置が狂う
+ 表示については、[少しの修正で](https://github.com/hymkor/study-go-readchar/commit/f4dd61cab3c17023bffabe3f38514602f0ba7a31) 絵文字もいけた

[SetConsoleMode]: https://learn.microsoft.com/ja-jp/windows/console/setconsolemode
