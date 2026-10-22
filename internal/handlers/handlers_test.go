package handlers

import (
	"bytes"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestHandleMain(t *testing.T) {
	t.Chdir("../..")

	want, err := os.ReadFile(htmlFileName)
	if err != nil {
		t.Fatalf("read %s: %v", htmlFileName, err)
	}

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()

	HandleMain(rec, req)

	res := rec.Result()
	defer func() { _ = res.Body.Close() }()

	if res.StatusCode != http.StatusOK {
		t.Fatalf("status = %d, want %d", res.StatusCode, http.StatusOK)
	}

	if ct := res.Header.Get("Content-Type"); !strings.Contains(ct, "text/html") {
		t.Fatalf("Content-Type = %q, want it to contain %q", ct, "text/html")
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		t.Fatalf("read body: %v", err)
	}

	if !bytes.Equal(body, want) {
		t.Fatalf("body mismatch:\n got: %q\nwant: %q", body, want)
	}
}

func TestHandleUpload(t *testing.T) {
	t.Chdir("../..")

	tests := []struct {
		name     string
		filename string
		content  string
		want     string
	}{
		{
			name:     "text to morse",
			filename: "test12",
			content:  "Привет",
			want:     ".--. .-. .. .-- . -",
		},
		{
			name:     "text with extension to morse",
			filename: "test12.txt",
			content:  "ПРИВЕТ",
			want:     ".--. .-. .. .-- . -",
		},
		{
			name:     "morse to text",
			filename: "test",
			content:  ".--. .-. .. .-- . -",
			want:     "ПРИВЕТ",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			before := listTmpFiles(t)

			req := newUploadRequest(t, tt.filename, tt.content)
			rec := httptest.NewRecorder()

			HandleUpload(rec, req)

			res := rec.Result()
			defer func() { _ = res.Body.Close() }()

			if res.StatusCode != http.StatusOK {
				t.Fatalf("status = %d, want %d", res.StatusCode, http.StatusOK)
			}

			body, err := io.ReadAll(res.Body)
			if err != nil {
				t.Fatalf("read body: %v", err)
			}

			if got := string(body); got != tt.want {
				t.Fatalf("body = %q, want %q", got, tt.want)
			}

			created := newTmpFiles(t, before)
			if len(created) != 1 {
				t.Fatalf("created %d files in tmp, want 1 (%v)", len(created), created)
			}

			if ext, wantExt := filepath.Ext(created[0]), filepath.Ext(tt.filename); !strings.HasSuffix(created[0], wantExt) || (wantExt != "" && ext != wantExt) {
				t.Fatalf("created file %q ext = %q, want %q", created[0], ext, wantExt)
			}
		})
	}
}

func TestHandleUploadNoFile(t *testing.T) {
	t.Chdir("../..")

	var buf bytes.Buffer
	w := multipart.NewWriter(&buf)
	if err := w.Close(); err != nil {
		t.Fatalf("close writer: %v", err)
	}

	req := httptest.NewRequest(http.MethodPost, "/upload", &buf)
	req.Header.Set("Content-Type", w.FormDataContentType())
	rec := httptest.NewRecorder()

	HandleUpload(rec, req)

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusInternalServerError)
	}
}

func newUploadRequest(t *testing.T, filename, content string) *http.Request {
	t.Helper()

	var buf bytes.Buffer
	w := multipart.NewWriter(&buf)

	fw, err := w.CreateFormFile("myFile", filename)
	if err != nil {
		t.Fatalf("create form file: %v", err)
	}

	if _, err := io.WriteString(fw, content); err != nil {
		t.Fatalf("write form file: %v", err)
	}

	if err := w.Close(); err != nil {
		t.Fatalf("close writer: %v", err)
	}

	req := httptest.NewRequest(http.MethodPost, "/upload", &buf)
	req.Header.Set("Content-Type", w.FormDataContentType())

	return req
}

func listTmpFiles(t *testing.T) map[string]struct{} {
	t.Helper()

	entries, err := os.ReadDir("tmp")
	if err != nil {
		t.Fatalf("read tmp dir: %v", err)
	}

	files := make(map[string]struct{}, len(entries))
	for _, e := range entries {
		if !e.IsDir() {
			files[e.Name()] = struct{}{}
		}
	}

	return files
}

func newTmpFiles(t *testing.T, before map[string]struct{}) []string {
	t.Helper()

	var created []string
	for name := range listTmpFiles(t) {
		if _, ok := before[name]; !ok {
			created = append(created, name)
		}
	}

	t.Cleanup(func() {
		for _, name := range created {
			_ = os.Remove(filepath.Join("tmp", name))
		}
	})

	return created
}
