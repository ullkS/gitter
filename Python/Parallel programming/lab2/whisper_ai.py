import speech_recognition as sr
import os
from jiwer import wer
from nltk.translate.bleu_score import sentence_bleu
import nltk
import tkinter as tk
from tkinter import filedialog, messagebox
import threading
import whisper
import numpy as np
from pydub import AudioSegment
import io
import re
from pydub import AudioSegment
from rouge import Rouge
import time
# Загрузка необходимых ресурсов NLTK
nltk.download('punkt', quiet=True)

# Загрузка модели Whisper
model = whisper.load_model("base")

def convert_mp3_to_wav(mp3_path):
    audio = AudioSegment.from_mp3(mp3_path)
    wav_data = io.BytesIO()
    audio.export(wav_data, format="wav")
    wav_data.seek(0)
    return wav_data

def transcribe_audio(audio_file):
    try:
        # Загрузка аудио
        audio = AudioSegment.from_file(audio_file)
        
        # Длительность аудио в миллисекундах
        audio_duration = len(audio)
        
        # Длительность чанка в миллисекундах (25 секунд)
        chunk_duration = 25 * 1000
        
        # Инициализируем переменные
        transcribed_text = ""
        detected_language = None

        # Обрабатываем аудио по частям
        for start_time in range(0, audio_duration, chunk_duration):
            end_time = min(start_time + chunk_duration, audio_duration)
            
            # Выделяем часть аудио
            audio_chunk = audio[start_time:end_time]
            
            # Сохраняем чанк в отдельный файл
            chunk_filename = f"chunk_{start_time // 1000}_{end_time // 1000}.wav"
            audio_chunk.export(chunk_filename, format="wav")
            
            # Загружаем чанк с помощью whisper
            audio_data = whisper.load_audio(chunk_filename)
            audio_data = whisper.pad_or_trim(audio_data)
            
            # Получаем мел-спектрограмму
            mel = whisper.log_mel_spectrogram(audio_data).to(model.device)

            # Определение языка (только для первого чанка)
            if detected_language is None:
                _, probs = model.detect_language(mel)
                detected_language = max(probs, key=probs.get)
                print(f"Detected language: {detected_language}")

            # Транскрибация
            options = whisper.DecodingOptions(language=detected_language, without_timestamps=True)
            result = model.decode(mel, options)

            transcribed_text += result.text + " "
            
            # Удаляем временный файл
            os.remove(chunk_filename)

        # Удаляем лишние пробелы
        transcribed_text = transcribed_text.strip()
        
        # Дополнительная обработка текста
        transcribed_text = re.sub(r'\s+', ' ', transcribed_text)  # Удаляем лишние пробелы
        transcribed_text = transcribed_text.capitalize()  # Делаем первую букву заглавной

        return transcribed_text
    except Exception as e:
        print(f"Error in transcribe_audio: {e}")
        return ""

def simple_tokenize(text):
    return re.findall(r'\w+|[^\w\s]', text.lower())

def evaluate_transcription(reference, hypothesis):
    wer_score = wer(reference, hypothesis)
    reference_tokens = simple_tokenize(reference)
    hypothesis_tokens = simple_tokenize(hypothesis)
    bleu_score = sentence_bleu([reference_tokens], hypothesis_tokens)

    # Calculate ROUGE-L
    rouge = Rouge()
    rouge_scores = rouge.get_scores(hypothesis, reference, avg=True)
    rouge_l_score = rouge_scores['rouge-l']['f'] * 100

    wer_score *= 100
    bleu_score *= 100

    return wer_score, bleu_score, rouge_l_score

def process_audio_file(audio_file, expected_text):
    start_time = time.time()
    # Распознаем текст из аудио
    transcribed_text = transcribe_audio(audio_file)
    
    # Определяем путь для сохранения файла с распознанным текстом
    output_file = os.path.join(os.path.dirname(audio_file), "transcribed_output.txt")
    
    # Сохраняем распознанный текст в файл
    with open(output_file, 'w', encoding='utf-8') as f:
        f.write(transcribed_text)
    
    processing_time = time.time() - start_time

    # Оцениваем качество распознавания
    wer_score, bleu_score, rouge_l_score = evaluate_transcription(expected_text, transcribed_text)
    return {
        'transcribed_text': transcribed_text,
        'wer_score': wer_score,
        'bleu_score': bleu_score,
        'rouge_l_score': rouge_l_score,
        'processing_time': processing_time
    }

