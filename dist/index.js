"use strict";
var __awaiter = (this && this.__awaiter) || function (thisArg, _arguments, P, generator) {
    function adopt(value) { return value instanceof P ? value : new P(function (resolve) { resolve(value); }); }
    return new (P || (P = Promise))(function (resolve, reject) {
        function fulfilled(value) { try { step(generator.next(value)); } catch (e) { reject(e); } }
        function rejected(value) { try { step(generator["throw"](value)); } catch (e) { reject(e); } }
        function step(result) { result.done ? resolve(result.value) : adopt(result.value).then(fulfilled, rejected); }
        step((generator = generator.apply(thisArg, _arguments || [])).next());
    });
};
Object.defineProperty(exports, "__esModule", { value: true });
exports.setState = exports.getState = void 0;
const getState = (key) => __awaiter(void 0, void 0, void 0, function* () {
    const res = yield fetch("https://bima-global.bimasoft.workers.dev/?_=/get", {
        method: "POST",
        headers: {
            Accept: "application/json",
            "Content-Type": "application/json",
        },
        body: JSON.stringify({
            key,
        }),
    });
    try {
        const data = yield res.json();
        return data;
    }
    catch (error) {
        return error;
    }
});
exports.getState = getState;
const setState = (key, value) => __awaiter(void 0, void 0, void 0, function* () {
    const res = yield fetch("https://bima-global.bimasoft.workers.dev/?_=/set", {
        method: "POST",
        headers: {
            Accept: "application/json",
            "Content-Type": "application/json",
        },
        body: JSON.stringify({
            key,
            value,
        }),
    });
    try {
        const data = yield res.json();
        return data;
    }
    catch (error) {
        return error;
    }
});
exports.setState = setState;
