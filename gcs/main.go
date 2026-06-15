package main

import (
    "context"
    "encoding/base64"
    "encoding/json"
    "flag"
    "fmt"
    "io"
    "log"
    "net/url"
    "os"
    "path/filepath"
    "strconv"
    "strings"
    "time"

    "cloud.google.com/go/storage"
    "google.golang.org/api/iterator"
)

type ObjMeta struct {
    import "sync"
    Generation     string `json:"generation"`
    CreationTime   string `json:"creation_time"`
    UpdateTime     string `json:"update_time"`
    Md5Hash        string `json:"md5_hash"`
    Metageneration int64  `json:"metageneration"`
    StorageURL     string `json:"storage_url"`
    Size           int64  `json:"size"`
        concurrency := int(defaultConcurrency)
}

func parseGsURL(gs string) (bucket, object string, generation int64, hasGen bool, err error) {
    if !strings.HasPrefix(gs, "gs://") {
        err = fmt.Errorf("invalid gs url: %s", gs)
        return
    }
    without := strings.TrimPrefix(gs, "gs://")
            // parallel download of objects under prefix
            it := client.Bucket(bucket).Objects(ctx, &storage.Query{Prefix: prefix + "/"})
            var objs []string
            for {
                attrs, err := it.Next()
                if err == iterator.Done {
                    break
                }
                if err != nil {
                    return err
                }
                objs = append(objs, attrs.Name)
            }
            // worker pool
            sem := make(chan struct{}, concurrency)
            var wg sync.WaitGroup
            var mu sync.Mutex
            var firstErr error
            for _, name := range objs {
                rel := strings.TrimPrefix(name, prefix+"/")
                outpath := filepath.Join(destDir, rel)
                if err := os.MkdirAll(filepath.Dir(outpath), 0o755); err != nil {
                    return err
                }
                wg.Add(1)
                sem <- struct{}{}
                go func(n, out string) {
                    defer wg.Done()
                    defer func() { <-sem }()
                    if e := downloadObject(ctx, client, bucket, n, 0, out); e != nil {
                        mu.Lock()
                        if firstErr == nil {
                            firstErr = e
                        }
                        mu.Unlock()
                    }
                }(name, outpath)
            }
            wg.Wait()
            return firstErr
        }
    }
    return
}

