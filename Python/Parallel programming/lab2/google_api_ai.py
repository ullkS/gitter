import os, nltk, threading, speech_recognition as sr, tkinter as tk
from jiwer import wer
from nltk.translate.bleu_score import sentence_bleu
from tkinter import filedialog, messagebox

# Загрузка необходимых ресурсов NLTK
nltk.download('punkt', quiet=True)

def transcribe_audio_chunks(audio_file, chunk_length=60, max_attempts=3):
    recognizer = sr.Recognizer()
    transcribed_text = []

    with sr.AudioFile(audio_file) as source:
        audio_length = int(source.DURATION)
        for start in range(0, audio_length, chunk_length):
            source_offset = start
            audio_chunk = recognizer.record(source, offset=source_offset, duration=chunk_length)
            attempt = 0
            success = False
            while attempt < max_attempts and not success:
                try:
                    text = recognizer.recognize_google(audio_chunk, language="ru-RU")
                    transcribed_text.append(text)
                    success = True
                except sr.UnknownValueError:
                    attempt += 1
                    print(f"Не удалось распознать фрагмент с {start} до {start + chunk_length} секунд. Попытка {attempt}/{max_attempts}.")
                except sr.RequestError as e:
                    transcribed_text.append(f"Error: {e}")
                    print(f"Ошибка запроса: {e}")
                    break  # Прерываем попытки при ошибке запроса

            if not success:
                transcribed_text.append("...")
                print(f"Фрагмент с {start} до {start + chunk_length} секунд не распознан после {max_attempts} попыток.")               

    return " ".join(transcribed_text)

def evaluate_transcription(reference, hypothesis):
    wer_score = wer(reference, hypothesis)
    reference_tokens = nltk.word_tokenize(reference.lower())
    hypothesis_tokens = nltk.word_tokenize(hypothesis.lower())
    bleu_score = sentence_bleu([reference_tokens], hypothesis_tokens)

    wer_score *= 100
    bleu_score *= 100

    return wer_score, bleu_score

def process_audio_file(audio_file, expected_text):
    # Распознаем текст из аудио
    transcribed_text = transcribe_audio_chunks(audio_file)
    
    # Определяем путь для сохранения файла с распознанным текстом
    output_file = os.path.join(os.path.dirname(audio_file), "transcribed_output.txt")
    
    # Сохраняем распознанный текст в файл
    with open(output_file, 'w', encoding='utf-8') as f:
        f.write(transcribed_text)

    # Оцениваем качество распознавания
    wer_score, bleu_score = evaluate_transcription(expected_text, transcribed_text)
    return {
        'transcribed_text': transcribed_text,
        'wer_score': wer_score,
        'bleu_score': bleu_score
    }

class SpeechRecognitionApp:
    def __init__(self, master):
        self.master = master
        master.title("Speech Recognition App")

        self.audio_path = tk.StringVar()
        self.text_path = tk.StringVar()
        self.is_recording = False
        self.is_paused = False

        tk.Label(master, text="Путь к аудио файлу:").grid(row=0, column=0, sticky="e")
        tk.Entry(master, textvariable=self.audio_path, width=50).grid(row=0, column=1)
        tk.Button(master, text="Выбрать", command=self.select_audio).grid(row=0, column=2)

        tk.Label(master, text="Путь к текстовому файлу:").grid(row=1, column=0, sticky="e")
        tk.Entry(master, textvariable=self.text_path, width=50).grid(row=1, column=1)
        tk.Button(master, text="Выбрать", command=self.select_text).grid(row=1, column=2)

        tk.Button(master, text="Обработать", command=self.process).grid(row=2, column=1)

        self.result_text = tk.Text(master, height=20, width=70)
        self.result_text.grid(row=3, column=0, columnspan=3, padx=10, pady=10)

        tk.Button(master, text="Запись речи", command=self.record_speech).grid(row=4, column=0)
        self.pause_button = tk.Button(master, text="Пауза", command=self.toggle_pause, state=tk.DISABLED)
        self.pause_button.grid(row=4, column=2)

    def select_audio(self):
        filename = filedialog.askopenfilename(filetypes=[("Audio Files", "*.wav *.mp3")])
        self.audio_path.set(filename)

    def select_text(self):
        filename = filedialog.askopenfilename(filetypes=[("Text Files", "*.txt")])
        self.text_path.set(filename)

    def process(self):
        audio_file = self.audio_path.get()
        text_file = self.text_path.get()

        if not audio_file or not text_file:
            messagebox.showerror("Ошибка", "Пожалуйста, выберите аудио и текстовый файл.")
            return

        with open(text_file, 'r', encoding='utf-8') as f:
            expected_text = f.read()

        result = process_audio_file(audio_file, expected_text)

        self.result_text.delete(1.0, tk.END)
        self.result_text.insert(tk.END, f"Распознанный текст:\n{result['transcribed_text']}\n")
        self.result_text.insert(tk.END, f"\nWER: {result['wer_score']:.2f}%\n")
        self.result_text.insert(tk.END, f"BLEU: {result['bleu_score']:.2f}%\n")

    def toggle_pause(self):
        if self.is_paused:
            self.is_paused = False
            self.pause_button.config(text="Пауза")
        else:
            self.is_paused = True
            self.pause_button.config(text="Продолжить")

    def record_speech(self):
        def record():
            self.is_recording = True
            self.pause_button.config(state=tk.NORMAL)
            r = sr.Recognizer()
            with sr.Microphone() as source:
                self.result_text.delete(1.0, tk.END)
                self.result_text.insert(tk.END, "Говорите...")
    
                while self.is_recording:
                    if not self.is_paused:
                        audio = r.listen(source, timeout=1, phrase_time_limit=29)  
                        try:
                            text = r.recognize_google(audio, language="ru-RU")
                            self.result_text.insert(tk.END, f"\n{text}")
                        except sr.UnknownValueError:
                            self.result_text.insert(tk.END, "\nНе удалось распознать аудио")
                        except sr.RequestError as e:
                            self.result_text.insert(tk.END, f"\nОшибка запроса: {e}")
    
        threading.Thread(target=record).start()

root = tk.Tk()
app = SpeechRecognitionApp(root)
root.mainloop()