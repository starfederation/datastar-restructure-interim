package datastar

import "time"

const (
	Version                        = "0.20.0-beta"
	VersionClientByteSize          = 43677
	VersionClientByteSizeGzip      = 14902
	VersionClientByteSizeGzipHuman = "15 KiB"

	DefaultSettleTime = 300 * time.Millisecond
	DefaultSseSendRetry = 1000 * time.Millisecond
	DefaultFragmentMergeMode = FragmentMergeMode("morph")
)