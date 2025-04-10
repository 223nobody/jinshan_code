exports.add = function(numbers) {
    return numbers.reduce((acc, num) => acc + num, 0);
};

exports.subtract = function(numbers) {
    return numbers.slice(1).reduce((acc, num) => acc - num, numbers[0]);
};

exports.multiply = function(numbers) {
    return numbers.reduce((acc, num) => acc * num, 1);
};

exports.divide = function(numbers) {
    return numbers.slice(1).reduce((acc, num) => {
        if (num === 0) throw new Error("被除数不能为0");
        return acc / num;
    }, numbers[0]);
};

