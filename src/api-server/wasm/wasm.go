package wasm

import (
	"bytes"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"strings"

	"github.com/tricorder/src/utils/errors"
	"github.com/tricorder/src/utils/file"
	"github.com/tricorder/src/utils/uuid"
)

const (
	defaultWASISDKPath         = "/opt/tricorder/wasm/wasi-sdk"
	defaultWASIClang           = defaultWASISDKPath + "/bin/clang"
	defaultWASICFlags          = "--sysroot=" + defaultWASISDKPath + "/share/wasi-sysroot"
	defaultWASIStarshipInclude = "/opt/tricorder/wasm/include"
	defaultBuildTmpDir         = "/tmp"
	cJSONDotC                  = "cJSON.c"
)

type WASICompiler struct {
	clangPath   string
	cFlags      string
	includesDir string
	buildTmpDir string
}

func NewWASICompiler(wasiSDKPath string, includeDir string, buildTmpDir string) *WASICompiler {
	return &WASICompiler{
		clangPath:   path.Join(wasiSDKPath, "bin", "clang"),
		cFlags:      "--sysroot=" + path.Join(wasiSDKPath, "share", "wasi-sysroot"),
		includesDir: includeDir,
		buildTmpDir: buildTmpDir,
	}
}

func NewWASICompilerWithDefaults() *WASICompiler {
	return NewWASICompiler(defaultWASISDKPath, defaultWASIStarshipInclude, defaultBuildTmpDir)
}

func (w *WASICompiler) BuildC(code string) ([]byte, error) {
	srcID := strings.Replace(uuid.New(), "-", "_", -1)
	const (
		cExt    = ".c"
		wasmExt = ".wasm"
	)
	srcFilePath := path.Join(w.buildTmpDir, srcID+cExt)
	dstFilePath := path.Join(w.buildTmpDir, srcID+wasmExt)

	// write code string to tmp file
	phase := "write code to " + srcFilePath
	_, err := os.Stat(srcFilePath)
	if errors.Is(err, os.ErrNotExist) {
		content := []byte(code)
		err = ioutil.WriteFile(srcFilePath, content, 0o644)
		if err != nil {
			return nil, errors.Wrap("compile wasm code", phase, err)
		}
	} else if err == nil {
		return nil, errors.New("compile wasm code", phase+" error: File already exists.")
	} else {
		return nil, errors.Wrap("compile wasm code", phase, err)
	}

	// compile code
	phase = "compile " + srcFilePath + " to " + dstFilePath
	cmd := exec.Command(w.clangPath, w.cFlags,
		path.Join(w.includesDir, cJSONDotC), "-I"+w.includesDir, srcFilePath,
		"-Wl,--export-all", "-Wall", "-Wextra", "-o", dstFilePath)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	out, err := cmd.Output()
	if err != nil {
		return nil, errors.Wrap("compile wasm code", phase+" error cc output:\n"+stderr.String(), err)
	}

	if len(out) > 0 {
		return nil, errors.New("compile wasm code", phase+" error: cc output:\n"+string(out))
	}

	// check compiled file if exists
	phase = "check compiled file " + dstFilePath
	_, err = os.Stat(dstFilePath)
	if err != nil {
		return nil, errors.Wrap("compile wasm code", phase, err)
	}

	// check comiled file fmt
	phase = "check compiled file format " + dstFilePath
	if !file.IsWasmELF(dstFilePath) {
		return nil, errors.New("compile wasm code", phase+" error: File is not a wasm file.")
	}

	// read compiled file
	phase = "read compiled file " + dstFilePath
	data, err := ioutil.ReadFile(dstFilePath)
	if err != nil {
		return nil, errors.Wrap("compile wasm code", phase, err)
	}

	// delete tmp files
	phase = "delete tmp files"
	err = os.Remove(srcFilePath)
	if err != nil {
		return nil, errors.Wrap("compile wasm code", phase+" "+srcFilePath, err)
	}
	err = os.Remove(dstFilePath)
	if err != nil {
		return nil, errors.Wrap("compile wasm code", phase+" "+dstFilePath, err)
	}
	return data, nil
}
