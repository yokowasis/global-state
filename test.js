import { getState, setState } from "./index.js";

(async () => {
  // setState("test", "test");
  const a = await getState("test");
  console.log(a);
})();
