const crypto = require('crypto');
require('dotenv').config({ path: '../.env' });

const SECRET_KEY = process.env.SECRET_KEY;

const verifyToken = (token) =>{
    if (!token || typeof token !== 'string') return false;
     
    // Ensure token has 3 parts only
    const temp = token.split('.');
    if (temp.length !== 3) return false;

    const unsignedToken = `${temp[0]}.${temp[1]}`;

    const signature = crypto
        .createHmac('sha256', SECRET_KEY)
        .update(unsignedToken)
        .digest('base64url'); 
    
    // timingSafeEqual compares raw binary data
    const signatureBuffer = Buffer.from(signature);
    const providedSignatureBuffer = Buffer.from(temp[2]);

    // timingSafeEqual requires buffers of the exact same length
    if (signatureBuffer.length !== providedSignatureBuffer.length) {
        return false;
    }

    return crypto.timingSafeEqual(signatureBuffer, providedSignatureBuffer);
}
module.exports = verifyToken;