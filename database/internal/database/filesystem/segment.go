package filesystem

import (
	"fmt"
	"os"
	"time"
)

var now = time.Now

type Segment struct {
	file      *os.File
	directory string

	segmentSize    int
	maxSegmentSize int
}

func NewSegment(directory string, maxSegmentSize int) *Segment {
	return &Segment{
		directory:      directory,
		maxSegmentSize: maxSegmentSize,
	}
}

func (s *Segment) Write(data []byte) error {
	if s.file == nil || s.segmentSize >= s.maxSegmentSize {
		if err := s.rotateSegment(); err != nil {
			return fmt.Errorf("failed to rotate segment file: %w", err)
		}
	}

	writtenBytes, err := s.file.Write(data)
	if err != nil {
		return fmt.Errorf("failed to data to segment file: %w", err)
	}

	s.segmentSize += writtenBytes
	if err = s.file.Sync(); err != nil {
		return fmt.Errorf("failed to sync segment file: %w", err)
	}

	return nil
}

func (s *Segment) rotateSegment() error {
	segmentName := fmt.Sprintf("%s/wal_%d.log", s.directory, now().UnixMilli())

	flags := os.O_CREATE | os.O_WRONLY
	file, err := os.OpenFile(segmentName, flags, 0644)
	if err != nil {
		return err
	}

	s.file = file
	s.segmentSize = 0
	return nil
}
