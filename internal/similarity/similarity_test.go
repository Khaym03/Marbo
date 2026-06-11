package similarity

import "testing"

func TestDotProductSameVector(t *testing.T) {

	v := []float32{
		1, 0, 0,
	}

	score, err := DotProduct(v, v)

	if err != nil {
		t.Fatal(err)
	}

	if score != 1 {
		t.Fatalf(
			"expected 1 got %f",
			score,
		)
	}
}

func TestDotProductOpposite(t *testing.T) {

	a := []float32{
		1, 0,
	}

	b := []float32{
		-1, 0,
	}

	score, err := DotProduct(a, b)

	if err != nil {
		t.Fatal(err)
	}

	if score != -1 {
		t.Fatalf(
			"expected -1 got %f",
			score,
		)
	}
}

func TestDotProductDimensionMismatch(
	t *testing.T,
) {

	a := []float32{1, 2}
	b := []float32{1}

	_, err := DotProduct(a, b)

	if err == nil {
		t.Fatal("expected error")
	}
}
