/**
 * window.localStorage
 * @method set
 * @method get
 * @method remove
 * @method clear remove all
 */
export const Local = {
	set(key: string, val: any) {
		window.localStorage.setItem(key, JSON.stringify(val));
	},
	get(key: string) {
		let json: any = window.localStorage.getItem(key);
		return JSON.parse(json);
	},
	remove(key: string) {
		window.localStorage.removeItem(key);
	},
	clear() {
		window.localStorage.clear();
	},
};