const getState = async (key: string): Promise<any> => {
  const res = await fetch("https://bima-global.bimasoft.workers.dev/?_=/get", {
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
    const data = await res.json();
    return data;
  } catch (error) {
    return error;
  }
};

const setState = async (key: string, value: any): Promise<any> => {
  const res = await fetch("https://bima-global.bimasoft.workers.dev/?_=/set", {
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
    const data = await res.json();
    return data;
  } catch (error) {
    return error;
  }
};

export { getState, setState };
