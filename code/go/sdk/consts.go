package datastar

import "time"

const (
	Version                        = "0.20.0-beta"
	VersionClientByteSize          = 43974
	VersionClientByteSizeGzip      = 15017
	VersionClientByteSizeGzipHuman = "15 KiB"

	DefaultSettleTime = 300 * time.Millisecond
	DefaultSseSendRetry = 1000 * time.Millisecond
	DefaultFragmentMergeMode = FragmentMergeMode("morph")
)