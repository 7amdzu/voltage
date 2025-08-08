package handler

import (
	"io"
	"net/http"

	"github.com/7amdzu/voltage/internal/mask"
	"github.com/gin-gonic/gin"
)

const smallThreshold = 1 << 20 // 1 MiB

// Ingest accepts JSON, masks PANs, and writes masked JSON back.
// Uses zero-alloc MaskBytes for small bodies, StreamMasker for large.
func Ingest(c *gin.Context) {
	c.Header("Content-Type", "application/json")

	// If Content-Length known and small, read fully and zero-alloc mask:
	if c.Request.ContentLength > 0 && c.Request.ContentLength < smallThreshold {
		data, err := io.ReadAll(c.Request.Body)
		if err != nil {
			c.Status(http.StatusInternalServerError)
			return
		}
		masked := mask.MaskBytes(data)
		c.Data(200, "application/json", masked)
		return
	}

	// Otherwise, stream-mask (handles any size, splits across chunks):
	reader := mask.NewStreamMasker(c.Request.Body)

	if _, err := io.Copy(c.Writer, reader); err != nil {
		// Log or handle error as needed
		c.Status(http.StatusInternalServerError)
	}
}
