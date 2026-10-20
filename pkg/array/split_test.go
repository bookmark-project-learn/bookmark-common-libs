package array

import "testing"

func TestSplitIntoBatches(t *testing.T) {
	testCases := []struct {
		name      string
		input     []int
		batchSize int
		expected  [][]int
	}{
		{
			name:      "empty input",
			input:     []int{},
			batchSize: 3,
			expected:  [][]int{},
		},
		{
			name:      "batch size larger than input",
			input:     []int{1, 2},
			batchSize: 5,
			expected:  [][]int{{1, 2}},
		},
		{
			name:      "batch size equal to input length",
			input:     []int{1, 2, 3},
			batchSize: 3,
			expected:  [][]int{{1, 2, 3}},
		},
		{
			name:      "batch size smaller than input",
			input:     []int{1, 2, 3, 4, 5},
			batchSize: 2,
			expected:  [][]int{{1, 2}, {3, 4}, {5}},
		},
		{
			name:      "batch size of one",
			input:     []int{1, 2, 3},
			batchSize: 1,
			expected:  [][]int{{1}, {2}, {3}},
		},
		{
			name:      "batch size of zero",
			input:     []int{1, 2, 3},
			batchSize: 0,
			expected:  [][]int{},
		},
		{
			name:      "batch size of negative number",
			input:     []int{1, 2, 3},
			batchSize: -1,
			expected:  [][]int{},
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			result := SplitIntoBatch(tc.input, tc.batchSize)

			// Gọi hàm helper để kiểm tra tính chính xác của kết quả
			assertBatchesEqual(t, tc.expected, result)
		})
	}
}

// assertBatchesEqual là hàm bổ trợ giúp so sánh kết quả mong đợi và thực tế.
// Sử dụng t.Helper() để khi test fail sẽ trỏ chính xác về dòng gọi hàm trong TestSplitIntoBatches.
func assertBatchesEqual(t *testing.T, expected, result [][]int) {
	t.Helper()

	if len(result) != len(expected) {
		t.Fatalf("expected %d batches, got %d", len(expected), len(result))
	}

	for i := range result {
		if len(result[i]) != len(expected[i]) {
			t.Fatalf("expected batch %d to have %d elements, got %d", i, len(expected[i]), len(result[i]))
		}
		for j := range result[i] {
			if result[i][j] != expected[i][j] {
				t.Fatalf("expected batch %d element %d to be %d, got %d", i, j, expected[i][j], result[i][j])
			}
		}
	}
}
