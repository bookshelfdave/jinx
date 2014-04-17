package jinx

import (
    //"bufio"
    "errors"
    "io"
    "strings"
)

type JinxReader struct {
    R io.ReadSeeker
    Offset int64
}

func NewJinxReaderFromString(s string) *JinxReader {
    string_reader := strings.NewReader(s)
    rs := io.ReadSeeker(string_reader)
    return &JinxReader{rs, 0}
}

func (r *JinxReader) Read(n int) ([]byte, error) {
    buf := make([]byte, n)
    retn, err := r.R.Read(buf)
    if err != nil {
        return nil, err
    }
    if retn != n {
        return nil, errors.New("Not enough data available to Peek")
    }
    r.Offset += int64(n)
    return buf, nil
}

func (r *JinxReader) Peek(n int) ([]byte, error) {
    buf := make([]byte, n)
    retn, err := r.R.Read(buf)
    if err != nil {
        return nil, err
    }
    if retn != n {
        return nil, errors.New("Not enough data available to Peek")
    }
    r.R.Seek(r.Offset, 0)
    return buf, nil
}

func (r *JinxReader) Seek(n int64) error {
    r.Offset = n
    _, err := r.R.Seek(n, 0);
    return err
}
