import "./style.css";

// Кнопка конвертации
const convertButton = document.getElementById("convert");
const resetButton = document.getElementById("reset");

const asideContent = document.querySelector(".aside-content");
const mainContent = document.querySelector(".main-content");

const radios = document.querySelectorAll('input[name="unit-type"]');

const fromDropdown = document.getElementById("from");
const toDropdown = document.getElementById("to");

const options = {
  length: [
    "millimeter",
    "centimeter",
    "meter",
    "kilometer",
    "inch",
    "foot",
    "yard",
    "mile",
  ],
  weight: ["milligram", "gram", "kilogram", "ounce", "pound"],
  temperature: ["Celsius", "Fahrenheit", "Kelvin"],
};
// Основные элементы

function updateDropdown(selectedType) {
  // Очищаем текущее содержимое dropdown
  fromDropdown.innerHTML = "";
  toDropdown.innerHTML = "";

  // Получаем опции для выбранного типа
  const selectedOptions = options[selectedType];

  // Заполняем dropdown
  selectedOptions.forEach((unit) => {
    const fromOption = document.createElement("option");
    fromOption.value = unit;
    fromOption.textContent = unit;

    const toOption = document.createElement("option");
    toOption.value = unit;
    toOption.textContent = unit;

    fromDropdown.appendChild(fromOption);
    toDropdown.appendChild(toOption);
  });
}

// Навешиваем обработчик на радиокнопки
radios.forEach((radio) => {
  radio.addEventListener("change", (event) => {
    updateDropdown(event.target.value); // Обновляем dropdown при изменении
  });
});

// Инициализация - загружаем варианты для "Length" по умолчанию
updateDropdown("length");

// Асинхронная функция для получения данных
async function getData() {
  // Получаем значения полей
  const length = document.getElementById("length").value;
  const resultElement = document.getElementById("result");
  const type = document.querySelector('input[name="unit-type"]:checked').value;

  const fromUnit = fromDropdown.value;
  const toUnit = toDropdown.value;

  // resultElement.textContent = "Resuilt:";
  // asideContent.classList.remove("close");

  // Проверка заполнения полей
  if (!length) {
    alert("Пожалуйста, заполните все поля!");
    return;
  }

  // URL вашего API
  const url = "/api/convert";

  try {
    // Запрос к серверу
    const response = await fetch(url, {
      method: "POST",
      headers: {
        "Content-Type": "application/json;charset=utf-8",
      },
      body: JSON.stringify({
        type: type,
        value: parseFloat(length),
        from: fromUnit,
        to: toUnit,
      }),
    });

    // Проверка ответа
    if (response.ok) {
      const data = await response.json();
      resultElement.textContent = `${parseFloat(length)} ${
        fromDropdown.value
      } = ${data.result} ${toDropdown.value}`;
      asideContent.classList.remove("close"); // Показать результат
      mainContent.classList.add("close");
    } else {
      alert(`Ошибка ${response.status}: ${response.statusText}`);
    }
  } catch (error) {
    console.error("Ошибка запроса:", error);
    alert("Произошла ошибка. Пожалуйста, попробуйте позже.");
  }
}

// Функция сброса состояния
function resetState() {
  asideContent.classList.add("close"); // Убираем класс
  mainContent.classList.remove("close");
}

// Привязка события к кнопке
convertButton.addEventListener("click", () => {
  getData(); // Запуск функции конвертации
});

resetButton.addEventListener("click", () => {
  resetState(); // Запуск функции конвертации
});
