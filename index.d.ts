/**
 *
 * @param {string} key
 * @returns {any}
 */
declare const getState: (key: any) => Promise<any>;
/**
 *
 * @param {string} key
 * @param {any} value
 */
declare const setState: (key: any, value: any) => Promise<any>;
export { getState, setState };
