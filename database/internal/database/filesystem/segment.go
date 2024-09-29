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

	writtenBytes, err := WriteToFile(s.file, data)
	if err != nil {
		return fmt.Errorf("failed to write data to segment file: %w", err)
	}

	s.segmentSize += writtenBytes
	return nil
}

func (s *Segment) rotateSegment() error {
	segmentName := fmt.Sprintf("%s/wal_%d.log", s.directory, now().UnixMilli())
	file, err := CreateFile(segmentName)
	if err != nil {
		return err
	}

	s.file = file
	s.segmentSize = 0
	return nil
}
