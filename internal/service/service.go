package service

import (
	"context"
	"fmt"
	"io"
	"log"
	"multi-output-data-processor/internal/config"
	"multi-output-data-processor/internal/entity"
	"os"
	"sync"
	"time"
)

type Pipeliner interface {
	ValidateInput(ctx context.Context, input entity.InputData) error
	SelectOutputCh(ctx context.Context, config config.Config, input entity.InputData) []entity.OutputChannel
	Process(ctx context.Context, input entity.InputData, outputChannels []entity.OutputChannel)
	WriteToOutputCh(ctx context.Context, input entity.InputData, channelName interface{}) error
}

type PipelineService struct {
	config config.Config
}

func NewPipelineService(config config.Config) *PipelineService {
	srv := &PipelineService{
		config: config,
	}

	return srv
}

// ValidateInput validates the input data fields (tag and data).
func (p PipelineService) ValidateInput(ctx context.Context, input entity.InputData) error {

	// Check the tag field (it mustn't be empty).
	if input.Tag == "" {
		return entity.ErrEmptyTagParameter
	}

	validTag := false

	for _, cfg := range p.config.Conf {
		if cfg.Tag == input.Tag {
			validTag = true
			break
		}
	}

	if !validTag {
		return entity.ErrInvalidTagParameter
	}

	// Check the data field (it mustn't be empty).
	if input.Data == "" {
		return entity.ErrEmptyDataParameter
	}

	return nil
}

// SelectOutputCh selects the output channels for the input data.
func (p PipelineService) SelectOutputCh(ctx context.Context, config config.Config, input entity.InputData) []entity.OutputChannel {
	var outputCh []entity.OutputChannel

	for _, conf := range config.Conf {
		if conf.Tag == input.Tag {
			for _, output := range conf.Outputs {
				outputCh = append(outputCh, entity.OutputChannel{Name: output})
			}
		}
	}

	return outputCh
}

// Process processes data to the output channels concurrently.
func (p PipelineService) Process(ctx context.Context, input entity.InputData, outputCh []entity.OutputChannel) {
	var wg sync.WaitGroup
	wg.Add(len(outputCh))

	for _, ch := range outputCh {
		go func(ch entity.OutputChannel) {
			defer wg.Done()

			// Retry 3 times if the output channel is not working or there are errors writing into it.
			for i := 0; i < 3; i++ {
				err := p.WriteToOutputCh(ctx, input, ch.Name)
				if err == nil {
					return
				}

				log.Printf("error writing into %v: %v", ch.Name, err.Error())
				time.Sleep(3 * time.Second)
			}

			// Log and add data to the dead-letter queue
			// after failing to write to the output channel after retries.
			log.Printf("failed to write to %v after retries", ch.Name)
			file, err := os.OpenFile("./examples/dead-letter-queue.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
			if err != nil {
				return
			}
			defer file.Close()

			_, err = fmt.Fprintln(file, input.Data)
			if err != nil {
				return
			}

		}(ch)
	}

	wg.Wait()
}

// WriteToOutputCh writes data to the output channels.
func (p PipelineService) WriteToOutputCh(ctx context.Context, input entity.InputData, channelName interface{}) error {
	switch channelName {
	case "stdout":
		fmt.Println(input.Data)
	case "stderr":
		fmt.Fprintln(os.Stderr, input.Data)
	case "file":
		file, err := os.OpenFile("./file.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			return err
		}
		defer file.Close()

		_, err = fmt.Fprintln(file, input.Data)
		if err != nil {
			return err
		}
	case "null":
		_, err := fmt.Fprintln(io.Discard, input.Data)
		if err != nil {
			return err
		}
	case nil:
		return fmt.Errorf("output is nil")
	default:
		return fmt.Errorf("invalid output channel: %v", channelName)
	}

	return nil
}
