import { defineStore } from 'pinia';
import { Local } from "../utils/storage";
import { ref } from 'vue';

// VUE3 風格
export const useCounterStore = defineStore('counter', () => {
  const count = ref(Local.get("count") || 0);
  const getCount = (): number => {
    return count.value;
  }  
  const increment = (): void => {
    count.value++;
    Local.set("count", count.value);
  }
  const decrement = (): void => {
    count.value--;
    Local.set("count", count.value);
  }
  return {getCount, increment, decrement}
});

// VUE2 風格
export const useCounterStore2 = defineStore('counter', {
  state: () => ({
    count: 0,
  }),
  actions: {
    increment() {
      this.count++;
    },
    decrement() {
      this.count--;
    },
  },
});