# rutty-slack
[Rutty](https://github.com/yantene/rutty) はチャットツールで投稿されたコード片を環境で実行し、結果を表示する bot です。
rutty-slackは、Slackから入力を[Rutty](https://github.com/yantene/rutty)に投げ、その結果をSlackへ投稿するbotとなっています。

# 利用方法
以下で立ち上げる。
初回起動時に各環境のDockerイメージのpullが走るため、気長に待つ

```
docker-compose up -d
```


# 必要なSlackアプリの権限
現在は hubot と同様の権限で利用できます。