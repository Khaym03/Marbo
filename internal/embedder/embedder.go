// Package embedder provides functionality for generating text embeddings using ONNX models.
package embedder

import (
	"fmt"
	"math"

	"github.com/sugarme/tokenizer"
	"github.com/sugarme/tokenizer/pretrained"
	ort "github.com/yalue/onnxruntime_go"
)

const (
	EmbeddingSize = 384
	MaxSeqLength  = 512
)

// E5Embedder implements the Embedder interface using an E5 model.
type E5Embedder struct {
	tokenizer *tokenizer.Tokenizer
	options   *ort.SessionOptions
	session   *ort.AdvancedSession

	// We keep tensors to feed the session
	inputIDsTensor      *ort.Tensor[int64]
	attentionMaskTensor *ort.Tensor[int64]
	tokenTypeTensor     *ort.Tensor[int64]
	outputTensor        *ort.Tensor[float32]

	inputIDsData      []int64
	attentionMaskData []int64
	tokenTypeData     []int64
	outputData        []float32
}

func New(
	modelPath string,
	tokenizerPath string,
) (*E5Embedder, error) {

	tk, err := pretrained.FromFile(tokenizerPath)
	if err != nil {
		return nil, fmt.Errorf("load tokenizer: %w", err)
	}

	opts, err := ort.NewSessionOptions()
	if err != nil {
		return nil, fmt.Errorf("create session options: %w", err)
	}

	// Pre-allocate buffers
	inputIDsData := make([]int64, MaxSeqLength)
	attentionMaskData := make([]int64, MaxSeqLength)
	tokenTypeData := make([]int64, MaxSeqLength)
	outputData := make([]float32, MaxSeqLength*EmbeddingSize)

	// Pre-create tensors
	shape := ort.NewShape(1, MaxSeqLength)
	outShape := ort.NewShape(1, MaxSeqLength, EmbeddingSize)

	inputIDsTensor, err := ort.NewTensor(shape, inputIDsData)
	if err != nil {
		return nil, err
	}
	attentionMaskTensor, err := ort.NewTensor(shape, attentionMaskData)
	if err != nil {
		return nil, err
	}
	tokenTypeTensor, err := ort.NewTensor(shape, tokenTypeData)
	if err != nil {
		return nil, err
	}
	outputTensor, err := ort.NewTensor(outShape, outputData)
	if err != nil {
		return nil, err
	}

	session, err := ort.NewAdvancedSession(
		modelPath,
		[]string{"input_ids", "attention_mask", "token_type_ids"},
		[]string{"last_hidden_state"},
		[]ort.Value{inputIDsTensor, attentionMaskTensor, tokenTypeTensor},
		[]ort.Value{outputTensor},
		opts,
	)
	if err != nil {
		return nil, err
	}

	return &E5Embedder{
		tokenizer:           tk,
		options:             opts,
		session:             session,
		inputIDsTensor:      inputIDsTensor,
		attentionMaskTensor: attentionMaskTensor,
		tokenTypeTensor:     tokenTypeTensor,
		outputTensor:        outputTensor,
		inputIDsData:        inputIDsData,
		attentionMaskData:   attentionMaskData,
		tokenTypeData:       tokenTypeData,
		outputData:          outputData,
	}, nil
}

func (e *E5Embedder) Close() {
	if e.session != nil {
		e.session.Destroy()
	}
	e.inputIDsTensor.Destroy()
	e.attentionMaskTensor.Destroy()
	e.tokenTypeTensor.Destroy()
	e.outputTensor.Destroy()
	if e.options != nil {
		e.options.Destroy()
	}
}

func (e *E5Embedder) Embed(text string) ([]float32, error) {

	en, err := e.tokenizer.EncodeSingle(text, true)
	if err != nil {
		return nil, err
	}

	seqLen := len(en.Ids)
	if seqLen > MaxSeqLength {
		seqLen = MaxSeqLength
	}

	// Reset buffers
	for i := 0; i < MaxSeqLength; i++ {
		e.inputIDsData[i] = 0
		e.attentionMaskData[i] = 0
		e.tokenTypeData[i] = 0
	}

	for i := 0; i < seqLen; i++ {
		e.inputIDsData[i] = int64(en.Ids[i])
		e.attentionMaskData[i] = int64(en.AttentionMask[i])
	}

	// Run inference (using pre-allocated tensors with MaxSeqLength)
	// Note: We are not changing tensor shape here as it's tricky with this library.
	// We'll rely on attention masking to ignore padding tokens.
	if err := e.session.Run(); err != nil {
		return nil, err
	}

	// Pool data from persistent buffer
	return meanPoolAndNormalize(
		e.outputData,
		e.attentionMaskData[:seqLen],
		int64(seqLen),
	), nil
}

func meanPoolAndNormalize(
	output []float32,
	mask []int64,
	seqLen int64,
) []float32 {

	embedding := make([]float32, EmbeddingSize)
	var valid float32

	for t := int64(0); t < seqLen; t++ {
		if mask[t] == 0 {
			continue
		}
		valid++
		offset := t * EmbeddingSize
		for dim := 0; dim < EmbeddingSize; dim++ {
			embedding[dim] += output[int(offset)+dim]
		}
	}

	var norm float64
	for i := range embedding {
		embedding[i] /= valid
		norm += float64(embedding[i] * embedding[i])
	}

	l2 := float32(math.Sqrt(norm))
	for i := range embedding {
		embedding[i] /= l2
	}

	return embedding
}
