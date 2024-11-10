module StarFederation.Datastar.Defaults

open System

let [<Literal>] Version = "1.0.0-beta.1"
let [<Literal>] VersionClientByteSize = 43972
let [<Literal>] VersionClientByteSizeGzip = 15021
let [<Literal>] VersionClientByteSizeGzipHuman = "15 KiB"

let DefaultSettleTime = TimeSpan.FromMilliseconds 300
let DefaultSseSendRetry = TimeSpan.FromMilliseconds 1000
let DefaultFragmentMergeMode = "morph"
