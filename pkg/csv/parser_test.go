package csv

import (
	"bytes"
	"mime/multipart"
	"testing"

	"github.com/stretchr/testify/assert"
)

// helper to build a *multipart.FileHeader from raw CSV content
func CreateMultipartFileHeader(filename, content string) *multipart.FileHeader {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, _ := writer.CreateFormFile("file", filename)
	part.Write([]byte(content))
	writer.Close()

	reader := multipart.NewReader(body, writer.Boundary())
	form, err := reader.ReadForm(1 << 20)
	if err != nil {
		return nil
	}
	return form.File["file"][0]
}

func TestParseFromMultipartFile(t *testing.T) {
	t.Parallel()

	type row struct {
		Description string `csv:"description"`
		URL         string `csv:"url"`
	}

	testCases := []struct {
		name string

		csvContent string

		expectedData  []row
		expectedError bool
	}{
		{
			name:       "valid csv - single row",
			csvContent: "description,url\nmy bookmark,https://example.com",
			expectedData: []row{
				{Description: "my bookmark", URL: "https://example.com"},
			},
		},
		{
			name:       "valid csv - multiple rows",
			csvContent: "description,url\nbookmark 1,https://example.com\nbookmark 2,https://google.com",
			expectedData: []row{
				{Description: "bookmark 1", URL: "https://example.com"},
				{Description: "bookmark 2", URL: "https://google.com"},
			},
		},
		{
			name:         "headers only - no data rows",
			csvContent:   "description,url\n",
			expectedData: nil,
		},
		{
			name:          "empty file",
			csvContent:    "",
			expectedError: true, // gocsv errors on empty content
		},
		{
			name:          "malformed csv - inconsistent columns",
			csvContent:    "description,url\nonly-one-column",
			expectedError: true,
		},
		{
			name:       "extra unknown columns - should be ignored",
			csvContent: "description,url,unknown_col\nbookmark,https://example.com,ignored",
			expectedData: []row{
				{Description: "bookmark", URL: "https://example.com"},
			},
		},
		{
			name:       "missing description column - zero value",
			csvContent: "url\nhttps://example.com",
			expectedData: []row{
				{Description: "", URL: "https://example.com"},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			fileHeader := CreateMultipartFileHeader("test.csv", tc.csvContent)

			var result []row
			err := ParseCSVFile(fileHeader, &result)

			if tc.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedData, result)
			}
		})
	}
}
