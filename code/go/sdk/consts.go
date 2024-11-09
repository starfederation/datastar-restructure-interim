package datastar

import "time"

const (
	Version                        = "1.0.0-beta.1"
	VersionClientByteSize          = 43972
	VersionClientByteSizeGzip      = 15021
	VersionClientByteSizeGzipHuman = "15 KiB"

	DefaultSettleTime = 300 * time.Millisecond
	DefaultSseSendRetry = 1000 * time.Millisecond
	DefaultFragmentMergeMode = FragmentMergeMode("morph")
)