package main

import "testing"

func TestParseGsURL(t *testing.T) {
    tests := []struct{
        in string
        bucket string
        object string
        gen int64
        hasGen bool
    }{
        {"gs://my-bucket/path/to/object", "my-bucket", "path/to/object", 0, false},
        {"gs://bucket/dir/", "bucket", "dir/", 0, false},
        {"gs://b/o#123", "b", "o", 123, true},
    }
    for _, tt := range tests {
        b,o,g,has,err := parseGsURL(tt.in)
        if err != nil {
            t.Fatalf("parseGsURL(%s) error: %v", tt.in, err)
        }
        if b != tt.bucket || o != tt.object || g != tt.gen || has != tt.hasGen {
            t.Fatalf("parseGsURL(%s) = %v,%v,%v,%v; want %v,%v,%v,%v", tt.in, b,o,g,has, tt.bucket,tt.object,tt.gen,tt.hasGen)
        }
    }
}
