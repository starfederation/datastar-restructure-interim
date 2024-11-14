package datastar

import "time"

const (
	Version                        = "0.20.0"
	VersionClientByteSize          = 43672
	VersionClientByteSizeGzip      = 14899
	VersionClientByteSizeGzipHuman = "14 KiB"

	DefaultSettleTime = 300 * time.Millisecond
	DefaultSseSendRetry = 1000 * time.Millisecond
	DefaultFragmentMergeMode = FragmentMergeMode("morph")
)