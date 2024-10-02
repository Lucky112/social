import os
import csv
import psycopg2
from datetime import datetime
from pathlib import Path

def connect_db():
    host=os.getenv('DB_HOST')
    dbname=os.getenv('DB_NAME')
    user=os.getenv('DB_USER')
    pwd_file=os.getenv('DB_PASSWORD_FILE')
    password= Path(pwd_file).read_text()

    print(f"{user}:{password}@{host}:{dbname}")

    return psycopg2.connect(
        host=host,
        dbname=dbname,
        user=user,
        password=password
    )


def insert_data(data, batch_size=1000):
    conn = connect_db()
    cursor = conn.cursor()

    insert_query = """
    INSERT INTO scl.profiles (user_id, name, surname, birthdate, sex, address, hobbies)
    VALUES %s
    """

    try:
        # Разбиваем данные на батчи
        for i in range(0, len(data), batch_size):
            batch = data[i:i + batch_size]

            # Формируем часть запроса с большим количеством значений
            values_str = ','.join(cursor.mogrify("(%s, %s, %s, %s, %s, %s, %s)", row).decode('utf-8') for row in batch)
            full_query = insert_query % values_str

            # Выполняем запрос на вставку батча
            cursor.execute(full_query)

        conn.commit()
        print(f"Данные успешно вставлены ({len(data)} строк)")
    except Exception as e:
        conn.rollback()
        print(f"Ошибка вставки данных: {e}")
    finally:
        cursor.close()
        conn.close()


def read_csv(file_path):
    data = []
    with open(file_path, mode='r', encoding='utf-8') as file:
        reader = csv.reader(file)
        for row in reader:
            # Разбираем строку формата "Абрамов,Тимофей, 1909-01-01, Лиски"
            surname, name, birthdate, address = row
            birthdate = datetime.strptime(birthdate, '%Y-%m-%d').date()
            sex = "female" if surname[-1] == "а" else "male"
            hobbies = "Чтение, сериалы"
            user_id = "1"

            data.append((user_id, name, surname, birthdate, sex, address, hobbies))
    return data


def main():
    csv_file = "data.csv"
    data = read_csv(csv_file)

    insert_data(data, batch_size=65000//7)

if __name__ == "__main__":
    main()