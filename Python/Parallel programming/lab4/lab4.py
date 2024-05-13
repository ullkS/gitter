import threading
import time
import queue
import random

# Массив кол-й
message_queue = queue.Queue()

# Параметры симуляции
NUM_WRITERS = 3
NUM_READERS = 4
NUM_MESSAGES = 7
WRITER_PRIORITY = 0.2  # Имитация приоритета писателя\читателя 
READER_PRIORITY = 0.1  # - таймер (меньше значение - выше приоритет)

# Статистика работы
stats = {
    'writer': {i: {'messages_written': 0, 'time_spent': 0} for i in range(NUM_WRITERS)},
    'reader': {i: {'messages_read': 0, 'time_spent': 0} for i in range(NUM_READERS)}
}

stats_lock = threading.Lock()

def writer(writer_id):
    for _ in range(NUM_MESSAGES):
        start_time = time.time()
        message = f'Сообщение от писателя {writer_id}'
        message_queue.put(message)
        time.sleep(WRITER_PRIORITY)
        end_time = time.time()
        with stats_lock:
            stats['writer'][writer_id]['messages_written'] += 1
            stats['writer'][writer_id]['time_spent'] += end_time - start_time

def reader(reader_id):
    while True:
        try:
            start_time = time.time()
            message = message_queue.get(timeout=0.1) # Контроль веса
            time.sleep(READER_PRIORITY)
            end_time = time.time()
            with stats_lock:
                stats['reader'][reader_id]['messages_read'] += 1
                stats['reader'][reader_id]['time_spent'] += end_time - start_time
            #print(f'Читатель {reader_id} прочитал: {message}')  -- Вывод сообщений
        except queue.Empty:
            break

# Создание и запуск потоков писателей\читателей
writer_threads = []
for i in range(NUM_WRITERS):
    thread = threading.Thread(target=writer, args=(i,))
    thread.start()
    writer_threads.append(thread)

reader_threads = []
for i in range(NUM_READERS):
    thread = threading.Thread(target=reader, args=(i,))
    thread.start()
    reader_threads.append(thread)

# Ожидание завершения работы писателей\писателей
for thread in writer_threads:
    thread.join()

for thread in reader_threads:
    thread.join()

print('Статистика работы писателей и читателей:')
for writer_id, data in stats['writer'].items():
    print(f'Писатель {writer_id}: сообщений написано {data["messages_written"]}, активное время работы {data["time_spent"]:.2f} сек')
for reader_id, data in stats['reader'].items():
    print(f'Читатель {reader_id}: сообщений прочитано {data["messages_read"]}, активное время работы {data["time_spent"]:.2f} сек')