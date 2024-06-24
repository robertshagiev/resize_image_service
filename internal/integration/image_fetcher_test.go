package integration

import (
	"bytes"
	"context"
	"errors"
	"io"
	"net/http"
	"testing"
	"time"

	"resize_image_service/internal/integration/mocks"
	"resize_image_service/internal/model"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestFetchImage(t *testing.T) {
	type input struct {
		url string
	}

	type expected struct {
		imageData *model.ImageData
		err       error
	}

	var testCases = []struct {
		name     string
		input    input
		expected expected
	}{
		{
			name: "valid image URL",
			input: input{
				url: "https://upload.wikimedia.org/wikipedia/commons/thumb/0/01/Bufo_bufo_03-clean.jpg/275px-Bufo_bufo_03-clean.jpg",
			},
			expected: expected{
				imageData: &model.ImageData{
					ContentType: "image/jpeg",
					URL:         "https://upload.wikimedia.org/wikipedia/commons/thumb/0/01/Bufo_bufo_03-clean.jpg/275px-Bufo_bufo_03-clean.jpg",
					Width:       275,
					Height:      183,
				},
				err: nil,
			},
		},
		{
			name: "invalid image URL",
			input: input{
				url: "https://popa.ru/popa.jpeg",
			},
			expected: expected{
				imageData: nil,
				err:       errors.New("error sending request: status code: 404, body: Not Found"),
			},
		},
	}

	timeout := 10 * time.Second

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			clientHttp := mocks.NewClientHttp(t)
			fetchImage := &ImageFetcher{
				client:  clientHttp,
				timeout: timeout,
			}

			if tc.expected.err == nil {
				resp, err := http.Get(tc.input.url)
				require.NoError(t, err)
				defer resp.Body.Close()

				imageData, err := io.ReadAll(resp.Body)
				require.NoError(t, err)

				mockedResp := &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewReader(imageData)),
					Header:     http.Header{"Content-Type": []string{tc.expected.imageData.ContentType}},
				}
				clientHttp.On("Do", mock.AnythingOfType("*http.Request")).Return(mockedResp, nil)
			} else {
				clientHttp.On("Do", mock.AnythingOfType("*http.Request")).Return(nil, tc.expected.err)
			}

			ctx := context.Background()
			result, err := fetchImage.FetchImage(ctx, tc.input.url)

			if tc.expected.err != nil {
				require.Error(t, err)
				require.Contains(t, err.Error(), "status code: 404, body: Not Found")
			} else {
				require.NoError(t, err)
				require.NotNil(t, result)
				require.Equal(t, tc.expected.imageData.ContentType, result.ContentType)
				require.Equal(t, tc.expected.imageData.URL, result.URL)
				require.Equal(t, tc.expected.imageData.Width, result.Width)
				require.Equal(t, tc.expected.imageData.Height, result.Height)
			}

		})
	}
}
