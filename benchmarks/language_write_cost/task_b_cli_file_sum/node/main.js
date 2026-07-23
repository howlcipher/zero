const fs = require("fs");

let content;
try {
  content = fs.readFileSync("names.txt", "utf8");
} catch (err) {
  console.log("Error: could not read names.txt");
  process.exit(0);
}

for (const line of content.split("\n")) {
  if (line !== "") {
    console.log("Hello,", line);
  } else {
    console.log("");
  }
}
