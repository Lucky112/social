
# Загрузите данные из CSV файла
file_path = '/home/ribragimov/work/otus/social/docs/resp_time_index'  # Путь к вашему CSV файлу

import pandas as pd
import matplotlib.pyplot as plt

# Загрузите данные из CSV файла
data = pd.read_csv(file_path)

# Преобразуйте timeStamp в datetime формат и в секунды от начала
data['timeStamp'] = pd.to_datetime(data['timeStamp'], unit='ms')
data['timeInSeconds'] = (data['timeStamp'] - data['timeStamp'].min()).dt.total_seconds()

# Для уменьшения количества точек на графике, выберем каждую n-ую точку
n = 100  # Показывать каждую 5-ю точку (измените это значение для контроля количества точек)
reduced_data = data.iloc[::n]

# Постройте график Latency over Time
plt.figure(figsize=(12, 6))
plt.plot(reduced_data['timeInSeconds'], reduced_data['Latency'], marker='o', linestyle='-', color='b', label='Latency', markersize=5)

# Настройки графика
plt.title('Latency Over Time')
plt.xlabel('Time (seconds since start)')
plt.ylabel('Latency (ms)')
plt.xticks(rotation=45)  # Поверните метки по оси X для удобочитаемости

# Настройка меток по оси X
# plt.xticks(range(int(data['timeInSeconds'].max()) + 1))  # Установка меток по каждой секунде

plt.grid()

# Сохраните график в файл
plt.tight_layout()
plt.savefig('/mnt/c/Users/mhr11/Downloads/latency_over_time_index.png')  # Сохранение в PNG
plt.close()  # Закрытие графика

