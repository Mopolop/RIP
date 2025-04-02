const readline = require('readline').createInterface({
    input: process.stdin,
    output: process.stdout
});

// Функция RLE-сжатия
function rleCompress(inputArray) {
    if (!inputArray || inputArray.length === 0) return [];
    let compressed = [];
    let currentValue = inputArray[0];
    let count = 1;

    for (let i = 1; i < inputArray.length; i++) {
        if (inputArray[i] === currentValue) {
            count++;
        } else {
            compressed.push([count, currentValue]);
            currentValue = inputArray[i];
            count = 1;
        }
    }
    compressed.push([count, currentValue]); // Добавляем последнюю группу
    return compressed;
}

// Запрос данных у пользователя
readline.question(
    'Введите элементы массива через пробел (например: 1 1 2 a a b): ',
    input => {
        const inputArray = input.trim().split(/\s+/); // Разделение по пробелам
        const compressedData = rleCompress(inputArray);
        console.log('Сжатые данные:', compressedData);
        readline.close();
    }
);