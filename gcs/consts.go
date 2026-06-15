package main

const (
    // defaultConcurrency used when not provided via env
    defaultConcurrency = 4
    // defaultRangeThreshold size in bytes above which we perform ranged parallel download
    defaultRangeThreshold = 10 * 1024 * 1024 // 10MB
    // defaultRangeChunk size of each range chunk in bytes
    defaultRangeChunk = 4 * 1024 * 1024 // 4MB
)
