const verifyToken = require('./verifyToken');
const crypto = require('crypto');
require('dotenv').config({ path: '../.env' });
const express = require('express');

// Create Redis Client
const {createClient} = require('redis');
const redisClient = createClient({
    url : "redis://localhost:6379"
})
redisClient.on("error", (err) => {
    console.log("Redis error", err);
});

// Load the Token Bucket Script
const fs = require('fs');
const luaScript = fs.readFileSync('../redis/token_bucket.lua', 'utf8');
let luaScriptSHA = null;

const app = express();

const PORT = process.env.PORT || 3000;


app.get("/", async (req,res)=>{
    // Check Wether Header Contain Bearer Token
    const authHeader = req.headers["authorization"]

    if (!authHeader || !authHeader.startsWith("Bearer ")){
        return res.status(401).send("Missing or Invalid Token Format")
    }
    const token = authHeader.split(" ")[1]

    // Verify The Token
    const decoded = verifyToken(token)
    if (!decoded){
        return res.status(401).send("Invalid Token Signature")
    }else{
        return res.status(200).send("Token Verified.")
    }

});


// Start the server
app.listen(PORT,()=>{
    console.log(`Server Running on port ${PORT}`)
})

