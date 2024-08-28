package main_test

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"os"

	main "a21hc3NpZ25tZW50"

	"github.com/joho/godotenv"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

type MockDoType func(req *http.Request) (*http.Response, error)

type MockClient struct {
	MockDo        MockDoType
	MockRoundTrip func(req *http.Request) (*http.Response, error)
}

func (m *MockClient) Do(req *http.Request) (*http.Response, error) {
	return m.MockDo(req)
}

func (m *MockClient) RoundTrip(req *http.Request) (*http.Response, error) {
	return m.MockRoundTrip(req)
}

var _ = Describe("Main", func() {
	Describe("csvToSlice", func() {
		It("converts CSV data to a slice 1", func() {
			data := `header1,header2
value1,value2`
			expected := map[string][]string{
				"header1": {"value1"},
				"header2": {"value2"},
			}

			result, err := main.CsvToSlice(data)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(result).Should(Equal(expected))
		})

		It("converts CSV data to a slice 2", func() {
			data := `header3,header4
value3,value4`
			expected := map[string][]string{
				"header3": {"value3"},
				"header4": {"value4"},
			}

			result, err := main.CsvToSlice(data)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(result).Should(Equal(expected))
		})
	})

	Describe("connectAIModel", func() {
		It("connects to the AI model and returns a response 1", func() {
			jsonData := `{"answer": "SUM", "coordinates": [[0, 0]], "cells": ["10"], "aggregator": "SUM"}`
			r := ioutil.NopCloser(bytes.NewReader([]byte(jsonData)))

			mockClient := &MockClient{
				MockDo: func(req *http.Request) (*http.Response, error) {
					return &http.Response{
						StatusCode: 200,
						Body:       r,
					}, nil
				},
				MockRoundTrip: func(req *http.Request) (*http.Response, error) {
					return &http.Response{
						StatusCode: 200,
						Body:       ioutil.NopCloser(bytes.NewReader([]byte(jsonData))),
					}, nil
				},
			}

			connector := &main.AIModelConnector{
				Client: &http.Client{
					Transport: mockClient,
				},
			}

			payload := main.Inputs{
				Table: map[string][]string{
					"header1": {"value1"},
					"header2": {"value2"},
				},
				Query: "What is the total?",
			}

			expected := main.Response{
				Answer:      "SUM",
				Coordinates: [][]int{{0, 0}},
				Cells:       []string{"10"},
				Aggregator:  "SUM",
			}

			err := godotenv.Load()
			Expect(err).ShouldNot(HaveOccurred())

			result, err := connector.ConnectAIModel(payload, os.Getenv("HUGGINGFACE_TOKEN"))
			Expect(err).ShouldNot(HaveOccurred())
			Expect(result).Should(Equal(expected))
		})

		It("connects to the AI model and returns a response 2", func() {
			jsonData := `{"answer": "AVG", "coordinates": [[0, 1]], "cells": ["5", "15"], "aggregator": "AVG"}`
			r := ioutil.NopCloser(bytes.NewReader([]byte(jsonData)))

			mockClient := &MockClient{
				MockDo: func(req *http.Request) (*http.Response, error) {
					return &http.Response{
						StatusCode: 200,
						Body:       r,
					}, nil
				},
				MockRoundTrip: func(req *http.Request) (*http.Response, error) {
					return &http.Response{
						StatusCode: 200,
						Body:       ioutil.NopCloser(bytes.NewReader([]byte(jsonData))),
					}, nil
				},
			}

			connector := &main.AIModelConnector{
				Client: &http.Client{
					Transport: mockClient,
				},
			}

			payload := main.Inputs{
				Table: map[string][]string{
					"header1": {"value3"},
					"header2": {"value4"},
				},
				Query: "What is the average?",
			}

			expected := main.Response{
				Answer:      "AVG",
				Coordinates: [][]int{{0, 1}},
				Cells:       []string{"5", "15"},
				Aggregator:  "AVG",
			}

			err := godotenv.Load()
			Expect(err).ShouldNot(HaveOccurred())

			result, err := connector.ConnectAIModel(payload, os.Getenv("HUGGINGFACE_TOKEN"))
			Expect(err).ShouldNot(HaveOccurred())
			Expect(result).Should(Equal(expected))
		})
	})
})
