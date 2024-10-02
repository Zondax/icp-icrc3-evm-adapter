import * as buffer from "buffer";

// Adding missing abob function, required by ethers BigNumber package (atob is present on browser runtimes only)
global.Buffer = global.Buffer || buffer.Buffer;

if (typeof btoa === 'undefined') {
    global.btoa = function (str) {
        return new Buffer(str, 'binary').toString('base64');
    };
}

if (typeof atob === 'undefined') {
    global.atob = function (b64Encoded) {
        return new Buffer(b64Encoded, 'base64').toString('binary');
    };
}

export * from "./mappings/mappingHandlers";