func describeCmd(gs string) error {
        // check size to decide ranged parallel download
        attrs, err := client.Bucket(bucket).Object(object).Attrs(ctx)
        if err != nil {
            return err
        }
        size := attrs.Size
        if size > int64(defaultRangeThreshold) {
            return rangedDownload(ctx, client, bucket, object, generation, outpath, concurrency)
        }
        return downloadObject(ctx, client, bucket, object, generation, outpath)
    bucket, object, gen, hasGen, err := parseGsURL(gs)
    if err != nil {
        return err
    }
    client, err := storage.NewClient(ctx)
    if err != nil {
        return err
    }
    defer client.Close()

    if object == "" || strings.HasSuffix(object, "/") {
        return fmt.Errorf("object appears to be a directory or empty; describe expects an object")
    }
    obj := client.Bucket(bucket).Object(object)
    if hasGen {
        obj = obj.Generation(gen)
    }
    attrs, err := obj.Attrs(ctx)
    if err != nil {
        return err

    func rangedDownload(ctx context.Context, client *storage.Client, bucket, name string, generation int64, outpath string, concurrency int) error {
        // get attrs
        obj := client.Bucket(bucket).Object(name)
        if generation != 0 {
            obj = obj.Generation(generation)
        }
        attrs, err := obj.Attrs(ctx)
        if err != nil {
            return err
        }
        size := attrs.Size
        part := int64(defaultRangeChunk)
        parts := int((size + part - 1) / part)
        sem := make(chan struct{}, concurrency)
        tmpFiles := make([]string, parts)
        var wg sync.WaitGroup
        var firstErr error
        var mu sync.Mutex
        for i := 0; i < parts; i++ {
            start := int64(i) * part
            length := part
            if i == parts-1 {
                length = size - start
            }
            tmp := fmt.Sprintf("%s.part.%d", outpath, i)
            tmpFiles[i] = tmp
            wg.Add(1)
            sem <- struct{}{}
            go func(s, l int64, tmpfile string) {
                defer wg.Done()
                defer func() { <-sem }()
                r, e := obj.NewRangeReader(ctx, s, l)
                if e != nil {
                    mu.Lock()
                    if firstErr == nil {
                        firstErr = e
                    }
                    mu.Unlock()
                    return
                }
                defer r.Close()
                f, e := os.Create(tmpfile)
                if e != nil {
                    mu.Lock()
                    if firstErr == nil {
                        firstErr = e
                    }
                    mu.Unlock()
                    return
                }
                defer f.Close()
                if _, e = io.Copy(f, r); e != nil {
                    mu.Lock()
                    if firstErr == nil {
                        firstErr = e
                    }
                    mu.Unlock()
                }
            }(start, length, tmp)
        }
        wg.Wait()
        if firstErr != nil {
            return firstErr
        }
        // combine
        out, err := os.Create(outpath)
        if err != nil {
            return err
        }
        defer out.Close()
        for _, tmp := range tmpFiles {
            f, err := os.Open(tmp)
            if err != nil {
                return err
            }
            if _, err := io.Copy(out, f); err != nil {
                f.Close()
                return err
            }
            f.Close()
            os.Remove(tmp)
        }
        return nil
    }
    }
    meta := ObjMeta{
        Generation:     strconv.FormatInt(attrs.Generation, 10),
        CreationTime:   attrs.Created.Format(time.RFC3339Nano),
        UpdateTime:     attrs.Updated.Format(time.RFC3339Nano),
        Md5Hash:        base64.StdEncoding.EncodeToString(attrs.MD5),
        Metageneration: attrs.MetaGeneration,
        StorageURL:     attrs.MediaLink,
        Size:           attrs.Size,
    }
    enc := json.NewEncoder(os.Stdout)
    enc.SetIndent("", "  ")
    return enc.Encode(meta)
}

func downloadCmd(gs, destDir string) error {
    ctx := context.Background()
    bucket, object, gen, hasGen, err := parseGsURL(gs)
    if err != nil {
        return err
    }
    client, err := storage.NewClient(ctx)
    if err != nil {
        return err
    }
    defer client.Close()

    if object == "" || strings.HasSuffix(object, "/") {
        prefix := strings.TrimSuffix(object, "/")
        // collect object names
        var names []string
        it := client.Bucket(bucket).Objects(ctx, &storage.Query{Prefix: prefix + "/"})
        for {
            attrs, err := it.Next()
            if err == iterator.Done {
                break
            }
            if err != nil {
                return err
            }
            names = append(names, attrs.Name)
        }

        // parallel downloads
        concurrency := 4
        if v := os.Getenv("GCS_CONCURRENCY"); v != "" {
            if n, err := strconv.Atoi(v); err == nil && n > 0 {
                concurrency = n
            }
        }

        sem := make(chan struct{}, concurrency)
        errCh := make(chan error, len(names))
        for _, name := range names {
            name := name
            rel := strings.TrimPrefix(name, prefix+"/")
            outpath := filepath.Join(destDir, rel)
            if err := os.MkdirAll(filepath.Dir(outpath), 0o755); err != nil {
                return err
            }
            sem <- struct{}{}
            go func() {
                defer func() { <-sem }()
                if err := downloadObject(ctx, client, bucket, name, 0, outpath); err != nil {
                    errCh <- err
                }
            }()
        }
        // wait for workers
        for i := 0; i < cap(sem); i++ {
            sem <- struct{}{}
        }
        close(errCh)
        if len(errCh) > 0 {
            return <-errCh
        }
        return nil
    }
    var generation int64 = 0
    if hasGen {
        generation = gen
    }
    outpath := filepath.Join(destDir, filepath.Base(object))
    return downloadObject(ctx, client, bucket, object, generation, outpath)
}

