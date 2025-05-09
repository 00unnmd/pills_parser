import { parse } from 'node-html-parser';
import { vildegra } from "./vildegra.js";
import { antistenMV } from './antistenMV.js';
import { artrafik } from './artrafik.js';
import { vismutaTD } from './vismutaTD.js';
import { detravenol } from './detravenol.js';
import { doksiciklin } from './doksiciklin.js';
import { izislip } from './izislip.js';
import { kontrazud } from './kontrazud.js';
import { mebespalin } from './mebespalin.js';
import { meladapt } from './meladapt.js';
import { mukocil } from './mukocil.js';
import { naftoderil } from './naftoderil.js';
import { noocil } from './noocil.js';
import { orungamin } from './orungamin.js';
import { safiston } from './safiston.js';
import { skinoklir } from './skinoklir.js';
import { sumatrolid } from './sumatrolid.js';
import { tiloram } from './tiloram.js';
import { tolizor } from './tolizor.js';
import { toraksol } from './toraksol.js';
import { ulblok } from './ulblok.js';
import { fazostabil } from './fazostabil.js';
import { flebofa } from './flebofa.js';
import { freimitus } from './freimitus.js';
import { ezlor } from './ezlor.js';
import { ekziter } from './ekziter.js';
import { ekurohol } from './ekurohol.js';
import { elufor } from './elufor.js';
import { esslial } from './esslial.js';

const getDiscount = (isInStock, price, oldPrice) => {
    if (!isInStock || oldPrice === 0) {
        return 0;
    } else {
        return oldPrice - price;
    }
}

const producerNames = ["Озон Фарм, Россия", "Озон, Россия", "Озон/Озон Фарм, Россия"]

const getProducer = (listItem) => {
    const temporaryProducer = listItem.querySelector("span.listing-card__manufacturer").closest("p");
    const producerName = temporaryProducer.querySelector("a").innerHTML;

    return producerNames.some(item => item === producerName);
}

const getItems = (listArr) => {
    const temporaryArr = listArr.map(item => {
        const isInStock = Boolean(item._attrs["data-oldma-item-serp-is-in-stock"]);
        const isOzonProducer = getProducer(item);
        const oldPriceEl = item.querySelector("span.listing-card__price-old");
        const oldPrice = oldPriceEl ? oldPriceEl._attrs["data-old-price"] : 0;
        const price = Number(item._attrs["data-oldma-item-serp-price"].replace(/\s/g, ""));
    
        if (!isOzonProducer) {
            return 
        } else {
            return {
                name: item._attrs["data-oldma-item-serp-name"],
                price: isInStock ? price : 0,
                discount: getDiscount(isInStock, price, oldPrice),
            }
        }
    })

    return temporaryArr.filter(item => item != undefined);
}

const getParsedArr = (arr) => {
    return arr.map(item => {
        const listArr = parse(item.html).querySelectorAll("article.listing-card");

        return {
            name: item.name,
            items: getItems(listArr),
        }
    })
}

const logParsedArray = (arr) => {
    for (let i = 0; i < arr.length; i++) {
        console.log(arr[i].name)
        console.table(arr[i].items)
    }
};

const parsedAntistenMV = getParsedArr(antistenMV);
const parsedArtrafik = getParsedArr(artrafik);
const parsedVildegra = getParsedArr(vildegra);
const parsedVismutaTD = getParsedArr(vismutaTD);
const parsedDetravenol = getParsedArr(detravenol);
const parsedDoksiciklin = getParsedArr(doksiciklin);
const parsedIzislip = getParsedArr(izislip);
const parsedKontrazud = getParsedArr(kontrazud);
const parsedMebespalin = getParsedArr(mebespalin);
const parsedMeladapt = getParsedArr(meladapt);
const parsedMukocil = getParsedArr(mukocil);
const parsedNaftoderil = getParsedArr(naftoderil);
const parsedNoocil = getParsedArr(noocil);
const parsedOrungamin = getParsedArr(orungamin);
const parsedSafiston = getParsedArr(safiston);
const parsedSkinoklir = getParsedArr(skinoklir);
const parsedSumatrolid = getParsedArr(sumatrolid);
const parsedTiloram = getParsedArr(tiloram);
const parsedTolizor = getParsedArr(tolizor);
const parsedToraksol = getParsedArr(toraksol);
const parsedUlblok = getParsedArr(ulblok);
const parsedFazostabil = getParsedArr(fazostabil);
const parsedFlebofa = getParsedArr(flebofa);
const parsedFreimitus = getParsedArr(freimitus);
const parsedEzlor = getParsedArr(ezlor);
const parsedEkziter = getParsedArr(ekziter);
const parsedEkurohol = getParsedArr(ekurohol);
const parsedElufor = getParsedArr(elufor);
const parsedEsslial = getParsedArr(esslial);


logParsedArray(parsedAntistenMV);
logParsedArray(parsedArtrafik);
logParsedArray(parsedVildegra);
logParsedArray(parsedVismutaTD);
logParsedArray(parsedDetravenol);
logParsedArray(parsedDoksiciklin);
logParsedArray(parsedIzislip);
logParsedArray(parsedKontrazud);
logParsedArray(parsedMebespalin);
logParsedArray(parsedMeladapt);
logParsedArray(parsedMukocil);
logParsedArray(parsedNaftoderil);
logParsedArray(parsedNoocil);
logParsedArray(parsedOrungamin);
logParsedArray(parsedSafiston);
logParsedArray(parsedSkinoklir);
logParsedArray(parsedSumatrolid);
logParsedArray(parsedTiloram);
logParsedArray(parsedTolizor);
logParsedArray(parsedToraksol);
logParsedArray(parsedUlblok);
logParsedArray(parsedFazostabil);
logParsedArray(parsedFlebofa);
logParsedArray(parsedFreimitus);
logParsedArray(parsedEzlor);
logParsedArray(parsedEkziter);
logParsedArray(parsedEkurohol);
logParsedArray(parsedElufor);
logParsedArray(parsedEsslial);

// export const eaptekaPillsArray = []

/*

{
    name: "Москва и область",
    items: [],
    html: ``
},
{
    name: "Санкт-Петербург и область",
    items: [],
    html: ``
},
{
    name: "Новосибирск",
    items: [],
    html: ``
},
{
    name: "Екатеринбург",
    items: [],
    html: ``
},
{
    name: "Казань",
    items: [],
    html: ``
},
{
    name: "Нижний Новгород",
    items: [],
    html: ``
},
{
    name: "Уфа",
    items: [],
    html: ``
},
{
    name: "Ростов-на-Дону",
    items: [],
    html: ``
}

*/