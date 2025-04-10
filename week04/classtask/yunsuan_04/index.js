const calculator = require('./calculator');

// 验证参数数量，如果少于4个报错
if (process.argv.length < 4) {
    console.error('Usage: node index.js <operator> <num1> <num2> [num3...]');
    process.exit(1);
}

const operator = process.argv[2];
const numArgs = process.argv.slice(3);

// 验证数字参数
const numbers = numArgs.map(arg => {
    const num = parseFloat(arg);
    if (isNaN(num)) {
        console.error(`Error: Invalid number '${arg}'`);
        process.exit(1);
    }
    return num;
});

// 执行对应操作
try {
    let result;
    switch (operator) {
        case '+':
            result = calculator.add(numbers);
            break;
        case '-':
            result = calculator.subtract(numbers);
            break;
        case '*':
            result = calculator.multiply(numbers);
            break;
        case '/':
            result = calculator.divide(numbers);
            break;
        default:
            console.error(`Error: Invalid operator '${operator}'`);
            process.exit(1);
    }
    console.log(`Result: ${result}`);
} catch (error) {
    console.error(error.message);
    process.exit(1);
}