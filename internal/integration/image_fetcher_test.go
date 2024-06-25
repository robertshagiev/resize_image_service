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

	"github.com/go-test/deep"
	"github.com/stretchr/testify/mock"
)

// func TestFetchImage(t *testing.T) {
// 	type input struct {
// 		url string
// 	}

// 	type expected struct {
// 		imageData *model.ImageData
// 		err       error
// 	}

// 	var testCases = []struct {
// 		name     string
// 		input    input
// 		expected expected
// 	}{
// 		{
// 			name: "valid image URL",
// 			input: input{
// 				url: "https://via.placeholder.com/50",
// 			},
// 			expected: expected{
// 				imageData: &model.ImageData{
// 					Data:        []byte{137, 80, 78, 71},
// 					ContentType: "image/png",
// 					URL:         "https://via.placeholder.com/50",
// 					Width:       50,
// 					Height:      50,
// 				},
// 				err: nil,
// 			},
// 		},
// 		{
// 			name: "invalid image URL",
// 			input: input{
// 				url: "https://popa.ru/popa.jpeg",
// 			},
// 			expected: expected{
// 				imageData: nil,
// 				err:       errors.New("error sending request: status code: 404, body: Not Found"),
// 			},
// 		},
// 	}

// 	timeout := 10 * time.Second
// 	resp := &http.Response{
// 		StatusCode: http.StatusOK,
// 		Body:       io.NopCloser(bytes.NewReader([]byte{137, 80, 78, 71})),
// 		Header:     make(http.Header),
// 	}
// 	resp.Header.Set("Content-Type", "image/png")
// 	ctx := context.Background()

// 	for _, tc := range testCases {
// 		t.Run(tc.name, func(t *testing.T) {
// 			clientHttp := mocks.NewClientHttp(t)

// 			clientHttp.EXPECT().Do(mock.Anything).Return(resp, tc.expected.err)

// 			fetchImage := NewImageFetcher(timeout)

// 			result, err := fetchImage.FetchImage(ctx, tc.input.url)
// 			require.Equal(t, tc.expected, expected{imageData: result, err: err})

// 		})
// 	}

func TestFetchImage(t *testing.T) {
	type input struct {
		url  string
		resp *http.Response
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
				url: "https://via.placeholder.com/50",
				resp: &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewReader([]byte{137, 80, 78, 71})),
					Header:     http.Header{"Content-Type": []string{"image/png"}},
				},
			},
			expected: expected{
				imageData: &model.ImageData{
					Data:        []byte{137, 80, 78, 71, 13, 10, 26, 10, 0, 0, 0, 13, 73, 72, 68, 82, 0, 0, 0, 50, 0, 0, 0, 50, 8, 2, 0, 0, 0, 145, 93, 31, 230, 0, 0, 1, 167, 73, 68, 65, 84, 120, 156, 237, 152, 205, 170, 234, 48, 20, 70, 211, 198, 212, 106, 39, 17, 20, 133, 130, 153, 40, 20, 156, 248, 254, 207, 224, 84, 17, 4, 43, 197, 128, 80, 173, 208, 162, 84, 201, 207, 25, 244, 194, 189, 232, 109, 206, 217, 163, 227, 96, 175, 81, 75, 246, 247, 177, 146, 64, 7, 245, 86, 171, 21, 249, 60, 252, 223, 22, 248, 63, 168, 5, 1, 181, 32, 160, 22, 4, 212, 130, 128, 90, 16, 80, 11, 2, 106, 65, 64, 45, 8, 168, 5, 1, 181, 32, 124, 168, 86, 199, 189, 188, 221, 110, 235, 186, 38, 132, 132, 97, 152, 36, 137, 148, 242, 114, 185, 116, 187, 93, 33, 68, 24, 134, 223, 182, 191, 196, 9, 33, 63, 108, 112, 105, 89, 107, 61, 207, 91, 46, 151, 205, 107, 85, 85, 247, 251, 125, 177, 88, 84, 85, 117, 60, 30, 103, 179, 153, 219, 233, 37, 14, 106, 112, 93, 226, 227, 241, 48, 198, 172, 215, 235, 205, 102, 83, 150, 229, 237, 118, 27, 12, 6, 148, 82, 206, 249, 243, 249, 180, 214, 54, 99, 231, 243, 57, 77, 83, 66, 200, 225, 112, 200, 243, 188, 45, 78, 8, 105, 107, 128, 105, 213, 117, 205, 24, 75, 146, 68, 8, 145, 101, 153, 214, 186, 211, 249, 115, 186, 190, 239, 27, 99, 154, 231, 225, 112, 104, 173, 221, 239, 247, 198, 152, 209, 104, 212, 22, 183, 214, 182, 53, 188, 227, 186, 68, 206, 57, 231, 156, 16, 18, 69, 17, 99, 140, 82, 170, 181, 110, 150, 180, 214, 190, 255, 119, 75, 227, 241, 120, 183, 219, 205, 231, 115, 71, 92, 41, 229, 104, 120, 193, 117, 90, 121, 158, 167, 105, 170, 148, 170, 170, 202, 24, 19, 69, 81, 81, 20, 74, 169, 235, 245, 26, 4, 129, 231, 121, 205, 152, 181, 86, 74, 25, 199, 177, 148, 242, 223, 123, 121, 137, 51, 198, 218, 26, 222, 241, 28, 191, 70, 140, 49, 89, 150, 149, 101, 25, 4, 193, 116, 58, 237, 247, 251, 82, 202, 162, 40, 24, 99, 66, 136, 94, 175, 215, 140, 73, 41, 41, 165, 147, 201, 228, 116, 58, 105, 173, 227, 56, 110, 139, 55, 195, 239, 13, 48, 173, 95, 228, 67, 63, 167, 168, 5, 1, 181, 32, 160, 22, 4, 212, 130, 128, 90, 16, 80, 11, 2, 106, 65, 64, 45, 8, 168, 5, 225, 11, 45, 83, 16, 141, 96, 87, 246, 5, 0, 0, 0, 0, 73, 69, 78, 68, 174, 66, 96, 130},
					ContentType: "image/png",
					URL:         "https://via.placeholder.com/50",
					Width:       50,
					Height:      50,
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

	ctx := context.Background()

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			clientHttp := mocks.NewClientHttp(t)

			tc.input.resp.Header.Set("Content-Type", "image/png")
			clientHttp.EXPECT().Do(mock.Anything).Return(tc.input.resp, tc.expected.err)

			fetchImage := NewImageFetcher(timeout)

			result, err := fetchImage.FetchImage(ctx, tc.input.url)

			if diff := deep.Equal(result, tc.expected.imageData); diff != nil {
				t.Error(diff)
			}

			if !errors.Is(err, tc.expected.err) {
				t.Errorf("expected error %v, got %v", tc.expected.err, err)
			}
		})
	}
}
