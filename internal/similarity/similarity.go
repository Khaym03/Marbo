package similarity

import "fmt"

func DotProduct(
	a []float32,
	b []float32,
) (float32, error) {

	if len(a) != len(b) {
		return 0, fmt.Errorf(
			"vector size mismatch: %d != %d",
			len(a),
			len(b),
		)
	}

	var score float32

	for i := range a {
		score += a[i] * b[i]
	}

	return score, nil
}
