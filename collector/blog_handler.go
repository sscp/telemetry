package collector

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"time"

	"github.com/sscp/telemetry/blog"

	"github.com/opentracing/opentracing-go"
)

// GetBlogFileName wraps GetBaseFileName by adding the .blog extension
func GetBlogFileName(runName string, startTime time.Time) string {
	return fmt.Sprintf("%v.blog", GetBaseFileName(runName, startTime))
}

// BlogWriter is a BinaryHandler (handlers.go) that writes to a .blog file in
// the folder that telemetry is run in
type BlogWriter struct {
	file       *os.File
	buffer     *bufio.Writer
	blogWriter *blog.Writer
}

// NewBlogWriter allocates a BlogWriter, which is returned as a BinaryHandler
// interface
func NewBlogWriter() BinaryHandler {
	return &BlogWriter{}
}

// HandleStartRun is called by collector when data collection starts and
// created the .blog file and the buffered writer
func (bw *BlogWriter) HandleStartRun(ctx context.Context, runName string, startTime time.Time) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "BlogWriter/HandleStartRun")
	defer span.Finish()

	bw.createFile(runName, startTime)
	bw.buffer = bufio.NewWriter(bw.file)
	bw.blogWriter = blog.NewWriter(bw.buffer)
}

// createFile simply creates a file in the current directory to write binary log data to
func (bw *BlogWriter) createFile(runName string, startTime time.Time) {
	filename := GetBlogFileName(runName, startTime)
	var err error
	bw.file, err = os.Create(filename)
	if err != nil {
		panic(err)
	}
}

// HandlePacket is called when collector passes off a packet to BlogWriter and
// simply writes the packet the blogWriter
func (bw *BlogWriter) HandlePacket(ctx context.Context, packet []byte) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "BlogWriter/HandleData")
	defer span.Finish()
	bw.blogWriter.Write(packet)
}

// HandleDroppedPacket is called when BlogWriter falls behind and cannot
// process an incomming packet, it currently does nothing
func (bw *BlogWriter) HandleDroppedPacket(ctx context.Context) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "BlogWriter/HandleDroppedPacket")
	defer span.Finish()

}

// HandleEndRun flushes the buffers and closes the file and is called when the
// collector stops recording packets and BlogWriter has cleared its input
// channel of packets.
func (bw *BlogWriter) HandleEndRun(ctx context.Context, endTime time.Time) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "BlogWriter/HandleEndRun")
	defer span.Finish()
	bw.buffer.Flush()
	bw.file.Close()
}
