package embedder

import (
	"math"
	"testing"
)

func TestMeanPoolReturns384Dimensions(t *testing.T) {

	output := make([]float32, 3*EmbeddingSize)

	mask := []int64{1, 1, 1}

	vec := meanPoolAndNormalize(output, mask, int64(len(mask)))

	if len(vec) != EmbeddingSize {
		t.Fatalf(
			"expected %d dims got %d",
			EmbeddingSize,
			len(vec),
		)
	}
}

func TestEmbeddingNormalized(t *testing.T) {

	output := make([]float32, 2*EmbeddingSize)

	for i := range output {
		output[i] = 1
	}

	mask := []int64{1, 1}

	vec := meanPoolAndNormalize(output, mask, int64(len(mask)))

	var norm float64

	for _, v := range vec {
		norm += float64(v * v)
	}

	norm = math.Sqrt(norm)

	if math.Abs(norm-1.0) > 0.0001 {
		t.Fatalf(
			"expected norm=1 got %f",
			norm,
		)
	}
}

func TestMeanPoolingDeterministic(t *testing.T) {

	output := make([]float32, 2*EmbeddingSize)

	for i := range output {
		output[i] = float32(i)
	}

	mask := []int64{1, 1}

	v1 := meanPoolAndNormalize(output, mask, int64(len(mask)))
	v2 := meanPoolAndNormalize(output, mask, int64(len(mask)))

	for i := range v1 {

		if v1[i] != v2[i] {
			t.Fatalf(
				"vectors differ at %d",
				i,
			)
		}
	}
}
