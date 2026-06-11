from optimum.onnxruntime import ORTModelForFeatureExtraction
from transformers import AutoTokenizer

model_id = "intfloat/multilingual-e5-small"
save_dir = "modelo_e5_onnx/"

print("Downloading and converting the model to ONNX natively...")
model = ORTModelForFeatureExtraction.from_pretrained(model_id, export=True)
tokenizer = AutoTokenizer.from_pretrained(model_id)

model.save_pretrained(save_dir)
tokenizer.save_pretrained(save_dir)
print(f"All save in {save_dir}")
