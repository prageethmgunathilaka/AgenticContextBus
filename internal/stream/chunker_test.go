package stream

import (
    "context"
    "testing"
)

func TestChunkAndReassemble(t *testing.T) {
    chunker := NewChunker()
    data := make([]byte, 10*1024) // 10KB of zeros

    chunks := chunker.Chunk(data)
    if len(chunks) == 0 {
        t.Fatal("expected chunks > 0")
    }

    reassembled := chunker.Reassemble(chunks)
    if len(reassembled) != len(data) {
        t.Fatalf("expected %d bytes, got %d", len(data), len(reassembled))
    }
}

func TestCalculateChecksum(t *testing.T) {
    a := []byte("hello")
    b := []byte("hello")
    c := []byte("world")

    if CalculateChecksum(a) != CalculateChecksum(b) {
        t.Fatal("checksums of identical data should match")
    }
    if CalculateChecksum(a) == CalculateChecksum(c) {
        t.Fatal("checksums of different data should not match")
    }
}

func TestInitAndCompleteStream(t *testing.T) {
    svc := NewStreamService()
    progress, err := svc.InitStream(context.Background(), &InitStreamRequest{TotalSize: 1234})
    if err != nil {
        t.Fatalf("InitStream error: %v", err)
    }
    if progress.TotalBytes != 1234 {
        t.Fatalf("expected TotalBytes=1234, got %d", progress.TotalBytes)
    }

    // UploadChunk currently returns nil (stub)
    if err := svc.UploadChunk(context.Background(), progress.StreamID, 0, []byte("abc"), false); err != nil {
        t.Fatalf("UploadChunk error: %v", err)
    }

    // CompleteStream should return an error per MVP
    if _, err := svc.CompleteStream(context.Background(), progress.StreamID); err == nil {
        t.Fatal("expected error from CompleteStream")
    }
}


