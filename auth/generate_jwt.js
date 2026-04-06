const crypto = require('crypto');
const fs = require('fs');

// Import from .env file
const SECRET_KEY = 'my_super_secret_benchmark_key';
const POOL_SIZE = 200;

function generateJWT(uID) {

    const header = Buffer.from(JSON.stringify({alg:"HS256",typ:"jwt"})).toString('base64url');
    const payload = Buffer.from(JSON.stringify({sub:uID})).toString('base64url');
    
    const unsignedToken = `${header}.${payload}`
    
    const signature = crypto
        .createHmac('sha256', SECRET_KEY)
        .update(unsignedToken)
        .digest('base64url'); 

    // 4. Return the fully assembled JWT
    return `${unsignedToken}.${signature}`;
}
    
// Create a writeStream (Internal Buffer + Open File Descripter)==>Better than fs.writfileSync
const writeStream = fs.createWriteStream('jwt.txt');

for (let i = 1; i <= POOL_SIZE; i++) {
    // Creates users like user_001, user_002
    const userId = `user_${i.toString().padStart(3,'0')}`;
    writeStream.write(generateJWT(userId)+'\n')
}

writeStream.end()

writeStream.on('finish', () => {
    console.log(`Generated ${POOL_SIZE} JWTs in jwt.txt`);
});