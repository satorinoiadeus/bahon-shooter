package main

import (
    "archive/zip"
    "embed"
    "fmt"
    "io"
    "os"
    "os/exec"
    "path/filepath"
    "strings"
    "syscall"
    "unsafe"
)

//go:embed payload.zip
var payloadFS embed.FS

var buildVersion = "dev"

func messageBox(title, text string) {
    user32 := syscall.NewLazyDLL("user32.dll")
    proc := user32.NewProc("MessageBoxW")
    t, _ := syscall.UTF16PtrFromString(title)
    m, _ := syscall.UTF16PtrFromString(text)
    proc.Call(0, uintptr(unsafe.Pointer(m)), uintptr(unsafe.Pointer(t)), 0x10)
}

func fail(err error) {
    messageBox("Bahon Shooter", err.Error())
    os.Exit(1)
}

func extractPayload(dst string) error {
    data, err := payloadFS.ReadFile("payload.zip")
    if err != nil {
        return err
    }

    tmp, err := os.CreateTemp("", "bahon-payload-*.zip")
    if err != nil {
        return err
    }
    tmpPath := tmp.Name()
    defer os.Remove(tmpPath)
    if _, err := tmp.Write(data); err != nil {
        tmp.Close()
        return err
    }
    if err := tmp.Close(); err != nil {
        return err
    }

    zr, err := zip.OpenReader(tmpPath)
    if err != nil {
        return err
    }
    defer zr.Close()

    if err := os.RemoveAll(dst); err != nil {
        return err
    }
    if err := os.MkdirAll(dst, 0755); err != nil {
        return err
    }

    cleanRoot := filepath.Clean(dst) + string(os.PathSeparator)
    for _, f := range zr.File {
        outPath := filepath.Join(dst, f.Name)
        if !strings.HasPrefix(filepath.Clean(outPath)+string(os.PathSeparator), cleanRoot) {
            return fmt.Errorf("invalid embedded path: %s", f.Name)
        }
        if f.FileInfo().IsDir() {
            if err := os.MkdirAll(outPath, 0755); err != nil {
                return err
            }
            continue
        }
        if err := os.MkdirAll(filepath.Dir(outPath), 0755); err != nil {
            return err
        }
        in, err := f.Open()
        if err != nil {
            return err
        }
        out, err := os.OpenFile(outPath, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, f.Mode())
        if err != nil {
            in.Close()
            return err
        }
        _, copyErr := io.Copy(out, in)
        closeErr1 := out.Close()
        closeErr2 := in.Close()
        if copyErr != nil {
            return copyErr
        }
        if closeErr1 != nil {
            return closeErr1
        }
        if closeErr2 != nil {
            return closeErr2
        }
    }
    return os.WriteFile(filepath.Join(dst, ".version"), []byte(buildVersion), 0644)
}

func main() {
    cache, err := os.UserCacheDir()
    if err != nil {
        fail(err)
    }
    runtimeDir := filepath.Join(cache, "BahonShooter", "runtime")
    versionFile := filepath.Join(runtimeDir, ".version")
    current, _ := os.ReadFile(versionFile)
    if string(current) != buildVersion {
        if err := extractPayload(runtimeDir); err != nil {
            fail(fmt.Errorf("could not unpack game runtime: %w", err))
        }
    }

    gzdoom := filepath.Join(runtimeDir, "gzdoom.exe")
    iwad := filepath.Join(runtimeDir, "freedoom2.wad")
    mod := filepath.Join(runtimeDir, "bahon-shooter.pk3")

    for _, p := range []string{gzdoom, iwad, mod} {
        if _, err := os.Stat(p); err != nil {
            fail(fmt.Errorf("missing game file: %s", p))
        }
    }

    cmd := exec.Command(gzdoom,
        "-iwad", iwad,
        "-file", mod,
        "+map", "MAP01",
        "-config", filepath.Join(runtimeDir, "bahon-shooter.ini"),
    )
    cmd.Dir = runtimeDir
    if err := cmd.Start(); err != nil {
        fail(fmt.Errorf("could not launch GZDoom: %w", err))
    }
}
