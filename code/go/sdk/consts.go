package datastar

import "time"

const (
	Version                        = "0.20.0-beta"
	VersionClientByteSize          = 43609
	VersionClientByteSizeGzip      = 14889
	VersionClientByteSizeGzipHuman = "14 KiB"

	DefaultSettleTime = 300 * time.Millisecond
	DefaultSseSendRetry = 1000 * time.Millisecond
	DefaultFragmentMergeMode = FragmentMergeMode("morph")
)