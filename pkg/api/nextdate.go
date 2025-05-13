package api

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const (
	DateFormat = "20060102"
)

// afterNow обнуляет время для сравнения дат
// date - параметр сравниваемый с точкой отсчета (now)
// now - точка отсчета
func afterNow(date, now time.Time) bool {
	dateDay := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())
	nowDay := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	return dateDay.After(nowDay)
}

// NextDate вычисляет следующую дату.
// now - время, от которого ищется ближайшая дата.
// dstart - исходное время в формате 20060102, от которого начинается отсчёт повторений.
// repeat - правило повторения в описанном выше формате.
// Функция возвращает следующую дату в формате 20060102 и ошибку.
func NextDate(now time.Time, dstart string, repeat string) (string, error) {

	if repeat == "" {
		return "", nil
	}

	startDate, err := time.Parse(DateFormat, dstart)
	if err != nil {
		return "ошибка выполнения time.Parse", err
	}

	currentDate := startDate

	switch {
	case repeat == "y":

		for {
			currentDate = currentDate.AddDate(1, 0, 0)
			if afterNow(currentDate, now) {
				break
			}
		}
		return currentDate.Format(DateFormat), nil

	case strings.HasPrefix(repeat, "d"):
		parts := strings.Split(repeat, " ")
		if len(parts) != 2 {
			return "указан не верный формат", err
		}
		daysStr := parts[1]
		days, err := strconv.Atoi(daysStr)
		if err != nil {
			return "ошибка преобразования", err
		}
		if days <= 0 || days > 400 {
			return "максимально допустимое число равно 400", err
		}

		for {
			currentDate = currentDate.AddDate(0, 0, days)
			if afterNow(currentDate, now) {
				break
			}
		}
		return currentDate.Format(DateFormat), err

	default:
		return "", fmt.Errorf("неподдерживаемый формат %s", repeat)
	}
}

// обработчик следующей даты
func nextDateHandler(w http.ResponseWriter, r *http.Request) {
	// Проверяем запрос на GET
	if r.Method != http.MethodGet {
		http.Error(w, "ошибка запроса", http.StatusMethodNotAllowed)
		return
	}

	// параметры из URL
	nowParam := r.FormValue("now")
	dateParam := r.FormValue("date")
	repeatParam := r.FormValue("repeat")

	//определяем формат now если параметр Now не определен тогда берем текущую дату
	var now time.Time
	if nowParam == "" {
		now = time.Now()
	} else {

		parsedNow, err := time.Parse(DateFormat, nowParam)
		if err != nil {
			http.Error(w, "Неверный формат now", http.StatusBadRequest)
			return
		}
		now = parsedNow
	}

	//вызываем функци nextDate
	nextDate, err := NextDate(now, dateParam, repeatParam)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// вернуть дату следующего выполнения
	w.Header().Set("Content-Type", "text/plain")
	fmt.Fprint(w, nextDate)
}
