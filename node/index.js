import * as XLSX from "xlsx/xlsx.mjs";
import * as fs from "fs";
import { fetchAllPills } from "./fetcher.js";

const date = new Date();
const filenameDate = date.toLocaleDateString();

XLSX.set_fs(fs);

const logParsedArray = (arr) => {
    for (let i = 0; i < arr.length; i++) {
        console.table(arr[i])
    }
};

const allPills = await fetchAllPills();
// logParsedArray(allPills);

let wb = XLSX.utils.book_new();
let ws = XLSX.utils.json_to_sheet(allPills);
XLSX.utils.book_append_sheet(wb, ws, "allPills");
XLSX.writeFileXLSX(wb, `parsing-${filenameDate}.xlsx`)