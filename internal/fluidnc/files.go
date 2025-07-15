package fluidnc

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

// ListFiles lists files on the FluidNC filesystem
func (c *Client) ListFiles() (*FileListResponse, error) {
	url := fmt.Sprintf("http://%s:%d/files", c.config.Host, c.config.Port)
	resp, err := c.client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to list files: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("list files failed with status %d: %s", resp.StatusCode, string(bodyBytes))
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	// Try to parse as JSON first, if that fails, parse as plain text
	var fileListResp FileListResponse
	if err := json.Unmarshal(bodyBytes, &fileListResp); err != nil {
		// Parse plain text response (common format from FluidNC)
		fileListResp = c.parseFileListText(string(bodyBytes))
	}

	return &fileListResp, nil
}

// parseFileListText parses plain text file listing response
func (c *Client) parseFileListText(response string) FileListResponse {
	var files []FileInfo
	lines := strings.Split(response, "\n")

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "[") {
			continue
		}

		// Parse typical format: "filename.ext (1234 bytes)"
		parts := strings.Fields(line)
		if len(parts) >= 1 {
			name := parts[0]
			var size int64

			// Extract size if present
			if len(parts) > 1 {
				sizeStr := strings.Trim(strings.Join(parts[1:], " "), "()")
				if strings.Contains(sizeStr, "bytes") {
					sizeParts := strings.Fields(sizeStr)
					if len(sizeParts) > 0 {
						if s, err := strconv.ParseInt(sizeParts[0], 10, 64); err == nil {
							size = s
						}
					}
				}
			}

			fileType := "file"
			if strings.HasSuffix(name, "/") {
				fileType = "directory"
				name = strings.TrimSuffix(name, "/")
			}

			files = append(files, FileInfo{
				Name: name,
				Size: size,
				Type: fileType,
			})
		}
	}

	return FileListResponse{
		Files: files,
		Path:  "/",
	}
}

// UploadFile uploads a file to FluidNC
func (c *Client) UploadFile(filePath, destination, endpoint string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	var body bytes.Buffer
	writer := multipart.NewWriter(&body)

	part, err := writer.CreateFormFile("file", filepath.Base(filePath))
	if err != nil {
		return fmt.Errorf("failed to create form file: %w", err)
	}

	if _, err := io.Copy(part, file); err != nil {
		return fmt.Errorf("failed to copy file: %w", err)
	}

	if destination != "" {
		if err := writer.WriteField("filename", destination); err != nil {
			return fmt.Errorf("failed to write filename field: %w", err)
		}
	}

	if err := writer.Close(); err != nil {
		return fmt.Errorf("failed to close writer: %w", err)
	}

	url := fmt.Sprintf("http://%s:%d/%s", c.config.Host, c.config.Port, endpoint)
	req, err := http.NewRequest("POST", url, &body)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())

	resp, err := c.client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to upload file: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("upload failed with status %d: %s", resp.StatusCode, string(bodyBytes))
	}

	if c.config.Verbose {
		fmt.Printf("File uploaded successfully to %s\n", endpoint)
	}

	return nil
}

// UploadToLocalFS uploads a file to the local filesystem via /files endpoint
func (c *Client) UploadToLocalFS(filePath string, destinationName string) error {
	return c.UploadFile(filePath, destinationName, "files")
}

// UploadToSD uploads a file directly to SD card via /upload endpoint
func (c *Client) UploadToSD(filePath string, destinationName string) error {
	return c.UploadFile(filePath, destinationName, "upload")
}

// UpdateFirmware uploads firmware file for update via /updatefw endpoint
func (c *Client) UpdateFirmware(firmwarePath string) error {
	file, err := os.Open(firmwarePath)
	if err != nil {
		return fmt.Errorf("failed to open firmware file: %w", err)
	}
	defer file.Close()

	var body bytes.Buffer
	writer := multipart.NewWriter(&body)

	part, err := writer.CreateFormFile("file", filepath.Base(firmwarePath))
	if err != nil {
		return fmt.Errorf("failed to create form file: %w", err)
	}

	if _, err := io.Copy(part, file); err != nil {
		return fmt.Errorf("failed to copy firmware file: %w", err)
	}

	if err := writer.Close(); err != nil {
		return fmt.Errorf("failed to close writer: %w", err)
	}

	url := fmt.Sprintf("http://%s:%d/updatefw", c.config.Host, c.config.Port)
	req, err := http.NewRequest("POST", url, &body)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())

	resp, err := c.client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to upload firmware: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("firmware update failed with status %d: %s", resp.StatusCode, string(bodyBytes))
	}

	if c.config.Verbose {
		fmt.Println("Firmware update started successfully")
	}

	return nil
}
