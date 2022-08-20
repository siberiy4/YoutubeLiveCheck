# YoutubeLiveCheck

### youtube data api

参考
https://zenn.dev/meihei/articles/1021b1a3f8c226

#### チャンネルの動画リストをとる
```
https://www.googleapis.com/youtube/v3/search?key=`APiキー`&part=id&channelId=`チャンネルのID`&order=date&maxResults=`取得したい件数`
```

#### ビデオの詳細情報をとる
```
https://www.googleapis.com/youtube/v3/videos?key=`APIキー`&id=`ビデオID`&part=liveStreamingDetails
```

参考記事にビデオIDはコンマ区切りで列挙できると書いてある。  
実験したところ、最大列挙数は50。
50個列挙した最後にコンマを付けると51個目を列挙するものとして ` The request specifies an invalid filter parameter.`


#### どうやって配信前、配信中、配信後か判定するか
ビデオの詳細情報のitemsの中にある`liveStreamingDetails`を見る
- ただの動画
    - liveStreamingDetailsが存在しない
- 配信終了後のLive
    - liveStreamingDetailsに scheduledStartTime、actualStartTime、actualEndTime
- 配信中のLive
    - liveStreamingDetailsに scheduledStartTime、actualStartTime 
- 配信前のLive
    - liveStreamingDetailsに scheduledStartTime

