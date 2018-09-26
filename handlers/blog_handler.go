package handlers

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/sscp/telemetry/blog"
	"github.com/sscp/telemetry/events"
	"github.com/sscp/telemetry/log"

	"github.com/opentracing/opentracing-go"
)

// GetBlogFileName wraps GetBaseFileName by adding the .blog extension
func GetBlogFileName(runName string, startTime time.Time) string {
	return fmt.Sprintf("%v.blog", GetBaseFileName(runName, startTime))
}

// BlogWriter is a BinaryHandler (handlers.go) that writes to a .blog file in
// the folder that telemetry is run in
type BlogWriter struct {
	folderPath string
	file       *os.File
	buffer     *bufio.Writer
	blogWriter *blog.Writer
}

// BlogConfig contains config info for the CSVWriter
type BlogConfig struct {
	Folder string
}

// NewBlogWriter allocates a BlogWriter, which is returned as a BinaryHandler
// interface
func NewBlogWriter(cfg BlogConfig) (BinaryHandler, error) {
	err := os.MkdirAll(cfg.Folder, os.ModePerm) // Create folder if it doesn't exist
	if err != nil {
		return nil, err
	}
	return &BlogWriter{folderPath: cfg.Folder}, nil
}

// HandleStartRun is called by collector when data collection starts and
// created the .blog file and the buffered writer
func (bw *BlogWriter) HandleStartRun(ctx context.Context, runName string, startTime time.Time) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "BlogWriter/HandleStartRun")
	defer span.Finish()

	bw.createFile(ctx, runName, startTime)
	bw.buffer = bufio.NewWriter(bw.file)
	bw.blogWriter = blog.NewWriter(bw.buffer)
}

// createFile simply creates a file in the current directory to write binary log data to
func (bw *BlogWriter) createFile(ctx context.Context, runName string, startTime time.Time) {
	filename := GetBlogFileName(runName, startTime)
	var err error
	bw.file, err = os.Create(filepath.Join(bw.folderPath, filename))
	if err != nil {
		log.Error(ctx, err, "Could not create blog file")
	}
}

// HandleRawEvent is called when collector passes off a packet to BlogWriter and
// simply writes the packet the blogWriter
func (bw *BlogWriter) HandleRawEvent(ctx context.Context, rawEvent events.RawEvent) {
	span, _ := opentracing.StartSpanFromContext(ctx, "BlogWriter/HandleData")
	defer span.Finish()
	_, err := bw.blogWriter.Write(rawEvent.Data)
	if err != nil {
		log.Error(ctx, err, "could not write to blog")
	}
}

// HandleDroppedPacket is called when BlogWriter falls behind and cannot
// process an incomming packet, it currently does nothing
func (bw *BlogWriter) HandleDroppedPacket(ctx context.Context) {
	span, _ := opentracing.StartSpanFromContext(ctx, "BlogWriter/HandleDroppedPacket")
	defer span.Finish()
}

// HandleEndRun flushes the buffers and closes the file and is called when the
// collector stops recording packets and BlogWriter has cleared its input
// channel of packets.
func (bw *BlogWriter) HandleEndRun(ctx context.Context, endTime time.Time) {
	span, _ := opentracing.StartSpanFromContext(ctx, "BlogWriter/HandleEndRun")
	defer span.Finish()
	bw.buffer.Flush()
	bw.file.Close()
}
