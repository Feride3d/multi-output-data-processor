package service

import (
	"context"
	"fmt"
	"multi-output-data-processor/internal/config"
	"multi-output-data-processor/internal/entity"
	"os"
	"path/filepath"
	"testing"

	"gopkg.in/yaml.v3"

	"github.com/stretchr/testify/require"
)

// TestValidateInput tests ValidateInput function.
func TestValidateInput(t *testing.T) {

	configFilePath := filepath.Join("..", "..", "config", "config.yml")
	configFile, err := os.ReadFile(configFilePath)
	if err != nil {
		t.Fatalf("failed to read config file: %s", err)
	}
	var config config.Config
	err = yaml.Unmarshal(configFile, &config)
	if err != nil {
		t.Fatalf("failed to parse config file: %s", err)
	}

	srv := NewPipelineService(config)

	tests := []struct {
		name   string
		input  entity.InputData
		output error
	}{
		{
			name: "Empty tag field",
			input: entity.InputData{
				Tag:  "",
				Data: "This is a test message from an empty tag test",
			},
			output: entity.ErrEmptyTagParameter,
		},
		{
			name: "Invalid tag field: doesn't match to the tags in config file",
			input: entity.InputData{
				Tag:  "warn",
				Data: "This is a test message from an invalid tag test",
			},
			output: entity.ErrInvalidTagParameter,
		},
		{
			name: "Empty data field",
			input: entity.InputData{
				Tag:  "info",
				Data: "",
			},
			output: entity.ErrEmptyDataParameter,
		},
		{
			name: "Positive test",
			input: entity.InputData{
				Tag:  "info",
				Data: "This is a test message from positive test",
			},
			output: nil,
		},
	}
	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			result := srv.ValidateInput(context.Background(), tc.input)
			require.Equal(t, tc.output, result)
		})
	}
}

// TestSelectOutputCh tests SelectOutputCh function.
func TestSelectOutputCh(t *testing.T) {

	configFilePath := filepath.Join("..", "..", "config", "config.yml")
	configFile, err := os.ReadFile(configFilePath)
	if err != nil {
		t.Fatalf("failed to read config file: %s", err)
	}
	var config config.Config
	err = yaml.Unmarshal(configFile, &config)
	if err != nil {
		t.Fatalf("failed to parse config file: %s", err)
	}

	srv := NewPipelineService(config)

	tests := []struct {
		name   string
		input  entity.InputData
		output []entity.OutputChannel
	}{
		{
			name: "Positive test of error tag",
			input: entity.InputData{
				Tag:  "error",
				Data: "This is a test message from error tag",
			},
			output: []entity.OutputChannel{entity.OutputChannel{Name: "stderr"}, entity.OutputChannel{Name: "file"}},
		},
		{
			name: "Positive test of info tag",
			input: entity.InputData{
				Tag:  "info",
				Data: "This is a test message from info tag",
			},
			output: []entity.OutputChannel{entity.OutputChannel{Name: "stdout"}},
		},
		{
			name: "Positive test of debug tag",
			input: entity.InputData{
				Tag:  "debug",
				Data: "This is a test message from debug tag",
			},
			output: []entity.OutputChannel{entity.OutputChannel{Name: "file"}},
		},
		{
			name: "Positive test of trase tag",
			input: entity.InputData{
				Tag:  "trace",
				Data: "This is a test message from info trace",
			},
			output: []entity.OutputChannel{entity.OutputChannel{Name: "null"}},
		},
	}
	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			result := srv.SelectOutputCh(context.Background(), config, tc.input)
			require.Equal(t, tc.output, result)
		})
	}
}

// TestProcess tests Process function.
func TestProcess(t *testing.T) {

	configFilePath := filepath.Join("..", "..", "config", "config.yml")
	configFile, err := os.ReadFile(configFilePath)
	if err != nil {
		t.Fatalf("failed to read config file: %s", err)
	}
	var config config.Config
	err = yaml.Unmarshal(configFile, &config)
	if err != nil {
		t.Fatalf("failed to parse config file: %s", err)
	}

	srv := NewPipelineService(config)

	tests := []struct {
		name     string
		input    entity.InputData
		outputCh []entity.OutputChannel
		expLogs  []string
		dlQueue  string
		output   error
	}{
		{
			name: "Positive test",
			input: entity.InputData{
				Tag:  "error",
				Data: "This is a test message from error tag",
			},
			outputCh: []entity.OutputChannel{entity.OutputChannel{Name: "stderr"}, entity.OutputChannel{Name: "file"}},
			expLogs:  nil,
			dlQueue:  "",
		},
	}
	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			srv.Process(context.Background(), tc.input, tc.outputCh)
			require.Equal(t, tc.output, nil)
		})
	}
}

// TestWriteToOutputCh tests WriteToOutputCh function.
func TestWriteToOutputCh(t *testing.T) {

	configFilePath := filepath.Join("..", "..", "config", "config.yml")
	configFile, err := os.ReadFile(configFilePath)
	if err != nil {
		t.Fatalf("failed to read config file: %s", err)
	}
	var config config.Config
	err = yaml.Unmarshal(configFile, &config)
	if err != nil {
		t.Fatalf("failed to parse config file: %s", err)
	}

	srv := NewPipelineService(config)

	tests := []struct {
		name   string
		input  entity.InputData
		chName interface{}
		output error
	}{
		{
			name: "Positive test for stdout output",
			input: entity.InputData{
				Data: "This is a test message",
			},
			chName: "stdout",
			output: nil,
		},
		{
			name: "Positive test for stderr output",
			input: entity.InputData{
				Data: "This is a test message",
			},
			chName: "stderr",
			output: nil,
		},
		{
			name: "Positive test for file output",
			input: entity.InputData{
				Data: "This is a test message",
			},
			chName: "file",
			output: nil,
		},
		{
			name: "Positive test for file output",
			input: entity.InputData{
				Data: "This is a test message",
			},
			chName: "file",
			output: nil,
		},
		{
			name: "Positive test for null output",
			input: entity.InputData{
				Data: "This is a test message",
			},
			chName: "null",
			output: nil,
		},
		{
			name: "Test for nil output channel",
			input: entity.InputData{
				Data: "This is a test message",
			},
			chName: nil,
			output: fmt.Errorf("output is nil"),
		},
		{
			name: "Test for invalid output channel name",
			input: entity.InputData{
				Data: "This is a test message",
			},
			chName: "files",
			output: fmt.Errorf("invalid output channel: files"),
		},
	}
	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			result := srv.WriteToOutputCh(context.Background(), tc.input, tc.chName)
			require.Equal(t, tc.output, result)
		})
	}
}
