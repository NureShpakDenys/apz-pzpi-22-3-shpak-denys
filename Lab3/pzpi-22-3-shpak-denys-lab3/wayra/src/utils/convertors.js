function cToF(celsius) {
   return (celsius * 9 / 5) + 32;
}

function fToC(fahrenheit) {
   return (fahrenheit - 32) * 5 / 9;
}

function kmToMiles(km) {
   return km * 0.621371;
}

function milesToKm(miles) {
   return miles / 0.621371;
}

function kgToLbs(kg) {
   return kg * 2.20462;
}

function lbsToKg(lbs) {
   return lbs / 2.20462;
}

function kmphToMph(kmh) {
   return kmh * 0.621371;
}

function mphToKmh(mph) {
   return mph / 0.621371;
}

export default function convert(value, type, to) {
   switch (type) {
      case "temperature":
         return to === "imperial" ? cToF(value) : fToC(value);
      case "distance":
         return to === "imperial" ? kmToMiles(value) : milesToKm(value);
      case "weight":
         return to === "imperial" ? kgToLbs(value) : lbsToKg(value);
      case "speed":
         return to === "imperial" ? kmphToMph(value) : mphToKmh(value);
      default:
         return value;
   }
}