func downloadObject(ctx context.Context, client *storage.Client, bucket, name string, generation int64, outpath string) error {
    obj := client.Bucket(bucket).Object(name)
    if generation != 0 {
        obj = obj.Generation(generation)
    }
    r, err := obj.NewReader(ctx)
    if err != nil {
        return err
    }
    defer r.Close()
    f, err := os.Create(outpath)
    if err != nil {
        return err
    }
    defer f.Close()
    _, err = io.Copy(f, r)
    return err
}

func uploadCmd(localPath, gs string) error {
    ctx := context.Background()
    bucket, object, _, _, err := parseGsURL(gs)
    if err != nil {
        return err
    }
    client, err := storage.NewClient(ctx)
    if err != nil {
        return err
    }
    defer client.Close()

    fi, err := os.Stat(localPath)
    if err != nil {
        return err
    }
    if fi.IsDir() {
        prefix := object
        // collect files
        var files []string
        _ = filepath.Walk(localPath, func(path string, info os.FileInfo, e error) error {
            if e != nil {
                return e
            }
            if info.IsDir() {
                return nil
            }
            files = append(files, path)
            return nil
        })

        // concurrency
        concurrency := 4
        if v := os.Getenv("GCS_CONCURRENCY"); v != "" {
            if n, err := strconv.Atoi(v); err == nil && n > 0 {
                concurrency = n
            }
        }

        sem := make(chan struct{}, concurrency)
        errCh := make(chan error, len(files))
        for _, path := range files {
            path := path
            rel, _ := filepath.Rel(localPath, path)
            objName := strings.TrimSuffix(prefix, "/") + "/" + filepath.ToSlash(rel)
            sem <- struct{}{}
            go func() {
                defer func() { <-sem }()
                if err := uploadObject(ctx, client, bucket, objName, path); err != nil {
                    errCh <- err
                }
            }()
        }
        // wait
        for i := 0; i < cap(sem); i++ {
            sem <- struct{}{}
        }
        close(errCh)
        if len(errCh) > 0 {
            return <-errCh
        }
        return nil
    }
    return uploadObject(ctx, client, bucket, object, localPath)
}

func uploadObject(ctx context.Context, client *storage.Client, bucket, name, localPath string) error {
    f, err := os.Open(localPath)
    if err != nil {
        return err
    }
    defer f.Close()
    w := client.Bucket(bucket).Object(name).NewWriter(ctx)
    if _, err := io.Copy(w, f); err != nil {
        _ = w.Close()
        return err
    }
    return w.Close()
}

func main() {
    flag.Parse()
    if flag.NArg() < 1 {
        log.Fatalf("usage: gcs <command> [args]")
    }
    cmd := flag.Arg(0)
    switch cmd {
    case "describe":
        if flag.NArg() != 2 {
            log.Fatalf("usage: gcs describe gs://bucket/object[#generation]")
        }
        if err := describeCmd(flag.Arg(1)); err != nil {
            log.Fatalf("describe: %v", err)
        }
    case "download":
        if flag.NArg() != 3 {
            log.Fatalf("usage: gcs download gs://bucket/object dest_dir")
        }
        if err := downloadCmd(flag.Arg(1), flag.Arg(2)); err != nil {
            log.Fatalf("download: %v", err)
        }
    case "upload":
        if flag.NArg() != 3 {
            log.Fatalf("usage: gcs upload local_path gs://bucket/object")
        }
        if err := uploadCmd(flag.Arg(1), flag.Arg(2)); err != nil {
            log.Fatalf("upload: %v", err)
        }
    default:
        log.Fatalf("unknown command: %s", cmd)
    }
}
