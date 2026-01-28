package dashgram

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"
)

// asyncTask represents a task to be executed asynchronously
type asyncTask struct {
	ctx      context.Context
	endpoint string
	data     any
}

// HttpClient is an interface that wraps the Do method
type HttpClient interface {
	Do(req *http.Request) (*http.Response, error)
}

// Dashgram represents the main client for interacting with the Dashgram API
type Dashgram struct {
	ProjectID int
	AccessKey string
	APIURL    string
	Origin    string
	client    HttpClient

	// Async worker
	useAsync     bool
	numWorkers   int
	workerCtx    context.Context
	workerCancel context.CancelFunc
	taskChan     chan asyncTask
	workerWg     sync.WaitGroup
}

// New creates a new Dashgram client instance
func New(projectID int, accessKey string, options ...Option) *Dashgram {
	ctx, cancel := context.WithCancel(context.Background())

	d := &Dashgram{
		ProjectID: projectID,
		AccessKey: accessKey,
		APIURL:    "https://api.dashgram.io/v1",
		Origin:    "Go + Dashgram SDK",
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
		useAsync:     false,
		numWorkers:   1,
		workerCtx:    ctx,
		workerCancel: cancel,
		taskChan:     make(chan asyncTask, 1000), // Buffer for 1000 tasks
	}

	// Apply options
	for _, option := range options {
		option(d)
	}

	// Set up API URL with project ID
	d.APIURL = fmt.Sprintf("%s/%d", d.APIURL, d.ProjectID)

	// Start the async worker
	d.StartWorker()

	return d
}

// Close stops the async worker and waits for pending tasks
func (d *Dashgram) Close() {
	d.workerCancel()
	d.workerWg.Wait()
}

// startWorker starts the background worker goroutine
func (d *Dashgram) StartWorker() {
	d.workerWg.Add(1)
	go func() {
		defer d.workerWg.Done()
		for {
			select {
			case task := <-d.taskChan:
				d.request(task.ctx, task.endpoint, task.data)
			case <-d.workerCtx.Done():
				return
			}
		}
	}()
}

// Option is a function type for configuring Dashgram client options
type Option func(*Dashgram)

// WithAPIURL sets a custom API URL
func WithAPIURL(apiURL string) Option {
	return func(d *Dashgram) {
		d.APIURL = apiURL
	}
}

// WithOrigin sets a custom origin string
func WithOrigin(origin string) Option {
	return func(d *Dashgram) {
		d.Origin = origin
	}
}

// WithHTTPClient sets a custom HTTP client
func WithHTTPClient(client HttpClient) Option {
	return func(d *Dashgram) {
		d.client = client
	}
}

// WithUseAsync enables asynchronous requests
func WithUseAsync() Option {
	return func(d *Dashgram) {
		d.useAsync = true
	}
}

// WithNumWorkers sets the number of workers for asynchronous requests
func WithNumWorkers(numWorkers int) Option {
	return func(d *Dashgram) {
		d.numWorkers = numWorkers
	}
}

// request makes an HTTP request to the Dashgram API
func (d *Dashgram) request(ctx context.Context, endpoint string, data any) error {
	// Prepare request body
	var body io.Reader
	if data != nil {
		jsonData, err := json.Marshal(data)
		if err != nil {
			return fmt.Errorf("failed to marshal request data: %w", err)
		}
		body = bytes.NewBuffer(jsonData)
	}

	// Create request
	req, err := http.NewRequestWithContext(ctx, "POST", fmt.Sprintf("%s/%s", d.APIURL, endpoint), body)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", d.AccessKey))
	req.Header.Set("Content-Type", "application/json")

	// Make request
	resp, err := d.client.Do(req)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	// Read response body
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode == http.StatusForbidden {
		return &InvalidCredentialsError{}
	}

	var response struct {
		Status  string `json:"status"`
		Details string `json:"details"`
	}

	if err := json.Unmarshal(respBody, &response); err != nil {
		return fmt.Errorf("failed to parse response: %w", err)
	}

	// Check if status code is in 2xx range (200-299)
	if resp.StatusCode < 200 || resp.StatusCode >= 300 || response.Status != "success" {
		return &DashgramAPIError{
			StatusCode: resp.StatusCode,
			Details:    response.Details,
		}
	}

	return nil
}
