package osint

import (
	"context"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/rwcarlsen/goexif/exif"
	"github.com/EmmmmDeee/One-Liner-OSINT/pkg/config"
	"github.com/EmmmmDeee/One-Liner-OSINT/pkg/logger"
	"github.com/EmmmmDeee/One-Liner-OSINT/pkg/output"
	"github.com/EmmmmDeee/One-Liner-OSINT/pkg/utils"
)

// ImageModule provides image OSINT capabilities
type ImageModule struct {
	config *config.Config
	logger *logger.Logger
	client *utils.HTTPClient
}

// NewImageModule creates a new image module
func NewImageModule(cfg *config.Config) *ImageModule {
	return &ImageModule{
		config: cfg,
		logger: logger.NewLogger(),
		client: utils.NewHTTPClient(cfg.Timeout, cfg.UserAgent),
	}
}

// Name returns the module name
func (m *ImageModule) Name() string {
	return "image"
}

// Description returns the module description
func (m *ImageModule) Description() string {
	return "Image metadata extraction and analysis"
}

// Execute performs image OSINT
func (m *ImageModule) Execute(ctx context.Context, target string) ([]output.Result, error) {
	var results []output.Result

	// Check if target is a URL or file path
	var reader io.ReadCloser
	var err error

	if strings.HasPrefix(target, "http://") || strings.HasPrefix(target, "https://") {
		// Download image from URL
		resp, err := m.client.Get(ctx, target)
		if err != nil {
			return nil, fmt.Errorf("failed to download image: %v", err)
		}
		reader = resp.Body
	} else {
		// Open local file
		file, err := os.Open(target)
		if err != nil {
			return nil, fmt.Errorf("failed to open image: %v", err)
		}
		reader = file
	}
	defer reader.Close()

	// Extract EXIF data
	exifResult, err := m.extractEXIF(reader, target)
	if err != nil {
		m.logger.Debugf("EXIF extraction failed: %v", err)
		return results, nil
	}

	results = append(results, exifResult)
	return results, nil
}

// extractEXIF extracts EXIF metadata from image
func (m *ImageModule) extractEXIF(reader io.Reader, target string) (output.Result, error) {
	x, err := exif.Decode(reader)
	if err != nil {
		return output.Result{}, err
	}

	data := make(map[string]interface{})

	// Extract common EXIF fields
	if cameraMake, err := x.Get(exif.Make); err == nil {
		if val, err := cameraMake.StringVal(); err == nil {
			data["camera_make"] = val
		}
	}

	if cameraModel, err := x.Get(exif.Model); err == nil {
		if val, err := cameraModel.StringVal(); err == nil {
			data["camera_model"] = val
		}
	}

	if dateTime, err := x.Get(exif.DateTime); err == nil {
		if val, err := dateTime.StringVal(); err == nil {
			data["date_time"] = val
		}
	}

	if dateTimeOriginal, err := x.Get(exif.DateTimeOriginal); err == nil {
		if val, err := dateTimeOriginal.StringVal(); err == nil {
			data["date_time_original"] = val
		}
	}

	if software, err := x.Get(exif.Software); err == nil {
		if val, err := software.StringVal(); err == nil {
			data["software"] = val
		}
	}

	if orientation, err := x.Get(exif.Orientation); err == nil {
		if val, err := orientation.Int(0); err == nil {
			data["orientation"] = val
		}
	}

	if width, err := x.Get(exif.PixelXDimension); err == nil {
		if val, err := width.Int(0); err == nil {
			data["width"] = val
		}
	}

	if height, err := x.Get(exif.PixelYDimension); err == nil {
		if val, err := height.Int(0); err == nil {
			data["height"] = val
		}
	}

	// Extract GPS coordinates
	lat, lon, err := x.LatLong()
	if err == nil {
		data["gps_latitude"] = lat
		data["gps_longitude"] = lon
		data["gps_coordinates"] = fmt.Sprintf("%f,%f", lat, lon)
		data["google_maps_url"] = fmt.Sprintf("https://www.google.com/maps?q=%f,%f", lat, lon)
	}

	return output.Result{
		Timestamp: time.Now(),
		Source:    "image",
		Type:      "exif",
		Target:    target,
		Data:      data,
	}, nil
}

// ExtractGPSCoordinates extracts GPS coordinates from image
func (m *ImageModule) ExtractGPSCoordinates(imagePath string) (float64, float64, error) {
	file, err := os.Open(imagePath)
	if err != nil {
		return 0, 0, err
	}
	defer file.Close()

	x, err := exif.Decode(file)
	if err != nil {
		return 0, 0, err
	}

	return x.LatLong()
}
