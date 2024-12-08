import pandas as pd
import torch
from transformers import AutoTokenizer, AutoModelForSequenceClassification, pipeline

# Model ve tokenizer yolu
model1_path = r"/scraper/Models/final_model_pytorch"

# Tokenizer ve modelleri yükle
tokenizer = AutoTokenizer.from_pretrained(model1_path)  

# Ensemble model sınıfı
class EnsembleModel:
    def __init__(self, model1, model2, pipe):
        self.model1 = model1
        self.model2 = model2
        self.pipe = pipe
        self.labels = {0: "Negatif", 1: "Nötr", 2: "Pozitif"}  # 3 sınıf etiketi
        self.pipe_labels = {"negative": "Negatif", "positive": "Pozitif"}  # Pipe modelinin 2 sınıf etiketi

    def predict(self, input_text):
        inputs = tokenizer(input_text, return_tensors="pt", padding=True, truncation=True)
        output1 = self.model1(**inputs).logits
        probs1 = torch.softmax(output1, dim=-1)

        output2 = self.model2(**inputs).logits
        probs2 = torch.softmax(output2, dim=-1)

        pipe_preds = self.pipe(input_text)
        label = pipe_preds[0]['label']

        if label == 'negative':
            probs3 = torch.tensor([[1.0, 0.0, 0.0]])
        elif label == 'positive':
            probs3 = torch.tensor([[0.0, 0.0, 1.0]])
        else:
            probs3 = torch.tensor([[0.0, 1.0, 0.0]])

        max_categories = 3
        if probs1.size(1) != max_categories:
            probs1 = torch.cat([probs1, torch.zeros(probs1.size(0), max_categories - probs1.size(1))], dim=1)

        if probs2.size(1) != max_categories:
            probs2 = torch.cat([probs2, torch.zeros(probs2.size(0), max_categories - probs2.size(1))], dim=1)

        if probs3.size(1) != max_categories:
            probs3 = torch.cat([probs3, torch.zeros(probs3.size(0), max_categories - probs3.size(1))], dim=1)

        final_probs = (probs1 + probs2 + probs3) / 3  
        final_category = torch.argmax(final_probs, dim=-1)

        return final_category.item(), final_probs


