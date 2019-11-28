package runner

import (
	"errors"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"sync"
	"time"
)

var (
	// ErrRequestFailed is returned when the response from the server has a non-successful HTTP status code.
	ErrRequestFailed = errors.New("received non-200 response from the server")
	// ErrResponseMalformed is returned when the response from the server could not be read.
	ErrResponseMalformed = errors.New("failed to read response from the server")
	// ErrZeroRate is returned when the runner has been configured to send requests at zero rate.
	ErrZeroRate = errors.New("cannot use 0 as request rate")
)

type httpClient interface {
	Do(*http.Request) (*http.Response, error)
}

// Runner send HTTP requests using provided HTTP request configuration and HTTP client
// at the given rate. If the runner is configured to be verbose, it will log failed requests.
type Runner struct {
	period       time.Duration
	verbose      bool
	request      *http.Request
	client       httpClient
	counterCh    chan interface{}
	errCounterCh chan interface{}
	closeCh      chan interface{}
}

// New returns a new Runner, configured with given request, client, rate and verbose parameters.
// If the rate is 0, an error is returned, since a zero-rate runner is a noop.
func New(
	rate uint,
	verbose bool,
	request *http.Request,
	client httpClient,
	closeCh chan interface{},
) (*Runner, error) {
	if rate < 1 {
		return nil, ErrZeroRate
	}

	return &Runner{
		period:       time.Duration(float64(time.Second) / float64(rate)),
		verbose:      verbose,
		request:      request,
		client:       client,
		counterCh:    make(chan interface{}, rate),
		errCounterCh: make(chan interface{}, rate),
		closeCh:      closeCh,
	}, nil
}

// Start starts sending requests at a given rate using timers.
// Start will print the results of the execution every second.
// Start will block until the runner is stopped using the close channel.
// Sending to close channel will prevent the runner from producing new requests,
// but it will wait for all in-flight requests before returning.
func (r *Runner) Start() {
	var (
		reqCounter uint
		errCounter uint
		wg         sync.WaitGroup
	)

	ticker := time.NewTicker(r.period)
	secondTicker := time.NewTicker(1 * time.Second)

	for {
		select {
		case <-r.counterCh:
			reqCounter++
		case <-r.errCounterCh:
			errCounter++
		case <-secondTicker.C:
			var errRate float64
			if reqCounter > 0 {
				errRate = float64(errCounter) / float64(reqCounter) * 100
			}

			log.Printf("Sent %d requests, %3.2f%% failed\n", reqCounter, errRate)

			reqCounter = 0
			errCounter = 0
		case <-ticker.C:
			wg.Add(1)
			go r.doRequest(&wg)
		case <-r.closeCh:
			ticker.Stop()
			secondTicker.Stop()
			wg.Wait()
			return
		}
	}
}

func (r *Runner) doRequest(wg *sync.WaitGroup) {
	defer wg.Done()

	res, err := r.client.Do(r.request)
	if err != nil {
		r.logRequest(err)
		return
	}

	defer res.Body.Close()
	if _, err := io.Copy(ioutil.Discard, res.Body); err != nil {
		r.logRequest(ErrResponseMalformed)
		return
	}

	// Response statuses between 2xx and 4xx are considered successful.
	if res.StatusCode < http.StatusOK || res.StatusCode > http.StatusPermanentRedirect {
		r.logRequest(ErrRequestFailed)
		return
	}

	r.logRequest(nil)
}

func (r *Runner) logRequest(err error) {
	r.counterCh <- struct{}{}

	if err != nil {
		r.errCounterCh <- struct{}{}

		if r.verbose {
			log.Println(err)
		}
	}
}