class SpeechRecognitionApp:
    def __init__(self, master):
        self.master = master
        master.title("Speech Recognition App")

        self.audio_path = tk.StringVar()
        self.text_path = tk.StringVar()
        self.recording = False  # Флаг для управления записью

        tk.Label(master, text="Путь к аудио файлу:").grid(row=0, column=0, sticky="e")
        tk.Entry(master, textvariable=self.audio_path, width=50).grid(row=0, column=1)
        tk.Button(master, text="Выбрать", command=self.select_audio).grid(row=0, column=2)

        tk.Label(master, text="Путь к текстовому файлу:").grid(row=1, column=0, sticky="e")
        tk.Entry(master, textvariable=self.text_path, width=50).grid(row=1, column=1)
        tk.Button(master, text="Выбрать", command=self.select_text).grid(row=1, column=2)

        tk.Button(master, text="Обработать", command=self.process).grid(row=2, column=1)

        self.result_text = tk.Text(master, height=20, width=70)
        self.result_text.grid(row=3, column=0, columnspan=3, padx=10, pady=10)

        tk.Button(master, text="Запись речи", command=self.start_recording).grid(row=4, column=0)
        tk.Button(master, text="Остановить запись", command=self.stop_recording).grid(row=4, column=2)

    def select_audio(self):
        file_path = filedialog.askopenfilename(filetypes=[("Audio Files", "*.wav *.mp3")])
        if file_path:
            self.audio_path.set(file_path)

    def select_text(self):
        file_path = filedialog.askopenfilename(filetypes=[("Text Files", "*.txt")])
        if file_path:
            self.text_path.set(file_path)

    def process(self):
        audio_file = self.audio_path.get()
        text_file = self.text_path.get()

        if not audio_file or not text_file:
            messagebox.showerror("Error", "Please select an audio and text file.")
            return

        with open(text_file, 'r', encoding='utf-8') as f:
            expected_text = f.read()

        result = process_audio_file(audio_file, expected_text)

        self.result_text.delete(1.0, tk.END)
        self.result_text.insert(tk.END, f"Transcribed Text: {result['transcribed_text']}\n")
        self.result_text.insert(tk.END, f"WER: {result['wer_score']:.2f}%\n")
        self.result_text.insert(tk.END, f"BLEU: {result['bleu_score']:.2f}%\n")
        self.result_text.insert(tk.END, f"ROUGE-L: {result['rouge_l_score']:.2f}%\n")
        self.result_text.insert(tk.END, f"Processing Time: {result['processing_time']:.2f} seconds\n")

    def start_recording(self):
        self.recording = True
        threading.Thread(target=self.record_speech).start()

    def stop_recording(self):
        self.recording = False

    def record_speech(self):
        r = sr.Recognizer()
        with sr.Microphone() as source:
            self.result_text.delete(1.0, tk.END)
            self.result_text.insert(tk.END, "Говорите...")

            while self.recording:
                audio = r.listen(source, timeout=1, phrase_time_limit=29)
                try:
                    # Сохраняем аудио во временный файл
                    with open("temp_audio.wav", "wb") as f:
                        f.write(audio.get_wav_data())
                    
                    # Используем Whisper для распознавания
                    text = transcribe_audio("temp_audio.wav")
                    
                    self.result_text.delete(1.0, tk.END)
                    self.result_text.insert(tk.END, f"Распознанный текст: {text}")
                    
                    # Удаляем временный файл
                    os.remove("temp_audio.wav")
                except Exception as e:
                    self.result_text.delete(1.0, tk.END)
                    self.result_text.insert(tk.END, f"Ошибка при распознавании речи: {e}")

# Запуск приложения
root = tk.Tk()
app = SpeechRecognitionApp(root)
root.mainloop()