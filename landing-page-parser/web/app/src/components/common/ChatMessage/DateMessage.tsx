import React from "react";
import './DateMessage.css'

const DateMessage: React.FC<{ date: Date }> = (props) => {
    const date = new Date(props.date)

    return (
        <p
            className="message-date"
            key={`date-message-${props.date}`}
        >
            {formatDate(date)}
        </p>
    )
}

const weekDays = new Map<number, string>([
    [1, "Вс"],
    [2, "Пн"],
    [3, "Вт"],
    [4, "Ср"],
    [5, "Чт"],
    [6, "Пт"],
    [7, "Сб"],
])

const months = new Map<number, string>([
    [1, "Январь"],
    [2, "Февраль"],
    [3, "Март"],
    [4, "Апрель"],
    [5, "Май"],
    [6, "Июнь"],
    [7, "Июль"],
    [8, "Август"],
    [9, "Сентябрь"],
    [10, "Октябрь"],
    [11, "Ноябрь"],
    [12, "Декабрь"],
])

const formatDate = (date: Date): string => {
    let month = months.get(date.getMonth() + 1)
    if (!month) {
        month = "???"
    }

    let weekday = weekDays.get(date.getDay() + 1)
    if (!weekday) {
        weekday = "???"
    }

    const day = date.getDate() < 10 ? "0" + date.getDate() : date.getDate().toString()
    const year = date.getFullYear().toString()

    return `${month}, ${weekday} ${day} ${year}`
}

export default DateMessage