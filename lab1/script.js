// файл script.js
window.onload = function(){ 

    let a = ''
    let b = ''
    let expressionResult = ''
    let selectedOperation = null
    
    // окно вывода результата
    outputElement = document.getElementById("result")
    
    // список объектов кнопок циферблата (id которых начинается с btn_digit_)
    digitButtons = document.querySelectorAll('[id ^= "btn_digit_"]')
    
    function onDigitButtonClicked(digit) {
        if (!selectedOperation) {
            if ((digit != '.') || (digit == '.' && !a.includes(digit))) { 
                a += digit
            }
            outputElement.innerHTML = a
        } else {
            if ((digit != '.') || (digit == '.' && !b.includes(digit))) { 
                b += digit
                outputElement.innerHTML = b        
            }
        }
    }
    
    // устанавка колбек-функций на кнопки циферблата по событию нажатия
    digitButtons.forEach(button => {
        button.onclick = function() {
            const digitValue = button.innerHTML
            onDigitButtonClicked(digitValue)
        }
    });
    
    // установка колбек-функций для кнопок операций
    document.getElementById("btn_op_mult").onclick = function() { 
        if (a === '') return
        selectedOperation = 'x'
    }
    document.getElementById("btn_op_plus").onclick = function() { 
        if (a === '') return
        selectedOperation = '+'
    }
    document.getElementById("btn_op_minus").onclick = function() { 
        if (a === '') return
        selectedOperation = '-'
    }
    document.getElementById("btn_op_div").onclick = function() { 
        if (a === '') return
        selectedOperation = '/'
    }
    
    // кнопка смены знака
    document.getElementById("btn_op_sign").onclick = function() {
        if (selectedOperation) {
            if (b !== '') {
                b = (-(+b)).toString(); // Меняем знак второго операнда
                outputElement.innerHTML = b;
            }
        } else {
            if (a !== '') {
                a = (-(+a)).toString(); // Меняем знак первого операнда
                outputElement.innerHTML = a;
            }
        }
    }

        // кнопка процента
        document.getElementById("btn_op_percent").onclick = function() {
            if (a === '' || b === '' || !selectedOperation) return; // Если нет двух операндов и операции
            
            switch (selectedOperation) {
                case '+':
                    expressionResult = (+a) + ((+a) * (+b) / 100);
                    break;
                case '-':
                    expressionResult = (+a) - ((+a) * (+b) / 100);
                    break;
                case 'x':
                    expressionResult = (+a) * ((+b) / 100);
                    break;
                case '/':
                    expressionResult = (+a) / ((+b) / 100);
                    break;
            }
            
            a = expressionResult.toString();
            b = '';
            selectedOperation = null;
            outputElement.innerHTML = a;
        }

    // кнопка очищения
    document.getElementById("btn_op_clear").onclick = function() { 
        a = ''
        b = ''
        selectedOperation = ''
        expressionResult = ''
        outputElement.innerHTML = 0
    }
    
    // кнопка расчёта результата
    document.getElementById("btn_op_equal").onclick = function() { 
        if (a === '' || b === '' || !selectedOperation)
            return
            
        switch(selectedOperation) { 
            case 'x':
                expressionResult = (+a) * (+b)
                break;
            case '+':
                expressionResult = (+a) + (+b)
                break;
            case '-':
                expressionResult = (+a) - (+b)
                break;
            case '/':
                expressionResult = (+a) / (+b)
                break;
        }
        
        a = expressionResult.toString()
        b = ''
        selectedOperation = null
    
        outputElement.innerHTML = a
    }

    // Кнопка для вычисления квадратного корня
    document.getElementById("btn_op_root").onclick = function() {
        if (selectedOperation) {
            // Если выбрана операция, вычисляем корень из b
            if (b !== '') {
                b = Math.sqrt(+b).toString();
                outputElement.innerHTML = b;
            }
        } else {
            // Если операция не выбрана, вычисляем корень из a
            if (a !== '') {
                a = Math.sqrt(+a).toString();
                outputElement.innerHTML = a;
            }
        }
    };


    // Функция для вычисления факториала
    function factorial(n) {
        if (n === 0 || n === 1) return 1;
        let result = 1;
        for (let i = 2; i <= n; i++) {
            result *= i;
        }
        return result;
    }

    // Кнопка для вычисления факториала
    document.getElementById("btn_op_factorial").onclick = function() {
        if (selectedOperation) {
            // Если выбрана операция, вычисляем факториал из второго числа (b)
            if (b !== '') {
                const numB = Math.floor(+b); // Преобразуем в целое число
                if (numB >= 0 && Number.isInteger(numB)) {
                    b = factorial(numB).toString(); // Вычисляем факториал
                    outputElement.innerHTML = b;
                } else {
                    outputElement.innerHTML = 'Ошибка'; // Ошибка для нецелых чисел
                }
            }
        } else {
            // Если операция не выбрана, вычисляем факториал из первого числа (a)
            if (a !== '') {
                const numA = Math.floor(+a); // Преобразуем в целое число
                if (numA >= 0 && Number.isInteger(numA)) {
                    a = factorial(numA).toString(); // Вычисляем факториал
                    outputElement.innerHTML = a;
                } else {
                    outputElement.innerHTML = 'Ошибка'; // Ошибка для нецелых чисел
                }
            }
        }
    };

    // Кнопка для возведения в квадрат
    document.getElementById("btn_op_square").onclick = function() {
        if (selectedOperation) {
            // Если выбрана операция, возводим в квадрат число b
            if (b !== '') {
                b = Math.pow(+b, 2).toString(); // Возводим в квадрат b
                outputElement.innerHTML = b;
            }
        } else {
            // Если операция не выбрана, возводим в квадрат число a
            if (a !== '') {
                a = Math.pow(+a, 2).toString(); // Возводим в квадрат a
                outputElement.innerHTML = a;
            }
        }
    };

    // Кнопка для удаления последнего символа (Backspace)
    document.getElementById("btn_op_back").onclick = function() {
        if (selectedOperation) {
            // Если выбрана операция, удаляем последний символ из второго числа (b)
            if (b !== '') {
                b = b.slice(0, -1); // Удаляем последний символ из b
                outputElement.innerHTML = b || '0'; // Если строка пустая, показываем '0'
            }
        } else {
            // Если операция не выбрана, удаляем последний символ из первого числа (a)
            if (a !== '') {
                a = a.slice(0, -1); // Удаляем последний символ из a
                outputElement.innerHTML = a || '0'; // Если строка пустая, показываем '0'
            }
        }
    };
};