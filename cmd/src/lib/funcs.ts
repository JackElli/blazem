const randomKeyVals = "0123456789abcdef";
const keyLength = 16;

export const generateKey = () => {
    let randomKey = "";
    for (let i = 0; i < keyLength; i++) {
        let randIndex = Math.floor(Math.random() * randomKeyVals.length);
        randomKey += randomKeyVals[randIndex];
    }
    return randomKey;
};