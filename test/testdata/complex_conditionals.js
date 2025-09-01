let age = 25;
let hasLicense = true;
let hasInsurance = false;

if (age >= 18) {
    console.log("Adult");
} else {
    console.log("Minor");
}

if (age >= 18) {
    if (hasLicense) {
        console.log("Can drive");
    } else {
        console.log("Cannot drive - no license");
    }
} else {
    console.log("Cannot drive - too young");
}

if (hasLicense && hasInsurance) {
    console.log("Fully qualified driver");
} else if (hasLicense && !hasInsurance) {
    console.log("Has license but no insurance");
} else if (!hasLicense && hasInsurance) {
    console.log("Has insurance but no license");
} else {
    console.log("Neither license nor insurance");
}

let score = 85;
if (score >= 90) {
    console.log("A");
} else if (score >= 80) {
    console.log("B");
} else if (score >= 70) {
    console.log("C");
} else {
    console.log("F");
}
