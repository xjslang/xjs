try {
  openDb();
} catch {
  console.error("failed to open db");
}

try {
  openDb();
} catch (e) {
  console.error("failed to open db: ", e);
}

try {
  openDb();
} finally {
  console.log("cleanup");
}

try {
  openFile();
} catch {
  console.log("failed to open file");
} finally {
  console.log("cleanup");
}

// with comments
try {
  openFile();
} catch /*c1*/ ( //c2
  e
) /*c3*/ {
  console.log("failed to open file: ", e);
} finally /*c4*/ {
  console.log("cleanup");
}
