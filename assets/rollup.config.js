import typescript from 'rollup-plugin-typescript';
import minify from "rollup-plugin-babel-minify";

export default {
    input: "./script/main.ts",
    plugins: [
        typescript(),
        minify(),
    ],
    output: {
        file: "../build/production/main.js",
        format: "iife",
        sourcemap: true,
    }
};