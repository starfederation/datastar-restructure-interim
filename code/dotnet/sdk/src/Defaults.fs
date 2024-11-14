module StarFederation.Datastar.Defaults

open System

let [<Literal>] Version = "1.0.0-beta.1"
let [<Literal>] VersionClientByteSize = 43972
let [<Literal>] VersionClientByteSizeGzip = 15021

let DefaultSettleDuration = TimeSpan.FromMilliseconds 300
let DefaultSseSendRetry = TimeSpan.FromMilliseconds 1000
let DefaultFragmentMergeMode = "morph"
