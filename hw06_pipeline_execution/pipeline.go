package hw06_pipeline_execution //nolint:golint,stylecheck
import (
	"sync"
)

type (
	I   = interface{}
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

type workerParams struct {
	stageFunc Stage
	outTemp   Out
	inCh      chan interface{}
	outCh     chan interface{}
	mu        sync.Mutex
}

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	stageCh := make([]chan I, len(stages)+1)
	for i := 0; i < len(stageCh); i++ {
		stageCh[i] = make(chan I)
	}
	tempOutChannels := make([]Out, len(stages))

	// Main output channel
	out := make(Bi)

	go func() {
		sourceWorker(done, in, stageCh[0])
		close(stageCh[0])
	}()
	go func() {
		sinkWorker(done, out, stageCh[len(stageCh)-1])
		close(out)
	}()

	for i := 0; i < len(stages); i++ {
		worker := &workerParams{
			stageFunc: stages[i],
			outTemp:   tempOutChannels[i],
			inCh:      stageCh[i],
			outCh:     stageCh[i+1],
		}

		go process(done, worker)

		// Sync out channel with next stage in
		go func(params *workerParams) {
			defer close(params.outCh)
			defer params.mu.Unlock()
			params.mu.Lock()

			for val := range params.outTemp {
				params.outCh <- val
			}
		}(worker)
	}

	return out
}

// Main worker: read from In-Channel and put Data to temporary Out channel of the stage.
func process(done In, params *workerParams) {
	input := make(Bi)
	defer close(input)

	// Need mutex because there is an error when we write in out channel,
	// which is being read in SYNC go-routine
	params.mu.Lock()
	params.outTemp = params.stageFunc(input)
	params.mu.Unlock()

	for {
		select {
		case <-done:
			return
		case val, ok := <-params.inCh:
			if !ok {
				return
			}
			input <- val
		}
	}
}

// Zero-stage: transfer In-data to the first stage of pipeline.
func sourceWorker(done In, in In, outCh chan interface{}) {
	for {
		select {
		case <-done:
			return
		case payload, ok := <-in:
			if !ok {
				return
			}
			outCh <- payload
		}
	}
}

func sinkWorker(done In, out chan<- interface{}, inCh <-chan interface{}) {
	for {
		select {
		case <-done:
			return
		case payload, ok := <-inCh:
			if !ok {
				return
			}
			out <- payload
		}
	}
}
